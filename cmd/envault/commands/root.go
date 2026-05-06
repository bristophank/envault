package commands

import (
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	root := &cobra.Command{
		Use:   "envault",
		Short: "Minimal secrets manager using age encryption",
		Long: `envault encrypts and decrypts .env files using age encryption,
enabling secure secret sharing across teams.`,
		SilenceUsage: true,
	}

	root.AddCommand(
		NewInitCmd(),
		NewSealCmd(),
		NewUnsealCmd(),
		NewRecipientsCmd(),
	)

	return root
}
