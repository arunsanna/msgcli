package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/arunsanna/msgcli/internal/auth"
	"github.com/spf13/cobra"
)

var authSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure the Azure app client ID (one-time setup)",
	Long: `Store your Azure app registration client ID.

You need to create an Azure app registration first:
1. Go to https://portal.azure.com → App registrations → New registration
2. Name: "msgcli" (or anything)
3. Supported account types: "Accounts in any organizational directory and personal Microsoft accounts"
4. Click Register
5. Copy the "Application (client) ID"
6. Go to Authentication → Enable "Allow public client flows" → Save
7. Run this command and paste the client ID`,
	RunE: runAuthSetup,
}

func init() {
	authCmd.AddCommand(authSetupCmd)
}

func runAuthSetup(cmd *cobra.Command, args []string) error {
	// Check if already configured
	existing, err := auth.LoadConfig()
	if err == nil && existing.ClientID != "" {
		Infof("Current client ID: %s", existing.ClientID)
		if IsNoInput() {
			return fmt.Errorf("already configured - use --no-input=false to reconfigure")
		}
		fmt.Fprint(os.Stderr, "Replace existing configuration? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			Infof("Keeping existing configuration")
			return nil
		}
	}

	if IsNoInput() {
		return fmt.Errorf("--no-input requires MSGCLI_CLIENT_ID environment variable")
	}

	// Check for environment variable first
	clientID := os.Getenv("MSGCLI_CLIENT_ID")
	if clientID == "" {
		fmt.Fprint(os.Stderr, "Enter your Azure app client ID: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		clientID = strings.TrimSpace(input)
	}

	if clientID == "" {
		return fmt.Errorf("client ID cannot be empty")
	}

	// Basic validation (GUID format)
	if len(clientID) != 36 || strings.Count(clientID, "-") != 4 {
		Infof("Warning: client ID doesn't look like a valid GUID")
	}

	config := &auth.Config{
		ClientID: clientID,
	}

	if err := auth.SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	Infof("Client ID saved successfully")
	Infof("Next: run 'msgcli auth add <alias>' to add an account")
	return nil
}
