package project_test

import (
	"testing"

	"github.com/neopilot-ai/neo/pkg/project"
)

func TestInferPythonTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name: "simple types",
			input: map[string]interface{}{
				"name":   "test",
				"age":    25,
				"height": 5.8,
				"active": true,
			},
			expected: "active: bool\nage: int\nheight: float\nname: str\n",
		},
		{
			name: "nested object",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
					"age":  30,
				},
				"count": 42,
			},
			expected: "count: int\nclass user:\n    age: int\n    name: str\n",
		},
		{
			name: "deeply nested",
			input: map[string]interface{}{
				"app": map[string]interface{}{
					"config": map[string]interface{}{
						"debug": true,
						"port":  8080,
					},
				},
			},
			expected: "class app:\n    class config:\n        debug: bool\n        port: int\n",
		},
		{
			name: "with indent",
			input: map[string]interface{}{
				"test": "value",
			},
			expected: "    test: str\n",
		},
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.name == "with indent" {
				result = project.InferPythonTypes(tt.input, "    ")
			} else {
				result = project.InferPythonTypes(tt.input)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestInferPythonTypeForValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"string", "test", "str"},
		{"int", 42, "int"},
		{"float64", 3.14, "float"},
		{"float32", float32(2.71), "float"},
		{"bool", true, "bool"},
		{"map", map[string]interface{}{}, "dict"},
		{"slice", []string{"a", "b"}, "Any"},
		{"nil", nil, "Any"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := project.InferPythonTypeForValue(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
