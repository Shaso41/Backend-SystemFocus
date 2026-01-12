package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// RESP (REdis Serialization Protocol) types
const (
	SimpleString = '+'
	Error        = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
)

// Parser handles RESP protocol parsing
type Parser struct {
	reader *bufio.Reader
}

// NewParser creates a new RESP parser
func NewParser(reader io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(reader),
	}
}

// Parse reads and parses a RESP message
func (p *Parser) Parse() (interface{}, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	if len(line) == 0 {
		return nil, fmt.Errorf("empty line")
	}

	switch line[0] {
	case SimpleString:
		return string(line[1:]), nil
	case Error:
		return nil, fmt.Errorf(string(line[1:]))
	case Integer:
		return strconv.ParseInt(string(line[1:]), 10, 64)
	case BulkString:
		return p.parseBulkString(line)
	case Array:
		return p.parseArray(line)
	default:
		// Handle inline commands (plain text)
		return p.parseInline(string(line))
	}
}

// parseBulkString parses a RESP bulk string
func (p *Parser) parseBulkString(line []byte) (interface{}, error) {
	length, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, fmt.Errorf("invalid bulk string length: %w", err)
	}

	if length == -1 {
		return nil, nil // Null bulk string
	}

	buf := make([]byte, length+2) // +2 for \r\n
	_, err = io.ReadFull(p.reader, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read bulk string: %w", err)
	}

	return string(buf[:length]), nil
}

// parseArray parses a RESP array
func (p *Parser) parseArray(line []byte) (interface{}, error) {
	count, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, fmt.Errorf("invalid array length: %w", err)
	}

	if count == -1 {
		return nil, nil // Null array
	}

	array := make([]interface{}, count)
	for i := 0; i < count; i++ {
		element, err := p.Parse()
		if err != nil {
			return nil, err
		}
		array[i] = element
	}

	return array, nil
}

// parseInline parses inline commands (plain text, space-separated)
func (p *Parser) parseInline(line string) (interface{}, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	result := make([]interface{}, len(parts))
	for i, part := range parts {
		result[i] = part
	}

	return result, nil
}

// readLine reads a line from the reader
func (p *Parser) readLine() ([]byte, error) {
	line, err := p.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// Remove \r\n
	if len(line) >= 2 && line[len(line)-2] == '\r' {
		line = line[:len(line)-2]
	} else if len(line) >= 1 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}

	return line, nil
}

// Encoder handles RESP protocol encoding
type Encoder struct {
	writer io.Writer
}

// NewEncoder creates a new RESP encoder
func NewEncoder(writer io.Writer) *Encoder {
	return &Encoder{writer: writer}
}

// WriteSimpleString writes a RESP simple string
func (e *Encoder) WriteSimpleString(s string) error {
	_, err := fmt.Fprintf(e.writer, "+%s\r\n", s)
	return err
}

// WriteError writes a RESP error
func (e *Encoder) WriteError(msg string) error {
	_, err := fmt.Fprintf(e.writer, "-%s\r\n", msg)
	return err
}

// WriteInteger writes a RESP integer
func (e *Encoder) WriteInteger(n int64) error {
	_, err := fmt.Fprintf(e.writer, ":%d\r\n", n)
	return err
}

// WriteBulkString writes a RESP bulk string
func (e *Encoder) WriteBulkString(s string) error {
	if s == "" {
		_, err := e.writer.Write([]byte("$-1\r\n"))
		return err
	}
	_, err := fmt.Fprintf(e.writer, "$%d\r\n%s\r\n", len(s), s)
	return err
}

// WriteNull writes a RESP null bulk string
func (e *Encoder) WriteNull() error {
	_, err := e.writer.Write([]byte("$-1\r\n"))
	return err
}

// WriteArray writes a RESP array
func (e *Encoder) WriteArray(items []string) error {
	_, err := fmt.Fprintf(e.writer, "*%d\r\n", len(items))
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := e.WriteBulkString(item); err != nil {
			return err
		}
	}

	return nil
}
