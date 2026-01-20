package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/robertusnegoro/k8ctl/internal/config"
	"github.com/robertusnegoro/k8ctl/internal/errors"
	"github.com/robertusnegoro/k8ctl/internal/k8s"
	"github.com/robertusnegoro/k8ctl/internal/output"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NewDescribeCommand creates a new describe command for showing resource details.
func NewDescribeCommand() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "describe [resource-type] [resource-name]",
		Short: "Show details of a specific resource",
		Long: `Show detailed information about a Kubernetes resource with
formatted YAML/JSON output and syntax highlighting.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescribe(cmd, args, namespace)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (overrides config)")

	return cmd
}

func runDescribe(_ *cobra.Command, args []string, namespace string) error {
	resourceType := args[0]
	resourceName := args[1]

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

	// Expand shortcuts (same as in get command)
	resourceShortcuts := map[string]string{
		"po":     "pods",
		"deploy": "deployments",
		"svc":    "services",
		"ing":    "ingresses",
		"cm":     "configmaps",
		"sec":    "secrets",
		"ns":     "namespaces",
		"no":     "nodes",
		"sa":     "serviceaccounts",
	}
	if expanded, ok := resourceShortcuts[strings.ToLower(resourceType)]; ok {
		resourceType = expanded
	}

	switch strings.ToLower(resourceType) {
	case ResourcePod, ResourcePods, ResourcePo:
		return describePod(client, namespace, resourceName)
	case ResourceDeployment, ResourceDeployments, ResourceDeploy:
		return describeDeployment(client, namespace, resourceName)
	case ResourceService, ResourceServices, ResourceSvc:
		return describeService(client, namespace, resourceName)
	case ResourceConfigMap, ResourceConfigMaps, ResourceCm:
		return describeConfigMap(client, namespace, resourceName)
	case ResourceSecret, ResourceSecrets, ResourceSec:
		return describeSecret(client, namespace, resourceName)
	case ResourceIngress, ResourceIngresses, ResourceIng:
		return describeIngress(client, namespace, resourceName)
	case ResourceServiceAccount, ResourceServiceAccounts, ResourceSa:
		return describeServiceAccount(client, namespace, resourceName)
	default:
		return errors.WrapError(
			fmt.Errorf("resource type %s not yet implemented", resourceType),
			fmt.Sprintf("Resource type '%s' is not yet supported for describe", resourceType),
		)
	}
}

func describePod(client kubernetes.Interface, namespace, name string) error {
	pod, err := client.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return errors.HandleKubernetesError(err, "pod", name, namespace)
	}

	// Format as YAML with sections
	return output.FormatYAML(pod)
}

func describeDeployment(client kubernetes.Interface, namespace, name string) error {
	deployment, err := client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return errors.HandleKubernetesError(err, "deployment", name, namespace)
	}

	return output.FormatYAML(deployment)
}

func describeService(client kubernetes.Interface, namespace, name string) error {
	service, err := client.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return errors.HandleKubernetesError(err, "service", name, namespace)
	}

	return output.FormatYAML(service)
}

func describeConfigMap(client kubernetes.Interface, namespace, name string) error {
	configMap, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return errors.HandleKubernetesError(err, "configmap", name, namespace)
	}

	return output.FormatYAML(configMap)
}

func describeSecret(client kubernetes.Interface, namespace, name string) error {
	secret, err := client.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return errors.HandleKubernetesError(err, "secret", name, namespace)
	}

	return output.FormatYAML(secret)
}

func describeIngress(client kubernetes.Interface, namespace, name string) error {
	ingress, err := client.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return errors.HandleKubernetesError(err, "ingress", name, namespace)
	}

	return output.FormatYAML(ingress)
}

func describeServiceAccount(client kubernetes.Interface, namespace, name string) error {
	serviceAccount, err := client.CoreV1().ServiceAccounts(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return errors.HandleKubernetesError(err, "serviceaccount", name, namespace)
	}

	return output.FormatYAML(serviceAccount)
}
