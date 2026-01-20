// Package output provides output formatting and colorization utilities for k8ctl.
package output

import (
	"github.com/fatih/color"
)

var (
	// StatusReady is the color for ready status.
	StatusReady = color.New(color.FgGreen, color.Bold)
	// StatusRunning is the color for running status.
	StatusRunning = color.New(color.FgGreen)
	// StatusPending is the color for pending status.
	StatusPending = color.New(color.FgYellow)
	// StatusFailed is the color for failed status.
	StatusFailed = color.New(color.FgRed, color.Bold)
	// StatusError is the color for error status.
	StatusError = color.New(color.FgRed)
	// StatusWarning is the color for warning status.
	StatusWarning = color.New(color.FgYellow)

	// HeaderColor is the color for table headers.
	HeaderColor = color.New(color.FgBlue, color.Bold)

	// LogError is the color for error log level.
	LogError = color.New(color.FgRed, color.Bold)
	// LogWarn is the color for warning log level.
	LogWarn = color.New(color.FgYellow)
	// LogInfo is the color for info log level.
	LogInfo = color.New(color.FgCyan)
	// LogDebug is the color for debug log level.
	LogDebug = color.New(color.FgMagenta)

	// Success is the color for success messages.
	Success = color.New(color.FgGreen)
	// Info is the color for info messages.
	Info = color.New(color.FgCyan)
	// Warning is the color for warning messages.
	Warning = color.New(color.FgYellow)
	// Error is the color for error messages.
	Error = color.New(color.FgRed)
)

// ColorizeStatus returns a colored string based on status
func ColorizeStatus(status string) string {
	switch status {
	case "Running", "Ready", "Active", "Succeeded", "Completed":
		return StatusRunning.Sprint(status)
	case "Pending", "ContainerCreating", "PodInitializing":
		return StatusPending.Sprint(status)
	case "Failed", "Error", "CrashLoopBackOff", "ImagePullBackOff":
		return StatusFailed.Sprint(status)
	case "Warning", "Unknown":
		return StatusWarning.Sprint(status)
	default:
		return status
	}
}

// IsColorEnabled checks if color output is enabled
func IsColorEnabled() bool {
	return !color.NoColor
}

// DisableColors disables color output
func DisableColors() {
	color.NoColor = true
}

// EnableColors enables color output
func EnableColors() {
	color.NoColor = false
}
