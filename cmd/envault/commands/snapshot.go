package commands

import (
	"fmt"
	"path/filepath"

	"github.com/envault/envault/internal/snapshot"
	"github.com/spf13/cobra"
)

// NewSnapshotCmd returns the parent 'snapshot' command with sub-commands.
func NewSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage snapshots of sealed vault files",
	}
	cmd.AddCommand(newSnapshotSaveCmd())
	cmd.AddCommand(newSnapshotListCmd())
	cmd.AddCommand(newSnapshotRestoreCmd())
	return cmd
}

func newSnapshotSaveCmd() *cobra.Command {
	var input string
	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save a snapshot of the current sealed file",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := filepath.Join(filepath.Dir(input), ".envault", "snapshots")
			m := snapshot.New(dir)
			path, err := m.Save(input)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "snapshot saved: %s\n", path)
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", ".env.age", "sealed file to snapshot")
	return cmd
}

func newSnapshotListCmd() *cobra.Command {
	var input string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := filepath.Join(filepath.Dir(input), ".envault", "snapshots")
			m := snapshot.New(dir)
			paths, err := m.List()
			if err != nil {
				return err
			}
			if len(paths) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no snapshots found")
				return nil
			}
			for _, p := range paths {
				fmt.Fprintln(cmd.OutOrStdout(), p)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", ".env.age", "sealed file whose snapshots to list")
	return cmd
}

func newSnapshotRestoreCmd() *cobra.Command {
	var input, snap string
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore a snapshot over the sealed file",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := filepath.Join(filepath.Dir(input), ".envault", "snapshots")
			m := snapshot.New(dir)
			if err := m.Restore(snap, input); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "restored %s -> %s\n", snap, input)
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", ".env.age", "destination sealed file")
	cmd.Flags().StringVarP(&snap, "snapshot", "s", "", "snapshot file to restore (required)")
	_ = cmd.MarkFlagRequired("snapshot")
	return cmd
}
