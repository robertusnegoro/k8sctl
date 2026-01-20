// Package errors provides error handling and user-friendly error messages for k8ctl.
package errors //nolint:revive // Package name is intentional and does not conflict in practice

import (
	"fmt"
	"strings"

	"github.com/robertusnegoro/k8ctl/internal/output"
	"k8s.io/apimachinery/pkg/api/errors"
)

// UserFriendlyError wraps errors with user-friendly messages
type UserFriendlyError struct {
	Message    string
	Original   error
	Suggestion string
}

func (e *UserFriendlyError) Error() string {
	return e.Message
}

// HandleKubernetesError converts Kubernetes API errors to user-friendly messages
func HandleKubernetesError(err error, resourceType, resourceName, namespace string) error {
	if err == nil {
		return nil
	}

	// Check if it's a Kubernetes API error
	if k8sErr, ok := err.(*errors.StatusError); ok {
		status := k8sErr.Status()

		switch status.Code {
		case 404:
			msg := fmt.Sprintf("Resource '%s/%s' not found", resourceType, resourceName)
			if namespace != "" {
				msg += fmt.Sprintf(" in namespace '%s'", namespace)
			}
			return &UserFriendlyError{
				Message:    msg,
				Original:   err,
				Suggestion: fmt.Sprintf("Check if the %s exists: k8ctl get %s %s -n %s", resourceType, resourceType, resourceName, namespace),
			}
		case 403:
			return &UserFriendlyError{
				Message:    "Access denied",
				Original:   err,
				Suggestion: "Check your RBAC permissions or kubeconfig credentials",
			}
		case 401:
			return &UserFriendlyError{
				Message:    "Unauthorized",
				Original:   err,
				Suggestion: "Check your kubeconfig credentials: k8ctl ctx",
			}
		case 500, 502, 503:
			return &UserFriendlyError{
				Message:    "Kubernetes API server error",
				Original:   err,
				Suggestion: "Check cluster connectivity: kubectl cluster-info",
			}
		}
	}

	// Check for common error patterns
	errMsg := err.Error()

	if strings.Contains(errMsg, "no such host") || strings.Contains(errMsg, "connection refused") {
		return &UserFriendlyError{
			Message:    "Cannot connect to Kubernetes cluster",
			Original:   err,
			Suggestion: "Check your kubeconfig: k8ctl ctx",
		}
	}

	if strings.Contains(errMsg, "context not found") {
		return &UserFriendlyError{
			Message:    "Kubernetes context not found",
			Original:   err,
			Suggestion: "List available contexts: k8ctl ctx",
		}
	}

	if strings.Contains(errMsg, "namespace not found") {
		return &UserFriendlyError{
			Message:    fmt.Sprintf("Namespace '%s' not found", namespace),
			Original:   err,
			Suggestion: "List available namespaces: k8ctl ns",
		}
	}

	if strings.Contains(errMsg, "resource type") && strings.Contains(errMsg, "not yet implemented") {
		return &UserFriendlyError{
			Message:    errMsg,
			Original:   err,
			Suggestion: "Use 'k8ctl get <resource-type>' to see supported resources, or use standard kubectl for other resources",
		}
	}

	if strings.Contains(errMsg, "invalid client type") {
		return &UserFriendlyError{
			Message:    "Internal error: invalid Kubernetes client",
			Original:   err,
			Suggestion: "Try reconnecting to your cluster: k8ctl ctx",
		}
	}

	// Return original error if no specific handling
	return err
}

// PrintError prints an error with color coding and suggestions
func PrintError(err error) {
	if err == nil {
		return
	}

	// Check if it's a user-friendly error
	if ufErr, ok := err.(*UserFriendlyError); ok {
		_, _ = output.Error.Println("Error:", ufErr.Message)
		if ufErr.Suggestion != "" {
			_, _ = output.Info.Println("Suggestion:", ufErr.Suggestion)
		}
		if ufErr.Original != nil && ufErr.Original != ufErr {
			_, _ = output.Warning.Printf("Details: %v\n", ufErr.Original)
		}
		return
	}

	// Default error output
	_, _ = output.Error.Println("Error:", err.Error())
}

// WrapError wraps an error with a user-friendly message
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return &UserFriendlyError{
		Message:  message,
		Original: err,
	}
}
