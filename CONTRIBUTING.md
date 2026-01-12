# Contributing to Redis Clone

First off, thank you for considering contributing to Redis Clone! 🎉

## 🤝 How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples**
- **Describe the behavior you observed and what you expected**
- **Include logs and error messages**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description of the suggested enhancement**
- **Explain why this enhancement would be useful**
- **List some examples of how it would be used**

### Pull Requests

1. Fork the repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run formatting (`make fmt`)
6. Run linter (`make vet`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## 📝 Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for formatting
- Run `go vet` before committing
- Write meaningful commit messages

### Code Organization

```
redis-clone/
├── cmd/           # Application entry points
├── internal/      # Private application code
├── pkg/           # Public library code
└── scripts/       # Build and utility scripts
```

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
feat: add pub/sub support
fix: handle nil pointer in parser
docs: update README with new examples
test: add benchmarks for concurrent operations
refactor: simplify command handler logic
```

## 🧪 Testing

- Write tests for all new features
- Maintain test coverage above 80%
- Run tests with race detector: `go test -race ./...`
- Add benchmarks for performance-critical code

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run benchmarks
make bench

# Run with race detector
go test -race ./...
```

## 📚 Documentation

- Update README.md for user-facing changes
- Add inline comments for complex logic
- Update architecture diagrams if needed
- Include examples in documentation

## 🔍 Code Review Process

1. All submissions require review
2. Reviewers will check:
   - Code quality and style
   - Test coverage
   - Documentation
   - Performance implications
3. Address review comments
4. Maintainers will merge approved PRs

## 🎯 Development Setup

### Prerequisites

- Go 1.21 or higher
- Docker (optional)
- Make (optional)

### Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/redis-clone.git
cd redis-clone

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build ./cmd/server
```

## 🐛 Debugging

### Running Locally

```bash
# Build and run
make run

# Or manually
go build -o redis-clone ./cmd/server
./redis-clone -addr :6379
```

### Connecting

```bash
# Using telnet
telnet localhost 6379

# Using redis-cli
redis-cli -p 6379
```

## 📋 Project Roadmap

See [Issues](https://github.com/yourusername/redis-clone/issues) for planned features and known issues.

### Potential Enhancements

- [ ] Additional data types (Lists, Sets, Sorted Sets)
- [ ] Persistence (RDB, AOF)
- [ ] Replication
- [ ] Pub/Sub
- [ ] Clustering
- [ ] Lua scripting
- [ ] Transactions

## 💬 Communication

- **Issues**: For bug reports and feature requests
- **Pull Requests**: For code contributions
- **Discussions**: For questions and ideas

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

## 🙏 Thank You!

Your contributions make this project better for everyone. Thank you for taking the time to contribute! ❤️
