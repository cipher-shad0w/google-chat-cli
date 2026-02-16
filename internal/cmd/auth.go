package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cipher-shad0w/gogchat/internal/auth"
	"github.com/cipher-shad0w/gogchat/internal/config"
)

// NewAuthCmd creates the top-level "auth" command with login, logout, and
// status subcommands.
func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication for Google Chat API",
		Long:  "Login, logout, and check authentication status for the Google Chat API.",
	}

	cmd.AddCommand(
		newLoginCmd(),
		newLogoutCmd(),
		newStatusCmd(),
	)

	return cmd
}

// resolveCredentials returns the client ID and secret to use for OAuth2.
// Resolution order: CLI flags → config file / env vars → built-in defaults.
// Returns an error if no credentials are available from any source.
func resolveCredentials(cmd *cobra.Command) (string, string, error) {
	clientID, _ := cmd.Flags().GetString("client-id")
	clientSecret, _ := cmd.Flags().GetString("client-secret")

	// Fall back to viper (config file / env vars) when flags are empty.
	if clientID == "" {
		clientID = viper.GetString("client_id")
	}
	if clientSecret == "" {
		clientSecret = viper.GetString("client_secret")
	}

	// Fall back to the built-in defaults shipped with the binary.
	if clientID == "" {
		clientID = auth.DefaultClientID
	}
	if clientSecret == "" {
		clientSecret = auth.DefaultClientSecret
	}

	if err := auth.ValidateCredentials(clientID, clientSecret); err != nil {
		return "", "", err
	}

	return clientID, clientSecret, nil
}

// tokenPath returns the path to the token file from the loaded config or the
// default location.
func tokenPath() string {
	cfg, err := config.Load()
	if err != nil || cfg.TokenFile == "" {
		return auth.DefaultTokenPath()
	}
	return cfg.TokenFile
}

// newLoginCmd creates the "auth login" subcommand.
func newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Google Chat API via OAuth2",
		Long:  "Run the interactive OAuth2 authorization flow, open a browser for consent, and save the resulting token locally.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientID, clientSecret, err := resolveCredentials(cmd)
			if err != nil {
				return err
			}

			path := tokenPath()

			// If the user is already logged in, ask before re-authenticating.
			if auth.TokenExists(path) {
				fmt.Println("You are already logged in.")
				fmt.Print("Do you want to re-authenticate? [y/N]: ")

				var answer string
				fmt.Scanln(&answer)
				if answer != "y" && answer != "Y" {
					fmt.Println("Login cancelled.")
					return nil
				}
			}

			token, err := auth.Login(clientID, clientSecret)
			if err != nil {
				return fmt.Errorf("login failed: %w", err)
			}

			if err := auth.SaveToken(path, token); err != nil {
				return fmt.Errorf("saving token: %w", err)
			}

			fmt.Println("✓ Successfully logged in!")
			fmt.Printf("  Token saved to: %s\n", path)
			return nil
		},
	}

	cmd.Flags().String("client-id", "", "Google OAuth2 client ID")
	cmd.Flags().String("client-secret", "", "Google OAuth2 client secret")

	return cmd
}

// newLogoutCmd creates the "auth logout" subcommand.
func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove stored authentication token",
		Long:  "Delete the locally stored OAuth2 token, effectively logging out.",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := tokenPath()

			if !auth.TokenExists(path) {
				fmt.Println("Not currently logged in — nothing to do.")
				return nil
			}

			if err := auth.DeleteToken(path); err != nil {
				return fmt.Errorf("logout failed: %w", err)
			}

			fmt.Println("✓ Successfully logged out.")
			return nil
		},
	}
}

// newStatusCmd creates the "auth status" subcommand.
func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current authentication status",
		Long:  "Check whether a valid OAuth2 token exists and display its expiry information.",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := tokenPath()

			if !auth.TokenExists(path) {
				fmt.Println("✗ Not logged in")
				fmt.Println("  Run 'gogchat auth login' to authenticate")
				return nil
			}

			token, err := auth.LoadToken(path)
			if err != nil {
				fmt.Println("✗ Not logged in (token file is corrupt)")
				fmt.Printf("  Error: %v\n", err)
				fmt.Println("  Run 'gogchat auth login' to re-authenticate")
				return nil
			}

			if token.Expiry.IsZero() {
				fmt.Println("✓ Logged in")
				fmt.Println("  Token expires: (no expiry set)")
				fmt.Printf("  Token file: %s\n", path)
			} else if token.Expiry.Before(time.Now()) {
				fmt.Println("✓ Logged in (token expired — will refresh on next use)")
				fmt.Printf("  Token expired: %s\n", token.Expiry.UTC().Format("2006-01-02 15:04:05 UTC"))
				fmt.Printf("  Token file: %s\n", path)
			} else {
				fmt.Println("✓ Logged in")
				fmt.Printf("  Token expires: %s\n", token.Expiry.UTC().Format("2006-01-02 15:04:05 UTC"))
				fmt.Printf("  Token file: %s\n", path)
			}

			return nil
		},
	}
}
