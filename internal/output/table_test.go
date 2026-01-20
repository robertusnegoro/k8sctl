package output

import (
	"os"
	"testing"
)

func TestNewTable(t *testing.T) {
	headers := []string{"Name", "Status", "Age"}
	table := NewTable(headers)

	if table == nil {
		t.Fatal("NewTable should not return nil")
	}
	if len(table.headers) != len(headers) {
		t.Errorf("Expected %d headers, got %d", len(headers), len(table.headers))
	}
	if len(table.rows) != 0 {
		t.Error("New table should have no rows")
	}
}

func TestAddRow(t *testing.T) {
	table := NewTable([]string{"Name", "Status"})

	table.AddRow([]string{"test-pod", "Running"})
	if len(table.rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(table.rows))
	}

	table.AddRow([]string{"test-pod-2", "Pending"})
	if len(table.rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(table.rows))
	}
}

func TestRender(_ *testing.T) {
	// Redirect stdout to avoid cluttering test output
	originalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = originalStdout
		_ = w.Close()
	}()

	table := NewTable([]string{"Name", "Status"})
	table.AddRow([]string{"test-pod", "Running"})
	table.AddRow([]string{"test-pod-2", "Pending"})

	// Should not panic
	table.Render()
}

func TestIsStatusField(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Running", "Running", true},
		{"Pending", "Pending", true},
		{"Failed", "Failed", true},
		{"Ready", "Ready", true},
		{"NotReady", "NotReady", true},
		{"Unknown", "Unknown", false},
		{"Custom", "Custom", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isStatusField(tt.input)
			if result != tt.expected {
				t.Errorf("isStatusField(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
