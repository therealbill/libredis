# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Libredis is a comprehensive Go library for Redis client functionality, forked from Goredis with significant API improvements. It provides a Python Redis Client-like API and includes advanced features like Sentinel support, Redis INFO parsing, and SSL connections.

## Development Commands

### Testing
```bash
# Standard tests
go test

# Coverage testing  
go test -cover

# HTML coverage report
go test -coverprofile=cover.out
go tool cover -html=cover.out

# Integration tests (requires Redis server running)
go test -tags integration ./client ./info
```

### Benchmarking
```bash
go test -test.run=none -test.bench="Benchmark.*"
```

## Architecture

### Package Structure
- **client/**: Core Redis client implementation organized by Redis data types
  - `redis.go` - Main client with connection handling
  - `connection.go` - Basic connection operations (Echo, Ping)
  - `strings.go`, `hashes.go`, `lists.go`, `sets.go`, `sorted_sets.go` - Redis data type operations
  - `server.go` - Server management commands
  - `sentinel.go` - Sentinel support
  - `pubsub.go` - Publish/Subscribe messaging
  - `transactions.go` - MULTI/EXEC transactions
  - `pipelining.go` - Batch operations
  - `scripting.go` - Lua script evaluation

- **info/**: Redis INFO command parsing utilities
- **structures/**: Common data structures and command definitions
- **examples/**: Usage examples including Redis info display and Sentinel CLI

### Key Features
- Full Redis command set with typed responses
- Connection pooling with configurable timeouts
- SSL/TLS support for secure connections
- URL-based connection configuration (`DialURL`)
- Pipelining and transactions
- Pub/Sub messaging
- Sentinel support for high availability
- HyperLogLog operations
- Lua scripting support

### API Design Notes
- **API Volatility Warning**: Commands with variadic parameters are being refactored (e.g., `ZAdd` will become simple form with `ZAddVariadic` for variadic calls)
- Moving toward typed `CommandError` returns for better error handling
- Python Redis Client-inspired API design patterns

## Development Setup

### Prerequisites
- Go 1.24+ (latest stable release, specified in go.mod)
- Redis server >= 2.8.13 for full feature compatibility
- Redis server must be running locally for integration tests

### Go Modules
- Project uses Go modules (go.mod) for dependency management
- Module path: `github.com/therealbill/libredis`
- Use `go mod tidy` to manage dependencies

### CI/CD
- Travis CI runs integration tests with Go 1.24 using `go test -tags integration ./client ./info`
- Code Climate integration with golint and duplication detection
- Coverage reporting integrated with Code Climate

### Testing Strategy
- Standard tests can run without Redis server
- Integration tests require `-tags integration` and running Redis server
- Each Redis data type has dedicated test files
- Comprehensive coverage including benchmarks

## Key Files for Development
- `client/redis.go` - Main client entry point and connection management
- `structures/command.go` - Redis command definitions and structures  
- `structures/info.go` - Redis server info structures
- `examples/` - Reference implementations for common use cases