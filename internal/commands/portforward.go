package commands

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/robertusnegoro/k8ctl/internal/config"
	"github.com/robertusnegoro/k8ctl/internal/k8s"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

func NewPortForwardCommand() *cobra.Command {
	var namespace string
	var localPort string
	var remotePort string

	cmd := &cobra.Command{
		Use:   "port-forward [resource-type/resource-name]",
		Short: "Forward one or more local ports to a pod or service",
		Long: `Forward one or more local ports to a pod or service with auto-discovery.
Simplified syntax compared to kubectl port-forward.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPortForward(cmd, args, namespace, localPort, remotePort)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (overrides config)")
	cmd.Flags().StringVarP(&localPort, "local-port", "l", "", "Local port (auto-assigned if not specified)")
	cmd.Flags().StringVarP(&remotePort, "remote-port", "r", "", "Remote port (auto-discovered if not specified)")

	return cmd
}

func runPortForward(_ *cobra.Command, args []string, namespace, localPort, remotePort string) error {
	resource := args[0]

	client, err := k8s.GetClient()
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes client: %w", err)
	}

	restConfig, err := k8s.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes config: %w", err)
	}

	// Get namespace
	if namespace == "" {
		namespace = config.GetCurrentNamespace()
		if namespace == "" {
			namespace = DefaultNamespace
		}
	}

	// Parse resource (format: type/name or just name for pod)
	parts := strings.Split(resource, "/")
	var podName string
	var resourceType string

	if len(parts) == 2 {
		resourceType = parts[0]
		podName = parts[1]
	} else {
		// Assume pod if no type specified
		resourceType = "pod"
		podName = parts[0]
	}

	// Handle different resource types
	switch strings.ToLower(resourceType) {
	case ResourcePod, ResourcePods, ResourcePo:
		return portForwardToPod(client, restConfig, namespace, podName, localPort, remotePort)
	case "service", "services", "svc":
		return portForwardToService(client, restConfig, namespace, podName, localPort, remotePort)
	default:
		return fmt.Errorf("resource type %s not supported for port-forward", resourceType)
	}
}

func portForwardToPod(client *kubernetes.Clientset, config *rest.Config, namespace, podName, localPort, remotePort string) error {
	// Get pod
	pod, err := client.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get pod: %w", err)
	}

	// Auto-discover port if not specified
	if remotePort == "" {
		if len(pod.Spec.Containers) > 0 && len(pod.Spec.Containers[0].Ports) > 0 {
			remotePort = strconv.Itoa(int(pod.Spec.Containers[0].Ports[0].ContainerPort))
			fmt.Printf("Auto-discovered remote port: %s\n", remotePort)
		} else {
			return fmt.Errorf("could not auto-discover port. Please specify --remote-port")
		}
	}

	// Auto-assign local port if not specified
	if localPort == "" {
		localPort = remotePort
	}

	return startPortForward(config, namespace, pod.Name, localPort, remotePort)
}

func portForwardToService(client *kubernetes.Clientset, config *rest.Config, namespace, serviceName, localPort, remotePort string) error {
	// Get service
	service, err := client.CoreV1().Services(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get service: %w", err)
	}

	// Auto-discover port if not specified
	if remotePort == "" {
		if len(service.Spec.Ports) > 0 {
			remotePort = strconv.Itoa(int(service.Spec.Ports[0].Port))
			fmt.Printf("Auto-discovered remote port: %s\n", remotePort)
		} else {
			return fmt.Errorf("could not auto-discover port. Please specify --remote-port")
		}
	}

	// Get a pod from the service
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: service.Spec.Selector,
		}),
	})
	if err != nil || len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for service %s", serviceName)
	}

	podName := pods.Items[0].Name

	// Auto-assign local port if not specified
	if localPort == "" {
		localPort = remotePort
	}

	return startPortForward(config, namespace, podName, localPort, remotePort)
}

func startPortForward(config *rest.Config, namespace, podName, localPort, remotePort string) error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespace, podName)
	hostIP := strings.TrimPrefix(config.Host, "https://")

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return fmt.Errorf("failed to create round tripper: %w", err)
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", &url.URL{
		Scheme: "https",
		Host:   hostIP,
		Path:   path,
	})

	ports := []string{fmt.Sprintf("%s:%s", localPort, remotePort)}
	pf, err := portforward.New(dialer, ports, make(chan struct{}), make(chan struct{}), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create port forward: %w", err)
	}

	fmt.Printf("Forwarding from 127.0.0.1:%s -> %s:%s\n", localPort, podName, remotePort)
	fmt.Println("Press Ctrl+C to stop")

	errChan := make(chan error, 1)
	go func() {
		errChan <- pf.ForwardPorts()
	}()

	return <-errChan
}
