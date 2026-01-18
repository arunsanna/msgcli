package cmd

import (
	"github.com/spf13/cobra"
)

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Manage Outlook calendar",
	Long:  `Commands to list, create, update, and manage calendar events.`,
}

func init() {
	rootCmd.AddCommand(calendarCmd)
}
