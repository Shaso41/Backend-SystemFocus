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

func TestEncoder_SimpleString(t *testing.T) {
	var buf bytes.Buffer
	encoder := NewEncoder(&buf)

	err := encoder.WriteSimpleString("OK")
	assert.NoError(t, err)
	assert.Equal(t, "+OK\r\n", buf.String())
}

func TestEncoder_Error(t *testing.T) {
	var buf bytes.Buffer
	encoder := NewEncoder(&buf)

	err := encoder.WriteError("Error message")
	assert.NoError(t, err)
	assert.Equal(t, "-Error message\r\n", buf.String())
}

func TestEncoder_Integer(t *testing.T) {
	var buf bytes.Buffer
	encoder := NewEncoder(&buf)

	err := encoder.WriteInteger(1000)
	assert.NoError(t, err)
	assert.Equal(t, ":1000\r\n", buf.String())
}

func TestEncoder_BulkString(t *testing.T) {
	var buf bytes.Buffer
	encoder := NewEncoder(&buf)

	err := encoder.WriteBulkString("foobar")
	assert.NoError(t, err)
	assert.Equal(t, "$6\r\nfoobar\r\n", buf.String())
}

func TestEncoder_Null(t *testing.T) {
	var buf bytes.Buffer
	encoder := NewEncoder(&buf)

	err := encoder.WriteNull()
	assert.NoError(t, err)
	assert.Equal(t, "$-1\r\n", buf.String())
}

func TestEncoder_Array(t *testing.T) {
	var buf bytes.Buffer
	encoder := NewEncoder(&buf)

	err := encoder.WriteArray([]string{"foo", "bar"})
	assert.NoError(t, err)
	assert.Equal(t, "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", buf.String())
}
