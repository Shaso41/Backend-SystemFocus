package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoder_SimpleString(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)

	err := enc.WriteSimpleString("OK")
	assert.NoError(t, err)
	assert.Equal(t, "+OK\r\n", buf.String())
}

func TestEncoder_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)

	err := enc.WriteError("ERR unknown command")
	assert.NoError(t, err)
	assert.Equal(t, "-ERR unknown command\r\n", buf.String())
}

func TestEncoder_Integer(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)

	err := enc.WriteInteger(123)
	assert.NoError(t, err)
	assert.Equal(t, ":123\r\n", buf.String())
}

func TestEncoder_BulkString(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)

	err := enc.WriteBulkString("foobar")
	assert.NoError(t, err)
	assert.Equal(t, "$6\r\nfoobar\r\n", buf.String())
}

func TestEncoder_Null(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)

	err := enc.WriteNull()
	assert.NoError(t, err)
	assert.Equal(t, "$-1\r\n", buf.String())
}

func TestEncoder_Array(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)

	err := enc.WriteArray([]string{"foo", "bar"})
	assert.NoError(t, err)
	assert.Equal(t, "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", buf.String())
}

func TestEncoder_EmptyArray(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)

	err := enc.WriteArray([]string{})
	assert.NoError(t, err)
	assert.Equal(t, "*0\r\n", buf.String())
}
