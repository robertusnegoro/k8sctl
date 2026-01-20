package commands

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Helper functions for pod operations

func getReadyContainers(pod *corev1.Pod) int {
	ready := 0
	for i := range pod.Status.ContainerStatuses {
		if pod.Status.ContainerStatuses[i].Ready {
			ready++
		}
	}
	return ready
}

func getPodStatus(pod *corev1.Pod) string {
	if pod.Status.Phase != "" {
		return string(pod.Status.Phase)
	}
	return "Unknown"
}

func getRestartCount(pod *corev1.Pod) int32 {
	var restarts int32
	for i := range pod.Status.ContainerStatuses {
		restarts += pod.Status.ContainerStatuses[i].RestartCount
	}
	return restarts
}

func getAge(creationTime metav1.Time) string {
	if creationTime.IsZero() {
		return "<unknown>"
	}

	duration := time.Since(creationTime.Time)

	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	case duration < time.Hour:
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%dh", int(duration.Hours()))
	default:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	}
}
