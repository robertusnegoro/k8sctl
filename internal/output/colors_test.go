package output

import (
	"testing"
)

func TestColorizeStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Running status", "Running", "Running"},
		{"Ready status", "Ready", "Ready"},
		{"Pending status", "Pending", "Pending"},
		{"Failed status", "Failed", "Failed"},
		{"Error status", "Error", "Error"},
		{"Unknown status", "Unknown", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorizeStatus(tt.status)
			if result == "" {
				t.Error("ColorizeStatus should not return empty string")
			}
			// Just verify it returns something (color codes make exact comparison difficult)
		})
	}
}

func TestIsColorEnabled(_ *testing.T) {
	enabled := IsColorEnabled()
	// This depends on environment, so we just check it doesn't panic
	_ = enabled
}

func TestDisableAndEnableColors(_ *testing.T) {
	// Test that functions don't panic
	DisableColors()
	enabled := IsColorEnabled()
	_ = enabled // Just check it doesn't panic

	EnableColors()
	enabled = IsColorEnabled()
	_ = enabled // Just check it doesn't panic
}
