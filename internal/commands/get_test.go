package commands

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetPods(t *testing.T) {
	// Create fake clientset
	clientset := fake.NewSimpleClientset()

	// Create test pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "container1"},
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			ContainerStatuses: []corev1.ContainerStatus{
				{Ready: true},
			},
		},
	}

	// Add pod to fake clientset
	_, err := clientset.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create test pod: %v", err)
	}

	// Test getPods function
	err = getPods(clientset, "default", "", "table", false)
	if err != nil {
		t.Errorf("getPods failed: %v", err)
	}
}

func TestGetPodsWithName(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{Name: "container1"}},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}

	_, err := clientset.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create test pod: %v", err)
	}

	err = getPods(clientset, "default", "test-pod", "table", false)
	if err != nil {
		t.Errorf("getPods with name failed: %v", err)
	}
}

func TestGetDeployments(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	replicas := int32(3)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas:     3,
			UpdatedReplicas:   3,
			AvailableReplicas: 3,
		},
	}

	_, err := clientset.AppsV1().Deployments("default").Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create test deployment: %v", err)
	}

	err = getDeployments(clientset, "default", "", "table", false)
	if err != nil {
		t.Errorf("getDeployments failed: %v", err)
	}
}

func TestGetServices(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "10.0.0.1",
			Ports: []corev1.ServicePort{
				{Port: 80, Protocol: corev1.ProtocolTCP},
			},
		},
	}

	_, err := clientset.CoreV1().Services("default").Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create test service: %v", err)
	}

	err = getServices(clientset, "default", "", "table", false)
	if err != nil {
		t.Errorf("getServices failed: %v", err)
	}
}

func TestGetNamespaces(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
		Status: corev1.NamespaceStatus{
			Phase: corev1.NamespaceActive,
		},
	}

	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create test namespace: %v", err)
	}

	err = getNamespaces(clientset, "", "table")
	if err != nil {
		t.Errorf("getNamespaces failed: %v", err)
	}
}
