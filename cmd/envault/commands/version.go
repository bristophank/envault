package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BuildVersion is set at build time via -ldflags.
var BuildVersion = "dev"

// BuildCommit is set at build time via -ldflags.
var BuildCommit = "none"

// BuildDate is set at build time via -ldflags.
var BuildDate = "unknown"

// NewVersionCmd returns a cobra command that prints version information.
func NewVersionCmd() *cobra.Command {
	var short bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `Print the version, commit hash, and build date of this envault binary.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(cmd, short)
		},
	}

	cmd.Flags().BoolVarP(&short, "short", "s", false, "Print only the version number")

	return cmd
}

func runVersion(cmd *cobra.Command, short bool) error {
	if short {
		fmt.Fprintln(cmd.OutOrStdout(), BuildVersion)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "envault %s\n", BuildVersion)
	fmt.Fprintf(cmd.OutOrStdout(), "  commit: %s\n", BuildCommit)
	fmt.Fprintf(cmd.OutOrStdout(), "  built:  %s\n", BuildDate)
	return nil
}
