package ledger

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ardasevinc/netbird-cli/internal/mutation"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

type StageInput struct {
	Profile        string
	ServerIdentity string
	AccountID      string
	Operation      string
	Request        json.RawMessage
	Before         json.RawMessage
	IntendedAfter  json.RawMessage
	Findings       []Finding
}

type Finding struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Stage struct {
	ID             string
	Revision       int
	Profile        string
	ServerIdentity string
	AccountID      string
	Operation      string
	Request        json.RawMessage
	Before         json.RawMessage
	IntendedAfter  json.RawMessage
	Digest         string
	Findings       []Finding
	Cancelled      bool
	CreatedAt      time.Time
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("ledger path is empty")
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, fmt.Errorf("create ledger directory: %w", err)
		}
		if err := os.Chmod(dir, 0o700); err != nil { // #nosec G302 -- ledger directories must be owner-only and need execute permission.
			return nil, fmt.Errorf("protect ledger directory: %w", err)
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open ledger: %w", err)
	}
	db.SetMaxOpenConns(1)
	store := &Store{db: db}
	if err := store.migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL);
INSERT INTO schema_version(version)
SELECT 1 WHERE NOT EXISTS (SELECT 1 FROM schema_version);
CREATE TABLE IF NOT EXISTS stages (
  id TEXT NOT NULL,
  revision INTEGER NOT NULL,
  profile TEXT NOT NULL,
  server_identity TEXT NOT NULL,
  account_id TEXT NOT NULL,
  operation TEXT NOT NULL,
  request_json BLOB NOT NULL,
  before_json BLOB NOT NULL,
  intended_after_json BLOB NOT NULL,
  digest TEXT NOT NULL,
  findings_json BLOB NOT NULL,
  cancelled INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  PRIMARY KEY (id, revision)
);`)
	if err != nil {
		return fmt.Errorf("migrate ledger: %w", err)
	}
	return nil
}

func (s *Store) Create(ctx context.Context, input StageInput) (Stage, error) {
	if input.Profile == "" || input.Operation == "" {
		return Stage{}, errors.New("profile and operation are required")
	}
	request, err := mutation.CanonicalJSON(input.Request)
	if err != nil {
		return Stage{}, fmt.Errorf("canonical request: %w", err)
	}
	before, err := mutation.CanonicalJSON(input.Before)
	if err != nil {
		return Stage{}, fmt.Errorf("canonical preimage: %w", err)
	}
	after, err := mutation.CanonicalJSON(input.IntendedAfter)
	if err != nil {
		return Stage{}, fmt.Errorf("canonical intended state: %w", err)
	}
	findings, err := json.Marshal(input.Findings)
	if err != nil {
		return Stage{}, fmt.Errorf("encode findings: %w", err)
	}
	digest, err := mutation.Digest(request, before, after, findings)
	if err != nil {
		return Stage{}, fmt.Errorf("digest stage: %w", err)
	}
	id, err := newID()
	if err != nil {
		return Stage{}, err
	}
	created := time.Now().UTC().Truncate(time.Microsecond)
	_, err = s.db.ExecContext(ctx, `INSERT INTO stages (id, revision, profile, server_identity, account_id, operation, request_json, before_json, intended_after_json, digest, findings_json, created_at) VALUES (?, 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, id, input.Profile, input.ServerIdentity, input.AccountID, input.Operation, request, before, after, digest, findings, created.Format(time.RFC3339Nano))
	if err != nil {
		return Stage{}, fmt.Errorf("persist stage: %w", err)
	}
	return Stage{ID: id, Revision: 1, Profile: input.Profile, ServerIdentity: input.ServerIdentity, AccountID: input.AccountID, Operation: input.Operation, Request: request, Before: before, IntendedAfter: after, Digest: digest, Findings: input.Findings, CreatedAt: created}, nil
}

func (s *Store) Get(ctx context.Context, id string, revision int) (Stage, error) {
	var stage Stage
	var findings []byte
	var created string
	var cancelled int
	err := s.db.QueryRowContext(ctx, `SELECT id, revision, profile, server_identity, account_id, operation, request_json, before_json, intended_after_json, digest, findings_json, cancelled, created_at FROM stages WHERE id = ? AND revision = ?`, id, revision).Scan(&stage.ID, &stage.Revision, &stage.Profile, &stage.ServerIdentity, &stage.AccountID, &stage.Operation, &stage.Request, &stage.Before, &stage.IntendedAfter, &stage.Digest, &findings, &cancelled, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return Stage{}, fmt.Errorf("stage %s@%d not found", id, revision)
	}
	if err != nil {
		return Stage{}, fmt.Errorf("read stage: %w", err)
	}
	if err := json.Unmarshal(findings, &stage.Findings); err != nil {
		return Stage{}, fmt.Errorf("decode stage findings: %w", err)
	}
	stage.Cancelled = cancelled != 0
	stage.CreatedAt, err = time.Parse(time.RFC3339Nano, created)
	if err != nil {
		return Stage{}, fmt.Errorf("decode stage timestamp: %w", err)
	}
	return stage, nil
}

func (s *Store) Cancel(ctx context.Context, id string, revision int) error {
	result, err := s.db.ExecContext(ctx, `UPDATE stages SET cancelled = 1 WHERE id = ? AND revision = ? AND cancelled = 0`, id, revision)
	if err != nil {
		return fmt.Errorf("cancel stage: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check stage cancellation: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("stage %s@%d not found or already cancelled", id, revision)
	}
	return nil
}

func newID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate stage id: %w", err)
	}
	return hex.EncodeToString(raw[:]), nil
}
