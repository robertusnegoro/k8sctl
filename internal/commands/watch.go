package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/robertusnegoro/k8ctl/internal/config"
	"github.com/robertusnegoro/k8ctl/internal/errors"
	"github.com/robertusnegoro/k8ctl/internal/k8s"
	"github.com/robertusnegoro/k8ctl/internal/output"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// NewWatchCommand creates a new watch command for watching Kubernetes resources.
func NewWatchCommand() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "watch [resource-type]",
		Short: "Watch resources for changes",
		Long: `Watch resources and display updates in real-time with
color changes on state transitions.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatch(cmd, args, namespace)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (overrides config)")

	return cmd
}

func runWatch(_ *cobra.Command, args []string, namespace string) error {
	resourceType := args[0]

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

	switch strings.ToLower(resourceType) {
	case ResourcePods, ResourcePo:
		return watchPods(client, namespace)
	case ResourceDeployments, ResourceDeploy:
		return watchDeployments(client, namespace)
	case ResourceServices, ResourceSvc:
		return watchServices(client, namespace)
	case ResourceConfigMaps, ResourceCm:
		return watchConfigMaps(client, namespace)
	case ResourceSecrets, ResourceSec:
		return watchSecrets(client, namespace)
	default:
		return errors.WrapError(
			fmt.Errorf("resource type %s not yet implemented", resourceType),
			fmt.Sprintf("Resource type '%s' is not yet supported for watch", resourceType),
		)
	}
}

func watchPods(client kubernetes.Interface, namespace string) error {
	watcher, err := client.CoreV1().Pods(namespace).Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WrapError(err, "Failed to create watcher")
	}
	defer watcher.Stop()

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initial display
	displayPods(client, namespace)

	fmt.Println("\nWatching for changes... (Ctrl+C to stop)")

	// Watch for events
	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return errors.WrapError(fmt.Errorf("watch channel closed"), "Watch connection lost")
			}

			switch event.Type {
			case watch.Added, watch.Modified, watch.Deleted:
				// Clear screen and redisplay
				fmt.Print("\033[2J\033[H")
				displayPods(client, namespace)
				pod, ok := event.Object.(*corev1.Pod)
				if ok {
					fmt.Printf("\nEvent: %s - %s\n", event.Type, pod.Name)
				}
			}

		case <-sigChan:
			fmt.Println("\nStopping watch...")
			return nil
		}
	}
}

func displayPods(client kubernetes.Interface, namespace string) {
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods: %v\n", err)
		return
	}

	table := output.NewTable([]string{"NAME", "READY", "STATUS", "RESTARTS", "AGE"})
	for i := range pods.Items {
		pod := &pods.Items[i]
		ready := fmt.Sprintf("%d/%d", getReadyContainers(pod), len(pod.Spec.Containers))
		status := getPodStatus(pod)
		restarts := getRestartCount(pod)
		age := getAge(pod.CreationTimestamp)

		table.AddRow([]string{
			pod.Name,
			ready,
			status,
			fmt.Sprintf("%d", restarts),
			age,
		})
	}
	table.Render()
}

func watchDeployments(client kubernetes.Interface, namespace string) error {
	watcher, err := client.AppsV1().Deployments(namespace).Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WrapError(err, "Failed to create watcher")
	}
	defer watcher.Stop()

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initial display
	displayDeployments(client, namespace)

	fmt.Println("\nWatching for changes... (Ctrl+C to stop)")

	// Watch for events
	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return errors.WrapError(fmt.Errorf("watch channel closed"), "Watch connection lost")
			}

			switch event.Type {
			case watch.Added, watch.Modified, watch.Deleted:
				// Clear screen and redisplay
				fmt.Print("\033[2J\033[H")
				displayDeployments(client, namespace)
				deployment, ok := event.Object.(*appsv1.Deployment)
				if ok {
					fmt.Printf("\nEvent: %s - %s\n", event.Type, deployment.Name)
				}
			}

		case <-sigChan:
			fmt.Println("\nStopping watch...")
			return nil
		}
	}
}

func displayDeployments(client kubernetes.Interface, namespace string) {
	deployments, err := client.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing deployments: %v\n", err)
		return
	}

	table := output.NewTable([]string{"NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE"})
	for i := range deployments.Items {
		deployment := &deployments.Items[i]
		ready := fmt.Sprintf("%d/%d", deployment.Status.ReadyReplicas, deployment.Status.Replicas)
		upToDate := fmt.Sprintf("%d", deployment.Status.UpdatedReplicas)
		available := fmt.Sprintf("%d", deployment.Status.AvailableReplicas)
		age := getAge(deployment.CreationTimestamp)

		table.AddRow([]string{
			deployment.Name,
			ready,
			upToDate,
			available,
			age,
		})
	}
	table.Render()
}

func watchServices(client kubernetes.Interface, namespace string) error {
	watcher, err := client.CoreV1().Services(namespace).Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WrapError(err, "Failed to create watcher")
	}
	defer watcher.Stop()

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initial display
	displayServices(client, namespace)

	fmt.Println("\nWatching for changes... (Ctrl+C to stop)")

	// Watch for events
	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return errors.WrapError(fmt.Errorf("watch channel closed"), "Watch connection lost")
			}

			switch event.Type {
			case watch.Added, watch.Modified, watch.Deleted:
				// Clear screen and redisplay
				fmt.Print("\033[2J\033[H")
				displayServices(client, namespace)
				service, ok := event.Object.(*corev1.Service)
				if ok {
					fmt.Printf("\nEvent: %s - %s\n", event.Type, service.Name)
				}
			}

		case <-sigChan:
			fmt.Println("\nStopping watch...")
			return nil
		}
	}
}

func displayServices(client kubernetes.Interface, namespace string) {
	services, err := client.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing services: %v\n", err)
		return
	}

	table := output.NewTable([]string{"NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE"})
	for i := range services.Items {
		service := &services.Items[i]
		serviceType := string(service.Spec.Type)
		clusterIP := service.Spec.ClusterIP
		externalIP := NoneValue
		if len(service.Spec.ExternalIPs) > 0 {
			externalIP = service.Spec.ExternalIPs[0]
		} else if len(service.Status.LoadBalancer.Ingress) > 0 {
			lb := service.Status.LoadBalancer.Ingress[0]
			if lb.IP != "" {
				externalIP = lb.IP
			} else if lb.Hostname != "" {
				externalIP = lb.Hostname
			}
		}

		ports := []string{}
		for _, port := range service.Spec.Ports {
			portStr := fmt.Sprintf("%d/%s", port.Port, port.Protocol)
			if port.NodePort != 0 {
				portStr = fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol)
			}
			ports = append(ports, portStr)
		}
		portsStr := NoneValue
		if len(ports) > 0 {
			portsStr = strings.Join(ports, ",")
		}

		age := getAge(service.CreationTimestamp)

		table.AddRow([]string{
			service.Name,
			serviceType,
			clusterIP,
			externalIP,
			portsStr,
			age,
		})
	}
	table.Render()
}

func watchConfigMaps(client kubernetes.Interface, namespace string) error {
	watcher, err := client.CoreV1().ConfigMaps(namespace).Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WrapError(err, "Failed to create watcher")
	}
	defer watcher.Stop()

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initial display
	displayConfigMaps(client, namespace)

	fmt.Println("\nWatching for changes... (Ctrl+C to stop)")

	// Watch for events
	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return errors.WrapError(fmt.Errorf("watch channel closed"), "Watch connection lost")
			}

			switch event.Type {
			case watch.Added, watch.Modified, watch.Deleted:
				// Clear screen and redisplay
				fmt.Print("\033[2J\033[H")
				displayConfigMaps(client, namespace)
				configMap, ok := event.Object.(*corev1.ConfigMap)
				if ok {
					fmt.Printf("\nEvent: %s - %s\n", event.Type, configMap.Name)
				}
			}

		case <-sigChan:
			fmt.Println("\nStopping watch...")
			return nil
		}
	}
}

func displayConfigMaps(client kubernetes.Interface, namespace string) {
	configMaps, err := client.CoreV1().ConfigMaps(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing configmaps: %v\n", err)
		return
	}

	table := output.NewTable([]string{"NAME", "DATA", "AGE"})
	for i := range configMaps.Items {
		cm := &configMaps.Items[i]
		dataCount := len(cm.Data)
		age := getAge(cm.CreationTimestamp)

		table.AddRow([]string{
			cm.Name,
			fmt.Sprintf("%d", dataCount),
			age,
		})
	}
	table.Render()
}

func watchSecrets(client kubernetes.Interface, namespace string) error {
	watcher, err := client.CoreV1().Secrets(namespace).Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.WrapError(err, "Failed to create watcher")
	}
	defer watcher.Stop()

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initial display
	displaySecrets(client, namespace)

	fmt.Println("\nWatching for changes... (Ctrl+C to stop)")

	// Watch for events
	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return errors.WrapError(fmt.Errorf("watch channel closed"), "Watch connection lost")
			}

			switch event.Type {
			case watch.Added, watch.Modified, watch.Deleted:
				// Clear screen and redisplay
				fmt.Print("\033[2J\033[H")
				displaySecrets(client, namespace)
				secret, ok := event.Object.(*corev1.Secret)
				if ok {
					fmt.Printf("\nEvent: %s - %s\n", event.Type, secret.Name)
				}
			}

		case <-sigChan:
			fmt.Println("\nStopping watch...")
			return nil
		}
	}
}

func displaySecrets(client kubernetes.Interface, namespace string) {
	secrets, err := client.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing secrets: %v\n", err)
		return
	}

	table := output.NewTable([]string{"NAME", "TYPE", "DATA", "AGE"})
	for i := range secrets.Items {
		secret := &secrets.Items[i]
		secretType := string(secret.Type)
		if secretType == "" {
			secretType = SecretTypeOpaque
		}
		dataCount := len(secret.Data)
		age := getAge(secret.CreationTimestamp)

		table.AddRow([]string{
			secret.Name,
			secretType,
			fmt.Sprintf("%d", dataCount),
			age,
		})
	}
	table.Render()
}
