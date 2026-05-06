package commands

import (
	"fmt"

	"github.com/envault/envault/internal/keystore"
	"github.com/envault/envault/internal/recipients"
	"github.com/spf13/cobra"
)

func NewRecipientsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recipients",
		Short: "Manage recipients who can decrypt vault files",
	}

	cmd.AddCommand(
		newRecipientsAddCmd(),
		newRecipientsListCmd(),
		newRecipientsRemoveCmd(),
	)

	return cmd
}

func newRecipientsAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [alias] [public-key]",
		Short: "Add a recipient by alias and age public key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := keystore.New("")
			if err != nil {
				return err
			}
			r := recipients.New(ks)
			if err := r.Add(args[0], args[1]); err != nil {
				return fmt.Errorf("add recipient: %w", err)
			}
			fmt.Printf("Added recipient %q\n", args[0])
			return nil
		},
	}
}

func newRecipientsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all recipients",
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := keystore.New("")
			if err != nil {
				return err
			}
			r := recipients.New(ks)
			list, err := r.List()
			if err != nil {
				return fmt.Errorf("list recipients: %w", err)
			}
			if len(list) == 0 {
				fmt.Println("No recipients configured.")
				return nil
			}
			for _, rec := range list {
				fmt.Printf("  %-20s %s\n", rec.Alias, rec.PublicKey)
			}
			return nil
		},
	}
}

func newRecipientsRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [alias]",
		Short: "Remove a recipient by alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := keystore.New("")
			if err != nil {
				return err
			}
			r := recipients.New(ks)
			if err := r.Remove(args[0]); err != nil {
				return fmt.Errorf("remove recipient: %w", err)
			}
			fmt.Printf("Removed recipient %q\n", args[0])
			return nil
		},
	}
}
