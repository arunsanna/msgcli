package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/skylarbpayne/msgcli/internal/auth"
	"github.com/spf13/cobra"
)

var authAddCmd = &cobra.Command{
	Use:   "add <alias>",
	Short: "Add a new account via device code authentication",
	Long: `Authenticate with a Microsoft account and store the credentials.

The alias is a friendly name you choose (e.g., "personal", "work").
You can have multiple accounts and switch between them with --account.`,
	Args: cobra.ExactArgs(1),
	RunE: runAuthAdd,
}

func init() {
	authCmd.AddCommand(authAddCmd)
}

func runAuthAdd(cmd *cobra.Command, args []string) error {
	alias := args[0]

	if IsNoInput() {
		return fmt.Errorf("auth add requires interactive input - cannot use --no-input")
	}

	config, err := auth.LoadConfig()
	if err != nil {
		return fmt.Errorf("run 'msgcli auth setup' first: %w", err)
	}

	// Check if account already exists
	existing, err := auth.LoadToken(alias)
	if err == nil && existing != nil {
		Infof("Account '%s' already exists (email: %s)", alias, existing.Email)
		Infof("Use 'msgcli auth remove %s' first to replace it", alias)
		return fmt.Errorf("account already exists")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	Infof("Starting device code authentication...")

	// Start device code flow
	dcr, err := auth.StartDeviceCodeFlow(ctx, config.ClientID)
	if err != nil {
		return fmt.Errorf("failed to start authentication: %w", err)
	}

	// Display instructions to user
	Infof("")
	Infof("To sign in, open a browser and go to:")
	Infof("  %s", dcr.VerificationURI)
	Infof("")
	Infof("Enter the code: %s", dcr.UserCode)
	Infof("")
	Infof("Waiting for authentication...")

	// Poll for token
	tr, err := auth.PollForToken(ctx, config.ClientID, dcr.DeviceCode, dcr.Interval)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	Infof("Authentication successful!")

	// Get user info
	userInfo, err := auth.GetUserInfo(ctx, tr.AccessToken)
	if err != nil {
		Infof("Warning: couldn't fetch user info: %v", err)
		userInfo = &auth.UserInfo{UserPrincipalName: "unknown"}
	}

	// Determine email to store
	email := userInfo.Mail
	if email == "" {
		email = userInfo.UserPrincipalName
	}

	// Save token
	token := &auth.TokenData{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		ExpiresAt:    time.Now().Unix() + int64(tr.ExpiresIn),
		Email:        email,
	}

	if err := auth.SaveToken(alias, token); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	Infof("Account '%s' added successfully", alias)
	Infof("Email: %s", email)
	return nil
}
