package config

import (
	"path/filepath"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	dir := GetConfigDir()
	if dir == "" {
		t.Error("Config directory should not be empty")
	}
	if !filepath.IsAbs(dir) {
		t.Error("Config directory should be absolute path")
	}
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()
	if path == "" {
		t.Error("Config path should not be empty")
	}
	if filepath.Base(path) != "config.yaml" {
		t.Error("Config file should be named config.yaml")
	}
}

func TestLoadAndSave(t *testing.T) {
	// Create temporary config directory
	tmpDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath

	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
		cfg = nil
	}()

	configDir = tmpDir
	configPath = filepath.Join(tmpDir, "config.yaml")

	// Test loading non-existent config (should create default)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if cfg == nil {
		t.Fatal("Config should not be nil")
	}
	if cfg.Colors.Enabled != true {
		t.Error("Default colors.enabled should be true")
	}
	if cfg.Output.Format != "table" {
		t.Error("Default output.format should be 'table'")
	}

	// Test saving config
	cfg.CurrentContext = "test-context"
	cfg.CurrentNS = "test-namespace"
	err = Save(cfg)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test loading saved config
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}
	if cfg2.CurrentContext != "test-context" {
		t.Errorf("Expected context 'test-context', got '%s'", cfg2.CurrentContext)
	}
	if cfg2.CurrentNS != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", cfg2.CurrentNS)
	}
}

func TestSetCurrentContext(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath

	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
		cfg = nil
	}()

	configDir = tmpDir
	configPath = filepath.Join(tmpDir, "config.yaml")

	err := SetCurrentContext("test-ctx")
	if err != nil {
		t.Fatalf("Failed to set context: %v", err)
	}

	ctx := GetCurrentContext()
	if ctx != "test-ctx" {
		t.Errorf("Expected context 'test-ctx', got '%s'", ctx)
	}
}

func TestSetCurrentNamespace(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath

	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
		cfg = nil
	}()

	configDir = tmpDir
	configPath = filepath.Join(tmpDir, "config.yaml")

	err := SetCurrentNamespace("test-ns")
	if err != nil {
		t.Fatalf("Failed to set namespace: %v", err)
	}

	ns := GetCurrentNamespace()
	if ns != "test-ns" {
		t.Errorf("Expected namespace 'test-ns', got '%s'", ns)
	}
}

func TestGetWithNilConfig(t *testing.T) {
	originalCfg := cfg
	defer func() {
		cfg = originalCfg
	}()

	cfg = nil
	result := Get()
	if result == nil {
		t.Error("Get() should return a config even if cfg is nil")
	}
}
