# libredis

[![Code Climate](https://codeclimate.com/github/therealbill/libredis/badges/gpa.svg)](https://codeclimate.com/github/therealbill/libredis)
[![GoDoc](https://godoc.org/github.com/therealbill/libredis?status.png)](https://godoc.org/github.com/therealbill/libredis)
[![Build Status](https://travis-ci.org/therealbill/libredis.svg?branch=master)](https://travis-ci.org/therealbill/libredis)

Libredis is a library for interacting with Redis. This includes (a '^'
indicates features planned and/or in dev):

- client connection 
- Sentinel interaction and management
	- Support for SENTINEL commands
	- Ability for client code to manage Sentinel
	- Dedicated control methods^
- Redis INFO parsing
- Inbuilt Sentinel discovery support^
- Inbuilt Redis Cluster support^
- Tested under Go 1.24+ and Redis >= 2.8.13


Libredis is intended to be more than a simple client connection library.
It will include Redis specific custom operations, Structures, and
capabilities suitable for integrating with any Go code which interacts
with Redis ranging from simple CRUD operations to service management.

## Features

### Redis Client

The client code is a fork from
[Goredis](https://github.com/xuyu/goredis). The code is being cleaned up
and rewritten for better performance and developer usability.

#### API Volatility

The original API, while better than other libraries, did not meet what
I felt the API should look like. The info command returned simple
strings, any command which supported variadic arguments required them,
and certain errors should be handled at a lower level and returned in a
more formal manner at the API level. 

As such, the API is currently under a bit of flux. Any command which
accepts variadic paramters such as `ZAdd` will be changing to be of the
simple form with Variadic added for the variadic call. For example, `ZAdd`
will take the key, score, and value while `ZAddVariadic` will take a
`map[string]float64` to variadic commands. This pattern will spread
through the code as I find them or find the time to pre-emptively fix
them.

Where appropriate certain known errors will be converted from raw Redis
errors to typed CommandErrors to provide simpler error handling. It is likely
all commands will return a typed CommandError to provide commonality
across the client API.

#### Notables

* Python Redis Client Like API
* Support [Pipeling](http://godoc.org/github.com/TheRealBill/libredis#Pipelined)
* Support [Transaction](http://godoc.org/github.com/TheRealBill/libredis#Transaction)
* Support [Publish Subscribe](http://godoc.org/github.com/TheRealBill/libredis#PubSub)
* Support [Lua Eval](http://godoc.org/github.com/TheRealBill/libredis#Redis.Eval)
* Support [Connection Pool](http://godoc.org/github.com/TheRealBill/libredis#ConnPool)
* Support [Dial URL-Like](http://godoc.org/github.com/TheRealBill/libredis#DialURL)
* Support for Sentinel commands
* Support Parsing Redis Info commands into Maps and structs
* Support [monitor](http://godoc.org/github.com/TheRealBill/libredis#MonitorCommand), [sort](http://godoc.org/github.com/TheRealBill/libredis#SortCommand), [scan](http://godoc.org/github.com/TheRealBill/libredis#Redis.Scan), [slowlog](http://godoc.org/github.com/TheRealBill/libredis#SlowLog) .etc
* SSL Support! If you have a provider or proxy providing an SSL endpoint you can now connect to it via libredis.
* **Redis Streams Support** - Complete implementation with consumer groups and stream management
* **Geospatial Operations** - Location-based operations with radius and area search capabilities
* **ACL Security Management** - Enterprise-grade user authentication and access control
* **Performance Monitoring** - Memory usage analysis and latency tracking
* **Enhanced Pub/Sub** - Information commands and sharded messaging for Redis 7.0+
* **Comprehensive Command Coverage** - 75 new commands added across all Redis feature categories

## Recent Major Updates

### Phase 2: Major Feature Categories (22 New Commands)
Libredis now includes comprehensive support for Redis's most advanced features:

#### Redis Streams
Complete implementation of Redis Streams including:
- **Basic Operations**: `XADD`, `XREAD`, `XRANGE`, `XREVRANGE`, `XLEN`, `XDEL`, `XTRIM`
- **Consumer Groups**: `XGROUP CREATE/DESTROY/SETID`, `XREADGROUP`, `XACK`, `XCLAIM`, `XPENDING`
- **Stream Information**: `XINFO STREAM/GROUPS/CONSUMERS`

Redis Streams provide powerful event streaming capabilities with consumer group semantics similar to Apache Kafka, enabling reliable message processing at scale.

#### Geospatial Operations
Full geospatial support for location-based applications:
- **Basic Operations**: `GEOADD`, `GEODIST`, `GEOHASH`, `GEOPOS`
- **Modern Search**: `GEOSEARCH`, `GEOSEARCHSTORE` (Redis 6.2+)
- **Legacy Search**: `GEORADIUS`, `GEORADIUSBYMEMBER` (deprecated but supported)

Geospatial operations enable proximity searches, distance calculations, and location-based queries with support for multiple distance units and search geometries.

### Phase 3: Security & Management Features (29 New Commands)
Enterprise-grade capabilities for production Redis deployments:

#### ACL Security Management
Complete access control implementation for Redis 6.0+:
- **User Management**: `ACL SETUSER`, `ACL GETUSER`, `ACL DELUSER`, `ACL USERS`
- **Permissions & Audit**: `ACL CAT`, `ACL WHOAMI`, `ACL LOG`, `ACL DRYRUN`
- **Configuration**: `ACL LOAD`, `ACL SAVE`, `ACL LIST`, `ACL GENPASS`

ACL features enable rule-based permissions, secure password generation, and comprehensive security audit trails for enterprise environments.

#### Enhanced Pub/Sub & Monitoring
Advanced messaging and performance monitoring:
- **Pub/Sub Information**: `PUBSUB CHANNELS`, `PUBSUB NUMSUB`, `PUBSUB NUMPAT`
- **Sharded Pub/Sub**: `SPUBLISH`, `SSUBSCRIBE`, `SUNSUBSCRIBE` (Redis 7.0+)
- **Memory Management**: `MEMORY USAGE`, `MEMORY STATS`, `MEMORY DOCTOR`, `MEMORY PURGE`
- **Latency Monitoring**: `LATENCY LATEST`, `LATENCY GRAPH`, `LATENCY RESET`

#### Database Administration
Professional database management tools:
- **Database Operations**: `SWAPDB`, `REPLICAOF` (modern replacement for SLAVEOF)
- **Module Management**: `MODULE LIST`

### Phase 1: Core Missing Commands (24 New Commands)
Added essential Redis commands across all major data types:
- **Lists**: `LMOVE`, `BLMOVE`, `LPOS`, `LMPOP`, `BLMPOP`
- **Sets**: `SMISMEMBER`
- **Sorted Sets**: `ZPOPMAX`, `ZPOPMIN`, `BZPOPMAX`, `BZPOPMIN`, `ZRANDMEMBER`, `ZMSCORE`
- **Hashes**: `HSTRLEN`, `HRANDFIELD`
- **Keys**: `COPY`, `TOUCH`, `UNLINK`, `WAIT`
- **Bitmaps**: `BITFIELD`, `BITFIELD_RO`, `BITPOS`
- **Connection**: Enhanced `AUTH` (ACL), `HELLO`, `RESET`

## Redis Info

The info package provides functions for parsing the string results of an
Redis info command. When using the libredis/client package these are
unnecessary. This package is useful for those using other Redis client
packages which return strings.

### Sentinel Info

The arguments for `INFO` in Redis when being issued against a Sentinel
returned null for `INFO all`. This has been subsequently fixed to match
the Redis `INFO all` pattern. If using a Redis version of 2.8.13 or
older, use `SentinelInfo` instead of Info to handle this scenario.


## Related Articles

- [Redis Commands](http://redis.io/commands)
- [Redis Protocol](http://redis.io/topics/protocol)
- [Sentinel](http://redis.io/topics/sentinel)
- [GoDoc](http://godoc.org/github.com/TheRealBill/libredis)



# Running Tests


normal test:

	go test

coverage test:

	go test -cover

coverage test with html result:

	go test -coverprofile=cover.out
	go tool cover -html=cover.out


# Running Benchmarks

	go test -test.run=none -test.bench="Benchmark.*"

