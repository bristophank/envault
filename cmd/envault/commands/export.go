package commands

import (
	"fmt"
	"os"

	"github.com/envault/envault/internal/export"
	"github.com/envault/envault/internal/keystore"
	"github.com/envault/envault/internal/vault"
	"github.com/spf13/cobra"
)

// NewExportCmd returns the export sub-command.
func NewExportCmd() *cobra.Command {
	var format string
	var sealedFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export decrypted variables to dotenv, shell, or json format",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(cmd, sealedFile, export.Format(format))
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Output format: dotenv, shell, json")
	cmd.Flags().StringVarP(&sealedFile, "input", "i", ".env.age", "Path to the sealed .env file")
	return cmd
}

func runExport(cmd *cobra.Command, sealedFile string, format export.Format) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home dir: %w", err)
	}

	ks := keystore.New(home + "/.config/envault")
	privKey, err := ks.LoadPrivateKey()
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	v := vault.New()
	content, err := v.Unseal(sealedFile, "", privKey)
	if err != nil {
		return fmt.Errorf("unseal: %w", err)
	}

	vars, err := parseEnvContent(content)
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	ex := export.New()
	out, err := ex.Export(vars, format)
	if err != nil {
		return err
	}

	fmt.Fprint(cmd.OutOrStdout(), out)
	return nil
}
