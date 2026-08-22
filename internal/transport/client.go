package transport

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strings"
	"time"
)

const defaultMaxBody = 8 << 20

type Config struct {
	BaseURL   string
	Token     string
	CAFile    string
	Timeout   time.Duration
	MaxBody   int64
	HTTP      *http.Client
	LogLevel  string
	LogWriter io.Writer
}

type Client struct {
	baseURL *url.URL
	token   string
	maxBody int64
	http    *http.Client
	logger  *slog.Logger
}

// Response is the bounded response envelope returned by a management request.
// Body is only populated after the response has been fully read.
type Response struct {
	StatusCode int
	Status     string
	Body       []byte
}

// RequestError preserves the one fact that matters for mutation replay: whether
// the request may have reached the remote server.
type RequestError struct {
	Method      string
	Path        string
	StatusCode  int
	Dispatched  bool
	Description string
	Err         error
}

func (e *RequestError) Error() string {
	context := ""
	if e.StatusCode > 0 || e.Method != "" || e.Path != "" {
		context = " ("
		if e.StatusCode > 0 {
			context += fmt.Sprintf("HTTP %d", e.StatusCode)
		}
		if e.Method != "" {
			if len(context) > 2 {
				context += " "
			}
			context += e.Method
		}
		if e.Path != "" {
			if len(context) > 2 {
				context += " "
			}
			context += e.Path
		}
		context += ")"
	}
	if e.Err != nil {
		return e.Description + context + ": " + e.Err.Error()
	}
	return e.Description + context
}

func (e *RequestError) Unwrap() error { return e.Err }

func (e *RequestError) DispatchedState() bool { return e.Dispatched }

func (e *RequestError) StatusCodeState() int { return e.StatusCode }

// Origin returns the normalized origin used for server identity pinning.
func (c *Client) Origin() string { return strings.TrimRight(c.baseURL.String(), "/") }

func New(cfg Config) (*Client, error) {
	u, err := url.Parse(cfg.BaseURL)
	if err != nil || u.Scheme == "" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("base URL must be an absolute origin without credentials, query, or fragment")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 20 * time.Second
	}
	if cfg.MaxBody <= 0 {
		cfg.MaxBody = defaultMaxBody
	}
	logger, err := newLogger(cfg.LogLevel, cfg.LogWriter)
	if err != nil {
		return nil, err
	}
	client := cfg.HTTP
	if client == nil {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		if cfg.CAFile != "" {
			pool, err := loadCAPool(cfg.CAFile)
			if err != nil {
				return nil, err
			}
			transport.TLSClientConfig.RootCAs = pool
		}
		client = &http.Client{
			Transport: transport,
			Timeout:   cfg.Timeout,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
	return &Client{baseURL: u, token: cfg.Token, maxBody: cfg.MaxBody, http: client, logger: logger}, nil
}

func (c *Client) GetJSON(ctx context.Context, path string, result any) error {
	_, err := c.DoJSON(ctx, http.MethodGet, path, nil, result)
	return err
}

// DoJSON performs exactly one bounded request. It never retries. Callers that
// dispatch consequential mutations must persist their intent before invoking it.
func (c *Client) DoJSON(ctx context.Context, method, path string, request any, result any) (Response, error) {
	started := time.Now()
	c.logDebug("request", "method", method, "path", path)
	u, err := c.requestURL(path)
	if err != nil {
		return Response{}, err
	}
	var body io.Reader
	if request != nil {
		encoded, err := json.Marshal(request)
		if err != nil {
			return Response{}, fmt.Errorf("encode request %s: %w", path, err)
		}
		body = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return Response{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if request != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	dispatched := false
	trace := &httptrace.ClientTrace{WroteRequest: func(httptrace.WroteRequestInfo) { dispatched = true }}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := c.http.Do(req)
	if err != nil {
		c.logDebug("request failed", "method", method, "path", path, "dispatched", dispatched, "duration_ms", time.Since(started).Milliseconds())
		return Response{}, &RequestError{Method: method, Path: path, Dispatched: dispatched, Description: "request failed", Err: err}
	}
	c.logDebug("response", "method", method, "path", path, "status", resp.StatusCode, "duration_ms", time.Since(started).Milliseconds())
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, c.maxBody+1))
	if err != nil {
		return Response{StatusCode: resp.StatusCode, Status: resp.Status, Body: responseBody}, &RequestError{Method: method, Path: path, StatusCode: resp.StatusCode, Dispatched: true, Description: "read response failed", Err: err}
	}
	response := Response{StatusCode: resp.StatusCode, Status: resp.Status, Body: responseBody}
	if int64(len(responseBody)) > c.maxBody {
		return response, &RequestError{Method: method, Path: path, StatusCode: resp.StatusCode, Dispatched: true, Description: fmt.Sprintf("response exceeds %d-byte limit", c.maxBody)}
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return response, &RequestError{Method: method, Path: path, StatusCode: resp.StatusCode, Dispatched: true, Description: "remote rejected request"}
	}
	if result != nil && len(responseBody) != 0 {
		if err := json.Unmarshal(responseBody, result); err != nil {
			return response, &RequestError{Method: method, Path: path, StatusCode: resp.StatusCode, Dispatched: true, Description: "decode response failed", Err: err}
		}
	}
	return response, nil
}

func newLogger(level string, writer io.Writer) (*slog.Logger, error) {
	if level == "" {
		level = "error"
	}
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		return nil, fmt.Errorf("log level must be one of debug, info, warn, or error")
	}
	if writer == nil {
		return slog.New(slog.NewTextHandler(io.Discard, nil)), nil
	}
	return slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slogLevel})), nil
}

func (c *Client) logDebug(message string, args ...any) {
	c.logger.Debug(message, args...)
}

func (c *Client) requestURL(path string) (*url.URL, error) {
	if path == "" || !strings.HasPrefix(path, "/") {
		return nil, errors.New("request path must start with /")
	}
	relative, err := url.Parse(path)
	if err != nil || relative.IsAbs() || relative.Host != "" {
		return nil, errors.New("request path must be relative")
	}
	u := *c.baseURL
	u.Path = strings.TrimRight(c.baseURL.Path, "/") + relative.Path
	u.RawPath = strings.TrimRight(c.baseURL.EscapedPath(), "/") + relative.EscapedPath()
	if u.RawPath == u.Path {
		u.RawPath = ""
	}
	u.RawQuery = relative.RawQuery
	return &u, nil
}

func loadCAPool(path string) (*x509.CertPool, error) {
	pem, err := os.ReadFile(path) // #nosec G304 -- the CA file path is an explicit local profile reference.
	if err != nil {
		return nil, fmt.Errorf("read custom CA file: %w", err)
	}
	pool, err := x509.SystemCertPool()
	if err != nil {
		pool = x509.NewCertPool()
	}
	if !pool.AppendCertsFromPEM(pem) {
		return nil, errors.New("custom CA file contains no certificates")
	}
	return pool, nil
}
