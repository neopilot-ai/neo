package runtime_test

import (
	"context"
	"os"
	"testing"

	"github.com/neopilot-ai/neo/pkg/runtime"
)

// MockRuntime for benchmarking
type BenchmarkMockRuntime struct {
	matchString string
}

func (m *BenchmarkMockRuntime) Match(runtime string) bool {
	return runtime == m.matchString
}

func (m *BenchmarkMockRuntime) Build(ctx context.Context, input *runtime.BuildInput) (*runtime.BuildOutput, error) {
	return &runtime.BuildOutput{
		Out:        input.Out(),
		Handler:    input.Handler,
		Errors:     []string{},
		Sourcemaps: []string{},
	}, nil
}

func (m *BenchmarkMockRuntime) Run(ctx context.Context, input *runtime.RunInput) (runtime.Worker, error) {
	return &BenchmarkMockWorker{}, nil
}

func (m *BenchmarkMockRuntime) ShouldRebuild(functionID string, path string) bool {
	return true
}

type BenchmarkMockWorker struct{}

func (m *BenchmarkMockWorker) Stop() {}

func (m *BenchmarkMockWorker) Logs() *os.File {
	return nil
}

func BenchmarkNewCollection(b *testing.B) {
	mockRuntime := &BenchmarkMockRuntime{matchString: "node"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runtime.NewCollection("/test/path", mockRuntime)
	}
}

func BenchmarkCollectionRuntime(b *testing.B) {
	mockRuntime := &BenchmarkMockRuntime{matchString: "node"}
	collection := runtime.NewCollection("/test/path", mockRuntime)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = collection.Runtime("node")
	}
}

func BenchmarkCollectionRuntimeMiss(b *testing.B) {
	mockRuntime := &BenchmarkMockRuntime{matchString: "node"}
	collection := runtime.NewCollection("/test/path", mockRuntime)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = collection.Runtime("python")
	}
}

func BenchmarkBuildInputOut(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-bench-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	
	input := &runtime.BuildInput{
		CfgPath:    cfgPath,
		FunctionID: "test-function",
		Dev:        false,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = input.Out()
	}
}

func BenchmarkBuildInputOutDev(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-bench-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	
	input := &runtime.BuildInput{
		CfgPath:    cfgPath,
		FunctionID: "test-function",
		Dev:        true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = input.Out()
	}
}

func BenchmarkCollectionBuild(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-bench-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	mockRuntime := &BenchmarkMockRuntime{matchString: "node"}
	collection := runtime.NewCollection(cfgPath, mockRuntime)

	input := &runtime.BuildInput{
		CfgPath:    cfgPath,
		FunctionID: "test-function",
		Runtime:    "node",
		Handler:    "handler",
		Dev:        true,
	}
	
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := collection.Build(ctx, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCollectionBuildWithBundle(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-bench-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	bundleDir := filepath.Join(tempDir, "bundle")
	err = os.MkdirAll(bundleDir, 0755)
	if err != nil {
		b.Fatal(err)
	}

	mockRuntime := &BenchmarkMockRuntime{matchString: "node"}
	collection := runtime.NewCollection(cfgPath, mockRuntime)

	input := &runtime.BuildInput{
		CfgPath:    cfgPath,
		FunctionID: "test-function",
		Runtime:    "node",
		Handler:    "handler",
		Bundle:     bundleDir,
	}
	
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := collection.Build(ctx, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCollectionShouldRebuild(b *testing.B) {
	mockRuntime := &BenchmarkMockRuntime{matchString: "node"}
	collection := runtime.NewCollection("/test/path", mockRuntime)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = collection.ShouldRebuild("node", "test-function", "/test/file.py")
	}
}

func BenchmarkCollectionAddTarget(b *testing.B) {
	mockRuntime := &BenchmarkMockRuntime{matchString: "node"}
	collection := runtime.NewCollection("/test/path", mockRuntime)
	
	input := &runtime.BuildInput{
		FunctionID: "test-function",
		Runtime:    "node",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collection.AddTarget(input)
	}
}
