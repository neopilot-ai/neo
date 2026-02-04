package cache_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/neopilot-ai/neo/pkg/cache"
)

func BenchmarkNewCache(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-cache-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cache.NewCache(tempDir, 30*time.Minute)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCacheSet(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-cache-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	c, err := cache.NewCache(tempDir, 30*time.Minute)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("key", "value")
	}
}

func BenchmarkCacheGet(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-cache-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	c, err := cache.NewCache(tempDir, 30*time.Minute)
	if err != nil {
		b.Fatal(err)
	}

	// Pre-populate cache
	c.Set("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.Get("key")
	}
}

func BenchmarkCacheGetMiss(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-cache-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	c, err := cache.NewCache(tempDir, 30*time.Minute)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.Get("nonexistent")
	}
}

func BenchmarkCacheSetGet(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-cache-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	c, err := cache.NewCache(tempDir, 30*time.Minute)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key"
		value := "value"
		c.Set(key, value)
		_, _ = c.Get(key)
	}
}

func BenchmarkCacheConcurrentSetGet(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-cache-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	c, err := cache.NewCache(tempDir, 30*time.Minute)
	if err != nil {
		b.Fatal(err)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "key"
			value := "value"
			if i%2 == 0 {
				c.Set(key, value)
			} else {
				_, _ = c.Get(key)
			}
			i++
		}
	})
}

func BenchmarkCacheDelete(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-cache-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	c, err := cache.NewCache(tempDir, 30*time.Minute)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key"
		c.Set(key, "value")
		c.Delete(key)
	}
}

func BenchmarkCacheStats(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-cache-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	c, err := cache.NewCache(tempDir, 30*time.Minute)
	if err != nil {
		b.Fatal(err)
	}

	// Pre-populate cache
	for i := 0; i < 100; i++ {
		c.Set("key", "value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Stats()
	}
}

func BenchmarkGenerateCacheKey(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cache.GenerateCacheKey("prefix", "param1", "param2", "param3")
	}
}

func BenchmarkCacheLargeObject(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-cache-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	c, err := cache.NewCache(tempDir, 30*time.Minute)
	if err != nil {
		b.Fatal(err)
	}

	// Create a large object
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("large-key", largeData)
		_, _ = c.Get("large-key")
	}
}

func BenchmarkCacheManyKeys(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "neo-cache-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	c, err := cache.NewCache(tempDir, 30*time.Minute)
	if err != nil {
		b.Fatal(err)
	}

	numKeys := 1000
	keys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = "key"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := keys[i%numKeys]
		c.Set(key, "value")
		_, _ = c.Get(key)
	}
}
