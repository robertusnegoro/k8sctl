// Package main provides the entry point for the k8ctl CLI application.
package main

import (
	"fmt"
	"os"

	"github.com/robertusnegoro/k8ctl/internal/commands"
	"github.com/robertusnegoro/k8ctl/internal/errors"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "k8ctl",
		Short: "Enhanced kubectl with colored tables and built-in context/namespace management",
		Long: `k8ctl is an enhanced version of kubectl with:
- Colored table output with real tables
- Built-in context and namespace switching (replaces kubectx and kubens)
- Enhanced commands for better user experience
- Advanced features like search, health dashboard, and more`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Initialize commands
	commands.InitCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		errors.PrintError(err)
		os.Exit(1)
	}
}
