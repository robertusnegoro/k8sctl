package output

import (
	"github.com/fatih/color"
)

var (
	// Status colors
	StatusReady   = color.New(color.FgGreen, color.Bold)
	StatusRunning = color.New(color.FgGreen)
	StatusPending = color.New(color.FgYellow)
	StatusFailed  = color.New(color.FgRed, color.Bold)
	StatusError   = color.New(color.FgRed)
	StatusWarning = color.New(color.FgYellow)

	// Header colors
	HeaderColor = color.New(color.FgBlue, color.Bold)

	// Log level colors
	LogError = color.New(color.FgRed, color.Bold)
	LogWarn  = color.New(color.FgYellow)
	LogInfo  = color.New(color.FgCyan)
	LogDebug = color.New(color.FgMagenta)

	// General colors
	Success = color.New(color.FgGreen)
	Info    = color.New(color.FgCyan)
	Warning = color.New(color.FgYellow)
	Error   = color.New(color.FgRed)
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
