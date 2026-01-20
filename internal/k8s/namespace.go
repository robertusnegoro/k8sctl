package k8s

import (
	"context"
	"fmt"

	"github.com/robertusnegoro/k8ctl/internal/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListNamespaces lists all available namespaces
func ListNamespaces() ([]string, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	var nsList []string
	for i := range namespaces.Items {
		nsList = append(nsList, namespaces.Items[i].Name)
	}

	return nsList, nil
}

// GetCurrentNamespace returns the current namespace
func GetCurrentNamespace() (string, error) {
	// Try to get namespace from config first
	if ns := config.GetCurrentNamespace(); ns != "" {
		return ns, nil
	}

	// Default to "default" namespace
	return "default", nil
}

// NamespaceExists checks if a namespace exists
func NamespaceExists(namespace string) (bool, error) {
	clientset, err := GetClient()
	if err != nil {
		return false, err
	}

	_, err = clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		return false, nil
	}

	return true, nil
}

// GetNamespaceDetails returns detailed information about a namespace
func GetNamespaceDetails(namespace string) (*corev1.Namespace, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}

	ns, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	return ns, nil
}
