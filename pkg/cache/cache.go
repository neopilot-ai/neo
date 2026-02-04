package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// CacheEntry represents a cached item with expiration
type CacheEntry struct {
	Data      interface{} `json:"data"`
	ExpiresAt time.Time   `json:"expiresAt"`
	CreatedAt time.Time   `json:"createdAt"`
}

// IsExpired checks if the cache entry has expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache represents a thread-safe file-based cache
type Cache struct {
	dir     string
	ttl     time.Duration
	mu      sync.RWMutex
	memory  map[string]*CacheEntry
	dirty   map[string]bool
}

// NewCache creates a new cache instance
func NewCache(dir string, ttl time.Duration) (*Cache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	cache := &Cache{
		dir:    dir,
		ttl:    ttl,
		memory: make(map[string]*CacheEntry),
		dirty:  make(map[string]bool),
	}

	// Load existing cache entries from disk
	if err := cache.loadFromDisk(); err != nil {
		return nil, fmt.Errorf("failed to load cache from disk: %w", err)
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache, nil
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.memory[key]
	if !exists || entry.IsExpired() {
		return nil, false
	}

	return entry.Data, true
}

// Set stores a value in cache
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &CacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(c.ttl),
		CreatedAt: time.Now(),
	}

	c.memory[key] = entry
	c.dirty[key] = true

	// Persist to disk asynchronously
	go c.persistToDisk(key)
}

// Delete removes a value from cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.memory, key)
	delete(c.dirty, key)

	// Remove from disk
	go os.Remove(c.cacheFilePath(key))
}

// Clear removes all entries from cache
func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.memory = make(map[string]*CacheEntry)
	c.dirty = make(map[string]bool)

	// Remove all cache files
	return os.RemoveAll(c.dir)
}

// Stats returns cache statistics
func (c *Cache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := len(c.memory)
	expired := 0
	for _, entry := range c.memory {
		if entry.IsExpired() {
			expired++
		}
	}

	return CacheStats{
		Total:   total,
		Expired: expired,
		Valid:   total - expired,
		Dir:     c.dir,
	}
}

// CacheStats represents cache statistics
type CacheStats struct {
	Total   int    `json:"total"`
	Expired int    `json:"expired"`
	Valid   int    `json:"valid"`
	Dir     string `json:"dir"`
}

// cacheFilePath returns the file path for a cache key
func (c *Cache) cacheFilePath(key string) string {
	hash := sha256.Sum256([]byte(key))
	return filepath.Join(c.dir, hex.EncodeToString(hash[:])+".json")
}

// loadFromDisk loads cache entries from disk
func (c *Cache) loadFromDisk() error {
	files, err := os.ReadDir(c.dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(c.dir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip corrupted files
		}

		var entry CacheEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue // Skip corrupted entries
		}

		if !entry.IsExpired() {
			// Use filename as key (hash of original key)
			key := file.Name()[:len(file.Name())-5] // Remove .json extension
			c.memory[key] = &entry
		} else {
			// Remove expired files
			os.Remove(filePath)
		}
	}

	return nil
}

// persistToDisk saves a cache entry to disk
func (c *Cache) persistToDisk(key string) {
	c.mu.RLock()
	entry, exists := c.memory[key]
	if !exists {
		c.mu.RUnlock()
		return
	}
	c.mu.RUnlock()

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	filePath := c.cacheFilePath(key)
	os.WriteFile(filePath, data, 0644)
}

// cleanup removes expired entries
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.memory {
			if entry.IsExpired() {
				delete(c.memory, key)
				delete(c.dirty, key)
				os.Remove(c.cacheFilePath(key))
			}
		}
		c.mu.Unlock()
	}
}

// GenerateCacheKey creates a cache key from parameters
func GenerateCacheKey(prefix string, params ...string) string {
	key := prefix
	for _, param := range params {
		key += ":" + param
	}
	return key
}
