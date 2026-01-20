package output

import (
	"os"
	"testing"
)

func TestFormatYAML(t *testing.T) {
	testData := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := FormatYAML(testData)
	w.Close()
	os.Stdout = originalStdout

	if err != nil {
		t.Fatalf("FormatYAML failed: %v", err)
	}

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	if n == 0 {
		t.Error("FormatYAML should write to stdout")
	}
}

func TestFormatJSON(t *testing.T) {
	testData := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := FormatJSON(testData)
	w.Close()
	os.Stdout = originalStdout

	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	if n == 0 {
		t.Error("FormatJSON should write to stdout")
	}
}

func TestFormatJSONWithComplexData(t *testing.T) {
	testData := map[string]interface{}{
		"nested": map[string]interface{}{
			"key": "value",
		},
		"array": []string{"a", "b", "c"},
	}

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := FormatJSON(testData)
	w.Close()
	os.Stdout = originalStdout

	if err != nil {
		t.Fatalf("FormatJSON failed with complex data: %v", err)
	}

	buf := make([]byte, 2048)
	n, _ := r.Read(buf)
	if n == 0 {
		t.Error("FormatJSON should write complex data to stdout")
	}
}
