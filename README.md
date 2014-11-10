# libredis

[![GoDoc](https://godoc.org/github.com/TheRealBill/libredis?status.png)](https://godoc.org/github.com/TheRealBill/libredis)

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
- Tested under Go 1.3 and Redis 2.8.13


Libredis is intended to be more than a simple client connection library.
It will include Redis specific custom operations, Structures, and
capabilities suitable for integrating with any Go code which interacts
with Redis ranging from simple CRUD operations to service management.

## Features

### Redis Client

The client code is a fork from
[Goredis](https://github.com/xuyu/goredis). The code is being cleaned up
and rewritten for better performance.

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
- [Sentinel](http://redis.io/topics/sentinl)
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

