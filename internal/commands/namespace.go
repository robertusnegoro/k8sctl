package commands

import (
	"fmt"

	"github.com/robertusnegoro/k8ctl/internal/config"
	"github.com/robertusnegoro/k8ctl/internal/k8s"
	"github.com/robertusnegoro/k8ctl/internal/output"
	"github.com/spf13/cobra"
)

func NewNamespaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ns [namespace]",
		Short: "List or switch Kubernetes namespaces",
		Long: `List all available namespaces or switch to a specific namespace.
If no namespace name is provided, lists all available namespaces.
If a namespace name is provided, switches to that namespace.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runNamespace,
	}

	return cmd
}

func runNamespace(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		// List namespaces
		namespaces, err := k8s.ListNamespaces()
		if err != nil {
			return fmt.Errorf("failed to list namespaces: %w", err)
		}

		currentNS := config.GetCurrentNamespace()

		table := output.NewTable([]string{"NAMESPACE", "CURRENT"})
		for _, ns := range namespaces {
			current := ""
			if ns == currentNS {
				current = "*"
			}
			table.AddRow([]string{ns, current})
		}
		table.Render()

		return nil
	}

	// Switch namespace
	namespace := args[0]

	// Check if namespace exists
	exists, err := k8s.NamespaceExists(namespace)
	if err != nil {
		return fmt.Errorf("failed to check namespace: %w", err)
	}
	if !exists {
		return fmt.Errorf("namespace %s does not exist", namespace)
	}

	// Update k8ctl config
	if err := config.SetCurrentNamespace(namespace); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	// Reset client cache
	k8s.ResetClient()

	fmt.Printf("Switched to namespace: %s\n", namespace)
	return nil
}
