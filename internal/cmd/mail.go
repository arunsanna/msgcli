package cmd

import (
	"github.com/spf13/cobra"
)

var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Manage Outlook mail",
	Long:  `Commands to list, read, send, and manage email messages.`,
}

func init() {
	rootCmd.AddCommand(mailCmd)
}
