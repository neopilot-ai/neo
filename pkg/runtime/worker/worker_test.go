package worker_test

import (
	"testing"

	"github.com/neopilot-ai/neo/pkg/runtime"
	"github.com/neopilot-ai/neo/pkg/runtime/worker"
)

func TestWorkerRuntimeMatch(t *testing.T) {
	runtime := &worker.Runtime{}

	tests := []struct {
		runtime string
		match   bool
	}{
		{"worker", true},
		{"Worker", true}, // Case insensitive
		{"WORKER", true}, // Case insensitive
		{"node", false},
		{"python", false},
		{"go", false},
		{"", false},
		{"cloudflare-worker", false},
		{"workers", false},
	}

	for _, tt := range tests {
		t.Run(tt.runtime, func(t *testing.T) {
			result := runtime.Match(tt.runtime)
			if result != tt.match {
				t.Errorf("Expected Match(%q) to be %v, got %v", tt.runtime, tt.match, result)
			}
		})
	}
}

func TestWorkerRuntimeShouldRebuild(t *testing.T) {
	runtime := worker.New()

	tests := []struct {
		name       string
		functionID string
		file       string
		expected   bool
	}{
		{
			name:       "Non-existent function ID",
			functionID: "non-existent",
			file:       "/path/to/file.js",
			expected:   false,
		},
		{
			name:       "Empty function ID",
			functionID: "",
			file:       "/path/to/file.js",
			expected:   false,
		},
		{
			name:       "Empty file path",
			functionID: "test-function",
			file:       "",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runtime.ShouldRebuild(tt.functionID, tt.file)
			if result != tt.expected {
				t.Errorf("Expected ShouldRebuild(%q, %q) to be %v, got %v", tt.functionID, tt.file, tt.expected, result)
			}
		})
	}
}

func TestWorkerRuntimeRun(t *testing.T) {
	runtime := &worker.Runtime{}

	// Test that Run returns an error since it's not implemented
	// This test ensures the interface contract is maintained
	result, err := runtime.Run(nil, nil) // We can pass nil since we expect an error immediately
	
	if err == nil {
		t.Error("Expected Run to return an error (not implemented)")
	}
	
	if result != nil {
		t.Error("Expected Run to return nil worker when not implemented")
	}
}

func TestNewRuntime(t *testing.T) {
	runtime := worker.New()
	
	if runtime == nil {
		t.Fatal("Expected New() to return a runtime instance")
	}
	
	// Test that the runtime can be used for matching
	if !runtime.Match("worker") {
		t.Error("Expected new runtime to match 'worker'")
	}
}

func TestGetFile(t *testing.T) {
	runtime := &worker.Runtime{}
	
	// Test with non-existent file
	input := &runtime.BuildInput{
		Handler: "/non/existent/handler.js",
		CfgPath: "/test/config",
	}
	
	// We can't test the actual file system interaction easily without setup,
	// but we can test that the method doesn't panic
	file, found := runtime.GetFile(input)
	
	// Should not find non-existent file
	if found {
		t.Error("Expected not to find non-existent file")
	}
	
	if file != "" {
		t.Errorf("Expected empty file path, got %q", file)
	}
}
