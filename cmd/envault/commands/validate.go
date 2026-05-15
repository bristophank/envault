package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/envault/envault/internal/env"
	"github.com/envault/envault/internal/validate"
	"github.com/spf13/cobra"
)

// NewValidateCmd returns the validate command which checks a .env file against a schema.
func NewValidateCmd() *cobra.Command {
	var schemaFile string

	cmd := &cobra.Command{
		Use:   "validate [env-file]",
		Short: "Validate a .env file against a schema",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envFile := ".env"
			if len(args) > 0 {
				envFile = args[0]
			}
			return runValidate(cmd, envFile, schemaFile)
		},
	}

	cmd.Flags().StringVarP(&schemaFile, "schema", "s", ".env.schema", "path to schema file")
	return cmd
}

func runValidate(cmd *cobra.Command, envFile, schemaFile string) error {
	envData, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	entries, err := env.Parse(strings.NewReader(string(envData)))
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	envMap := make(map[string]string, len(entries))
	for _, e := range entries {
		envMap[e.Key] = e.Value
	}

	schemaData, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("reading schema file: %w", err)
	}

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(string(schemaData)))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	schema, err := validate.ParseSchema(lines)
	if err != nil {
		return fmt.Errorf("parsing schema: %w", err)
	}

	results := validate.Validate(envMap, schema)
	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✓ validation passed")
		return nil
	}

	for _, r := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "✗ %s\n", r.Error())
	}
	return fmt.Errorf("%d validation error(s) found", len(results))
}
