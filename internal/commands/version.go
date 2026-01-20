package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new version command for displaying version information.
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Println(cmd.Root().Version)
		},
	}
}
