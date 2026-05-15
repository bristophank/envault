package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/envault/envault/internal/keystore"
	"github.com/envault/envault/internal/merge"
	"github.com/envault/envault/internal/vault"
)

// NewMergeCmd returns a cobra command that merges two sealed .env files.
// The base file is updated with keys from the incoming file. By default,
// conflicts are resolved by keeping the base value; use --take-incoming to
// prefer the incoming value instead.
func NewMergeCmd() *cobra.Command {
	var (
		takeIncoming bool
		output       string
		keystoreDir  string
	)

	cmd := &cobra.Command{
		Use:   "merge <base.env.age> <incoming.env.age>",
		Short: "Merge two sealed .env files",
		Long: `Decrypt both sealed files, merge their key-value pairs, and write the
result back to a sealed file.

By default conflicts (keys present in both files with different values) are
resolved by keeping the base value. Pass --take-incoming to prefer the value
from the incoming file instead.

The merged output is written to <base.env.age> unless --output is specified.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMerge(args[0], args[1], output, keystoreDir, takeIncoming)
		},
	}

	cmd.Flags().BoolVar(&takeIncoming, "take-incoming", false, "resolve conflicts by preferring the incoming value")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write merged output to this file instead of overwriting base")
	cmd.Flags().StringVar(&keystoreDir, "keystore", "", "path to keystore directory (default: ~/.config/envault)")

	return cmd
}

func runMerge(basePath, incomingPath, outputPath, keystoreDir string, takeIncoming bool) error {
	ks, err := keystore.New(keystoreDir)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	privKey, err := ks.LoadPrivateKey()
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	recipients, err := ks.LoadRecipients()
	if err != nil {
		return fmt.Errorf("load recipients: %w", err)
	}

	v := vault.New(privKey, recipients)

	// Decrypt base file into a temp file so we can read its contents.
	baseTmp, err := os.CreateTemp("", "envault-merge-base-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(baseTmp.Name())
	baseTmp.Close()

	if err := v.Unseal(basePath, baseTmp.Name()); err != nil {
		return fmt.Errorf("unseal base file %q: %w", basePath, err)
	}

	baseContent, err := os.ReadFile(baseTmp.Name())
	if err != nil {
		return fmt.Errorf("read base content: %w", err)
	}

	// Decrypt incoming file.
	incomingTmp, err := os.CreateTemp("", "envault-merge-incoming-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(incomingTmp.Name())
	incomingTmp.Close()

	if err := v.Unseal(incomingPath, incomingTmp.Name()); err != nil {
		return fmt.Errorf("unseal incoming file %q: %w", incomingPath, err)
	}

	incomingContent, err := os.ReadFile(incomingTmp.Name())
	if err != nil {
		return fmt.Errorf("read incoming content: %w", err)
	}

	// Perform the merge.
	merged, err := merge.Merge(string(baseContent), string(incomingContent), takeIncoming)
	if err != nil {
		return fmt.Errorf("merge: %w", err)
	}

	// Write merged plaintext to a temp file, then seal it.
	mergedTmp, err := os.CreateTemp("", "envault-merge-result-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(mergedTmp.Name())

	if _, err := mergedTmp.WriteString(merged); err != nil {
		mergedTmp.Close()
		return fmt.Errorf("write merged content: %w", err)
	}
	mergedTmp.Close()

	dst := outputPath
	if dst == "" {
		dst = basePath
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if err := v.Seal(mergedTmp.Name(), dst); err != nil {
		return fmt.Errorf("seal merged file: %w", err)
	}

	fmt.Printf("merged %q + %q → %q\n", basePath, incomingPath, dst)
	return nil
}
