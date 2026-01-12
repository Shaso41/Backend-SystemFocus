package client

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Client represents a Redis client
type Client struct {
	conn   net.Conn
	reader *bufio.Reader
}

// New creates a new Redis client
func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// Ping sends a PING command
func (c *Client) Ping() (string, error) {
	if err := c.sendCommand("PING"); err != nil {
		return "", err
	}
	return c.readSimpleString()
}

// Set sets a key-value pair
func (c *Client) Set(key, value string) error {
	if err := c.sendCommand("SET", key, value); err != nil {
		return err
	}
	_, err := c.readSimpleString()
	return err
}

// SetEx sets a key-value pair with expiration
func (c *Client) SetEx(key, value string, seconds int) error {
	if err := c.sendCommand("SET", key, value, "EX", strconv.Itoa(seconds)); err != nil {
		return err
	}
	_, err := c.readSimpleString()
	return err
}

// Get gets a value by key
func (c *Client) Get(key string) (string, error) {
	if err := c.sendCommand("GET", key); err != nil {
		return "", err
	}
	return c.readBulkString()
}

// Delete deletes a key
func (c *Client) Delete(key string) (int64, error) {
	if err := c.sendCommand("DEL", key); err != nil {
		return 0, err
	}
	return c.readInteger()
}

// Exists checks if a key exists
func (c *Client) Exists(key string) (bool, error) {
	if err := c.sendCommand("EXISTS", key); err != nil {
		return false, err
	}
	n, err := c.readInteger()
	return n == 1, err
}

// Keys returns all keys matching a pattern
func (c *Client) Keys(pattern string) ([]string, error) {
	if err := c.sendCommand("KEYS", pattern); err != nil {
		return nil, err
	}
	return c.readArray()
}

// Expire sets expiration on a key
func (c *Client) Expire(key string, seconds int) (bool, error) {
	if err := c.sendCommand("EXPIRE", key, strconv.Itoa(seconds)); err != nil {
		return false, err
	}
	n, err := c.readInteger()
	return n == 1, err
}

// TTL gets time-to-live for a key
func (c *Client) TTL(key string) (int64, error) {
	if err := c.sendCommand("TTL", key); err != nil {
		return 0, err
	}
	return c.readInteger()
}

// sendCommand sends a RESP array command
func (c *Client) sendCommand(args ...string) error {
	// Build RESP array
	var cmd strings.Builder
	cmd.WriteString(fmt.Sprintf("*%d\r\n", len(args)))
	for _, arg := range args {
		cmd.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}

	_, err := c.conn.Write([]byte(cmd.String()))
	return err
}

// readSimpleString reads a RESP simple string
func (c *Client) readSimpleString() (string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return "", fmt.Errorf("empty response")
	}

	if line[0] == '+' {
		return line[1:], nil
	}

	if line[0] == '-' {
		return "", fmt.Errorf(line[1:])
	}

	return "", fmt.Errorf("unexpected response: %s", line)
}

// readInteger reads a RESP integer
func (c *Client) readInteger() (int64, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] != ':' {
		return 0, fmt.Errorf("invalid integer response")
	}

	return strconv.ParseInt(line[1:], 10, 64)
}

// readBulkString reads a RESP bulk string
func (c *Client) readBulkString() (string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] != '$' {
		return "", fmt.Errorf("invalid bulk string response")
	}

	length, err := strconv.Atoi(line[1:])
	if err != nil {
		return "", err
	}

	if length == -1 {
		return "", nil // Null
	}

	buf := make([]byte, length+2) // +2 for \r\n
	_, err = c.reader.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:length]), nil
}

// readArray reads a RESP array
func (c *Client) readArray() ([]string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] != '*' {
		return nil, fmt.Errorf("invalid array response")
	}

	count, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, err
	}

	result := make([]string, count)
	for i := 0; i < count; i++ {
		str, err := c.readBulkString()
		if err != nil {
			return nil, err
		}
		result[i] = str
	}

	return result, nil
}
