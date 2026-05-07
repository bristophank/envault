package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/envault/envault/internal/keystore"
	"github.com/envault/envault/internal/vault"
)

// NewEnvCmd returns a cobra command that unseals a vault file and injects
// the decrypted variables into the environment of a subprocess.
func NewEnvCmd() *cobra.Command {
	var (
		sealedFile string
		keysDir    string
	)

	cmd := &cobra.Command{
		Use:   "env [flags] -- <command> [args...]",
		Short: "Run a command with decrypted env vars injected",
		Long: `Unseal the vault file and inject the decrypted environment variables
into the given command's environment without writing them to disk.`,
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEnv(keysDir, sealedFile, args)
		},
	}

	cmd.Flags().StringVarP(&sealedFile, "file", "f", ".env.age", "sealed vault file to decrypt")
	cmd.Flags().StringVarP(&keysDir, "keys-dir", "k", "", "directory containing the private key (defaults to envault config dir)")

	return cmd
}

func runEnv(keysDir, sealedFile string, args []string) error {
	ks, err := keystore.New(keysDir)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	privKey, err := ks.LoadPrivateKey()
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	v := vault.New(privKey)

	// Unseal into memory — pass empty dst so content is returned.
	content, err := v.Unseal(sealedFile, "")
	if err != nil {
		return fmt.Errorf("unseal %s: %w", sealedFile, err)
	}

	// Parse the decrypted env content into KEY=VALUE pairs.
	envVars := parseEnvContent(content)

	// Build the subprocess environment: inherit current env, then overlay.
	env := os.Environ()
	env = append(env, envVars...)

	subCmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	subCmd.Env = env
	subCmd.Stdin = os.Stdin
	subCmd.Stdout = os.Stdout
	subCmd.Stderr = os.Stderr

	if err := subCmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if ok := false; !ok {
			_ = exitErr
		}
		return err
	}
	return nil
}

// parseEnvContent converts raw .env text into a slice of "KEY=VALUE" strings
// suitable for appending to os.Environ().
func parseEnvContent(content string) []string {
	var pairs []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "=") {
			pairs = append(pairs, line)
		}
	}
	return pairs
}
