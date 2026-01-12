package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_SimpleString(t *testing.T) {
	input := "+OK\r\n"
	parser := NewParser(bytes.NewBufferString(input))

	result, err := parser.Parse()
	assert.NoError(t, err)
	assert.Equal(t, "OK", result)
}

func TestParser_Error(t *testing.T) {
	input := "-Error message\r\n"
	parser := NewParser(bytes.NewBufferString(input))

	_, err := parser.Parse()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error message")
}

func TestParser_Integer(t *testing.T) {
	input := ":1000\r\n"
	parser := NewParser(bytes.NewBufferString(input))

	result, err := parser.Parse()
	assert.NoError(t, err)
	assert.Equal(t, int64(1000), result)
}

func TestParser_BulkString(t *testing.T) {
	input := "$6\r\nfoobar\r\n"
	parser := NewParser(bytes.NewBufferString(input))

	result, err := parser.Parse()
	assert.NoError(t, err)
	assert.Equal(t, "foobar", result)
}

func TestParser_NullBulkString(t *testing.T) {
	input := "$-1\r\n"
	parser := NewParser(bytes.NewBufferString(input))

	result, err := parser.Parse()
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestParser_Array(t *testing.T) {
	input := "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n"
	parser := NewParser(bytes.NewBufferString(input))

	result, err := parser.Parse()
	assert.NoError(t, err)

	array, ok := result.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(array))
	assert.Equal(t, "GET", array[0])
	assert.Equal(t, "key", array[1])
}

func TestParser_InlineCommand(t *testing.T) {
	input := "GET key\n"
	parser := NewParser(bytes.NewBufferString(input))

	result, err := parser.Parse()
	assert.NoError(t, err)

	array, ok := result.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(array))
	assert.Equal(t, "GET", array[0])
	assert.Equal(t, "key", array[1])
}
