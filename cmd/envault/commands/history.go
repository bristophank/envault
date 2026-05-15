package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/envault/envault/internal/history"
	"github.com/spf13/cobra"
)

// NewHistoryCmd returns the `envault history` command.
func NewHistoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history [file]",
		Short: "Show operation history for a vault file",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runHistory,
	}
	return cmd
}

func runHistory(cmd *cobra.Command, args []string) error {
	file := ".env"
	if len(args) == 1 {
		file = args[0]
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}
	histDir := filepath.Join(homeDir, ".config", "envault", "history")
	h := history.New(histDir)

	entries, err := h.List(file)
	if err != nil {
		return fmt.Errorf("reading history: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "No history recorded for %s\n", file)
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tOPERATION\tUSER\tNOTE")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Operation,
			e.User,
			e.Note,
		)
	}
	return w.Flush()
}
