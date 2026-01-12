package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStore_SetAndGet(t *testing.T) {
	store := New()
	defer store.Close()

	store.Set("key1", "value1", 0)
	
	val, exists := store.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", val)
}

func TestStore_GetNonExistent(t *testing.T) {
	store := New()
	defer store.Close()

	_, exists := store.Get("nonexistent")
	assert.False(t, exists)
}

func TestStore_Delete(t *testing.T) {
	store := New()
	defer store.Close()

	store.Set("key1", "value1", 0)
	deleted := store.Delete("key1")
	
	assert.True(t, deleted)
	
	_, exists := store.Get("key1")
	assert.False(t, exists)
}

func TestStore_DeleteNonExistent(t *testing.T) {
	store := New()
	defer store.Close()

	deleted := store.Delete("nonexistent")
	assert.False(t, deleted)
}

func TestStore_Exists(t *testing.T) {
	store := New()
	defer store.Close()

	store.Set("key1", "value1", 0)
	
	assert.True(t, store.Exists("key1"))
	assert.False(t, store.Exists("nonexistent"))
}

func TestStore_Expiration(t *testing.T) {
	store := New()
	defer store.Close()

	// Set key with 100ms expiration
	store.Set("key1", "value1", 100*time.Millisecond)
	
	// Should exist immediately
	val, exists := store.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", val)
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Should not exist after expiration
	_, exists = store.Get("key1")
	assert.False(t, exists)
}

func TestStore_Expire(t *testing.T) {
	store := New()
	defer store.Close()

	store.Set("key1", "value1", 0)
	
	// Set expiration
	success := store.Expire("key1", 100*time.Millisecond)
	assert.True(t, success)
	
	// Should exist immediately
	assert.True(t, store.Exists("key1"))
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Should not exist after expiration
	assert.False(t, store.Exists("key1"))
}

func TestStore_TTL(t *testing.T) {
	store := New()
	defer store.Close()

	// Key doesn't exist
	ttl := store.TTL("nonexistent")
	assert.Equal(t, int64(-1), ttl)
	
	// Key exists without expiration
	store.Set("key1", "value1", 0)
	ttl = store.TTL("key1")
	assert.Equal(t, int64(-2), ttl)
	
	// Key exists with expiration
	store.Set("key2", "value2", 10*time.Second)
	ttl = store.TTL("key2")
	assert.True(t, ttl > 0 && ttl <= 10)
}

func TestStore_Keys(t *testing.T) {
	store := New()
	defer store.Close()

	store.Set("key1", "value1", 0)
	store.Set("key2", "value2", 0)
	store.Set("key3", "value3", 0)
	
	keys := store.Keys("*")
	assert.Equal(t, 3, len(keys))
}

func TestStore_Count(t *testing.T) {
	store := New()
	defer store.Close()

	assert.Equal(t, 0, store.Count())
	
	store.Set("key1", "value1", 0)
	store.Set("key2", "value2", 0)
	
	assert.Equal(t, 2, store.Count())
	
	store.Delete("key1")
	assert.Equal(t, 1, store.Count())
}

func TestStore_ConcurrentAccess(t *testing.T) {
	store := New()
	defer store.Close()

	done := make(chan bool)
	
	// Multiple goroutines writing
	for i := 0; i < 10; i++ {
		go func(n int) {
			for j := 0; j < 100; j++ {
				key := "key"
				store.Set(key, "value", 0)
				store.Get(key)
				store.Delete(key)
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Benchmarks
func BenchmarkStore_Set(b *testing.B) {
	store := New()
	defer store.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Set("key", "value", 0)
	}
}

func BenchmarkStore_Get(b *testing.B) {
	store := New()
	defer store.Close()
	
	store.Set("key", "value", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get("key")
	}
}

func BenchmarkStore_ConcurrentSet(b *testing.B) {
	store := New()
	defer store.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Set("key", "value", 0)
		}
	})
}

func BenchmarkStore_ConcurrentGet(b *testing.B) {
	store := New()
	defer store.Close()
	
	store.Set("key", "value", 0)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Get("key")
		}
	})
}
