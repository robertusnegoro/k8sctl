package errors_test

import (
	stderrors "errors"
	"testing"

	"github.com/robertusnegoro/k8ctl/internal/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestHandleKubernetesError_404(t *testing.T) {
	err := apierrors.NewNotFound(schema.GroupResource{Resource: "pods"}, "test-pod")

	result := errors.HandleKubernetesError(err, "pod", "test-pod", "default")

	ufErr, ok := result.(*errors.UserFriendlyError)
	if !ok {
		t.Fatal("HandleKubernetesError should return UserFriendlyError for 404")
	}
	if ufErr.Message == "" {
		t.Error("UserFriendlyError should have a message")
	}
	if ufErr.Suggestion == "" {
		t.Error("UserFriendlyError should have a suggestion")
	}
}

func TestHandleKubernetesError_403(t *testing.T) {
	err := apierrors.NewForbidden(schema.GroupResource{Resource: "pods"}, "test-pod", stderrors.New("access denied"))

	result := errors.HandleKubernetesError(err, "pod", "test-pod", "default")

	ufErr, ok := result.(*errors.UserFriendlyError)
	if !ok {
		t.Fatal("HandleKubernetesError should return UserFriendlyError for 403")
	}
	if ufErr.Message != "Access denied" {
		t.Errorf("Expected 'Access denied', got '%s'", ufErr.Message)
	}
}

func TestHandleKubernetesError_401(t *testing.T) {
	err := apierrors.NewUnauthorized("unauthorized")

	result := errors.HandleKubernetesError(err, "pod", "test-pod", "default")

	ufErr, ok := result.(*errors.UserFriendlyError)
	if !ok {
		t.Fatal("HandleKubernetesError should return UserFriendlyError for 401")
	}
	if ufErr.Message != "Unauthorized" {
		t.Errorf("Expected 'Unauthorized', got '%s'", ufErr.Message)
	}
}

func TestHandleKubernetesError_ConnectionError(t *testing.T) {
	err := stderrors.New("no such host")

	result := errors.HandleKubernetesError(err, "pod", "test-pod", "default")

	ufErr, ok := result.(*errors.UserFriendlyError)
	if !ok {
		t.Fatal("HandleKubernetesError should return UserFriendlyError for connection errors")
	}
	if ufErr.Message == "" {
		t.Error("UserFriendlyError should have a message")
	}
}

func TestWrapError(t *testing.T) {
	originalErr := stderrors.New("original error")
	message := "wrapped message"

	result := errors.WrapError(originalErr, message)

	if result == nil {
		t.Fatal("WrapError should not return nil")
	}

	ufErr, ok := result.(*errors.UserFriendlyError)
	if !ok {
		t.Fatal("WrapError should return UserFriendlyError")
	}
	if ufErr.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, ufErr.Message)
	}
	if ufErr.Original != originalErr {
		t.Error("WrapError should preserve original error")
	}
}

func TestWrapError_Nil(t *testing.T) {
	result := errors.WrapError(nil, "message")
	if result != nil {
		t.Error("WrapError should return nil for nil error")
	}
}
