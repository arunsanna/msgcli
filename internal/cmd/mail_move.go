package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/arunsanna/msgcli/internal/auth"
	"github.com/arunsanna/msgcli/internal/graph"
	"github.com/spf13/cobra"
)

var mailMoveFolder string

var mailMoveCmd = &cobra.Command{
	Use:   "move <message-id>",
	Short: "Move an email message to a different folder",
	Long: `Move an email message to a different folder.

Use folder ID or well-known name (inbox, drafts, sentitems, deleteditems, archive, junkemail).`,
	Args: cobra.ExactArgs(1),
	RunE: runMailMove,
}

func init() {
	mailMoveCmd.Flags().StringVar(&mailMoveFolder, "folder", "", "Destination folder ID or name (required)")
	mailMoveCmd.MarkFlagRequired("folder")
	mailCmd.AddCommand(mailMoveCmd)
}

func runMailMove(cmd *cobra.Command, args []string) error {
	messageID := args[0]

	account, err := auth.ResolveAccount(GetAccountFlag())
	if err != nil {
		return err
	}

	client := graph.NewClient(account)
	ctx := context.Background()

	result, err := client.MoveMessage(ctx, messageID, mailMoveFolder)
	if err != nil {
		return fmt.Errorf("failed to move message: %w", err)
	}

	format := GetOutputFormat()
	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	Infof("Message moved to folder: %s", result.ParentFolderID)
	return nil
}
