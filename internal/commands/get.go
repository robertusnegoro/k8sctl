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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var resourceShortcuts = map[string]string{
	"po":     "pods",
	"deploy": "deployments",
	"svc":    "services",
	"ing":    "ingresses",
	"cm":     "configmaps",
	"sec":    "secrets",
	"ns":     "namespaces",
	"no":     "nodes",
	"pv":     "persistentvolumes",
	"pvc":    "persistentvolumeclaims",
	"sa":     "serviceaccounts",
	"ds":     "daemonsets",
	"rs":     "replicasets",
	"sts":    "statefulsets",
	"cj":     "cronjobs",
	"job":    "jobs",
}

// NewGetCommand creates a new get command for displaying Kubernetes resources.
func NewGetCommand() *cobra.Command {
	var outputFormat string
	var namespace string

	cmd := &cobra.Command{
		Use:   "get [resource-type] [resource-name]",
		Short: "Display one or many resources",
		Long: `Display one or many resources with enhanced colored table output.
Supports all standard Kubernetes resource types and shortcuts.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd, args, outputFormat, namespace)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json, yaml")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (overrides config)")

	return cmd
}

func runGet(cmd *cobra.Command, args []string, outputFormat, namespace string) error {
	resourceType := args[0]
	resourceName := ""
	if len(args) > 1 {
		resourceName = args[1]
	}

	// Expand shortcuts
	if expanded, ok := resourceShortcuts[resourceType]; ok {
		resourceType = expanded
	}

	client, err := k8s.GetClient()
	if err != nil {
		return errors.WrapError(err, "Failed to connect to Kubernetes cluster")
	}

	// Determine if namespace was explicitly provided
	namespaceExplicitlySet := cmd.Flags().Changed("namespace")
	showNamespaceColumn := !namespaceExplicitlySet

	// Get namespace - if not explicitly set, use all namespaces
	if !namespaceExplicitlySet {
		namespace = metav1.NamespaceAll // List from all namespaces
	} else if namespace == "" {
		namespace = config.GetCurrentNamespace()
		if namespace == "" {
			namespace = DefaultNamespace
		}
	}

	// Handle namespace-scoped vs cluster-scoped resources
	switch strings.ToLower(resourceType) {
	case ResourcePods, ResourcePo:
		return getPods(client, namespace, resourceName, outputFormat, showNamespaceColumn)
	case ResourceDeployments, ResourceDeploy:
		return getDeployments(client, namespace, resourceName, outputFormat, showNamespaceColumn)
	case ResourceServices, ResourceSvc:
		return getServices(client, namespace, resourceName, outputFormat, showNamespaceColumn)
	case "configmaps", "cm":
		return getConfigMaps(client, namespace, resourceName, outputFormat, showNamespaceColumn)
	case "secrets", "sec":
		return getSecrets(client, namespace, resourceName, outputFormat, showNamespaceColumn)
	case ResourceIngresses, ResourceIng:
		return getIngresses(client, namespace, resourceName, outputFormat, showNamespaceColumn)
	case ResourceServiceAccounts, ResourceSa:
		return getServiceAccounts(client, namespace, resourceName, outputFormat, showNamespaceColumn)
	case "namespaces", "ns":
		return getNamespaces(client, resourceName, outputFormat)
	case "nodes", "no":
		return getNodes(client, resourceName, outputFormat)
	default:
		return errors.WrapError(
			fmt.Errorf("resource type %s not yet implemented", resourceType),
			fmt.Sprintf("Resource type '%s' is not yet supported", resourceType),
		)
	}
}

func getPods(client interface{}, namespace, name, outputFormat string, showNamespace bool) error {
	clientset, ok := client.(kubernetes.Interface)
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	var pods []*corev1.Pod

	if name != "" {
		pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "pod", name, namespace)
		}
		pods = []*corev1.Pod{pod}
	} else {
		podList, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "pods", "", namespace)
		}
		for i := range podList.Items {
			pods = append(pods, &podList.Items[i])
		}
	}

	if outputFormat == OutputFormatJSON {
		return output.FormatJSON(pods)
	}
	if outputFormat == OutputFormatYAML {
		return output.FormatYAML(pods)
	}

	// Table output
	headers := []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE"}
	if showNamespace {
		headers = []string{"NAMESPACE", "NAME", "READY", "STATUS", "RESTARTS", "AGE"}
	}
	table := output.NewTable(headers)
	for _, pod := range pods {
		ready := fmt.Sprintf("%d/%d", getReadyContainers(pod), len(pod.Spec.Containers))
		status := getPodStatus(pod)
		restarts := getRestartCount(pod)
		age := getAge(pod.CreationTimestamp)

		row := []string{
			pod.Name,
			ready,
			status,
			fmt.Sprintf("%d", restarts),
			age,
		}
		if showNamespace {
			row = []string{
				pod.Namespace,
				pod.Name,
				ready,
				status,
				fmt.Sprintf("%d", restarts),
				age,
			}
		}
		table.AddRow(row)
	}
	table.Render()
	return nil
}

func getDeployments(client interface{}, namespace, name, outputFormat string, showNamespace bool) error {
	clientset, ok := client.(kubernetes.Interface)
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	var deployments []interface{}

	if name != "" {
		deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "deployment", name, namespace)
		}
		deployments = []interface{}{deployment}
	} else {
		deploymentList, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "deployments", "", namespace)
		}
		for i := range deploymentList.Items {
			deployments = append(deployments, &deploymentList.Items[i])
		}
	}

	if outputFormat == OutputFormatJSON {
		return output.FormatJSON(deployments)
	}
	if outputFormat == OutputFormatYAML {
		return output.FormatYAML(deployments)
	}

	// Table output
	headers := []string{"NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE"}
	if showNamespace {
		headers = []string{"NAMESPACE", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE"}
	}
	table := output.NewTable(headers)
	for _, dep := range deployments {
		d, ok := dep.(*appsv1.Deployment)
		if !ok {
			continue
		}
		ready := fmt.Sprintf("%d/%d", d.Status.ReadyReplicas, *d.Spec.Replicas)
		upToDate := fmt.Sprintf("%d", d.Status.UpdatedReplicas)
		available := fmt.Sprintf("%d", d.Status.AvailableReplicas)
		age := getAge(d.CreationTimestamp)

		row := []string{
			d.Name,
			ready,
			upToDate,
			available,
			age,
		}
		if showNamespace {
			row = []string{
				d.Namespace,
				d.Name,
				ready,
				upToDate,
				available,
				age,
			}
		}
		table.AddRow(row)
	}
	table.Render()
	return nil
}

func getServices(client interface{}, namespace, name, outputFormat string, showNamespace bool) error {
	clientset, ok := client.(kubernetes.Interface)
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	var services []interface{}

	if name != "" {
		service, err := clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "service", name, namespace)
		}
		services = []interface{}{service}
	} else {
		serviceList, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "services", "", namespace)
		}
		for i := range serviceList.Items {
			services = append(services, &serviceList.Items[i])
		}
	}

	if outputFormat == OutputFormatJSON {
		return output.FormatJSON(services)
	}
	if outputFormat == OutputFormatYAML {
		return output.FormatYAML(services)
	}

	// Table output
	headers := []string{"NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE"}
	if showNamespace {
		headers = []string{"NAMESPACE", "NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE"}
	}
	table := output.NewTable(headers)
	for _, svc := range services {
		s, ok := svc.(*corev1.Service)
		if !ok {
			continue
		}
		ports := ""
		for i, port := range s.Spec.Ports {
			if i > 0 {
				ports += ","
			}
			ports += fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		}
		if ports == "" {
			ports = NoneValue
		}

		externalIP := NoneValue
		if len(s.Status.LoadBalancer.Ingress) > 0 {
			externalIP = s.Status.LoadBalancer.Ingress[0].IP
		}

		clusterIP := s.Spec.ClusterIP
		if clusterIP == "" {
			clusterIP = NoneValue
		}

		age := getAge(s.CreationTimestamp)

		row := []string{
			s.Name,
			string(s.Spec.Type),
			clusterIP,
			externalIP,
			ports,
			age,
		}
		if showNamespace {
			row = []string{
				s.Namespace,
				s.Name,
				string(s.Spec.Type),
				clusterIP,
				externalIP,
				ports,
				age,
			}
		}
		table.AddRow(row)
	}
	table.Render()
	return nil
}

func getNamespaces(client interface{}, name, outputFormat string) error {
	clientset, ok := client.(kubernetes.Interface)
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	var namespaces []interface{}

	if name != "" {
		namespace, err := clientset.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "namespace", name, "")
		}
		namespaces = []interface{}{namespace}
	} else {
		namespaceList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "namespaces", "", "")
		}
		for i := range namespaceList.Items {
			namespaces = append(namespaces, &namespaceList.Items[i])
		}
	}

	if outputFormat == OutputFormatJSON {
		return output.FormatJSON(namespaces)
	}
	if outputFormat == OutputFormatYAML {
		return output.FormatYAML(namespaces)
	}

	// Table output
	table := output.NewTable([]string{"NAME", "STATUS", "AGE"})
	for _, ns := range namespaces {
		n, ok := ns.(*corev1.Namespace)
		if !ok {
			continue
		}
		status := "Active"
		if n.Status.Phase != "" {
			status = string(n.Status.Phase)
		}
		age := getAge(n.CreationTimestamp)

		table.AddRow([]string{
			n.Name,
			status,
			age,
		})
	}
	table.Render()
	return nil
}

func getNodes(client interface{}, name, outputFormat string) error {
	clientset, ok := client.(kubernetes.Interface)
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	var nodes []interface{}

	if name != "" {
		node, err := clientset.CoreV1().Nodes().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "node", name, "")
		}
		nodes = []interface{}{node}
	} else {
		nodeList, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "nodes", "", "")
		}
		for i := range nodeList.Items {
			nodes = append(nodes, &nodeList.Items[i])
		}
	}

	if outputFormat == OutputFormatJSON {
		return output.FormatJSON(nodes)
	}
	if outputFormat == OutputFormatYAML {
		return output.FormatYAML(nodes)
	}

	// Table output
	table := output.NewTable([]string{"NAME", "STATUS", "ROLES", "AGE", "VERSION"})
	for _, node := range nodes {
		n, ok := node.(*corev1.Node)
		if !ok {
			continue
		}
		status := StatusReady
		for _, condition := range n.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				if condition.Status != corev1.ConditionTrue {
					status = StatusNotReady
				}
				break
			}
		}

		roles := []string{}
		for label := range n.Labels {
			if strings.HasPrefix(label, "node-role.kubernetes.io/") {
				roles = append(roles, strings.TrimPrefix(label, "node-role.kubernetes.io/"))
			}
		}
		roleStr := strings.Join(roles, ",")
		if roleStr == "" {
			roleStr = NoneValue
		}

		age := getAge(n.CreationTimestamp)
		version := n.Status.NodeInfo.KubeletVersion

		table.AddRow([]string{
			n.Name,
			status,
			roleStr,
			age,
			version,
		})
	}
	table.Render()
	return nil
}

func getConfigMaps(client interface{}, namespace, name, outputFormat string, showNamespace bool) error {
	clientset, ok := client.(kubernetes.Interface)
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	var configMaps []interface{}

	if name != "" {
		cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "configmap", name, namespace)
		}
		configMaps = []interface{}{cm}
	} else {
		cmList, err := clientset.CoreV1().ConfigMaps(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "configmaps", "", namespace)
		}
		for i := range cmList.Items {
			configMaps = append(configMaps, &cmList.Items[i])
		}
	}

	if outputFormat == OutputFormatJSON {
		return output.FormatJSON(configMaps)
	}
	if outputFormat == OutputFormatYAML {
		return output.FormatYAML(configMaps)
	}

	// Table output
	headers := []string{"NAME", "DATA", "AGE"}
	if showNamespace {
		headers = []string{"NAMESPACE", "NAME", "DATA", "AGE"}
	}
	table := output.NewTable(headers)
	for _, cm := range configMaps {
		c, ok := cm.(*corev1.ConfigMap)
		if !ok {
			continue
		}
		dataCount := len(c.Data)
		age := getAge(c.CreationTimestamp)

		row := []string{
			c.Name,
			fmt.Sprintf("%d", dataCount),
			age,
		}
		if showNamespace {
			row = []string{
				c.Namespace,
				c.Name,
				fmt.Sprintf("%d", dataCount),
				age,
			}
		}
		table.AddRow(row)
	}
	table.Render()
	return nil
}

func getSecrets(client interface{}, namespace, name, outputFormat string, showNamespace bool) error {
	clientset, ok := client.(kubernetes.Interface)
	if !ok {
		return fmt.Errorf("invalid client type")
	}

	var secrets []interface{}

	if name != "" {
		secret, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "secret", name, namespace)
		}
		secrets = []interface{}{secret}
	} else {
		secretList, err := clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "secrets", "", namespace)
		}
		for i := range secretList.Items {
			secrets = append(secrets, &secretList.Items[i])
		}
	}

	if outputFormat == OutputFormatJSON {
		return output.FormatJSON(secrets)
	}
	if outputFormat == OutputFormatYAML {
		return output.FormatYAML(secrets)
	}

	// Table output
	headers := []string{"NAME", "TYPE", "DATA", "AGE"}
	if showNamespace {
		headers = []string{"NAMESPACE", "NAME", "TYPE", "DATA", "AGE"}
	}
	table := output.NewTable(headers)
	for _, secret := range secrets {
		s, ok := secret.(*corev1.Secret)
		if !ok {
			continue
		}
		secretType := string(s.Type)
		if secretType == "" {
			secretType = SecretTypeOpaque
		}
		dataCount := len(s.Data)
		age := getAge(s.CreationTimestamp)

		row := []string{
			s.Name,
			secretType,
			fmt.Sprintf("%d", dataCount),
			age,
		}
		if showNamespace {
			row = []string{
				s.Namespace,
				s.Name,
				secretType,
				fmt.Sprintf("%d", dataCount),
				age,
			}
		}
		table.AddRow(row)
	}
	table.Render()
	return nil
}
