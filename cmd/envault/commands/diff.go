package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/envault/internal/crypto"
	"github.com/yourusername/envault/internal/diff"
	"github.com/yourusername/envault/internal/env"
	"github.com/yourusername/envault/internal/keystore"
)

// NewDiffCmd returns a command that decrypts two sealed files and shows
// a human-readable diff of their environment variables.
func NewDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff <file-a> <file-b>",
		Short: "Show differences between two sealed .env files",
		Args:  cobra.ExactArgs(2),
		RunE:  runDiff,
	}
	return cmd
}

func runDiff(cmd *cobra.Command, args []string) error {
	ks := keystore.New(defaultKeystoreDir())
	privKey, err := ks.LoadPrivateKey()
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	decryptFile := func(path string) (map[string]string, error) {
		ciphertext, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		plaintext, err := crypto.Decrypt(string(ciphertext), privKey)
		if err != nil {
			return nil, fmt.Errorf("decrypt %s: %w", path, err)
		}
		entries, err := env.Parse(strings.NewReader(plaintext))
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		m := make(map[string]string, len(entries))
		for _, e := range entries {
			m[e.Key] = e.Value
		}
		return m, nil
	}

	oldEnv, err := decryptFile(args[0])
	if err != nil {
		return err
	}
	newEnv, err := decryptFile(args[1])
	if err != nil {
		return err
	}

	changes := diff.Compare(oldEnv, newEnv)
	fmt.Fprintln(cmd.OutOrStdout(), diff.Summary(changes))
	return nil
}
