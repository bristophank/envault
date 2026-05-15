package commands

import (
	"fmt"
	"os"

	"github.com/envault/envault/internal/keystore"
	"github.com/envault/envault/internal/search"
	"github.com/envault/envault/internal/vault"
	"github.com/spf13/cobra"
)

// NewSearchCmd returns the search command.
func NewSearchCmd() *cobra.Command {
	var searchValues bool

	cmd := &cobra.Command{
		Use:   "search <query> [sealed-file]",
		Short: "Search for keys or values in a sealed vault file",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			sealedFile := ".env.sealed"
			if len(args) == 2 {
				sealedFile = args[1]
			}
			return runSearch(cmd, query, sealedFile, searchValues)
		},
	}

	cmd.Flags().BoolVarP(&searchValues, "values", "v", false, "search within values instead of keys")
	return cmd
}

func runSearch(cmd *cobra.Command, query, sealedFile string, searchValues bool) error {
	ks := keystore.New(defaultKeystoreDir())
	privKey, err := ks.LoadPrivateKey()
	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	v := vault.New()
	content, err := v.Unseal(sealedFile, "", privKey)
	if err != nil {
		return fmt.Errorf("unseal %s: %w", sealedFile, err)
	}

	env, err := parseEnvContent(content)
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	s := search.New()
	var results []search.Result
	if searchValues {
		results = s.SearchValues(env, sealedFile, query)
	} else {
		results = s.SearchKeys(env, sealedFile, query)
	}

	fmt.Fprintln(os.Stdout, search.Format(results))
	return nil
}
