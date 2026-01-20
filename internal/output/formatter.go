package output

import (
	"encoding/json"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// FormatYAML formats data as YAML
func FormatYAML(data interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	_, err = fmt.Fprint(os.Stdout, string(yamlData))
	if err != nil {
		return fmt.Errorf("failed to write YAML: %w", err)
	}
	return nil
}

// FormatJSON formats data as JSON
func FormatJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	_, err = fmt.Fprint(os.Stdout, string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}
	return nil
}
