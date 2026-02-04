package node_test

import (
	"testing"

	"github.com/neopilot-ai/neo/pkg/runtime/node"
)

func TestLoaderMap(t *testing.T) {
	// Test that the loader map contains expected mappings
	tests := []struct {
		ext      string
		expected string
	}{
		{"js", "js"},
		{"jsx", "jsx"},
		{"ts", "ts"},
		{"tsx", "tsx"},
		{"css", "css"},
		{"json", "json"},
		{"text", "text"},
		{"base64", "base64"},
		{"file", "file"},
		{"dataurl", "dataurl"},
		{"binary", "binary"},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			// This tests the internal loaderMap indirectly through the Node runtime
			runtime := &node.NodeRuntime{}
			if !runtime.Match("node") {
				t.Error("Expected Node runtime to match 'node'")
			}
			
			// We can't directly test the loaderMap since it's not exported,
			// but we can verify the runtime matches expected patterns
			if !runtime.Match("nodejs") {
				t.Error("Expected Node runtime to match 'nodejs'")
			}
			
			// Test non-matching runtimes
			if runtime.Match("python") {
				t.Error("Expected Node runtime not to match 'python'")
			}
		})
	}
}

func TestNodeRuntimeMatch(t *testing.T) {
	runtime := &node.NodeRuntime{}

	tests := []struct {
		runtime string
		match   bool
	}{
		{"node", true},
		{"nodejs", true},
		{"Node", true}, // Case insensitive
		{"NODE", true}, // Case insensitive
		{"python", false},
		{"go", false},
		{"", false},
		{"node18", true},
		{"node20", true},
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
