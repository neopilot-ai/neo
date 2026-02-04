package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neopilot-ai/neo/pkg/project"
)

func TestDiscover(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "neo-discover-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test when no config file exists
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	_, err = project.Discover()
	if err == nil {
		t.Error("Expected error when no config file exists")
	}

	// Create a config file
	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	err = os.WriteFile(cfgPath, []byte("export default {}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	discoveredPath, err := project.Discover()
	if err != nil {
		t.Fatalf("Failed to discover config: %v", err)
	}

	if discoveredPath != cfgPath {
		t.Errorf("Expected %q, got %q", cfgPath, discoveredPath)
	}
}

func TestResolveWorkingDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "neo-working-dir-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	workingDir := project.ResolveWorkingDir(cfgPath)
	expected := filepath.Join(tempDir, ".neo")

	if workingDir != expected {
		t.Errorf("Expected %q, got %q", expected, workingDir)
	}
}

func TestResolvePlatformDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "neo-platform-dir-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	platformDir := project.ResolvePlatformDir(cfgPath)
	expected := filepath.Join(tempDir, ".neo", "platform")

	if platformDir != expected {
		t.Errorf("Expected %q, got %q", expected, platformDir)
	}
}

func TestResolveLogDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "neo-log-dir-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	logDir := project.ResolveLogDir(cfgPath)
	expected := filepath.Join(tempDir, ".neo", "log")

	if logDir != expected {
		t.Errorf("Expected %q, got %q", expected, logDir)
	}
}
