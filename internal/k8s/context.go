package k8s

import (
	"fmt"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// ListContexts lists all available contexts
func ListContexts() ([]string, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	var contexts []string
	for name := range config.Contexts {
		contexts = append(contexts, name)
	}

	return contexts, nil
}

// GetCurrentContext returns the current context from kubeconfig
func GetCurrentContext() (string, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return "", fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	return config.CurrentContext, nil
}

// SetContext sets the current context in kubeconfig
func SetContext(contextName string) error {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	if _, exists := config.Contexts[contextName]; !exists {
		return fmt.Errorf("context %s does not exist", contextName)
	}

	config.CurrentContext = contextName
	return clientcmd.WriteToFile(*config, kubeconfig)
}
