package commands

import (
	"github.com/spf13/cobra"
)

// Root returns the top-level cobra command for envault.
func Root() *cobra.Command {
	root := &cobra.Command{
		Use:   "envault",
		Short: "Minimal secrets manager using age encryption",
		Long: `envault encrypts .env files with age so teams can share secrets safely.

Get started:
  envault init          Initialise a new vault in the current directory
  envault keys generate Generate your personal age key pair
  envault seal          Encrypt .env → .env.sealed
  envault unseal        Decrypt .env.sealed → .env
  envault rotate        Re-encrypt the vault for the current recipients
  envault recipients    Manage who can access the vault`,
	}

	root.AddCommand(
		NewInitCmd(),
		NewSealCmd(),
		NewUnsealCmd(),
		NewRecipientsCmd(),
		NewKeysCmd(),
		NewRotateCmd(),
	)

	return root
}
