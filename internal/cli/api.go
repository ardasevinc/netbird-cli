package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/ardasevinc/netbird-cli/internal/catalog"
	"github.com/spf13/cobra"
)

func apiCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "api", Short: "invoke admitted read-only management operations"}
	command.AddCommand(apiGetCommand(state, stdout))
	return command
}

func apiGetCommand(state *commandState, stdout io.Writer) *cobra.Command {
	var queries []string
	command := &cobra.Command{
		Use:   "get <operation-id> [path-value ...]",
		Short: "invoke one manifest-backed GET operation",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			operation, err := catalog.ReadOperation(args[0])
			if err != nil {
				return fail(2, err)
			}
			path, err := resolveOperationPath(operation.Path, args[1:])
			if err != nil {
				return fail(2, err)
			}
			path, err = addQuery(path, queries)
			if err != nil {
				return fail(2, err)
			}
			client, err := managementClient(state)
			if err != nil {
				return err
			}
			payload, err := client.GetRaw(cmd.Context(), path)
			if err != nil {
				return fail(3, err)
			}
			var payloadValue any
			if err := json.Unmarshal(payload, &payloadValue); err != nil {
				payloadValue = string(payload)
			}
			data := map[string]any{
				"operation_id": operation.ID,
				"path":         path,
				"payload":      payloadValue,
				"completeness": map[string]any{"state": "unknown", "reason": "raw_endpoint"},
			}
			if state.json {
				return writeJSON(stdout, map[string]any{"schema": "nb/v1/raw-get-result", "ok": true, "operation": "api.get", "data": data})
			}
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, payload, "", "  "); err != nil {
				_, err = fmt.Fprintln(stdout, string(payload))
				return err
			}
			_, err = fmt.Fprintln(stdout, pretty.String())
			return err
		},
	}
	command.Flags().StringArrayVar(&queries, "query", nil, "query parameter in key=value form (repeatable)")
	return command
}

func resolveOperationPath(template string, values []string) (string, error) {
	parts := strings.Split(template, "/")
	valueIndex := 0
	for index, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			if valueIndex >= len(values) {
				return "", fmt.Errorf("operation path requires %d path values", countPathValues(template))
			}
			parts[index] = url.PathEscape(values[valueIndex])
			valueIndex++
		}
	}
	if valueIndex != len(values) {
		return "", fmt.Errorf("operation path accepts %d path values", valueIndex)
	}
	return strings.Join(parts, "/"), nil
}

func countPathValues(template string) int {
	count := 0
	for _, part := range strings.Split(template, "/") {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			count++
		}
	}
	return count
}

func addQuery(path string, queries []string) (string, error) {
	if len(queries) == 0 {
		return path, nil
	}
	values := url.Values{}
	for _, query := range queries {
		key, value, ok := strings.Cut(query, "=")
		if !ok || strings.TrimSpace(key) == "" {
			return "", fmt.Errorf("query must use key=value syntax")
		}
		values.Add(key, value)
	}
	return path + "?" + values.Encode(), nil
}
