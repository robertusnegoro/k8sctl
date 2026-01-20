// Package aliases provides functionality for managing command aliases in k8ctl.
package aliases

import (
	"fmt"

	"github.com/robertusnegoro/k8ctl/internal/config"
	"github.com/spf13/cobra"
)

// ApplyAliases applies command aliases from config
func ApplyAliases(rootCmd *cobra.Command) {
	cfg := config.Get()
	if cfg == nil || cfg.Aliases == nil {
		return
	}

	for alias, command := range cfg.Aliases {
		// Create alias command
		aliasCmd := &cobra.Command{
			Use:   alias,
			Short: fmt.Sprintf("Alias for '%s'", command),
			RunE: func(_ *cobra.Command, args []string) error {
				// Find the actual command
				targetCmd, _, err := rootCmd.Find([]string{command})
				if err != nil {
					return err
				}
				// Execute the target command with remaining args
				targetCmd.SetArgs(args)
				return targetCmd.Execute()
			},
		}
		rootCmd.AddCommand(aliasCmd)
	}
}
