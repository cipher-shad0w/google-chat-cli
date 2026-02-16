package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cipher-shad0w/gogchat/internal/api"
	"github.com/spf13/viper"
)

// knownErrors maps specific API error signatures to user-friendly hints.
// Each entry is checked against the error and if matched, the hint is displayed.
var knownErrors = []struct {
	// match criteria
	code        int
	status      string
	msgContains string
	// hint to display
	hint string
}{
	{
		code:        404,
		status:      "NOT_FOUND",
		msgContains: "Google Chat app not found",
		hint: `Your Google Cloud project has the Chat API enabled, but the Chat app
is not configured. This is required by Google even for user-authenticated CLI tools.

To fix this:
  1. Open: https://console.cloud.google.com/apis/api/chat.googleapis.com/hangouts-chat
  2. Fill in the required fields (App name, Avatar URL, Description)
  3. You can disable Interactive Features if you don't need bot functionality
  4. Click Save
  5. Re-run your command`,
	},
	{
		code:        403,
		status:      "PERMISSION_DENIED",
		msgContains: "insufficient authentication scopes",
		hint: `Your access token is missing the required scopes for this operation.

To fix this:
  1. Run: gogchat auth logout
  2. Run: gogchat auth login
  3. Re-authorize when prompted in your browser`,
	},
	{
		code:        403,
		status:      "PERMISSION_DENIED",
		msgContains: "Chat API has not been used",
		hint: `The Google Chat API is not enabled in your Google Cloud project.

To fix this:
  1. Open: https://console.cloud.google.com/apis/library/chat.googleapis.com
  2. Click "Enable"
  3. Wait a few minutes for the change to propagate
  4. Re-run your command`,
	},
	{
		code:        401,
		status:      "UNAUTHENTICATED",
		msgContains: "",
		hint: `Your authentication token is invalid or expired.

To fix this:
  1. Run: gogchat auth logout
  2. Run: gogchat auth login`,
	},
	{
		code:        429,
		status:      "RESOURCE_EXHAUSTED",
		msgContains: "",
		hint: `You've exceeded the API rate limit. Wait a moment and try again.
If this persists, check your quota at:
  https://console.cloud.google.com/apis/api/chat.googleapis.com/quotas`,
	},
	{
		code:        403,
		status:      "PERMISSION_DENIED",
		msgContains: "not allowed to manage this resource",
		hint: `You don't have permission to perform this operation.
If this is a Workspace admin operation, try adding --admin flag.
Make sure you have the required role in Google Workspace admin console.`,
	},
}

// findHint searches for an actionable hint matching the given API error.
func findHint(apiErr *api.APIError) string {
	for _, ke := range knownErrors {
		if ke.code != 0 && ke.code != apiErr.Code {
			continue
		}
		if ke.status != "" && ke.status != apiErr.Status {
			continue
		}
		if ke.msgContains != "" && !strings.Contains(strings.ToLower(apiErr.Message), strings.ToLower(ke.msgContains)) {
			continue
		}
		return ke.hint
	}
	return ""
}

// printRichError prints a detailed, user-friendly error message to stderr.
// It handles both regular errors and *api.APIError with extended details.
func printRichError(err error) {
	var apiErr *api.APIError
	if !errors.As(err, &apiErr) {
		// Not an API error – print as-is.
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Header line
	fmt.Fprintf(os.Stderr, "\n✗ API Error %d (%s)\n", apiErr.Code, apiErr.Status)
	fmt.Fprintf(os.Stderr, "  %s\n", apiErr.Message)

	// Check for a known error hint
	if hint := findHint(apiErr); hint != "" {
		fmt.Fprintf(os.Stderr, "\n  Hint:\n")
		for _, line := range strings.Split(hint, "\n") {
			fmt.Fprintf(os.Stderr, "  %s\n", line)
		}
	}

	// Show help links from the API response details
	links := apiErr.HelpLinks()
	if len(links) > 0 {
		fmt.Fprintf(os.Stderr, "\n  Help:\n")
		for _, link := range links {
			if link.Description != "" {
				fmt.Fprintf(os.Stderr, "  • %s\n", link.Description)
			}
			if link.URL != "" {
				fmt.Fprintf(os.Stderr, "    %s\n", link.URL)
			}
		}
	}

	// Show reason if available and verbose
	if viper.GetBool("verbose") {
		if reason := apiErr.ErrorReason(); reason != "" {
			fmt.Fprintf(os.Stderr, "\n  Reason: %s\n", reason)
		}
		// Show metadata from details
		for _, d := range apiErr.Details {
			if len(d.Metadata) > 0 {
				fmt.Fprintf(os.Stderr, "  Detail [%s]:\n", d.Type)
				for k, v := range d.Metadata {
					fmt.Fprintf(os.Stderr, "    %s: %s\n", k, v)
				}
			}
		}
		if apiErr.RawBody != "" {
			fmt.Fprintf(os.Stderr, "\n  Raw response:\n  %s\n", apiErr.RawBody)
		}
	}

	fmt.Fprintln(os.Stderr) // trailing newline for readability
}
