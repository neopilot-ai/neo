package python_test

import (
	"testing"

	"github.com/neopilot-ai/neo/pkg/runtime/python"
)

func TestPythonRuntimeMatch(t *testing.T) {
	runtime := &python.PythonRuntime{}

	tests := []struct {
		runtime string
		match   bool
	}{
		{"python", true},
		{"python3", true},
		{"Python", true}, // Case insensitive
		{"PYTHON", true}, // Case insensitive
		{"py", false},
		{"node", false},
		{"go", false},
		{"", false},
		{"python3.9", true},
		{"python3.11", true},
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

func TestPythonRuntimeShouldRebuild(t *testing.T) {
	runtime := &python.PythonRuntime{}

	tests := []struct {
		name       string
		functionID string
		file       string
		expected   bool
	}{
		{
			name:       "Python file change",
			functionID: "test-function",
			file:       "/path/to/handler.py",
			expected:   true,
		},
		{
			name:       "Python file in subdirectory",
			functionID: "test-function",
			file:       "/path/to/utils/helper.py",
			expected:   true,
		},
		{
			name:       "Non-Python file",
			functionID: "test-function",
			file:       "/path/to/config.json",
			expected:   false,
		},
		{
			name:       "Empty file path",
			functionID: "test-function",
			file:       "",
			expected:   false,
		},
		{
			name:       "Different function ID",
			functionID: "other-function",
			file:       "/path/to/handler.py",
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
