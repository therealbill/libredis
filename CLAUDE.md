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
  - `streams.go` - Redis Streams operations (Phase 2)
  - `geospatial.go` - Geospatial operations (Phase 2)
  - `bitmaps.go` - Bitmap operations (Phase 1)
  - `acl.go` - Access Control Lists and security management (Phase 3)
  - `server.go` - Server management, memory, and latency commands
  - `sentinel.go` - Sentinel support
  - `pubsub.go` - Publish/Subscribe messaging with sharded pub/sub (Phase 3)
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
- Pub/Sub messaging with sharded pub/sub (Phase 3)
- Sentinel support for high availability
- HyperLogLog operations
- Lua scripting support
- Redis Streams with consumer groups (Phase 2)
- Geospatial operations and location search (Phase 2)
- Bitmap operations and bitfield manipulation (Phase 1)
- ACL security management and user authentication (Phase 3)
- Memory management and performance monitoring (Phase 3)
- Database administration and module management (Phase 3)

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

## Recent Updates

### Phase 2 Completed (22 New Commands) - Major Feature Categories
Successfully implemented Redis Streams and Geospatial commands, adding major new capabilities:

**Redis Streams (14 commands):**
- Basic Operations: `XADD`, `XREAD`, `XRANGE`, `XREVRANGE`, `XLEN`, `XDEL`, `XTRIM`
- Consumer Groups: `XGROUP CREATE/DESTROY/SETID`, `XREADGROUP`, `XACK`, `XCLAIM`, `XPENDING`
- Stream Information: `XINFO STREAM/GROUPS/CONSUMERS`

**Geospatial (8 commands):**
- Basic Operations: `GEOADD`, `GEODIST`, `GEOHASH`, `GEOPOS`
- Modern Search: `GEOSEARCH`, `GEOSEARCHSTORE` (Redis 6.2+)
- Legacy Search: `GEORADIUS`, `GEORADIUSBYMEMBER` (deprecated but supported)

### Phase 1 Completed (24 New Commands) - Core Missing Commands
Successfully implemented 24 core missing Redis commands across 6 categories:

**Lists (5 commands):**
- `LMOVE`, `BLMOVE` - Move elements between lists
- `LPOS` (with options) - Find element positions
- `LMPOP`, `BLMPOP` - Pop from multiple lists

**Sets (1 command):**
- `SMISMEMBER` - Check multiple member existence

**Sorted Sets (6 commands):**
- `ZPOPMAX`, `ZPOPMIN` - Pop highest/lowest scored members
- `BZPOPMAX`, `BZPOPMIN` - Blocking variants
- `ZRANDMEMBER` - Get random members
- `ZMSCORE` - Get multiple member scores

**Hashes (2 commands):**
- `HSTRLEN` - Get field value string length
- `HRANDFIELD` - Get random fields

**Keys (4 commands):**
- `COPY` - Copy keys with options
- `TOUCH` - Update last access time
- `UNLINK` - Non-blocking deletion
- `WAIT` - Wait for replication sync

**Bitmaps (3 commands):**
- `BITFIELD`, `BITFIELD_RO` - Bitfield operations
- `BITPOS` - Find bit positions

**Connection (3 commands):**
- `AUTH` with username (ACL support)
- `HELLO` - Protocol handshake
- `RESET` - Reset connection state

### Implementation Notes
- All commands follow existing codebase patterns using `packArgs()` and `ExecuteCommand()`
- Proper Go types and structured options for complex commands
- Comprehensive test coverage with modern Redis compatibility
- Redis version requirements documented for each command
- Phase 2 adds major new data type support (Streams) and location-based operations (Geospatial)

### Phase 3 Completed (29 New Commands) - Security & Management Features
Successfully implemented enterprise-grade security, monitoring, and management capabilities:

**ACL (Access Control) Commands (12 commands):**
- User Management: `ACL SETUSER`, `ACL GETUSER`, `ACL DELUSER`, `ACL USERS`
- Permissions: `ACL CAT`, `ACL WHOAMI`, `ACL LOG`
- Configuration: `ACL LOAD`, `ACL SAVE`, `ACL LIST`
- Utilities: `ACL GENPASS`, `ACL DRYRUN`

**Enhanced Pub/Sub Commands (6 commands):**
- Information: `PUBSUB CHANNELS`, `PUBSUB NUMSUB`, `PUBSUB NUMPAT`
- Sharded Pub/Sub: `SPUBLISH`, `SSUBSCRIBE`, `SUNSUBSCRIBE` (Redis 7.0+)

**Server Management Enhancements (11 commands):**
- Memory Management: `MEMORY USAGE`, `MEMORY STATS`, `MEMORY DOCTOR`, `MEMORY PURGE`
- Latency Monitoring: `LATENCY LATEST`, `LATENCY GRAPH`, `LATENCY RESET`
- Database Management: `SWAPDB`, `REPLICAOF`
- Module Management: `MODULE LIST`

### Implementation Notes
- Enterprise-grade security with ACL user management and rule-based permissions
- Production monitoring with memory usage analysis and latency tracking
- Sharded pub/sub for distributed messaging in Redis 7.0+
- Database administration tools for maintenance operations
- Complete test coverage with Redis version compatibility (2.8.13+ through 7.0+)

### Phase 4 Completed (45 New Commands) - Structured Data & Module Support
Successfully implemented comprehensive Redis Stack module support for advanced data operations:

**JSON Commands (15 commands) - RedisJSON Module:**
- Basic Operations: `JSON.SET`, `JSON.GET`, `JSON.DEL`, `JSON.TYPE`
- Numeric Operations: `JSON.NUMINCRBY`, `JSON.NUMMULTBY`
- String Operations: `JSON.STRAPPEND`, `JSON.STRLEN`
- Array Operations: `JSON.ARRAPPEND`, `JSON.ARRINDEX`, `JSON.ARRINSERT`, `JSON.ARRLEN`, `JSON.ARRPOP`, `JSON.ARRTRIM`
- Object Operations: `JSON.OBJKEYS`, `JSON.OBJLEN`

**Search Commands (8 commands) - RediSearch Module:**
- Index Management: `FT.CREATE`, `FT.DROPINDEX`, `FT.INFO`
- Search Operations: `FT.SEARCH`, `FT.AGGREGATE`, `FT.EXPLAIN`
- Document Management: `FT.ADD`, `FT.DEL` (deprecated in RediSearch 2.0+)

**Time Series Commands (10 commands) - RedisTimeSeries Module:**
- Basic Operations: `TS.CREATE`, `TS.ADD`, `TS.MADD`, `TS.INCRBY`, `TS.DECRBY`
- Query Operations: `TS.RANGE`, `TS.REVRANGE`, `TS.MRANGE`, `TS.MREVRANGE`
- Metadata Operations: `TS.INFO`

**Probabilistic Data Structures (12 commands) - RedisBloom Module:**
- Bloom Filter: `BF.RESERVE`, `BF.ADD`, `BF.MADD`, `BF.EXISTS`, `BF.MEXISTS`
- Cuckoo Filter: `CF.RESERVE`, `CF.ADD`, `CF.EXISTS`, `CF.DEL`
- Count-Min Sketch: `CMS.INITBYDIM`, `CMS.INCRBY`, `CMS.QUERY`

### Implementation Notes
- **Module Dependencies**: Requires Redis Stack or individual modules (RedisJSON, RediSearch, RedisTimeSeries, RedisBloom)
- **Graceful Degradation**: All tests include module availability checks with graceful skipping when modules unavailable
- **API Design**: Follows existing libredis patterns with proper Go typing, structured options, and comprehensive error handling
- **Advanced Features**: Support for complex JSON path operations, full-text search with aggregations, time series with automatic compression, and memory-efficient probabilistic data structures
- **Production Ready**: Complete test coverage with Redis Stack compatibility and real-world usage examples

## Key Files for Development
- `client/redis.go` - Main client entry point and connection management
- `client/lists.go`, `client/sets.go`, `client/sorted_sets.go` - Data structure commands
- `client/hashes.go`, `client/keys.go` - Core Redis operations
- `client/bitmaps.go` - Bitmap operations (new in Phase 1)
- `client/streams.go` - Redis Streams implementation (new in Phase 2)
- `client/geospatial.go` - Geospatial operations (new in Phase 2)
- `client/acl.go` - ACL security and user management (new in Phase 3)
- `client/pubsub.go` - Enhanced pub/sub with sharded messaging (updated in Phase 3)
- `client/server.go` - Server management, memory, and latency monitoring (updated in Phase 3)
- `client/json.go` - JSON document operations with RedisJSON (new in Phase 4)
- `client/search.go` - Full-text search with RediSearch (new in Phase 4)
- `client/timeseries.go` - Time series data with RedisTimeSeries (new in Phase 4)
- `client/probabilistic.go` - Probabilistic data structures with RedisBloom (new in Phase 4)
- `client/connection.go` - Connection and authentication commands
- `structures/command.go` - Redis command definitions and structures  
- `structures/info.go` - Redis server info structures
- `examples/` - Reference implementations for common use cases