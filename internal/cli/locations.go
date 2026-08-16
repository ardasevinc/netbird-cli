package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func locationsCommand(state *commandState, stdout io.Writer) *cobra.Command {
	command := &cobra.Command{Use: "locations", Short: "inspect NetBird geo locations"}
	command.AddCommand(&cobra.Command{Use: "countries", Short: "list country codes", RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		countries, err := client.ListCountries(cmd.Context())
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"countries": countries, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/locations-countries-result", "ok": true, "operation": "locations.countries", "data": data})
		}
		for _, country := range countries {
			if _, err := fmt.Fprintln(stdout, country); err != nil {
				return err
			}
		}
		return nil
	}}, &cobra.Command{Use: "cities <country>", Args: cobra.ExactArgs(1), Short: "list cities for a country code", RunE: func(cmd *cobra.Command, args []string) error {
		client, err := managementClient(state)
		if err != nil {
			return err
		}
		cities, err := client.ListCountryCities(cmd.Context(), args[0])
		if err != nil {
			return fail(3, err)
		}
		data := map[string]any{"country": args[0], "cities": cities, "completeness": completeCompleteness()}
		if state.json {
			return writeJSON(stdout, map[string]any{"schema": "nb/v1/locations-cities-result", "ok": true, "operation": "locations.cities", "data": data})
		}
		for _, city := range cities {
			if _, err := fmt.Fprintln(stdout, city); err != nil {
				return err
			}
		}
		return nil
	}})
	return command
}
