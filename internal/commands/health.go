package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/robertusnegoro/k8ctl/internal/k8s"
	"github.com/robertusnegoro/k8ctl/internal/output"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewHealthCommand() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "health",
		Short: "Show cluster and resource health dashboard",
		Long: `Display a comprehensive health dashboard showing:
- Cluster node status
- Pod health by namespace
- Resource status summary`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runHealth(nil, namespace)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (overrides config, empty for all)")

	return cmd
}

func runHealth(_ *cobra.Command, namespace string) error {
	client, err := k8s.GetClient()
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes client: %w", err)
	}

	// Display nodes
	fmt.Println("=== Node Status ===")
	if err := displayNodes(client); err != nil {
		fmt.Printf("Error displaying nodes: %v\n", err)
	}

	fmt.Println("\n=== Pod Health ===")
	if namespace == "" {
		// Show all namespaces
		namespaces, err := k8s.ListNamespaces()
		if err != nil {
			return fmt.Errorf("failed to list namespaces: %w", err)
		}
		for _, ns := range namespaces {
			if err := displayPodHealth(client, ns); err != nil {
				fmt.Printf("Error displaying pods in %s: %v\n", ns, err)
			}
		}
	} else {
		if err := displayPodHealth(client, namespace); err != nil {
			return fmt.Errorf("failed to display pod health: %w", err)
		}
	}

	return nil
}

func displayNodes(client *kubernetes.Clientset) error {
	nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	table := output.NewTable([]string{"NAME", "STATUS", "ROLES", "AGE", "VERSION"})
	for i := range nodes.Items {
		node := &nodes.Items[i]
		status := "Ready"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				if condition.Status != corev1.ConditionTrue {
					status = "NotReady"
				}
				break
			}
		}

		roles := []string{}
		for label := range node.Labels {
			if strings.HasPrefix(label, "node-role.kubernetes.io/") {
				roles = append(roles, strings.TrimPrefix(label, "node-role.kubernetes.io/"))
			}
		}
		roleStr := strings.Join(roles, ",")
		if roleStr == "" {
			roleStr = NoneValue
		}

		age := getAge(node.CreationTimestamp)
		version := node.Status.NodeInfo.KubeletVersion

		table.AddRow([]string{
			node.Name,
			status,
			roleStr,
			age,
			version,
		})
	}
	table.Render()
	return nil
}

func displayPodHealth(client *kubernetes.Clientset, namespace string) error {
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(pods.Items) == 0 {
		return nil
	}

	fmt.Printf("\nNamespace: %s\n", namespace)
	table := output.NewTable([]string{"NAME", "READY", "STATUS", "RESTARTS"})

	running := 0
	pending := 0
	failed := 0

	for i := range pods.Items {
		pod := &pods.Items[i]
		ready := getReadyContainers(pod)
		total := len(pod.Spec.Containers)
		status := getPodStatus(pod)
		restarts := getRestartCount(pod)

		table.AddRow([]string{
			pod.Name,
			fmt.Sprintf("%d/%d", ready, total),
			status,
			fmt.Sprintf("%d", restarts),
		})

		switch status {
		case "Running":
			running++
		case "Pending":
			pending++
		case "Failed", "Error":
			failed++
		}
	}

	table.Render()
	fmt.Printf("Summary: Running: %d, Pending: %d, Failed: %d\n", running, pending, failed)

	return nil
}
