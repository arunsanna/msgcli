package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/arunsanna/msgcli/internal/auth"
	"github.com/arunsanna/msgcli/internal/graph"
	"github.com/spf13/cobra"
)

var (
	mailListFolder string
	mailListLimit  int
	mailListQuery  string
)

var mailListCmd = &cobra.Command{
	Use:   "list",
	Short: "List email messages",
	Long: `List email messages in a folder.

Well-known folder names: inbox, drafts, sentitems, deleteditems, archive, junkemail`,
	RunE: runMailList,
}

func init() {
	mailListCmd.Flags().StringVarP(&mailListFolder, "folder", "f", "inbox", "Folder to list (inbox, drafts, sentitems, etc.)")
	mailListCmd.Flags().IntVarP(&mailListLimit, "limit", "l", 25, "Maximum number of messages to return")
	mailListCmd.Flags().StringVarP(&mailListQuery, "query", "q", "", "Search query (KQL syntax)")
	mailCmd.AddCommand(mailListCmd)
}

func runMailList(cmd *cobra.Command, args []string) error {
	account, err := auth.ResolveAccount(GetAccountFlag())
	if err != nil {
		return err
	}

	client := graph.NewClient(account)
	ctx := context.Background()

	var result *graph.ListResponse[graph.Message]

	if mailListQuery != "" {
		result, err = client.SearchMessages(ctx, mailListQuery, mailListLimit)
	} else {
		params := &graph.QueryParams{
			Top:     mailListLimit,
			OrderBy: "receivedDateTime desc",
			Select:  []string{"id", "subject", "from", "receivedDateTime", "isRead", "hasAttachments", "bodyPreview"},
		}
		result, err = client.ListMessages(ctx, mailListFolder, params)
	}

	if err != nil {
		return err
	}

	format := GetOutputFormat()
	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result.Value)
	}

	// Table format
	if len(result.Value) == 0 {
		Infof("No messages found")
		return nil
	}

	for _, msg := range result.Value {
		readMarker := "•"
		if msg.IsRead {
			readMarker = " "
		}
		attachMarker := " "
		if msg.HasAttachments {
			attachMarker = "📎"
		}

		from := ""
		if msg.From != nil {
			from = msg.From.EmailAddress.Address
			if msg.From.EmailAddress.Name != "" {
				from = msg.From.EmailAddress.Name
			}
		}

		subject := msg.Subject
		if len(subject) > 50 {
			subject = subject[:47] + "..."
		}

		timeStr := msg.ReceivedDateTime.Local().Format("Jan 02 15:04")

		fmt.Printf("%s%s %-20s %-50s %s\n", readMarker, attachMarker, truncate(from, 20), subject, timeStr)
		fmt.Printf("   ID: %s\n", msg.ID)
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s + strings.Repeat(" ", maxLen-len(s))
	}
	return s[:maxLen-3] + "..."
}

// Helper for formatting relative time
func relativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(diff.Hours()/24))
	default:
		return t.Format("Jan 02")
	}
}
