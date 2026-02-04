package npm_test

import (
	"os"
	"testing"

	"github.com/neopilot-ai/neo/pkg/npm"
)

func BenchmarkGet(b *testing.B) {
	b.Skip("Skipping network-dependent benchmark in CI")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := npm.Get("aws", "6.27.0")
		if err != nil {
			b.Skip("Network unavailable")
			return
		}
	}
}

func BenchmarkCachedGet(b *testing.B) {
	b.Skip("Skipping network-dependent benchmark in CI")
	
	// First call to populate cache
	_, err := npm.CachedGet("aws", "6.27.0")
	if err != nil {
		b.Skip("Network unavailable")
		return
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := npm.CachedGet("aws", "6.27.0")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCachedGetMiss(b *testing.B) {
	b.Skip("Skipping network-dependent benchmark in CI")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := npm.CachedGet("nonexistent-package", "1.0.0")
		if err != nil {
			// Expected to fail for non-existent package
		}
	}
}

func BenchmarkDetectPackageManager(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-npm-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = npm.DetectPackageManager(tempDir)
	}
}

func BenchmarkCachedDetectPackageManager(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-npm-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// First call to populate cache
	npm.CachedDetectPackageManager(tempDir)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = npm.CachedDetectPackageManager(tempDir)
	}
}

func BenchmarkConcurrentCachedGet(b *testing.B) {
	b.Skip("Skipping network-dependent benchmark in CI")
	
	// Populate cache first
	_, err := npm.CachedGet("aws", "6.27.0")
	if err != nil {
		b.Skip("Network unavailable")
		return
	}
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := npm.CachedGet("aws", "6.27.0")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCacheStats(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = npm.GetCacheStats()
	}
}

func BenchmarkWarmCache(b *testing.B) {
	packages := []npm.PackageInfo{
		{Name: "aws", Version: "6.27.0"},
		{Name: "aws", Version: "6.26.0"},
		{Name: "gcp", Version: "7.25.0"},
		{Name: "azure", Version: "5.85.0"},
		{Name: "random", Version: "4.13.0"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := npm.WarmCache(packages)
		if err != nil {
			b.Skip("Network unavailable")
			return
		}
	}
}
