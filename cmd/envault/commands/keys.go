package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/envault/envault/internal/crypto"
	"github.com/envault/envault/internal/keystore"
)

func NewKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage encryption keys",
	}

	cmd.AddCommand(newKeysGenerateCmd())
	cmd.AddCommand(newKeysShowCmd())

	return cmd
}

func newKeysGenerateCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new age key pair",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runKeysGenerate(force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing key if present")
	return cmd
}

func runKeysGenerate(force bool) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}

	ks := keystore.New(filepath.Join(homeDir, ".envault"))

	if !force {
		if _, err := ks.LoadPrivateKey(); err == nil {
			return fmt.Errorf("key already exists; use --force to overwrite")
		}
	}

	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	if err := ks.SavePrivateKey(priv); err != nil {
		return fmt.Errorf("failed to save private key: %w", err)
	}

	fmt.Printf("Generated new key pair.\nPublic key: %s\n", pub)
	fmt.Println("Private key saved to ~/.envault/key.txt")
	return nil
}

func newKeysShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show the current public key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runKeysShow()
		},
	}
}

func runKeysShow() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}

	ks := keystore.New(filepath.Join(homeDir, ".envault"))

	priv, err := ks.LoadPrivateKey()
	if err != nil {
		return fmt.Errorf("no key found; run 'envault keys generate' first: %w", err)
	}

	pub, err := crypto.PublicKeyFromPrivate(priv)
	if err != nil {
		return fmt.Errorf("failed to derive public key: %w", err)
	}

	fmt.Printf("Public key: %s\n", pub)
	return nil
}
