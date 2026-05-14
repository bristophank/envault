package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envault/envault/internal/lint"
)

// NewLintCmd returns the cobra command for linting a .env file.
func NewLintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint [file]",
		Short: "Check a .env file for common issues",
		Long: `Lint analyses a plaintext .env file and reports issues such as:
  - empty values
  - duplicate keys
  - placeholder values (changeme, TODO, FIXME, xxx)
  - keys that are not UPPER_SNAKE_CASE`,
		Args:    cobra.MaximumNArgs(1),
		RunE:    runLint,
		SilenceUsage: true,
	}
	return cmd
}

func runLint(cmd *cobra.Command, args []string) error {
	path := ".env"
	if len(args) == 1 {
		path = args[0]
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	results, err := lint.Lint(string(data))
	if err != nil {
		return fmt.Errorf("linting %s: %w", path, err)
	}

	if len(results) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "✓ %s looks good\n", path)
		return nil
	}

	for _, r := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "%s:%d: %s: %s\n", path, r.Line, r.Severity, r.Message)
	}

	// Exit with a non-zero status so CI pipelines can detect problems.
	return fmt.Errorf("%d issue(s) found in %s", len(results), path)
}
