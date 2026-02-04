package project_test

import (
	"context"
	"testing"

	"github.com/neopilot-ai/neo/pkg/project"
)

func TestEnvFor(t *testing.T) {
	// This is a basic test structure - the actual implementation requires
	// complex setup with AWS providers and complete events
	// For now, we'll test the basic structure and edge cases
	
	ctx := context.Background()
	
	// Test with nil complete event
	p := &project.Project{} // This would need proper initialization in real tests
	complete := &project.CompleteEvent{
		Devs: make(map[string]project.Dev),
	}
	
	// Test with non-existent dev
	env, err := p.EnvFor(ctx, complete, "non-existent")
	if err != nil {
		t.Fatalf("Expected no error for non-existent dev, got %v", err)
	}
	
	// Should still have basic app info
	if appEnv, exists := env["NEO_RESOURCE_App"]; !exists {
		t.Error("Expected NEO_RESOURCE_App to be set")
	} else if appEnv == "" {
		t.Error("Expected NEO_RESOURCE_App to have content")
	}
}

func TestEnvForWithLinks(t *testing.T) {
	ctx := context.Background()
	
	complete := &project.CompleteEvent{
		Devs: map[string]project.Dev{
			"test-resource": {
				Links: []string{"linked-resource"},
				Environment: map[string]string{
					"CUSTOM_VAR": "custom_value",
				},
			},
		},
		Links: map[string]project.Link{
			"linked-resource": {
				Properties: map[string]interface{}{
					"url":   "https://example.com",
					"port":  8080,
					"debug": true,
				},
			},
		},
	}
	
	p := &project.Project{} // Would need proper initialization
	env, err := p.EnvFor(ctx, complete, "test-resource")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Check linked resource environment variable
	if linkedEnv, exists := env["NEO_RESOURCE_linked-resource"]; !exists {
		t.Error("Expected NEO_RESOURCE_linked-resource to be set")
	} else if linkedEnv == "" {
		t.Error("Expected NEO_RESOURCE_linked-resource to have content")
	}
	
	// Check custom environment variable
	if customEnv, exists := env["CUSTOM_VAR"]; !exists {
		t.Error("Expected CUSTOM_VAR to be set")
	} else if customEnv != "custom_value" {
		t.Errorf("Expected CUSTOM_VAR to be 'custom_value', got %q", customEnv)
	}
}
