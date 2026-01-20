package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/robertusnegoro/k8ctl/internal/config"
	"github.com/robertusnegoro/k8ctl/internal/k8s"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

func NewDiffCommand() *cobra.Command {
	var namespace1 string
	var namespace2 string

	cmd := &cobra.Command{
		Use:   "diff [resource-type] [resource-name]",
		Short: "Compare resources",
		Long: `Compare resources and show differences.
Can compare resources across namespaces or contexts.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiff(cmd, args, namespace1, namespace2)
		},
	}

	cmd.Flags().StringVarP(&namespace1, "namespace1", "", "", "First namespace (defaults to current)")
	cmd.Flags().StringVarP(&namespace2, "namespace2", "", "", "Second namespace (for cross-namespace diff)")

	return cmd
}

func runDiff(_ *cobra.Command, args []string, namespace1, namespace2 string) error {
	resourceType := args[0]
	resourceName := args[1]

	client, err := k8s.GetClient()
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes client: %w", err)
	}

	// Get namespaces
	if namespace1 == "" {
		namespace1 = config.GetCurrentNamespace()
		if namespace1 == "" {
			namespace1 = "default"
		}
	}

	if namespace2 == "" {
		namespace2 = namespace1
	}

	switch strings.ToLower(resourceType) {
	case ResourcePod, ResourcePods, ResourcePo:
		return diffPods(client, namespace1, namespace2, resourceName)
	default:
		return fmt.Errorf("resource type %s not yet implemented for diff", resourceType)
	}
}

func diffPods(client *kubernetes.Clientset, namespace1, namespace2, name string) error {
	pod1, err := client.CoreV1().Pods(namespace1).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get pod in %s: %w", namespace1, err)
	}

	pod2, err := client.CoreV1().Pods(namespace2).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get pod in %s: %w", namespace2, err)
	}

	// Convert to YAML for comparison
	yaml1, err := resourceToYAML(pod1)
	if err != nil {
		return fmt.Errorf("failed to convert pod1 to YAML: %w", err)
	}

	yaml2, err := resourceToYAML(pod2)
	if err != nil {
		return fmt.Errorf("failed to convert pod2 to YAML: %w", err)
	}

	// Simple diff (can be enhanced with proper diff library)
	if yaml1 == yaml2 {
		fmt.Println("No differences found")
		return nil
	}

	fmt.Printf("=== Differences between %s/%s and %s/%s ===\n", namespace1, name, namespace2, name)
	fmt.Println("\n--- First resource ---")
	fmt.Println(yaml1)
	fmt.Println("\n--- Second resource ---")
	fmt.Println(yaml2)

	return nil
}

func resourceToYAML(obj runtime.Object) (string, error) {
	yamlBytes, err := yaml.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}
