package commands

import (
	"testing"

	"github.com/Shaso41/Backend-SystemFocus/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Ping(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	// PING without argument
	result, err := h.Execute([]interface{}{"PING"})
	assert.NoError(t, err)
	assert.Equal(t, SimpleString("PONG"), result)

	// PING with argument
	result, err = h.Execute([]interface{}{"PING", "hello"})
	assert.NoError(t, err)
	assert.Equal(t, BulkString("hello"), result)
}

func TestHandler_SetAndGet(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	// SET
	result, err := h.Execute([]interface{}{"SET", "key1", "value1"})
	assert.NoError(t, err)
	assert.Equal(t, SimpleString("OK"), result)

	// GET
	result, err = h.Execute([]interface{}{"GET", "key1"})
	assert.NoError(t, err)
	assert.Equal(t, BulkString("value1"), result)
}

func TestHandler_GetNonExistent(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	result, err := h.Execute([]interface{}{"GET", "nonexistent"})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestHandler_Delete(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	h.Execute([]interface{}{"SET", "key1", "value1"})

	result, err := h.Execute([]interface{}{"DEL", "key1"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result)

	// Verify deleted
	result, err = h.Execute([]interface{}{"GET", "key1"})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestHandler_Exists(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	h.Execute([]interface{}{"SET", "key1", "value1"})

	result, err := h.Execute([]interface{}{"EXISTS", "key1"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result)

	result, err = h.Execute([]interface{}{"EXISTS", "nonexistent"})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result)
}

func TestHandler_Keys(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	h.Execute([]interface{}{"SET", "key1", "value1"})
	h.Execute([]interface{}{"SET", "key2", "value2"})

	result, err := h.Execute([]interface{}{"KEYS", "*"})
	assert.NoError(t, err)

	keys, ok := result.([]string)
	assert.True(t, ok)
	assert.Equal(t, 2, len(keys))
}

func TestHandler_SetWithExpiration(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	result, err := h.Execute([]interface{}{"SET", "key1", "value1", "EX", "10"})
	assert.NoError(t, err)
	assert.Equal(t, SimpleString("OK"), result)

	// Verify TTL is set
	result, err = h.Execute([]interface{}{"TTL", "key1"})
	assert.NoError(t, err)
	ttl, ok := result.(int64)
	assert.True(t, ok)
	assert.True(t, ttl > 0 && ttl <= 10)
}

func TestHandler_Expire(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	h.Execute([]interface{}{"SET", "key1", "value1"})

	result, err := h.Execute([]interface{}{"EXPIRE", "key1", "10"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result)
}

func TestHandler_TTL(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	// Non-existent key
	result, err := h.Execute([]interface{}{"TTL", "nonexistent"})
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), result)

	// Key without expiration
	h.Execute([]interface{}{"SET", "key1", "value1"})
	result, err = h.Execute([]interface{}{"TTL", "key1"})
	assert.NoError(t, err)
	assert.Equal(t, int64(-2), result)
}

func TestHandler_Info(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	result, err := h.Execute([]interface{}{"INFO"})
	assert.NoError(t, err)

	info, ok := result.(BulkString)
	assert.True(t, ok)
	assert.Contains(t, string(info), "redis_version")
}

func TestHandler_UnknownCommand(t *testing.T) {
	s := store.New()
	defer s.Close()
	h := NewHandler(s)

	_, err := h.Execute([]interface{}{"UNKNOWN"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown command")
}
