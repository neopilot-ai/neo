package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/neopilot-ai/neo/pkg/cache"
)

var (
	// Global cache instance for configuration
	configCache *cache.Cache
	cacheOnce   sync.Once
)

// ConfigCacheEntry represents a cached configuration
type ConfigCacheEntry struct {
	Content    string            `json:"content"`
	Hash       string            `json:"hash"`
	Metadata   map[string]string `json:"metadata"`
	ModifiedAt time.Time         `json:"modifiedAt"`
}

// GetConfigCache returns the configuration cache instance
func GetConfigCache() (*cache.Cache, error) {
	var err error
	cacheOnce.Do(func() {
		cacheDir := filepath.Join(os.TempDir(), "neo-cache", "config")
		configCache, err = cache.NewCache(cacheDir, 1*time.Hour) // Configs cached for 1 hour
		if err != nil {
			err = fmt.Errorf("failed to create config cache: %w", err)
		}
	})
	
	return configCache, err
}

// CachedConfigFile represents a cached configuration file
type CachedConfigFile struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Hash     string `json:"hash"`
	Modified int64  `json:"modified"`
}

// GetCachedConfig reads a configuration file with caching
func GetCachedConfig(path string) (*CachedConfigFile, error) {
	cache, err := GetConfigCache()
	if err != nil {
		// Fallback to direct file read
		return readConfigFile(path)
	}

	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("config:%s:%d", path, fileInfo.ModTime().Unix())
	
	// Try cache first
	if cached, found := cache.Get(cacheKey); found {
		if config, ok := cached.(*CachedConfigFile); ok {
			return config, nil
		}
	}

	// Cache miss - read file
	config, err := readConfigFile(path)
	if err != nil {
		return nil, err
	}

	// Store in cache
	cache.Set(cacheKey, config)
	return config, nil
}

// readConfigFile reads a configuration file directly
func readConfigFile(path string) (*CachedConfigFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &CachedConfigFile{
		Path:     path,
		Content:  string(content),
		Hash:     computeHash(content),
		Modified: fileInfo.ModTime().Unix(),
	}, nil
}

// computeHash computes a simple hash for file content
func computeHash(content []byte) string {
	// Simple hash computation - could use crypto/sha256 for better security
	hash := 0
	for _, b := range content {
		hash = hash*31 + int(b)
	}
	return fmt.Sprintf("%x", hash)
}

// InvalidateConfigCache invalidates the configuration cache
func InvalidateConfigCache() error {
	cache, err := GetConfigCache()
	if err != nil {
		return err
	}
	return cache.Clear()
}

// GetConfigCacheStats returns cache statistics
func GetConfigCacheStats() *cache.CacheStats {
	cache, err := GetConfigCache()
	if err != nil {
		return nil
	}
	stats := cache.Stats()
	return &stats
}

// CacheProjectConfig caches project configuration metadata
func CacheProjectConfig(projectPath string, config map[string]interface{}) error {
	cache, err := GetConfigCache()
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("project-config:%s", projectPath)
	cache.Set(cacheKey, config)
	return nil
}

// GetCachedProjectConfig retrieves cached project configuration
func GetCachedProjectConfig(projectPath string) (map[string]interface{}, bool) {
	cache, err := GetConfigCache()
	if err != nil {
		return nil, false
	}

	cacheKey := fmt.Sprintf("project-config:%s", projectPath)
	if cached, found := cache.Get(cacheKey); found {
		if config, ok := cached.(map[string]interface{}); ok {
			return config, true
		}
	}
	return nil, false
}

// CacheProviderMetadata caches provider metadata
func CacheProviderMetadata(providerName string, metadata map[string]interface{}) error {
	cache, err := GetConfigCache()
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("provider-metadata:%s", providerName)
	cache.Set(cacheKey, metadata)
	return nil
}

// GetCachedProviderMetadata retrieves cached provider metadata
func GetCachedProviderMetadata(providerName string) (map[string]interface{}, bool) {
	cache, err := GetConfigCache()
	if err != nil {
		return nil, false
	}

	cacheKey := fmt.Sprintf("provider-metadata:%s", providerName)
	if cached, found := cache.Get(cacheKey); found {
		if metadata, ok := cached.(map[string]interface{}); ok {
			return metadata, true
		}
	}
	return nil, false
}
