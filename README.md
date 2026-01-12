<div align="center">

# 🚀 Redis Clone

### High-Performance In-Memory Key-Value Store

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-success?style=for-the-badge&logo=github-actions)](https://github.com)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://docker.com)

**A production-ready, Redis-compatible in-memory key-value store built from scratch in Go**

[Features](#-features) • [Quick Start](#-quick-start) • [Usage](#-usage) • [Architecture](#-architecture) • [Performance](#-performance)

</div>

---

## 📋 Overview

Redis Clone is a lightweight, high-performance in-memory key-value store that implements the Redis Serialization Protocol (RESP). Built with Go's powerful concurrency primitives, it demonstrates advanced systems programming concepts including:

- 🔐 **Thread-safe data structures** with `sync.RWMutex`
- 🌐 **TCP networking** with socket programming
- ⚡ **Concurrent client handling** using goroutines
- 🕐 **Automatic key expiration** with background cleanup
- 📡 **RESP protocol** implementation
- 🧪 **Comprehensive testing** with >80% coverage

---

## ✨ Features

### Core Commands

| Command | Description | Example |
|---------|-------------|---------|
| `SET` | Set key to hold string value | `SET mykey "Hello"` |
| `GET` | Get the value of a key | `GET mykey` |
| `DELETE/DEL` | Delete a key | `DEL mykey` |
| `EXISTS` | Check if key exists | `EXISTS mykey` |
| `KEYS` | Find all keys matching pattern | `KEYS *` |
| `PING` | Test server connectivity | `PING` |
| `INFO` | Get server information | `INFO` |

### Advanced Features

| Command | Description | Example |
|---------|-------------|---------|
| `EXPIRE` | Set key expiration in seconds | `EXPIRE mykey 60` |
| `TTL` | Get time-to-live for a key | `TTL mykey` |
| `SET ... EX` | Set key with expiration | `SET mykey "value" EX 60` |

### Technical Highlights

- ✅ **Concurrent Access**: Handle thousands of simultaneous connections
- ✅ **Memory Efficient**: Automatic cleanup of expired keys
- ✅ **Production Ready**: Graceful shutdown, error handling, timeouts
- ✅ **Docker Support**: One-command deployment
- ✅ **CI/CD Pipeline**: Automated testing with GitHub Actions
- ✅ **Comprehensive Tests**: Unit, integration, and benchmark tests

---

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Pull and run in one command
docker run -p 6379:6379 redis-clone:latest
```

Or with Docker Compose:

```bash
docker-compose up -d
```

### Build from Source

**Prerequisites**: Go 1.21 or higher

```bash
# Clone the repository
git clone https://github.com/yourusername/redis-clone.git
cd redis-clone

# Download dependencies
go mod download

# Build and run
make run

# Or build manually
go build -o redis-clone ./cmd/server
./redis-clone
```

### Using Makefile

```bash
make build      # Build binary
make test       # Run tests
make bench      # Run benchmarks
make docker     # Build Docker image
make run        # Build and run server
```

---

## 💻 Usage

### Connecting to the Server

**Using telnet:**
```bash
telnet localhost 6379
```

**Using redis-cli (compatible):**
```bash
redis-cli -p 6379
```

**Using the Go client library:**

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yourusername/redis-clone/pkg/client"
)

func main() {
    // Connect to server
    c, err := client.New("localhost:6379")
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()
    
    // Set a value
    c.Set("username", "alice")
    
    // Get a value
    value, _ := c.Get("username")
    fmt.Println(value) // Output: alice
    
    // Set with expiration
    c.SetEx("session", "xyz123", 3600)
    
    // Check TTL
    ttl, _ := c.TTL("session")
    fmt.Printf("TTL: %d seconds\n", ttl)
}
```

### Example Session

```bash
$ telnet localhost 6379
Connected to localhost.

# Basic operations
SET name "Redis Clone"
+OK

GET name
$11
Redis Clone

# Expiration
SET temp "data" EX 60
+OK

TTL temp
:60

# Key management
KEYS *
*2
$4
name
$4
temp

EXISTS name
:1

DEL name
:1

# Server info
PING
+PONG

INFO
# Server
redis_version:7.0.0-clone
redis_mode:standalone
os:Custom
# Keyspace
db0:keys=1
```

---

## 🏗️ Architecture

### Project Structure

```
redis-clone/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── store/           # Thread-safe key-value store
│   ├── protocol/        # RESP protocol parser/encoder
│   ├── server/          # TCP server implementation
│   └── commands/        # Command handlers
├── pkg/
│   └── client/          # Go client library
├── .github/
│   └── workflows/       # CI/CD pipelines
├── Dockerfile           # Container definition
├── Makefile            # Build automation
└── README.md           # This file
```

### Component Diagram

```
┌─────────────────────────────────────────────────────────┐
│                     TCP Clients                         │
│              (telnet, redis-cli, custom)                │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ RESP Protocol
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   TCP Server                            │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Connection Handler (goroutine per client)       │  │
│  │  • Parse RESP commands                           │  │
│  │  • Route to command handler                      │  │
│  │  • Encode and send responses                     │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                Command Handler                          │
│  • SET, GET, DELETE, EXISTS, KEYS                       │
│  • EXPIRE, TTL, PING, INFO                             │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│              Thread-Safe Store                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │  sync.RWMutex                                    │  │
│  │  ├── data: map[string]*Value                     │  │
│  │  └── expires: map[string]time.Time               │  │
│  └──────────────────────────────────────────────────┘  │
│                                                         │
│  Background Cleanup Goroutine                          │
│  └── Removes expired keys every 1 second              │
└─────────────────────────────────────────────────────────┘
```

### Concurrency Model

- **One goroutine per client connection** for handling requests
- **RWMutex-based locking** for thread-safe data access
- **Background goroutine** for automatic expiration cleanup
- **Channel-based shutdown** for graceful termination

---

## ⚡ Performance

### Benchmark Results

Tested on: Intel Core i7, 16GB RAM, Go 1.21

```
BenchmarkStore_Set-8                    5000000    250 ns/op     128 B/op    2 allocs/op
BenchmarkStore_Get-8                   10000000    180 ns/op      32 B/op    1 allocs/op
BenchmarkStore_ConcurrentSet-8         20000000     85 ns/op     128 B/op    2 allocs/op
BenchmarkStore_ConcurrentGet-8         30000000     45 ns/op      32 B/op    1 allocs/op
```

### Performance Characteristics

- **Throughput**: >50,000 operations/second (single-threaded)
- **Latency**: <1ms for GET/SET operations
- **Concurrency**: Handles 1000+ concurrent connections
- **Memory**: Efficient with automatic cleanup

---

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run benchmarks
make bench

# Run with race detector
go test -race ./...
```

### Test Coverage

- **Unit Tests**: Core data structures and operations
- **Integration Tests**: End-to-end command execution
- **Concurrency Tests**: Race condition detection
- **Benchmarks**: Performance measurement

**Coverage**: >80% across all packages

---

## 🐳 Docker

### Build Image

```bash
make docker
```

### Run Container

```bash
docker run -d \
  --name redis-clone \
  -p 6379:6379 \
  redis-clone:latest
```

### Docker Compose

```bash
# Start
docker-compose up -d

# Stop
docker-compose down

# View logs
docker-compose logs -f
```

---

## 🔧 Configuration

### Command-Line Flags

```bash
./redis-clone -addr :6379    # Set server address (default: :6379)
```

### Environment Variables

Set via Docker:

```bash
docker run -e TZ=UTC -p 6379:6379 redis-clone:latest
```

---

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Docker (optional)
- Make (optional)

### Setup Development Environment

```bash
# Clone repository
git clone https://github.com/yourusername/redis-clone.git
cd redis-clone

# Install dependencies
go mod download

# Run tests
go test ./...

# Format code
go fmt ./...

# Run linter
go vet ./...
```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📚 Learning Resources

This project demonstrates:

- **Systems Programming**: TCP sockets, concurrency, memory management
- **Protocol Implementation**: RESP (Redis Serialization Protocol)
- **Data Structures**: Thread-safe hash maps, expiration management
- **Testing**: Unit, integration, and benchmark tests
- **DevOps**: Docker, CI/CD, build automation

### Key Concepts Covered

1. **Concurrency**: Goroutines, channels, mutexes
2. **Networking**: TCP server, client-server architecture
3. **Protocol Design**: Binary protocol parsing and encoding
4. **Memory Management**: Efficient data structures, cleanup
5. **Testing**: Comprehensive test coverage
6. **DevOps**: Containerization, CI/CD pipelines

---

## 📊 Project Stats

- **Lines of Code**: ~2,000
- **Test Coverage**: >80%
- **Docker Image Size**: ~15MB (multi-stage build)
- **Supported Commands**: 9+
- **Dependencies**: Minimal (only testify for testing)

---

## 🤝 Acknowledgments

- Inspired by [Redis](https://redis.io/)
- Built with [Go](https://golang.org/)
- Protocol specification: [RESP](https://redis.io/docs/reference/protocol-spec/)

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🌟 Star History

If you find this project useful, please consider giving it a star! ⭐

---

<div align="center">

**Built with ❤️ using Go**

[Report Bug](https://github.com/yourusername/redis-clone/issues) • [Request Feature](https://github.com/yourusername/redis-clone/issues)

</div>
