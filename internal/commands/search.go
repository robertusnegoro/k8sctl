package commands

import (
	"context"
	"fmt"
	"strings"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/robertusnegoro/k8ctl/internal/config"
	"github.com/robertusnegoro/k8ctl/internal/errors"
	"github.com/robertusnegoro/k8ctl/internal/k8s"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NewSearchCommand creates a new search command for fuzzy searching Kubernetes resources.
func NewSearchCommand() *cobra.Command {
	var namespace string
	var resourceType string

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Fuzzy search for Kubernetes resources",
		Long: `Search for Kubernetes resources using fuzzy finding.
Supports searching across multiple resource types by name, labels, and annotations.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := ""
			if len(args) > 0 {
				query = args[0]
			}
			return runSearch(cmd, query, namespace, resourceType)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (overrides config)")
	cmd.Flags().StringVarP(&resourceType, "type", "t", "", "Resource type to search (pods, services, etc.)")

	return cmd
}

func runSearch(_ *cobra.Command, query, namespace, resourceType string) error {
	client, err := k8s.GetClient()
	if err != nil {
		return errors.WrapError(err, "Failed to connect to Kubernetes cluster")
	}

	// Get namespace
	if namespace == "" {
		namespace = config.GetCurrentNamespace()
		if namespace == "" {
			namespace = DefaultNamespace
		}
	}

	// Expand shortcuts
	resourceShortcuts := map[string]string{
		"po":     "pods",
		"deploy": "deployments",
		"svc":    "services",
		"ing":    "ingresses",
		"cm":     "configmaps",
		"sec":    "secrets",
		"sa":     "serviceaccounts",
	}
	if expanded, ok := resourceShortcuts[strings.ToLower(resourceType)]; ok {
		resourceType = expanded
	}

	// Search pods by default
	if resourceType == "" {
		resourceType = "pods"
	}

	switch strings.ToLower(resourceType) {
	case ResourcePods, ResourcePo:
		return searchPods(client, namespace, query)
	case ResourceDeployments, ResourceDeploy:
		return searchDeployments(client, namespace, query)
	case ResourceServices, ResourceSvc:
		return searchServices(client, namespace, query)
	case ResourceConfigMaps, ResourceCm:
		return searchConfigMaps(client, namespace, query)
	case ResourceSecrets, ResourceSec:
		return searchSecrets(client, namespace, query)
	default:
		return errors.WrapError(
			fmt.Errorf("resource type %s not yet implemented", resourceType),
			fmt.Sprintf("Resource type '%s' is not yet supported for search", resourceType),
		)
	}
}

func searchPods(client kubernetes.Interface, namespace, _ string) error {
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WrapError(err, "Failed to list pods")
	}

	// Prepare items for fuzzy finder
	items := make([]string, len(pods.Items))
	for i := range pods.Items {
		items[i] = pods.Items[i].Name
	}

	// Use fuzzy finder
	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return items[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
			if i < 0 || i >= len(pods.Items) {
				return ""
			}
			pod := pods.Items[i]
			return fmt.Sprintf("Pod: %s\nNamespace: %s\nStatus: %s\nReady: %d/%d",
				pod.Name,
				pod.Namespace,
				pod.Status.Phase,
				getReadyContainers(&pod),
				len(pod.Spec.Containers))
		}),
	)

	if err != nil {
		return nil // User cancelled
	}

	if idx >= 0 && idx < len(pods.Items) {
		fmt.Printf("Selected: %s\n", pods.Items[idx].Name)
	}

	return nil
}

func searchDeployments(client kubernetes.Interface, namespace, _ string) error {
	deployments, err := client.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WrapError(err, "Failed to list deployments")
	}

	// Prepare items for fuzzy finder
	items := make([]string, len(deployments.Items))
	for i := range deployments.Items {
		items[i] = deployments.Items[i].Name
	}

	// Use fuzzy finder
	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return items[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
			if i < 0 || i >= len(deployments.Items) {
				return ""
			}
			deployment := deployments.Items[i]
			return fmt.Sprintf("Deployment: %s\nNamespace: %s\nReady: %d/%d\nAvailable: %d",
				deployment.Name,
				deployment.Namespace,
				deployment.Status.ReadyReplicas,
				deployment.Status.Replicas,
				deployment.Status.AvailableReplicas)
		}),
	)

	if err != nil {
		return nil // User cancelled
	}

	if idx >= 0 && idx < len(deployments.Items) {
		fmt.Printf("Selected: %s\n", deployments.Items[idx].Name)
	}

	return nil
}

func searchServices(client kubernetes.Interface, namespace, _ string) error {
	services, err := client.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WrapError(err, "Failed to list services")
	}

	// Prepare items for fuzzy finder
	items := make([]string, len(services.Items))
	for i := range services.Items {
		items[i] = services.Items[i].Name
	}

	// Use fuzzy finder
	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return items[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
			if i < 0 || i >= len(services.Items) {
				return ""
			}
			service := services.Items[i]
			clusterIP := service.Spec.ClusterIP
			if clusterIP == "" {
				clusterIP = "<none>"
			}
			return fmt.Sprintf("Service: %s\nNamespace: %s\nType: %s\nCluster IP: %s",
				service.Name,
				service.Namespace,
				service.Spec.Type,
				clusterIP)
		}),
	)

	if err != nil {
		return nil // User cancelled
	}

	if idx >= 0 && idx < len(services.Items) {
		fmt.Printf("Selected: %s\n", services.Items[idx].Name)
	}

	return nil
}

func searchConfigMaps(client kubernetes.Interface, namespace, _ string) error {
	configMaps, err := client.CoreV1().ConfigMaps(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WrapError(err, "Failed to list configmaps")
	}

	// Prepare items for fuzzy finder
	items := make([]string, len(configMaps.Items))
	for i := range configMaps.Items {
		items[i] = configMaps.Items[i].Name
	}

	// Use fuzzy finder
	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return items[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
			if i < 0 || i >= len(configMaps.Items) {
				return ""
			}
			cm := configMaps.Items[i]
			return fmt.Sprintf("ConfigMap: %s\nNamespace: %s\nData Keys: %d",
				cm.Name,
				cm.Namespace,
				len(cm.Data))
		}),
	)

	if err != nil {
		return nil // User cancelled
	}

	if idx >= 0 && idx < len(configMaps.Items) {
		fmt.Printf("Selected: %s\n", configMaps.Items[idx].Name)
	}

	return nil
}

func searchSecrets(client kubernetes.Interface, namespace, _ string) error {
	secrets, err := client.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WrapError(err, "Failed to list secrets")
	}

	// Prepare items for fuzzy finder
	items := make([]string, len(secrets.Items))
	for i := range secrets.Items {
		items[i] = secrets.Items[i].Name
	}

	// Use fuzzy finder
	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return items[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
			if i < 0 || i >= len(secrets.Items) {
				return ""
			}
			secret := secrets.Items[i]
			secretType := string(secret.Type)
			if secretType == "" {
				secretType = SecretTypeOpaque
			}
			return fmt.Sprintf("Secret: %s\nNamespace: %s\nType: %s\nData Keys: %d",
				secret.Name,
				secret.Namespace,
				secretType,
				len(secret.Data))
		}),
	)

	if err != nil {
		return nil // User cancelled
	}

	if idx >= 0 && idx < len(secrets.Items) {
		fmt.Printf("Selected: %s\n", secrets.Items[idx].Name)
	}

	return nil
}
