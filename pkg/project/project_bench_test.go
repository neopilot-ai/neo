package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neopilot-ai/neo/pkg/project"
)

func BenchmarkResolveWorkingDir(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-bench-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = project.ResolveWorkingDir(cfgPath)
	}
}

func BenchmarkResolvePlatformDir(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-bench-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = project.ResolvePlatformDir(cfgPath)
	}
}

func BenchmarkLoadPersonalStage(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-bench-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	stageFile := filepath.Join(tempDir, ".neo", "stage")
	
	// Create stage file
	err = os.MkdirAll(filepath.Dir(stageFile), 0755)
	if err != nil {
		b.Fatal(err)
	}
	err = os.WriteFile(stageFile, []byte("test-stage"), 0644)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = project.LoadPersonalStage(cfgPath)
	}
}

func BenchmarkSetPersonalStage(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-bench-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := project.SetPersonalStage(cfgPath, "test-stage")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDiscover(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-bench-test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp dir
	originalWd, err := os.Getwd()
	if err != nil {
		b.Fatal(err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(tempDir)
	if err != nil {
		b.Fatal(err)
	}

	// Create config file
	cfgPath := filepath.Join(tempDir, "neo.config.ts")
	err = os.WriteFile(cfgPath, []byte("export default {}"), 0644)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := project.Discover()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInferPythonTypes(b *testing.B) {
	input := map[string]interface{}{
		"name":   "test",
		"age":    25,
		"height": 5.8,
		"active": true,
		"user": map[string]interface{}{
			"name": "John",
			"age":  30,
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = project.InferPythonTypes(input)
	}
}

func BenchmarkInferPythonTypeForValue(b *testing.B) {
	values := []interface{}{
		"test string",
		42,
		3.14,
		true,
		map[string]interface{}{},
		[]string{"a", "b"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value := values[i%len(values)]
		_ = project.InferPythonTypeForValue(value)
	}
}
