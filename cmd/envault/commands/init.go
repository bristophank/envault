package commands

import (
	"fmt"

	"github.com/envault/envault/internal/keystore"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialise envault and generate a keypair",
		RunE:  runInit,
	}
}

func runInit(cmd *cobra.Command, args []string) error {
	ks, err := keystore.New("")
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	pub, err := ks.Init()
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}

	fmt.Printf("Initialised envault.\nPublic key (share with teammates):\n\n  %s\n", pub)
	return nil
}
