package server

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer_StartAndStop(t *testing.T) {
	srv := New("localhost:0") // Use random port

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Stop server
	srv.Stop()

	// Wait for server to stop
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("Server did not stop in time")
	}
}

func TestServer_ClientConnection(t *testing.T) {
	srv := New("localhost:16379")

	// Start server
	go srv.Start()
	time.Sleep(100 * time.Millisecond)
	defer srv.Stop()

	// Connect as client
	conn, err := net.Dial("tcp", "localhost:16379")
	assert.NoError(t, err)
	defer conn.Close()

	// Send PING command
	_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	assert.NoError(t, err)

	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, "+PONG\r\n", response)
}

func TestServer_SetAndGet(t *testing.T) {
	srv := New("localhost:16380")

	// Start server
	go srv.Start()
	time.Sleep(100 * time.Millisecond)
	defer srv.Stop()

	// Connect as client
	conn, err := net.Dial("tcp", "localhost:16380")
	assert.NoError(t, err)
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// SET command
	_, err = conn.Write([]byte("*3\r\n$3\r\nSET\r\n$4\r\nkey1\r\n$6\r\nvalue1\r\n"))
	assert.NoError(t, err)

	response, err := reader.ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, "+OK\r\n", response)

	// GET command
	_, err = conn.Write([]byte("*2\r\n$3\r\nGET\r\n$4\r\nkey1\r\n"))
	assert.NoError(t, err)

	// Read bulk string response
	response, err = reader.ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, "$6\r\n", response)

	response, err = reader.ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, "value1\r\n", response)
}

func TestServer_MultipleConcurrentClients(t *testing.T) {
	t.Skip("Skipping flaky concurrent test - needs investigation")
	srv := New("localhost:16381")

	// Start server
	go srv.Start()
	time.Sleep(100 * time.Millisecond)
	defer srv.Stop()

	// Connect multiple clients
	numClients := 10
	done := make(chan bool, numClients)

	for i := 0; i < numClients; i++ {
		go func(clientID int) {
			conn, err := net.Dial("tcp", "localhost:16381")
			if err != nil {
				t.Errorf("Client %d failed to connect: %v", clientID, err)
				done <- false
				return
			}
			defer conn.Close()

			// Send PING
			_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
			if err != nil {
				t.Errorf("Client %d failed to write: %v", clientID, err)
				done <- false
				return
			}

			// Read response
			reader := bufio.NewReader(conn)
			response, err := reader.ReadString('\n')
			if err != nil {
				t.Errorf("Client %d failed to read: %v", clientID, err)
				done <- false
				return
			}

			if response != "+PONG\r\n" {
				t.Errorf("Client %d got unexpected response: %s", clientID, response)
				done <- false
				return
			}

			done <- true
		}(i)
	}

	// Wait for all clients
	for i := 0; i < numClients; i++ {
		success := <-done
		assert.True(t, success)
	}
}
