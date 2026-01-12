package protocol

import (
	"bufio"
	"fmt"
	"io"
)

// Encoder encodes RESP (REdis Serialization Protocol) responses
type Encoder struct {
	writer *bufio.Writer
}

// NewEncoder creates a new RESP encoder
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: bufio.NewWriter(w),
	}
}

// WriteSimpleString writes a RESP simple string (+OK\r\n)
func (e *Encoder) WriteSimpleString(s string) error {
	if _, err := e.writer.WriteString(fmt.Sprintf("+%s\r\n", s)); err != nil {
		return err
	}
	return e.writer.Flush()
}

// WriteError writes a RESP error (-ERR message\r\n)
func (e *Encoder) WriteError(msg string) error {
	if _, err := e.writer.WriteString(fmt.Sprintf("-%s\r\n", msg)); err != nil {
		return err
	}
	return e.writer.Flush()
}

// WriteInteger writes a RESP integer (:123\r\n)
func (e *Encoder) WriteInteger(n int64) error {
	if _, err := e.writer.WriteString(fmt.Sprintf(":%d\r\n", n)); err != nil {
		return err
	}
	return e.writer.Flush()
}

// WriteBulkString writes a RESP bulk string ($6\r\nfoobar\r\n)
func (e *Encoder) WriteBulkString(s string) error {
	if _, err := e.writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)); err != nil {
		return err
	}
	return e.writer.Flush()
}

// WriteNull writes a RESP null bulk string ($-1\r\n)
func (e *Encoder) WriteNull() error {
	if _, err := e.writer.WriteString("$-1\r\n"); err != nil {
		return err
	}
	return e.writer.Flush()
}

// WriteArray writes a RESP array (*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n)
func (e *Encoder) WriteArray(arr []string) error {
	if _, err := e.writer.WriteString(fmt.Sprintf("*%d\r\n", len(arr))); err != nil {
		return err
	}
	for _, s := range arr {
		if _, err := e.writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)); err != nil {
			return err
		}
	}
	return e.writer.Flush()
}
