package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/arunsanna/msgcli/internal/auth"
	"github.com/arunsanna/msgcli/internal/graph"
	"github.com/spf13/cobra"
)

var mailDeleteForce bool

var mailDeleteCmd = &cobra.Command{
	Use:   "delete <message-id>",
	Short: "Delete an email message",
	Args:  cobra.ExactArgs(1),
	RunE:  runMailDelete,
}

func init() {
	mailDeleteCmd.Flags().BoolVarP(&mailDeleteForce, "force", "f", false, "Skip confirmation prompt")
	mailCmd.AddCommand(mailDeleteCmd)
}

func runMailDelete(cmd *cobra.Command, args []string) error {
	messageID := args[0]

	account, err := auth.ResolveAccount(GetAccountFlag())
	if err != nil {
		return err
	}

	client := graph.NewClient(account)
	ctx := context.Background()

	if !mailDeleteForce && !IsNoInput() {
		// Fetch message to show what we're deleting
		msg, err := client.GetMessage(ctx, messageID)
		if err != nil {
			return fmt.Errorf("failed to get message: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Delete message: \"%s\"? [y/N]: ", msg.Subject)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			Infof("Cancelled")
			return nil
		}
	}

	if err := client.DeleteMessage(ctx, messageID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	Infof("Message deleted")
	return nil
}
