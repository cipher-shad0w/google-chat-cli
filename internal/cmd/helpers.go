package cmd

import (
	"fmt"

	"github.com/cipher-shad0w/gogchat/internal/api"
	"github.com/cipher-shad0w/gogchat/internal/auth"
	"github.com/cipher-shad0w/gogchat/internal/output"
	"github.com/spf13/viper"
)

// newAPIClient creates a new API client using the loaded configuration and
// stored OAuth2 token. It is shared by all command files in the cmd package.
func newAPIClient() (*api.Client, error) {
	clientID := Cfg.ClientID
	clientSecret := Cfg.ClientSecret

	// Fall back to the built-in defaults when the config has no credentials.
	if clientID == "" {
		clientID = auth.DefaultClientID
	}
	if clientSecret == "" {
		clientSecret = auth.DefaultClientSecret
	}

	if err := auth.ValidateCredentials(clientID, clientSecret); err != nil {
		return nil, err
	}

	tokenPath := Cfg.TokenFile
	if tokenPath == "" {
		tokenPath = auth.DefaultTokenPath()
	}

	token, err := auth.LoadToken(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("loading token (run 'gogchat auth login' first): %w", err)
	}

	httpClient := auth.HTTPClient(clientID, clientSecret, token)
	client := api.NewClient(httpClient)
	client.Verbose = viper.GetBool("verbose")
	return client, nil
}

// getFormatter returns a Formatter configured from the current CLI flags.
func getFormatter() *output.Formatter {
	return output.NewFormatter(viper.GetBool("json"), viper.GetBool("quiet"))
}
