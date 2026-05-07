package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envault/envault/internal/config"
	"github.com/envault/envault/internal/keystore"
	"github.com/envault/envault/internal/recipients"
	"github.com/envault/envault/internal/vault"
)

// NewRotateCmd returns a command that re-encrypts a vault file with the
// current set of recipients. This is useful after adding or removing a
// recipient so that the sealed file reflects the updated access list.
func NewRotateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate",
		Short: "Re-encrypt the vault with the current recipients",
		Long: `Decrypt the sealed vault using your private key, then re-encrypt it
for all currently registered recipients. Run this after adding or removing
a recipient to apply the access change to the vault file.`,
		RunE: runRotate,
	}

	cmd.Flags().StringP("src", "s", ".env.sealed", "Path to the sealed vault file")
	cmd.Flags().StringP("dst", "d", ".env.sealed", "Path for the rotated vault file (defaults to src)")
	return cmd
}

func runRotate(cmd *cobra.Command, args []string) error {
	src, _ := cmd.Flags().GetString("src")
	dst, _ := cmd.Flags().GetString("dst")

	cfgManager := config.New("")
	cfg, err := cfgManager.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks := keystore.New(cfg.KeystoreDir)
	privKey, err := ks.LoadPrivateKey()
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	rm := recipients.New(cfg.RecipientsFile)
	recips, err := rm.List()
	if err != nil {
		return fmt.Errorf("list recipients: %w", err)
	}
	if len(recips) == 0 {
		return fmt.Errorf("no recipients registered; add at least one with 'envault recipients add'")
	}

	v := vault.New()

	// Unseal to a temp file, then re-seal for current recipients.
	tmp, err := os.CreateTemp("", "envault-rotate-*.env")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)

	if err := v.Unseal(src, tmpPath, privKey); err != nil {
		return fmt.Errorf("unseal vault: %w", err)
	}

	if err := v.Seal(tmpPath, dst, recips); err != nil {
		return fmt.Errorf("re-seal vault: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Vault rotated for %d recipient(s) → %s\n", len(recips), dst)
	return nil
}
