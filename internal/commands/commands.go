package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yourusername/redis-clone/internal/store"
)

// Handler processes commands and returns responses
type Handler struct {
	store *store.Store
}

// NewHandler creates a new command handler
func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

// Execute processes a command and returns a response
func (h *Handler) Execute(args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("ERR empty command")
	}

	// Convert command to uppercase
	cmd, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("ERR invalid command type")
	}
	cmd = strings.ToUpper(cmd)

	// Route to appropriate handler
	switch cmd {
	case "PING":
		return h.handlePing(args)
	case "SET":
		return h.handleSet(args)
	case "GET":
		return h.handleGet(args)
	case "DELETE", "DEL":
		return h.handleDelete(args)
	case "EXISTS":
		return h.handleExists(args)
	case "KEYS":
		return h.handleKeys(args)
	case "EXPIRE":
		return h.handleExpire(args)
	case "TTL":
		return h.handleTTL(args)
	case "INFO":
		return h.handleInfo(args)
	default:
		return nil, fmt.Errorf("ERR unknown command '%s'", cmd)
	}
}

// handlePing handles PING command
func (h *Handler) handlePing(args []interface{}) (interface{}, error) {
	if len(args) == 1 {
		return "PONG", nil
	}
	if len(args) == 2 {
		msg, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("ERR invalid argument")
		}
		return msg, nil
	}
	return nil, fmt.Errorf("ERR wrong number of arguments for 'ping' command")
}

// handleSet handles SET command
// SET key value [EX seconds]
func (h *Handler) handleSet(args []interface{}) (interface{}, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'set' command")
	}

	key, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("ERR invalid key")
	}

	value, ok := args[2].(string)
	if !ok {
		return nil, fmt.Errorf("ERR invalid value")
	}

	var expiration time.Duration

	// Parse optional EX parameter
	if len(args) >= 5 {
		exFlag, ok := args[3].(string)
		if !ok || strings.ToUpper(exFlag) != "EX" {
			return nil, fmt.Errorf("ERR syntax error")
		}

		seconds, ok := args[4].(string)
		if !ok {
			return nil, fmt.Errorf("ERR invalid expire time")
		}

		sec, err := strconv.Atoi(seconds)
		if err != nil {
			return nil, fmt.Errorf("ERR invalid expire time")
		}

		expiration = time.Duration(sec) * time.Second
	}

	h.store.Set(key, value, expiration)
	return "OK", nil
}

// handleGet handles GET command
func (h *Handler) handleGet(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'get' command")
	}

	key, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("ERR invalid key")
	}

	value, exists := h.store.Get(key)
	if !exists {
		return nil, nil // Return null
	}

	return value, nil
}

// handleDelete handles DELETE/DEL command
func (h *Handler) handleDelete(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'del' command")
	}

	key, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("ERR invalid key")
	}

	deleted := h.store.Delete(key)
	if deleted {
		return int64(1), nil
	}
	return int64(0), nil
}

// handleExists handles EXISTS command
func (h *Handler) handleExists(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'exists' command")
	}

	key, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("ERR invalid key")
	}

	if h.store.Exists(key) {
		return int64(1), nil
	}
	return int64(0), nil
}

// handleKeys handles KEYS command
func (h *Handler) handleKeys(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'keys' command")
	}

	pattern, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("ERR invalid pattern")
	}

	keys := h.store.Keys(pattern)
	return keys, nil
}

// handleExpire handles EXPIRE command
func (h *Handler) handleExpire(args []interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'expire' command")
	}

	key, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("ERR invalid key")
	}

	seconds, ok := args[2].(string)
	if !ok {
		return nil, fmt.Errorf("ERR invalid expire time")
	}

	sec, err := strconv.Atoi(seconds)
	if err != nil {
		return nil, fmt.Errorf("ERR invalid expire time")
	}

	success := h.store.Expire(key, time.Duration(sec)*time.Second)
	if success {
		return int64(1), nil
	}
	return int64(0), nil
}

// handleTTL handles TTL command
func (h *Handler) handleTTL(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'ttl' command")
	}

	key, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("ERR invalid key")
	}

	ttl := h.store.TTL(key)
	return ttl, nil
}

// handleInfo handles INFO command
func (h *Handler) handleInfo(args []interface{}) (interface{}, error) {
	info := fmt.Sprintf("# Server\r\n"+
		"redis_version:7.0.0-clone\r\n"+
		"redis_mode:standalone\r\n"+
		"os:Custom\r\n"+
		"# Keyspace\r\n"+
		"db0:keys=%d\r\n",
		h.store.Count())

	return info, nil
}
