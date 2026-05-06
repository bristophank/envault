package commands

import (
	"fmt"

	"github.com/envault/envault/internal/keystore"
	"github.com/envault/envault/internal/vault"
	"github.com/spf13/cobra"
)

func NewSealCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seal [src] [dst]",
		Short: "Encrypt a .env file into a vault file",
		Args:  cobra.ExactArgs(2),
		RunE:  runSeal,
	}
	return cmd
}

func runSeal(cmd *cobra.Command, args []string) error {
	ks, err := keystore.New("")
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	recipients, err := ks.LoadRecipients()
	if err != nil {
		return fmt.Errorf("load recipients: %w", err)
	}

	v, err := vault.New(ks)
	if err != nil {
		return fmt.Errorf("create vault: %w", err)
	}

	if err := v.Seal(args[0], args[1], recipients); err != nil {
		return fmt.Errorf("seal: %w", err)
	}

	fmt.Printf("Sealed %s → %s\n", args[0], args[1])
	return nil
}

func NewUnsealCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unseal [src] [dst]",
		Short: "Decrypt a vault file into a .env file",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  runUnseal,
	}
	return cmd
}

func runUnseal(cmd *cobra.Command, args []string) error {
	ks, err := keystore.New("")
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	v, err := vault.New(ks)
	if err != nil {
		return fmt.Errorf("create vault: %w", err)
	}

	dst := ""
	if len(args) == 2 {
		dst = args[1]
	}

	out, err := v.Unseal(args[0], dst)
	if err != nil {
		return fmt.Errorf("unseal: %w", err)
	}

	if dst == "" {
		fmt.Print(out)
	} else {
		fmt.Printf("Unsealed %s → %s\n", args[0], dst)
	}
	return nil
}
