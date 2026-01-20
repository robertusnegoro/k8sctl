// Package k8s provides Kubernetes client management and utilities.
package k8s

import (
	"os"
	"path/filepath"

	"github.com/robertusnegoro/k8ctl/internal/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	clientset  *kubernetes.Clientset
	restConfig *rest.Config
)

// GetClient returns the Kubernetes clientset
func GetClient() (*kubernetes.Clientset, error) {
	if clientset != nil {
		return clientset, nil
	}

	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	clientset, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// GetConfig returns the Kubernetes REST config
func GetConfig() (*rest.Config, error) {
	if restConfig != nil {
		return restConfig, nil
	}

	// Try to get kubeconfig path
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if envKubeconfig := os.Getenv("KUBECONFIG"); envKubeconfig != "" {
		kubeconfig = envKubeconfig
	}

	// Use the current context in kubeconfig
	configLoadingRules := &clientcmd.ClientConfigLoadingRules{
		Precedence: []string{
			kubeconfig,
		},
	}

	configOverrides := &clientcmd.ConfigOverrides{}

	// Check if we have a context override from config
	currentContext := config.GetCurrentContext()
	if currentContext != "" {
		configOverrides.CurrentContext = currentContext
	}

	// Check if we have a namespace override from config
	currentNamespace := config.GetCurrentNamespace()
	if currentNamespace != "" {
		configOverrides.Context.Namespace = currentNamespace
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		configLoadingRules,
		configOverrides,
	)

	cfg, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	restConfig = cfg
	return restConfig, nil
}

// ResetClient resets the cached client (useful after context/namespace changes)
func ResetClient() {
	clientset = nil
	restConfig = nil
}
