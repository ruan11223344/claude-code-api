package claude

import (
	"sync"
	"time"
)

// CachedClient wraps Client with caching capabilities
type CachedClient struct {
	*Client
	cache      map[string]cacheEntry
	mu         sync.RWMutex
	ttl        time.Duration
}

type cacheEntry struct {
	response   string
	expireTime time.Time
}

// NewCachedClient creates a new cached client
func NewCachedClient(ttl time.Duration, opts ...Option) *CachedClient {
	return &CachedClient{
		Client: New(opts...),
		cache:  make(map[string]cacheEntry),
		ttl:    ttl,
	}
}

// Ask checks cache before calling Claude
func (c *CachedClient) Ask(question string) (string, error) {
	// Check cache first
	c.mu.RLock()
	if entry, exists := c.cache[question]; exists && entry.expireTime.After(time.Now()) {
		c.mu.RUnlock()
		return entry.response, nil
	}
	c.mu.RUnlock()
	
	// Not in cache or expired, call Claude
	response, err := c.Client.Ask(question)
	if err != nil {
		return "", err
	}
	
	// Store in cache
	c.mu.Lock()
	c.cache[question] = cacheEntry{
		response:   response,
		expireTime: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
	
	// Start cleanup goroutine if needed
	go c.cleanupExpired()
	
	return response, nil
}

// cleanupExpired removes expired entries from cache
func (c *CachedClient) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	now := time.Now()
	for key, entry := range c.cache {
		if entry.expireTime.Before(now) {
			delete(c.cache, key)
		}
	}
}

// ClearCache clears all cached entries
func (c *CachedClient) ClearCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]cacheEntry)
}