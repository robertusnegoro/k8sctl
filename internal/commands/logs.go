package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/robertusnegoro/k8ctl/internal/config"
	"github.com/robertusnegoro/k8ctl/internal/errors"
	"github.com/robertusnegoro/k8ctl/internal/k8s"
	"github.com/robertusnegoro/k8ctl/internal/output"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	logLevelRegex = regexp.MustCompile(`(?i)(ERROR|WARN|WARNING|INFO|DEBUG|FATAL)`)
)

func NewLogsCommand() *cobra.Command {
	var namespace string
	var follow bool
	var tailLines int64
	var container string

	cmd := &cobra.Command{
		Use:   "logs [pod-name]",
		Short: "Print the logs for a pod",
		Long: `Print the logs for a pod with enhanced formatting including:
- Log level color coding (ERROR, WARN, INFO, DEBUG)
- Better timestamp formatting
- Real-time streaming`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogs(cmd, args, namespace, follow, tailLines, container)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace (overrides config)")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	cmd.Flags().Int64Var(&tailLines, "tail", 10, "Lines of recent log file to display")
	cmd.Flags().StringVarP(&container, "container", "c", "", "Container name")

	return cmd
}

func runLogs(_ *cobra.Command, args []string, namespace string, follow bool, tailLines int64, container string) error {
	podName := args[0]

	client, err := k8s.GetClient()
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes client: %w", err)
	}

	// Get namespace
	if namespace == "" {
		namespace = config.GetCurrentNamespace()
		if namespace == "" {
			namespace = DefaultNamespace
		}
	}

	// Get pod to find container name if not specified
	if container == "" {
		pod, err := client.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "pod", podName, namespace)
		}
		if len(pod.Spec.Containers) > 0 {
			container = pod.Spec.Containers[0].Name
		}
	}

	opts := &corev1.PodLogOptions{
		Container: container,
		Follow:    follow,
		TailLines: &tailLines,
	}

	req := client.CoreV1().Pods(namespace).GetLogs(podName, opts)
	stream, err := req.Stream(context.Background())
	if err != nil {
		return errors.WrapError(err, "Failed to retrieve pod logs")
	}
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Text()
		colorizeLogLine(line)
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return fmt.Errorf("error reading logs: %w", err)
	}

	return nil
}

func colorizeLogLine(line string) {
	// Check for log levels and colorize
	if logLevelRegex.MatchString(line) {
		matches := logLevelRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			level := strings.ToUpper(matches[1])
			switch level {
			case "ERROR", "FATAL":
				output.LogError.Println(line)
				return
			case "WARN", "WARNING":
				output.LogWarn.Println(line)
				return
			case "INFO":
				output.LogInfo.Println(line)
				return
			case "DEBUG":
				output.LogDebug.Println(line)
				return
			}
		}
	}

	// Default output
	fmt.Println(line)
}
