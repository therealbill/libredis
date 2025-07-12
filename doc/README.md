# LibRedis Documentation

This directory contains comprehensive documentation for the LibRedis Go client library.

## Documentation Index

- **[commands.md](commands.md)** - Complete Redis command reference organized by data type
- **[api-reference.md](api-reference.md)** - Detailed API documentation with examples
- **[getting-started.md](getting-started.md)** - Quick start guide and basic usage
- **[advanced-features.md](advanced-features.md)** - Advanced features like pipelining, transactions, pub/sub
- **[sentinel.md](sentinel.md)** - Redis Sentinel configuration and usage

## Quick Links

### For New Users
Start with [getting-started.md](getting-started.md) to learn basic connection and command usage.

### For Command Reference
See [commands.md](commands.md) for a complete list of all supported Redis commands organized by data type.

### For Advanced Usage
Check [advanced-features.md](advanced-features.md) for transactions, pipelining, pub/sub, and Lua scripting.

### For High Availability
Review [sentinel.md](sentinel.md) for Redis Sentinel setup and failover handling.

## Redis Version Compatibility

LibRedis supports Redis versions 2.8.13+ with full compatibility for modern Redis features:

- **Redis 2.8+**: Core data types and commands
- **Redis 3.0+**: Cluster and replication improvements
- **Redis 3.2+**: BitField operations, Touch command
- **Redis 4.0+**: Unlink command, modules support
- **Redis 5.0+**: ZPOP commands, streams (planned)
- **Redis 6.0+**: ACL authentication, Hello protocol
- **Redis 6.2+**: Copy command, random field operations
- **Redis 7.0+**: Multi-pop operations
- **Redis 8.0+**: Latest features (tested)

## Recent Updates

### Phase 1 Completed (24 New Commands)
LibRedis has been significantly enhanced with 24 core missing Redis commands:

- **Lists**: LMOVE, BLMOVE, LPOS, LMPOP, BLMPOP
- **Sets**: SMISMEMBER
- **Sorted Sets**: ZPOPMAX, ZPOPMIN, BZPOPMAX, BZPOPMIN, ZRANDMEMBER, ZMSCORE
- **Hashes**: HSTRLEN, HRANDFIELD
- **Keys**: COPY, TOUCH, UNLINK, WAIT
- **Bitmaps**: BITFIELD, BITFIELD_RO, BITPOS
- **Connection**: AUTH (with username), HELLO, RESET

All commands follow Go best practices with proper typing and comprehensive test coverage.