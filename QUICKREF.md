# Redis Clone - Quick Reference

## 🚀 Quick Start Commands

```bash
# Build
go build -o redis-clone ./cmd/server

# Run
./redis-clone

# Run with custom port
./redis-clone -addr :6380

# Run tests
go test ./...

# Run with Docker
docker build -t redis-clone .
docker run -p 6379:6379 redis-clone
```

## 📝 Supported Commands

### Basic Operations

| Command | Syntax | Example | Description |
|---------|--------|---------|-------------|
| PING | `PING [message]` | `PING` | Test connection |
| SET | `SET key value [EX seconds]` | `SET name "Alice"` | Set key-value |
| GET | `GET key` | `GET name` | Get value |
| DEL | `DEL key` | `DEL name` | Delete key |
| EXISTS | `EXISTS key` | `EXISTS name` | Check if exists |

### Key Management

| Command | Syntax | Example | Description |
|---------|--------|---------|-------------|
| KEYS | `KEYS pattern` | `KEYS *` | List all keys |
| EXPIRE | `EXPIRE key seconds` | `EXPIRE session 3600` | Set expiration |
| TTL | `TTL key` | `TTL session` | Get time-to-live |

### Server

| Command | Syntax | Example | Description |
|---------|--------|---------|-------------|
| INFO | `INFO` | `INFO` | Server stats |

## 🔌 Connection Examples

### Telnet

```bash
$ telnet localhost 6379
SET mykey "Hello World"
+OK
GET mykey
$11
Hello World
```

### Redis-CLI

```bash
$ redis-cli -p 6379
127.0.0.1:6379> SET user:1 "Alice"
OK
127.0.0.1:6379> GET user:1
"Alice"
127.0.0.1:6379> EXPIRE user:1 60
(integer) 1
127.0.0.1:6379> TTL user:1
(integer) 58
```

### Go Client

```go
client, _ := client.New("localhost:6379")
defer client.Close()

client.Set("key", "value")
value, _ := client.Get("key")
client.SetEx("temp", "data", 60)
```

## 🐳 Docker Commands

```bash
# Build image
docker build -t redis-clone .

# Run container
docker run -d -p 6379:6379 --name redis-clone redis-clone

# View logs
docker logs -f redis-clone

# Stop container
docker stop redis-clone

# Remove container
docker rm redis-clone

# Docker Compose
docker-compose up -d
docker-compose down
```

## 🧪 Testing Commands

```bash
# All tests
make test

# With coverage
make coverage

# Benchmarks
make bench

# Race detector
go test -race ./...

# Specific package
go test ./internal/store -v
```

## 📊 Response Types (RESP)

| Type | Prefix | Example |
|------|--------|---------|
| Simple String | `+` | `+OK\r\n` |
| Error | `-` | `-ERR unknown command\r\n` |
| Integer | `:` | `:1000\r\n` |
| Bulk String | `$` | `$5\r\nhello\r\n` |
| Array | `*` | `*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n` |

## 🔧 Development

```bash
# Format code
make fmt

# Lint
make vet

# Clean
make clean

# All checks
make check
```

## 📈 Performance Tips

- Use connection pooling for multiple requests
- Batch operations when possible
- Set appropriate expiration times
- Monitor with INFO command

## 🐛 Troubleshooting

### Server won't start
- Check if port 6379 is already in use
- Verify Go version (1.21+)

### Connection refused
- Ensure server is running
- Check firewall settings
- Verify correct port

### Tests failing
- Run `go mod download`
- Check Go version
- Ensure no other instance running on test ports

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [RESP Protocol](https://redis.io/docs/reference/protocol-spec/)
- [Redis Commands](https://redis.io/commands/)
