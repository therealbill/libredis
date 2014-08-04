# libredis

[![GoDoc](https://godoc.org/github.com/therealbill/libredis?status.png)](https://godoc.org/github.com/therealbill/libredis)

Libredis is a library for interacting with Redis. This includes (a '^'
indicates features planned and/or in dev): are scheduled):

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
with Redis ranging from sinple CRUD operations to service management.

## Features

### Redis Client

The client code is a fork from
[Goredis](https://github.com/xuyu/goredis). The code is being cleaned up
and rewritten for better performance.

* Python Redis Client Like API
* Support [Pipeling](http://godoc.org/github.com/therealbill/libredis#Pipelined)
* Support [Transaction](http://godoc.org/github.com/therealbill/libredis#Transaction)
* Support [Publish Subscribe](http://godoc.org/github.com/therealbill/libredis#PubSub)
* Support [Lua Eval](http://godoc.org/github.com/therealbill/libredis#Redis.Eval)
* Support [Connection Pool](http://godoc.org/github.com/therealbill/libredis#ConnPool)
* Support [Dial URL-Like](http://godoc.org/github.com/therealbill/libredis#DialURL)
* Support for Sentinel command
* Support Parsing Redis Info commands into Maps and structs
* Support [monitor](http://godoc.org/github.com/therealbill/libredis#MonitorCommand), [sort](http://godoc.org/github.com/therealbill/libredis#SortCommand), [scan](http://godoc.org/github.com/therealbill/libredis#Redis.Scan), [slowlog](http://godoc.org/github.com/therealbill/libredis#SlowLog) .etc



## Related Articles

- [Redis Commands](http://redis.io/commands)
- [Redis Protocol](http://redis.io/topics/protocol)
- [Sentinel](http://redis.io/topics/sentinl)
- [GoDoc](http://godoc.org/github.com/therealbill/libredis)



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

