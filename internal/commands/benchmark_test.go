package commands

import (
	"context"
	"testing"

	"github.com/robertusnegoro/k8ctl/internal/output"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// BenchmarkGetPods benchmarks the getPods function
func BenchmarkGetPods(b *testing.B) {
	client := fake.NewSimpleClientset()
	
	// Create test pods
	for i := 0; i < 100; i++ {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		}
		_, _ = client.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getPods(client, "default", "", "table", false)
	}
}

// BenchmarkGetAge benchmarks the getAge helper function
func BenchmarkGetAge(b *testing.B) {
	timestamp := metav1.Now()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getAge(timestamp)
	}
}

// BenchmarkTableRender benchmarks table rendering
func BenchmarkTableRender(b *testing.B) {
	table := output.NewTable([]string{"NAME", "STATUS", "AGE"})
	for i := 0; i < 100; i++ {
		table.AddRow([]string{"test-pod", "Running", "1h"})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		table.Render()
	}
}
