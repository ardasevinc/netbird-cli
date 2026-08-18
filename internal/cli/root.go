package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ardasevinc/netbird-cli/internal/catalog"
	"github.com/ardasevinc/netbird-cli/internal/config"
	"github.com/ardasevinc/netbird-cli/internal/exit"
	"github.com/ardasevinc/netbird-cli/internal/version"
	"github.com/spf13/cobra"
)

type commandState struct {
	json            bool
	jsonl           bool
	configPath      string
	profileName     string
	statePath       string
	timeout         time.Duration
	logLevel        string
	logWriter       io.Writer
	timeoutEnvErr   error
	logLevelEnvErr  error
	timeoutExplicit bool
}

func Execute(ctx context.Context, args []string, stdout, stderr io.Writer, info version.Info) int {
	state := &commandState{json: wantsJSON(args) || wantsJSONL(args), jsonl: wantsJSONL(args)}
	commandStdout := stdout
	var streamOutput bytes.Buffer
	if state.jsonl {
		commandStdout = &streamOutput
	}
	root := newRoot(state, commandStdout, stderr, info)
	root.SetArgs(args)
	if err := root.ExecuteContext(ctx); err != nil {
		if state.jsonl {
			_ = writeJSONLFailure(stdout, "cli", err)
		} else if state.json {
			_ = writeJSON(stdout, map[string]any{
				"schema":    "nb/v1/error",
				"ok":        false,
				"operation": "cli",
				"error": map[string]any{
					"code":    "invalid_input",
					"class":   "input",
					"message": err.Error(),
				},
			})
		}
		_, _ = fmt.Fprintln(stderr, err)
		if coded, ok := err.(interface{ ExitCode() int }); ok {
			return coded.ExitCode()
		}
		return int(exit.InvalidInput)
	}
	if state.jsonl {
		if err := writeJSONL(stdout, streamOutput.Bytes()); err != nil {
			_, _ = fmt.Fprintln(stderr, err)
			return int(exit.Internal)
		}
	}
	return 0
}

func wantsJSON(args []string) bool {
	for _, arg := range args {
		if arg == "--json" || arg == "--json=true" {
			return true
		}
	}
	return false
}

func wantsJSONL(args []string) bool {
	for _, arg := range args {
		if arg == "--jsonl" || arg == "--jsonl=true" {
			return true
		}
	}
	return false
}

func newRoot(state *commandState, stdout, stderr io.Writer, info version.Info) *cobra.Command {
	root := &cobra.Command{
		Use:           "nb",
		Short:         "an agent-first NetBird management CLI",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	state.logWriter = stderr
	root.PersistentFlags().BoolVar(&state.json, "json", state.json, "emit one machine-readable JSON document")
	root.PersistentFlags().BoolVar(&state.jsonl, "jsonl", state.jsonl, "emit a bounded machine-readable JSONL stream")
	if state.timeout == 0 {
		timeout, _, err := config.TimeoutFromEnvironment()
		if err != nil {
			state.timeoutEnvErr = err
			state.timeout = config.DefaultTimeout
		} else {
			state.timeout = timeout
		}
	} else {
		state.timeoutExplicit = true
	}
	if state.logLevel == "" {
		level, _, err := config.LogLevelFromEnvironment()
		if err != nil {
			state.logLevelEnvErr = err
			state.logLevel = config.DefaultLogLevel
		} else {
			state.logLevel = level
		}
	}
	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		if cmd.Flags().Changed("timeout") {
			if _, err := config.ParseTimeout(state.timeout.String()); err != nil {
				return err
			}
			state.timeoutExplicit = true
		} else if state.timeoutEnvErr != nil {
			return state.timeoutEnvErr
		}
		if cmd.Flags().Changed("log-level") {
			level, err := config.ParseLogLevel(state.logLevel)
			if err != nil {
				return err
			}
			state.logLevel = level
		} else if state.logLevelEnvErr != nil {
			return state.logLevelEnvErr
		}
		return nil
	}
	configPath := state.configPath
	if configPath == "" {
		configPath = config.DefaultPath()
	}
	state.configPath = configPath
	profileName := state.profileName
	if profileName == "" {
		profileName = os.Getenv("NB_PROFILE")
		if profileName == "" {
			profileName = "default"
		}
	}
	state.profileName = profileName
	statePath := state.statePath
	if statePath == "" {
		statePath = config.DefaultStatePath()
	}
	state.statePath = statePath
	root.PersistentFlags().StringVar(&state.configPath, "config", configPath, "path to the TOML configuration")
	root.PersistentFlags().StringVar(&state.profileName, "profile", profileName, "named profile to use")
	root.PersistentFlags().StringVar(&state.statePath, "state", statePath, "path to the local mutation ledger")
	root.PersistentFlags().DurationVar(&state.timeout, "timeout", state.timeout, "maximum time to wait for one server request")
	root.PersistentFlags().StringVar(&state.logLevel, "log-level", state.logLevel, "diagnostic level written to stderr")
	root.AddCommand(versionCommand(state, stdout, info))
	root.AddCommand(schemaCommand(state, stdout))
	root.AddCommand(skillsCommand(state, stdout))
	root.AddCommand(coverageCommand(state, stdout))
	root.AddCommand(apiCommand(state, stdout))
	root.AddCommand(profileCommand(state, stdout))
	root.AddCommand(capabilitiesCommand(state, stdout))
	root.AddCommand(analyzeCommand(state, stdout))
	root.AddCommand(stageCommand(state, stdout))
	root.AddCommand(applyCommand(state, stdout))
	root.AddCommand(groupsCommand(state, stdout))
	root.AddCommand(peersCommand(state, stdout))
	root.AddCommand(policiesCommand(state, stdout))
	root.AddCommand(accountsCommand(state, stdout))
	root.AddCommand(usersCommand(state, stdout))
	root.AddCommand(routesCommand(state, stdout))
	root.AddCommand(networksCommand(state, stdout))
	root.AddCommand(dnsCommand(state, stdout))
	root.AddCommand(identityProvidersCommand(state, stdout))
	root.AddCommand(postureChecksCommand(state, stdout))
	root.AddCommand(eventsCommand(state, stdout))
	root.AddCommand(setupCommand(state, stdout))
	root.AddCommand(setupKeysCommand(state, stdout))
	root.AddCommand(locationsCommand(state, stdout))
	root.AddCommand(ingressCommand(state, stdout))
	root.AddCommand(peerIngressPortsCommand(state, stdout))
	return root
}

func versionCommand(state *commandState, stdout io.Writer, info version.Info) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "show the installed nb version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/version-result", "ok": true, "operation": "version", "data": info})
			}
			_, err := fmt.Fprintln(stdout, info.String())
			return err
		},
	}
}

func schemaCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "schema", Short: "inspect versioned machine schemas"}
	command.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list available schemas",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ids, err := catalog.SchemaIDs()
			if err != nil {
				return err
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/schema-list-result", "ok": true, "operation": "schema.list", "data": map[string]any{"schemas": ids}})
			}
			for _, id := range ids {
				if _, err := fmt.Fprintln(stdout, id); err != nil {
					return err
				}
			}
			return nil
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "show <schema-id>",
		Short: "show one schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := catalog.Schema(args[0])
			if err != nil {
				return err
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/schema-show-result", "ok": true, "operation": "schema.show", "data": map[string]any{"id": args[0], "document": json.RawMessage(data)}})
			}
			_, err = stdout.Write(data)
			if err == nil && !strings.HasSuffix(string(data), "\n") {
				_, err = io.WriteString(stdout, "\n")
			}
			return err
		},
	})
	return command
}

func skillsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "skills", Short: "inspect runtime agent skills"}
	command.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list available skills",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ids, err := catalog.SkillIDs()
			if err != nil {
				return err
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/skills-list-result", "ok": true, "operation": "skills.list", "data": map[string]any{"skills": ids}})
			}
			for _, id := range ids {
				if _, err := fmt.Fprintln(stdout, id); err != nil {
					return err
				}
			}
			return nil
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "get <skill-id>",
		Short: "show one runtime skill",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := catalog.Skill(args[0])
			if err != nil {
				return err
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/skill-result", "ok": true, "operation": "skills.get", "data": map[string]any{"id": args[0], "content": string(data)}})
			}
			_, err = stdout.Write(data)
			return err
		},
	})
	return command
}

func coverageCommand(state *commandState, stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "coverage",
		Short: "show the declared management API coverage manifest",
		RunE: func(cmd *cobra.Command, _ []string) error {
			manifest, err := catalog.CoverageManifest()
			if err != nil {
				return err
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/coverage-result", "ok": true, "operation": "coverage", "data": json.RawMessage(manifest)})
			}
			_, err = stdout.Write(manifest)
			if err == nil && !strings.HasSuffix(string(manifest), "\n") {
				_, err = io.WriteString(stdout, "\n")
			}
			return err
		},
	}
}

func writeJSON(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

const streamSchema = "nb/v1/stream-event"

func writeJSONL(w io.Writer, document []byte) error {
	var envelope struct {
		Operation string          `json:"operation"`
		Data      json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(document, &envelope); err != nil {
		return fmt.Errorf("decode finite JSON result for JSONL: %w", err)
	}
	var value any
	if err := json.Unmarshal(envelope.Data, &value); err != nil {
		return fmt.Errorf("decode finite result data for JSONL: %w", err)
	}

	items, completeness := streamItems(value)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	for _, item := range items {
		if err := encoder.Encode(map[string]any{"schema": streamSchema, "type": "record", "operation": envelope.Operation, "data": item}); err != nil {
			return err
		}
	}
	return encoder.Encode(map[string]any{
		"schema":    streamSchema,
		"type":      "complete",
		"operation": envelope.Operation,
		"data": map[string]any{
			"count":        len(items),
			"completeness": completeness,
		},
	})
}

func writeJSONLFailure(w io.Writer, operation string, err error) error {
	return writeJSON(w, map[string]any{
		"schema":    streamSchema,
		"type":      "error",
		"operation": operation,
		"error": map[string]any{
			"code":    "invalid_input",
			"class":   "input",
			"message": err.Error(),
		},
	})
}

func streamItems(value any) ([]any, any) {
	completeness := any(map[string]any{"state": "complete", "reason": nil})
	object, ok := value.(map[string]any)
	if !ok {
		return []any{value}, completeness
	}
	if current, ok := object["completeness"]; ok {
		completeness = current
	}
	for _, key := range []string{"items", "records", "events", "groups", "peers", "policies", "users", "accounts", "routes", "networks", "zones", "logs", "schemas", "skills", "operations"} {
		if collection, ok := object[key].([]any); ok {
			return collection, completeness
		}
	}
	return []any{value}, completeness
}

type commandFailure struct {
	code int
	err  error
}

func (e commandFailure) Error() string { return e.err.Error() }

func (e commandFailure) Unwrap() error { return e.err }

func (e commandFailure) ExitCode() int { return e.code }

func fail(code int, err error) error {
	return commandFailure{code: code, err: err}
}
