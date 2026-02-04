package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neopilot-ai/neo/pkg/project"
)

func TestLoadPersonalStage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "neo-stage-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")

	// Test when no stage file exists
	stage := project.LoadPersonalStage(cfgPath)
	if stage != "" {
		t.Errorf("Expected empty stage when file doesn't exist, got %q", stage)
	}

	// Test with valid stage file
	stageFile := filepath.Join(tempDir, "stage")
	err = os.WriteFile(stageFile, []byte("  test-stage  "), 0644)
	if err != nil {
		t.Fatalf("Failed to write stage file: %v", err)
	}

	stage = project.LoadPersonalStage(cfgPath)
	if stage != "test-stage" {
		t.Errorf("Expected 'test-stage', got %q", stage)
	}
}

func TestSetPersonalStage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "neo-stage-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")

	// Test setting a stage
	err = project.SetPersonalStage(cfgPath, "  my-stage  ")
	if err != nil {
		t.Fatalf("Failed to set stage: %v", err)
	}

	// Verify the stage was set correctly
	stage := project.LoadPersonalStage(cfgPath)
	if stage != "my-stage" {
		t.Errorf("Expected 'my-stage', got %q", stage)
	}

	// Test reading the file directly to verify content
	stageFile := filepath.Join(tempDir, "stage")
	data, err := os.ReadFile(stageFile)
	if err != nil {
		t.Fatalf("Failed to read stage file: %v", err)
	}

	if string(data) != "my-stage" {
		t.Errorf("Expected file content 'my-stage', got %q", string(data))
	}
}

func TestResolveStageFile(t *testing.T) {
	// This is an internal function, but we can test it indirectly
	// by checking if SetPersonalStage creates the file in the right location
	tempDir, err := os.MkdirTemp("", "neo-stage-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	
	err = project.SetPersonalStage(cfgPath, "test")
	if err != nil {
		t.Fatalf("Failed to set stage: %v", err)
	}

	expectedStageFile := filepath.Join(tempDir, "stage")
	if _, err := os.Stat(expectedStageFile); os.IsNotExist(err) {
		t.Errorf("Stage file was not created at expected location: %s", expectedStageFile)
	}
}
