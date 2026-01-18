package cmd

import (
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication and accounts",
	Long:  `Commands to set up credentials, add accounts, and manage authentication.`,
}

func init() {
	rootCmd.AddCommand(authCmd)
}
