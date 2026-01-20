package commands

import (
	"fmt"
	"os"

	"github.com/robertusnegoro/k8ctl/internal/config"
	"github.com/robertusnegoro/k8ctl/internal/k8s"
	"github.com/robertusnegoro/k8ctl/internal/output"
	"github.com/spf13/cobra"
)

// NewContextCommand creates a new context command for managing Kubernetes contexts.
func NewContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ctx [context-name]",
		Short: "List or switch Kubernetes contexts",
		Long: `List all available contexts or switch to a specific context.
If no context name is provided, lists all available contexts.
If a context name is provided, switches to that context.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runContext,
	}

	return cmd
}

func runContext(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		// List contexts
		contexts, err := k8s.ListContexts()
		if err != nil {
			return fmt.Errorf("failed to list contexts: %w", err)
		}

		currentContext, err := k8s.GetCurrentContext()
		if err != nil {
			return fmt.Errorf("failed to get current context: %w", err)
		}

		// Also check k8ctl config
		configContext := config.GetCurrentContext()
		if configContext != "" {
			currentContext = configContext
		}

		table := output.NewTable([]string{"CONTEXT", "CURRENT"})
		for _, ctx := range contexts {
			current := ""
			if ctx == currentContext {
				current = "*"
			}
			table.AddRow([]string{ctx, current})
		}
		table.Render()

		return nil
	}

	// Switch context
	contextName := args[0]

	// Update kubeconfig
	if err := k8s.SetContext(contextName); err != nil {
		return fmt.Errorf("failed to set context: %w", err)
	}

	// Update k8ctl config
	if err := config.SetCurrentContext(contextName); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to update k8ctl config: %v\n", err)
	}

	// Reset client cache
	k8s.ResetClient()

	fmt.Printf("Switched to context: %s\n", contextName)
	return nil
}
