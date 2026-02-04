package runtime_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/neopilot-ai/neo/pkg/runtime"
)

// MockRuntime implements the Runtime interface for testing
type MockRuntime struct {
	matchString string
	shouldBuild bool
}

func (m *MockRuntime) Match(runtime string) bool {
	return runtime == m.matchString
}

func (m *MockRuntime) Build(ctx context.Context, input *runtime.BuildInput) (*runtime.BuildOutput, error) {
	return &runtime.BuildOutput{
		Out:        input.Out(),
		Handler:    input.Handler,
		Errors:     []string{},
		Sourcemaps: []string{},
	}, nil
}

func (m *MockRuntime) Run(ctx context.Context, input *runtime.RunInput) (runtime.Worker, error) {
	return &MockWorker{}, nil
}

func (m *MockRuntime) ShouldRebuild(functionID string, path string) bool {
	return m.shouldBuild
}

// MockWorker implements the Worker interface for testing
type MockWorker struct{}

func (m *MockWorker) Stop() {}

func (m *MockWorker) Logs() *os.File {
	return nil
}

func TestNewCollection(t *testing.T) {
	mockRuntime := &MockRuntime{matchString: "node"}
	collection := runtime.NewCollection("/test/path", mockRuntime)

	if collection == nil {
		t.Fatal("Expected collection to be created")
	}
}

func TestCollectionRuntime(t *testing.T) {
	mockRuntime := &MockRuntime{matchString: "node"}
	collection := runtime.NewCollection("/test/path", mockRuntime)

	// Test matching runtime
	r, ok := collection.Runtime("node")
	if !ok {
		t.Error("Expected to find node runtime")
	}
	if r == nil {
		t.Error("Expected runtime to be returned")
	}

	// Test non-matching runtime
	_, ok = collection.Runtime("python")
	if ok {
		t.Error("Expected not to find python runtime")
	}
}

func TestBuildInputOut(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "neo-build-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")

	// Test production build
	input := &runtime.BuildInput{
		CfgPath:    cfgPath,
		FunctionID: "test-function",
		Dev:        false,
	}

	expected := filepath.Join(tempDir, ".neo", "artifacts", "test-function-src")
	if input.Out() != expected {
		t.Errorf("Expected %q, got %q", expected, input.Out())
	}

	// Test dev build
	input.Dev = true
	expected = filepath.Join(tempDir, ".neo", "artifacts", "test-function-dev")
	if input.Out() != expected {
		t.Errorf("Expected %q, got %q", expected, input.Out())
	}
}

func TestCollectionBuild(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "neo-collection-build-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	mockRuntime := &MockRuntime{matchString: "node"}
	collection := runtime.NewCollection(cfgPath, mockRuntime)

	input := &runtime.BuildInput{
		CfgPath:    cfgPath,
		FunctionID: "test-function",
		Runtime:    "node",
		Handler:    "handler",
		Dev:        true,
	}

	ctx := context.Background()
	output, err := collection.Build(ctx, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("Expected output to be returned")
	}

	if output.Handler != "handler" {
		t.Errorf("Expected handler to be 'handler', got %q", output.Handler)
	}

	// Verify output directory exists
	if _, err := os.Stat(output.Out); os.IsNotExist(err) {
		t.Errorf("Expected output directory to exist: %s", output.Out)
	}
}

func TestCollectionBuildWithBundle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "neo-collection-bundle-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	bundleDir := filepath.Join(tempDir, "bundle")
	err = os.MkdirAll(bundleDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create bundle dir: %v", err)
	}

	mockRuntime := &MockRuntime{matchString: "node"}
	collection := runtime.NewCollection(cfgPath, mockRuntime)

	input := &runtime.BuildInput{
		CfgPath:    cfgPath,
		FunctionID: "test-function",
		Runtime:    "node",
		Handler:    "handler",
		Bundle:     bundleDir,
	}

	ctx := context.Background()
	output, err := collection.Build(ctx, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if output.Out != bundleDir {
		t.Errorf("Expected output to be bundle directory %q, got %q", bundleDir, output.Out)
	}
}

func TestCollectionBuildWithUnknownRuntime(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "neo-collection-unknown-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	mockRuntime := &MockRuntime{matchString: "node"}
	collection := runtime.NewCollection(cfgPath, mockRuntime)

	input := &runtime.BuildInput{
		CfgPath:    cfgPath,
		FunctionID: "test-function",
		Runtime:    "unknown",
		Handler:    "handler",
	}

	ctx := context.Background()
	_, err = collection.Build(ctx, input)
	if err == nil {
		t.Error("Expected error for unknown runtime")
	}
}

func TestCollectionRun(t *testing.T) {
	mockRuntime := &MockRuntime{matchString: "node"}
	collection := runtime.NewCollection("/test/path", mockRuntime)

	input := &runtime.RunInput{
		Runtime:    "node",
		FunctionID: "test-function",
		Build: &runtime.BuildOutput{
			Out:     "/test/output",
			Handler: "handler",
		},
	}

	ctx := context.Background()
	worker, err := collection.Run(ctx, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if worker == nil {
		t.Fatal("Expected worker to be returned")
	}
}

func TestCollectionRunWithUnknownRuntime(t *testing.T) {
	mockRuntime := &MockRuntime{matchString: "node"}
	collection := runtime.NewCollection("/test/path", mockRuntime)

	input := &runtime.RunInput{
		Runtime:    "unknown",
		FunctionID: "test-function",
	}

	ctx := context.Background()
	_, err := collection.Run(ctx, input)
	if err == nil {
		t.Error("Expected error for unknown runtime")
	}
}

func TestCollectionShouldRebuild(t *testing.T) {
	mockRuntime := &MockRuntime{matchString: "node", shouldBuild: true}
	collection := runtime.NewCollection("/test/path", mockRuntime)

	// Test with matching runtime
	shouldRebuild := collection.ShouldRebuild("node", "test-function", "/test/file")
	if !shouldRebuild {
		t.Error("Expected should rebuild to be true")
	}

	// Test with non-matching runtime
	shouldRebuild = collection.ShouldRebuild("unknown", "test-function", "/test/file")
	if shouldRebuild {
		t.Error("Expected should rebuild to be false for unknown runtime")
	}
}

func TestCollectionAddTarget(t *testing.T) {
	cfgPath := "/test/path"
	mockRuntime := &MockRuntime{matchString: "node"}
	collection := runtime.NewCollection(cfgPath, mockRuntime)

	input := &runtime.BuildInput{
		FunctionID: "test-function",
		Runtime:    "node",
	}

	collection.AddTarget(input)

	// Verify the cfgPath was set
	if input.CfgPath != cfgPath {
		t.Errorf("Expected cfgPath to be set to %q, got %q", cfgPath, input.CfgPath)
	}
}
