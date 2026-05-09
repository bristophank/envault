package commands

import (
	"fmt"
	"path/filepath"
	"text/tabwriter"

	"github.com/envault/envault/internal/audit"
	"github.com/spf13/cobra"
)

// NewLogCmd returns the "log" subcommand that prints the audit log.
func NewLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Show the audit log of envault operations",
		RunE:  runLog,
	}
	return cmd
}

func runLog(cmd *cobra.Command, _ []string) error {
	cfgDir, err := defaultConfigDir()
	if err != nil {
		return fmt.Errorf("log: resolve config dir: %w", err)
	}

	logPath := filepath.Join(cfgDir, "audit.log")
	l, err := audit.New(logPath)
	if err != nil {
		return fmt.Errorf("log: init logger: %w", err)
	}

	entries, err := l.ReadAll()
	if err != nil {
		return fmt.Errorf("log: read entries: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No audit entries found.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tOPERATION\tTARGET\tUSER\tDETAILS")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Operation,
			e.Target,
			e.User,
			e.Details,
		)
	}
	return w.Flush()
}
