package store

import (
	"sync"
	"time"
)

// Value represents a stored value with metadata
type Value struct {
	Data      string
	CreatedAt time.Time
}

// Store is a thread-safe in-memory key-value store
type Store struct {
	mu      sync.RWMutex
	data    map[string]*Value
	expires map[string]time.Time
	stopCh  chan struct{}
}

// New creates a new Store instance and starts the cleanup goroutine
func New() *Store {
	s := &Store{
		data:    make(map[string]*Value),
		expires: make(map[string]time.Time),
		stopCh:  make(chan struct{}),
	}
	
	// Start background cleanup goroutine
	go s.cleanupExpired()
	
	return s
}

// Set stores a key-value pair with optional expiration
func (s *Store) Set(key, value string, expiration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.data[key] = &Value{
		Data:      value,
		CreatedAt: time.Now(),
	}
	
	if expiration > 0 {
		s.expires[key] = time.Now().Add(expiration)
	} else {
		delete(s.expires, key)
	}
}

// Get retrieves a value by key
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Check if key is expired
	if expireTime, exists := s.expires[key]; exists {
		if time.Now().After(expireTime) {
			return "", false
		}
	}
	
	val, exists := s.data[key]
	if !exists {
		return "", false
	}
	
	return val.Data, true
}

// Delete removes a key from the store
func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	_, exists := s.data[key]
	if exists {
		delete(s.data, key)
		delete(s.expires, key)
	}
	
	return exists
}

// Exists checks if a key exists and is not expired
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Check if key is expired
	if expireTime, exists := s.expires[key]; exists {
		if time.Now().After(expireTime) {
			return false
		}
	}
	
	_, exists := s.data[key]
	return exists
}

// Keys returns all non-expired keys matching a simple pattern
// For simplicity, only supports "*" wildcard
func (s *Store) Keys(pattern string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	keys := make([]string, 0)
	now := time.Now()
	
	for key := range s.data {
		// Check if expired
		if expireTime, exists := s.expires[key]; exists {
			if now.After(expireTime) {
				continue
			}
		}
		
		// Simple pattern matching (only supports "*")
		if pattern == "*" {
			keys = append(keys, key)
		}
	}
	
	return keys
}

// Expire sets an expiration time on an existing key
func (s *Store) Expire(key string, duration time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.data[key]; !exists {
		return false
	}
	
	s.expires[key] = time.Now().Add(duration)
	return true
}

// TTL returns the time-to-live for a key in seconds
// Returns -1 if key doesn't exist, -2 if key exists but has no expiration
func (s *Store) TTL(key string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if _, exists := s.data[key]; !exists {
		return -1
	}
	
	expireTime, hasExpiration := s.expires[key]
	if !hasExpiration {
		return -2
	}
	
	ttl := time.Until(expireTime).Seconds()
	if ttl < 0 {
		return -1
	}
	
	return int64(ttl)
}

// Count returns the number of keys in the store
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return len(s.data)
}

// Close stops the cleanup goroutine
func (s *Store) Close() {
	close(s.stopCh)
}

// cleanupExpired runs in the background and removes expired keys
func (s *Store) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.removeExpiredKeys()
		case <-s.stopCh:
			return
		}
	}
}

// removeExpiredKeys removes all expired keys from the store
func (s *Store) removeExpiredKeys() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	now := time.Now()
	for key, expireTime := range s.expires {
		if now.After(expireTime) {
			delete(s.data, key)
			delete(s.expires, key)
		}
	}
}
