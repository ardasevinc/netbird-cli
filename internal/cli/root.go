package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ardasevinc/netbird-cli/internal/catalog"
	"github.com/ardasevinc/netbird-cli/internal/config"
	"github.com/ardasevinc/netbird-cli/internal/exit"
	"github.com/ardasevinc/netbird-cli/internal/version"
	"github.com/spf13/cobra"
)

type commandState struct {
	json        bool
	configPath  string
	profileName string
	statePath   string
}

func Execute(ctx context.Context, args []string, stdout, stderr io.Writer, info version.Info) int {
	state := &commandState{}
	root := newRoot(state, stdout, stderr, info)
	root.SetArgs(args)
	if err := root.ExecuteContext(ctx); err != nil {
		if state.json {
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
	return 0
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
	root.PersistentFlags().BoolVar(&state.json, "json", state.json, "emit one machine-readable JSON document")
	configPath := state.configPath
	if configPath == "" {
		configPath = config.DefaultPath()
	}
	state.configPath = configPath
	profileName := state.profileName
	if profileName == "" {
		profileName = "default"
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
	root.AddCommand(versionCommand(state, stdout, info))
	root.AddCommand(schemaCommand(state, stdout))
	root.AddCommand(skillsCommand(state, stdout))
	root.AddCommand(coverageCommand(state, stdout))
	root.AddCommand(profileCommand(state, stdout))
	root.AddCommand(capabilitiesCommand(state, stdout))
	root.AddCommand(stageCommand(state, stdout))
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
