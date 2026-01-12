package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/redis-clone/internal/commands"
	"github.com/yourusername/redis-clone/internal/protocol"
	"github.com/yourusername/redis-clone/internal/store"
)

// Server represents the Redis-like TCP server
type Server struct {
	address  string
	listener net.Listener
	store    *store.Store
	handler  *commands.Handler
	stopCh   chan struct{}
}

// New creates a new server instance
func New(address string) *Server {
	s := store.New()
	return &Server{
		address: address,
		store:   s,
		handler: commands.NewHandler(s),
		stopCh:  make(chan struct{}),
	}
}

// Start starts the TCP server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.listener = listener
	log.Printf("🚀 Redis Clone server started on %s", s.address)
	log.Printf("📊 Ready to accept connections...")

	// Handle graceful shutdown
	go s.handleShutdown()

	// Accept connections
	for {
		select {
		case <-s.stopCh:
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-s.stopCh:
					return nil
				default:
					log.Printf("Error accepting connection: %v", err)
					continue
				}
			}

			// Handle connection in a new goroutine
			go s.handleConnection(conn)
		}
	}
}

// handleConnection handles a single client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("✅ New client connected: %s", clientAddr)

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(5 * time.Minute))

	parser := protocol.NewParser(conn)
	encoder := protocol.NewEncoder(conn)

	for {
		// Reset deadline on each command
		conn.SetDeadline(time.Now().Add(5 * time.Minute))

		// Parse command
		data, err := parser.Parse()
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("❌ Parse error from %s: %v", clientAddr, err)
				encoder.WriteError(fmt.Sprintf("ERR %v", err))
			}
			break
		}

		// Convert to command arguments
		args, ok := data.([]interface{})
		if !ok {
			encoder.WriteError("ERR invalid command format")
			continue
		}

		// Execute command
		result, err := s.handler.Execute(args)
		if err != nil {
			encoder.WriteError(err.Error())
			continue
		}

		// Send response
		if err := s.writeResponse(encoder, result); err != nil {
			log.Printf("❌ Write error to %s: %v", clientAddr, err)
			break
		}
	}

	log.Printf("👋 Client disconnected: %s", clientAddr)
}

// writeResponse writes the appropriate RESP response based on result type
func (s *Server) writeResponse(encoder *protocol.Encoder, result interface{}) error {
	switch v := result.(type) {
	case nil:
		return encoder.WriteNull()
	case string:
		return encoder.WriteSimpleString(v)
	case int64:
		return encoder.WriteInteger(v)
	case []string:
		return encoder.WriteArray(v)
	default:
		return encoder.WriteError("ERR unknown response type")
	}
}

// handleShutdown handles graceful shutdown on SIGINT/SIGTERM
func (s *Server) handleShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("\n🛑 Shutting down server...")

	close(s.stopCh)
	if s.listener != nil {
		s.listener.Close()
	}
	s.store.Close()

	log.Println("✅ Server stopped gracefully")
	os.Exit(0)
}

// Stop stops the server
func (s *Server) Stop() {
	close(s.stopCh)
	if s.listener != nil {
		s.listener.Close()
	}
	s.store.Close()
}
