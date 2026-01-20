package commands

import (
	"github.com/robertusnegoro/k8ctl/internal/aliases"
	"github.com/spf13/cobra"
)

// InitCommands initializes all commands
func InitCommands(rootCmd *cobra.Command) {
	// Core commands
	rootCmd.AddCommand(NewGetCommand())
	rootCmd.AddCommand(NewDescribeCommand())
	rootCmd.AddCommand(NewLogsCommand())
	rootCmd.AddCommand(NewWatchCommand())

	// Context and namespace management
	rootCmd.AddCommand(NewContextCommand())
	rootCmd.AddCommand(NewNamespaceCommand())

	// Advanced commands
	rootCmd.AddCommand(NewSearchCommand())
	rootCmd.AddCommand(NewHealthCommand())
	rootCmd.AddCommand(NewPortForwardCommand())
	rootCmd.AddCommand(NewDiffCommand())

	// Utility commands
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewCompletionCommand())

	// Apply aliases
	aliases.ApplyAliases(rootCmd)
}
