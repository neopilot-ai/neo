package npm

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/neopilot-ai/neo/internal/fs"
	"github.com/neopilot-ai/neo/pkg/cache"
)

var (
	// Global cache instance for NPM packages
	packageCache *cache.Cache
	cacheOnce    = func() func() {
		var once sync.Once
		return func() {
			once.Do(func() {
				cacheDir := filepath.Join(os.TempDir(), "neo-cache", "npm")
				var err error
				packageCache, err = cache.NewCache(cacheDir, 30*time.Minute)
				if err != nil {
					slog.Warn("Failed to create NPM cache", "error", err)
					packageCache = nil
				}
			})
		}
	}()
)

// CachedGet retrieves package information with caching
func CachedGet(name string, version string) (*Package, error) {
	cacheOnce()
	
	if packageCache == nil {
		// Fallback to non-cached version
		return Get(name, version)
	}

	cacheKey := cache.GenerateCacheKey("package", name, version)
	
	// Try to get from cache first
	if cached, found := packageCache.Get(cacheKey); found {
		if pkg, ok := cached.(*Package); ok {
			slog.Info("package cache hit", "name", name, "version", version)
			return pkg, nil
		}
	}

	// Cache miss - fetch from network
	slog.Info("package cache miss", "name", name, "version", version)
	pkg, err := Get(name, version)
	if err != nil {
		return nil, err
	}

	// Store in cache
	packageCache.Set(cacheKey, pkg)
	return pkg, nil
}

// CachedDetectPackageManager detects package manager with caching
func CachedDetectPackageManager(dir string) (string, string) {
	cacheOnce()
	
	if packageCache == nil {
		// Fallback to non-cached version
		return DetectPackageManager(dir)
	}

	cacheKey := cache.GenerateCacheKey("package-manager", dir)
	
	// Try to get from cache first
	if cached, found := packageCache.Get(cacheKey); found {
		if result, ok := cached.([]string); ok && len(result) == 2 {
			slog.Info("package manager cache hit", "dir", dir)
			return result[0], result[1]
		}
	}

	// Cache miss - detect
	slog.Info("package manager cache miss", "dir", dir)
	name, lock := DetectPackageManager(dir)
	
	// Store in cache
	packageCache.Set(cacheKey, []string{name, lock})
	return name, lock
}

// InvalidateCache invalidates the NPM package cache
func InvalidateCache() error {
	if packageCache == nil {
		return nil
	}
	return packageCache.Clear()
}

// GetCacheStats returns cache statistics
func GetCacheStats() *cache.CacheStats {
	if packageCache == nil {
		return nil
	}
	stats := packageCache.Stats()
	return &stats
}

// WarmCache pre-populates the cache with commonly used packages
func WarmCache(packages []PackageInfo) error {
	cacheOnce()
	
	if packageCache == nil {
		return fmt.Errorf("cache not available")
	}

	for _, pkg := range packages {
		cacheKey := cache.GenerateCacheKey("package", pkg.Name, pkg.Version)
		
		// Check if already cached
		if _, found := packageCache.Get(cacheKey); found {
			continue
		}

		// Fetch and cache
		_, err := CachedGet(pkg.Name, pkg.Version)
		if err != nil {
			slog.Warn("Failed to warm cache for package", "name", pkg.Name, "version", pkg.Version, "error", err)
		}
	}

	return nil
}

// PackageInfo represents a package for cache warming
type PackageInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
