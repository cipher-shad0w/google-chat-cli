package cmd

import (
	"fmt"
	"os"

	"github.com/cipher-shad0w/gogchat/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version is the current version of gogchat.
const Version = "0.1.0"

// Cfg holds the loaded application configuration, available after
// PersistentPreRun has executed.
var Cfg *config.Config

// usageTemplate is a customised usage template for the root command.
const usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

// rootCmd is the top-level cobra command for gogchat.
var rootCmd = &cobra.Command{
	Use:     "gogchat",
	Version: Version,
	Short:   "A CLI for the Google Chat API",
	Long: `gogchat is a command-line interface for the Google Chat API.

It provides comprehensive access to Google Chat resources including
spaces, messages, members, reactions, emoji, attachments, media,
space events, read state, and notification settings.

Authenticate once with 'gogchat auth login' and then interact with
the full Chat API from your terminal.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// If a custom config path was supplied, tell Viper about it.
		if cfgFile := viper.GetString("config"); cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		Cfg = cfg
		return nil
	},
}

func init() {
	// Persistent flags available to every sub-command.
	pflags := rootCmd.PersistentFlags()

	pflags.BoolP("json", "j", false, "Output in JSON format")
	pflags.Bool("admin", false, "Use admin access")
	pflags.BoolP("quiet", "q", false, "Suppress non-essential output")
	pflags.BoolP("verbose", "v", false, "Enable verbose/debug output")
	pflags.String("config", "", "Path to config file")

	// Bind each flag to Viper so env vars and config file values also work.
	_ = viper.BindPFlag("json", pflags.Lookup("json"))
	_ = viper.BindPFlag("admin", pflags.Lookup("admin"))
	_ = viper.BindPFlag("quiet", pflags.Lookup("quiet"))
	_ = viper.BindPFlag("verbose", pflags.Lookup("verbose"))
	_ = viper.BindPFlag("config", pflags.Lookup("config"))

	// Apply custom usage template.
	rootCmd.SetUsageTemplate(usageTemplate)

	// Register all sub-commands.
	rootCmd.AddCommand(
		NewAuthCmd(),
		NewSpacesCmd(),
		NewMessagesCmd(),
		NewMembersCmd(),
		NewReactionsCmd(),
		NewAttachmentsCmd(),
		NewEmojiCmd(),
		NewMediaCmd(),
		NewEventsCmd(),
		NewReadStateCmd(),
		NewNotificationsCmd(),
	)
}

// Execute runs the root command. It is the single entry point called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printRichError(err)
		os.Exit(1)
	}
}
