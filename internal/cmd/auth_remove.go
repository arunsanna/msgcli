package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/skylarbpayne/msgcli/internal/auth"
	"github.com/spf13/cobra"
)

var (
	forceRemove bool
)

var authRemoveCmd = &cobra.Command{
	Use:   "remove <alias>",
	Short: "Remove an account",
	Args:  cobra.ExactArgs(1),
	RunE:  runAuthRemove,
}

func init() {
	authRemoveCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Skip confirmation prompt")
	authCmd.AddCommand(authRemoveCmd)
}

func runAuthRemove(cmd *cobra.Command, args []string) error {
	alias := args[0]

	// Check if account exists
	token, err := auth.LoadToken(alias)
	if err != nil {
		return fmt.Errorf("account '%s' not found", alias)
	}

	if !forceRemove && !IsNoInput() {
		fmt.Fprintf(os.Stderr, "Remove account '%s' (%s)? [y/N]: ", alias, token.Email)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			Infof("Cancelled")
			return nil
		}
	}

	if err := auth.DeleteToken(alias); err != nil {
		return fmt.Errorf("failed to remove account: %w", err)
	}

	Infof("Account '%s' removed", alias)
	return nil
}
