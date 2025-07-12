# Redis Commands Reference

This document provides a comprehensive reference of all Redis commands supported by the LibRedis Go client, organized by data type and functionality.

## Table of Contents

- [String Operations](#string-operations)
- [List Operations](#list-operations)
- [Set Operations](#set-operations)
- [Sorted Set Operations](#sorted-set-operations)
- [Hash Operations](#hash-operations)
- [Key Operations](#key-operations)
- [Bitmap Operations](#bitmap-operations)
- [Redis Streams Operations](#redis-streams-operations) **NEW**
- [Geospatial Operations](#geospatial-operations) **NEW**
- [Connection Operations](#connection-operations)
- [Server Operations](#server-operations)
- [HyperLogLog Operations](#hyperloglog-operations)
- [Scripting Operations](#scripting-operations)
- [Pub/Sub Operations](#pubsub-operations)
- [Transaction Operations](#transaction-operations)
- [Sort Operations](#sort-operations)
- [Advanced Features](#advanced-features)

---

## String Operations

String operations work with the most basic Redis data type - binary-safe strings.

### Basic String Commands

| Command | Method | Description |
|---------|--------|-------------|
| GET | `Get(key)` | Gets string value as bytes |
| SET | `Set(key, value)` | Sets string value |
| APPEND | `Append(key, value)` | Appends value to string, returns new length |
| STRLEN | `StrLen(key)` | Returns string length |
| GETSET | `GetSet(key, value)` | Atomically sets and returns old value |
| SETNX | `Setnx(key, value)` | Sets only if key doesn't exist |

### String with Expiration

| Command | Method | Description |
|---------|--------|-------------|
| SETEX | `Setex(key, seconds, value)` | Sets with expiration in seconds |
| PSETEX | `PSetex(key, milliseconds, value)` | Sets with expiration in milliseconds |

### Multiple String Operations

| Command | Method | Description |
|---------|--------|-------------|
| MGET | `MGet(keys...)` | Gets multiple values |
| MSET | `MSet(pairs)` | Sets multiple key-value pairs |
| MSETNX | `MSetnx(pairs)` | Sets multiple only if none exist |

### Numeric Operations

| Command | Method | Description |
|---------|--------|-------------|
| INCR | `Incr(key)` | Increments by 1 |
| INCRBY | `IncrBy(key, increment)` | Increments by specified amount |
| INCRBYFLOAT | `IncrByFloat(key, increment)` | Increments by float |
| DECR | `Decr(key)` | Decrements by 1 |
| DECRBY | `DecrBy(key, decrement)` | Decrements by specified amount |

### Bit Operations

| Command | Method | Description |
|---------|--------|-------------|
| GETBIT | `GetBit(key, offset)` | Gets bit at offset |
| SETBIT | `SetBit(key, offset, value)` | Sets bit at offset |
| BITCOUNT | `BitCount(key, start, end)` | Counts set bits in range |
| BITOP | `BitOp(operation, destkey, keys...)` | Bitwise operations (AND, OR, XOR, NOT) |

### Range Operations

| Command | Method | Description |
|---------|--------|-------------|
| GETRANGE | `GetRange(key, start, end)` | Gets substring |
| SETRANGE | `SetRange(key, offset, value)` | Sets substring at offset |

---

## List Operations

Lists are ordered collections of strings, sorted by insertion order.

### Basic List Operations

| Command | Method | Description |
|---------|--------|-------------|
| LPUSH | `LPush(key, values...)` | Pushes to head |
| RPUSH | `RPush(key, values...)` | Pushes to tail |
| LPOP | `LPop(key)` | Pops from head |
| RPOP | `RPop(key)` | Pops from tail |
| LLEN | `LLen(key)` | Gets list length |
| LINDEX | `LIndex(key, index)` | Gets element at index |
| LRANGE | `LRange(key, start, end)` | Gets range of elements |

### Conditional Push Operations

| Command | Method | Description |
|---------|--------|-------------|
| LPUSHX | `LPushx(key, value)` | Pushes to head only if key exists |
| RPUSHX | `RPushx(key, value)` | Pushes to tail only if key exists |

### List Modification

| Command | Method | Description |
|---------|--------|-------------|
| LINSERT | `LInsert(key, position, pivot, value)` | Inserts before/after pivot |
| LSET | `LSet(key, index, value)` | Sets element at index |
| LREM | `LRem(key, count, value)` | Removes elements equal to value |
| LTRIM | `LTrim(key, start, stop)` | Trims list to range |

### Blocking Operations

| Command | Method | Description |
|---------|--------|-------------|
| BLPOP | `BLPop(keys, timeout)` | Blocking pop from head |
| BRPOP | `BRPop(keys, timeout)` | Blocking pop from tail |
| BRPOPLPUSH | `BRPopLPush(source, dest, timeout)` | Blocking right pop, left push |

### Advanced List Operations (Redis 6.0+)

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| RPOPLPUSH | `RPopLPush(source, dest)` | Atomic right pop, left push | 1.2+ |
| LMOVE | `LMove(source, dest, wherefrom, whereto)` | Atomic move with direction | 6.2+ |
| BLMOVE | `BLMove(source, dest, wherefrom, whereto, timeout)` | Blocking move | 6.2+ |
| LPOS | `LPos(key, element)` | Gets position of element | 6.0.6+ |
| LPOS | `LPosWithOptions(key, element, opts)` | Position with options | 6.0.6+ |
| LMPOP | `LMPop(keys, direction)` | Pops from first non-empty list | 7.0+ |
| BLMPOP | `BLMPop(timeout, keys, direction)` | Blocking multi-pop | 7.0+ |

---

## Set Operations

Sets are unordered collections of unique strings.

### Basic Set Operations

| Command | Method | Description |
|---------|--------|-------------|
| SADD | `SAdd(key, members...)` | Adds members to set |
| SCARD | `SCard(key)` | Gets set cardinality |
| SISMEMBER | `SIsMember(key, member)` | Checks membership |
| SMEMBERS | `SMembers(key)` | Gets all members |
| SREM | `SRem(key, members...)` | Removes members |
| SPOP | `SPop(key)` | Removes and returns random member |
| SRANDMEMBER | `SRandMember(key)` | Returns random member (no removal) |
| SRANDMEMBER | `SRandMemberCount(key, count)` | Returns multiple random members |

### Set Operations

| Command | Method | Description |
|---------|--------|-------------|
| SUNION | `SUnion(keys...)` | Union of sets |
| SUNIONSTORE | `SUnionStore(dest, keys...)` | Union stored to destination |
| SINTER | `SInter(keys...)` | Intersection of sets |
| SINTERSTORE | `SInterStore(dest, keys...)` | Intersection stored to destination |
| SDIFF | `SDiff(keys...)` | Difference of sets |
| SDIFFSTORE | `SDiffStore(dest, keys...)` | Difference stored to destination |
| SMOVE | `SMove(source, dest, member)` | Moves member between sets |

### Modern Set Operations

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| SMISMEMBER | `SMIsMember(key, members...)` | Checks multiple memberships | 6.2+ |

### Set Scanning

| Command | Method | Description |
|---------|--------|-------------|
| SSCAN | `SScan(key, cursor, pattern, count)` | Iterates set members |

---

## Sorted Set Operations

Sorted sets are collections of unique strings, each associated with a score for ordering.

### Basic Sorted Set Operations

| Command | Method | Description |
|---------|--------|-------------|
| ZADD | `ZAdd(key, score, member)` | Adds single member with score |
| ZADD | `ZAddVariadic(key, pairs)` | Adds multiple members with scores |
| ZCARD | `ZCard(key)` | Gets sorted set cardinality |
| ZSCORE | `ZScore(key, member)` | Gets member's score |
| ZINCRBY | `ZIncrBy(key, increment, member)` | Increments member's score |
| ZREM | `ZRem(key, members...)` | Removes members |

### Range Operations

| Command | Method | Description |
|---------|--------|-------------|
| ZRANGE | `ZRange(key, start, stop, withscores)` | Gets range by rank (low to high) |
| ZREVRANGE | `ZRevRange(key, start, stop, withscores)` | Gets range by rank (high to low) |
| ZRANGEBYSCORE | `ZRangeByScore(key, min, max, withscores, limit, offset, count)` | Gets range by score |
| ZREVRANGEBYSCORE | `ZRevRangeByScore(key, max, min, withscores, limit, offset, count)` | Gets range by score (reversed) |

### Rank Operations

| Command | Method | Description |
|---------|--------|-------------|
| ZRANK | `ZRank(key, member)` | Gets rank (0-based, low to high) |
| ZREVRANK | `ZRevRank(key, member)` | Gets rank (0-based, high to low) |

### Count Operations

| Command | Method | Description |
|---------|--------|-------------|
| ZCOUNT | `ZCount(key, min, max)` | Counts members in score range |
| ZLEXCOUNT | `ZLexCount(key, min, max)` | Counts members in lexicographical range |

### Modern Sorted Set Operations

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| ZPOPMAX | `ZPopMax(key)` | Pops member with highest score | 5.0+ |
| ZPOPMIN | `ZPopMin(key)` | Pops member with lowest score | 5.0+ |
| BZPOPMAX | `BZPopMax(keys, timeout)` | Blocking pop highest | 5.0+ |
| BZPOPMIN | `BZPopMin(keys, timeout)` | Blocking pop lowest | 5.0+ |
| ZRANDMEMBER | `ZRandMember(key)` | Returns random member | 6.2+ |
| ZMSCORE | `ZMScore(key, members...)` | Gets multiple scores | 6.2+ |

---

## Hash Operations

Hashes are maps between string fields and string values.

### Basic Hash Operations

| Command | Method | Description |
|---------|--------|-------------|
| HSET | `HSet(key, field, value)` | Sets field value |
| HGET | `HGet(key, field)` | Gets field value |
| HEXISTS | `HExists(key, field)` | Checks if field exists |
| HDEL | `HDel(key, fields...)` | Deletes fields |
| HLEN | `HLen(key)` | Gets number of fields |
| HSETNX | `HSetnx(key, field, value)` | Sets field only if it doesn't exist |

### Multiple Field Operations

| Command | Method | Description |
|---------|--------|-------------|
| HMSET | `HMSet(key, pairs)` | Sets multiple fields |
| HMGET | `HMGet(key, fields...)` | Gets multiple field values |
| HGETALL | `HGetAll(key)` | Gets all fields and values |
| HKEYS | `HKeys(key)` | Gets all field names |
| HVALS | `HVals(key)` | Gets all values |

### Numeric Operations

| Command | Method | Description |
|---------|--------|-------------|
| HINCRBY | `HIncrBy(key, field, increment)` | Increments field by integer |
| HINCRBYFLOAT | `HIncrByFloat(key, field, increment)` | Increments field by float |

### Modern Hash Operations

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| HSTRLEN | `HStrLen(key, field)` | Gets string length of field value | 3.2+ |
| HRANDFIELD | `HRandField(key)` | Returns random field | 6.2+ |

---

## Key Operations

Key operations work with Redis keys regardless of their data type.

### Basic Key Operations

| Command | Method | Description |
|---------|--------|-------------|
| EXISTS | `Exists(key)` | Checks if key exists |
| DEL | `Del(keys...)` | Deletes keys |
| TYPE | `Type(key)` | Gets key type |
| KEYS | `Keys(pattern)` | Gets keys matching pattern |
| RANDOMKEY | `RandomKey()` | Gets random key |

### Key Expiration

| Command | Method | Description |
|---------|--------|-------------|
| EXPIRE | `Expire(key, seconds)` | Sets expiration in seconds |
| EXPIREAT | `ExpireAt(key, timestamp)` | Sets expiration at Unix timestamp |
| PEXPIRE | `PExpire(key, milliseconds)` | Sets expiration in milliseconds |
| PEXPIREAT | `PExpireAt(key, timestamp)` | Sets expiration at Unix timestamp (ms) |
| TTL | `TTL(key)` | Gets time to live in seconds |
| PTTL | `PTTL(key)` | Gets time to live in milliseconds |
| PERSIST | `Persist(key)` | Removes expiration |

### Key Management

| Command | Method | Description |
|---------|--------|-------------|
| RENAME | `Rename(key, newkey)` | Renames key |
| RENAMENX | `Renamenx(key, newkey)` | Renames key only if newkey doesn't exist |
| MOVE | `Move(key, db)` | Moves key to another database |

### Modern Key Operations

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| TOUCH | `Touch(keys...)` | Updates last access time | 3.2.1+ |
| UNLINK | `Unlink(keys...)` | Non-blocking deletion | 4.0+ |
| COPY | `Copy(source, dest)` | Copies key | 6.2+ |
| WAIT | `Wait(numreplicas, timeout)` | Waits for replication | 3.0+ |

---

## Bitmap Operations

Bitmap operations allow manipulation of arbitrary bit fields within strings.

### BitField Operations

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| BITFIELD | `BitField(key, operations)` | Arbitrary bit field operations | 3.2+ |
| BITFIELD | `BitFieldWithOverflow(key, overflow, operations)` | BitField with overflow control | 3.2+ |
| BITFIELD_RO | `BitFieldRO(key, getOps)` | Read-only BitField | 6.0+ |

### Bit Position Operations

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| BITPOS | `BitPos(key, bit)` | Finds position of first bit set to 0/1 | 2.8.7+ |
| BITPOS | `BitPosWithRange(key, bit, opts)` | BitPos with start/end range | 2.8.7+ |

---

## Redis Streams Operations

Redis Streams provide a log-like data structure with consumer group capabilities for building event-driven applications.

### Basic Stream Operations

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| XADD | `XAdd(key, id, fields)` | Adds entry to stream | 5.0+ |
| XADD | `XAddWithOptions(key, id, fields, opts)` | Add with trimming/limits | 5.0+ |
| XREAD | `XRead(streams)` | Reads from streams | 5.0+ |
| XREAD | `XReadWithOptions(streams, opts)` | Read with blocking/count | 5.0+ |
| XRANGE | `XRange(key, start, end)` | Gets range of entries | 5.0+ |
| XRANGE | `XRangeWithOptions(key, start, end, opts)` | Range with count limit | 5.0+ |
| XREVRANGE | `XRevRange(key, end, start)` | Gets range in reverse | 5.0+ |
| XREVRANGE | `XRevRangeWithOptions(key, end, start, opts)` | Reverse range with count | 5.0+ |
| XLEN | `XLen(key)` | Gets stream length | 5.0+ |
| XDEL | `XDel(key, ids...)` | Deletes entries | 5.0+ |
| XTRIM | `XTrim(key, strategy, threshold)` | Trims stream | 5.0+ |
| XTRIM | `XTrimWithOptions(key, strategy, threshold, opts)` | Trim with options | 5.0+ |

### Consumer Group Operations

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| XGROUP CREATE | `XGroupCreate(key, group, id)` | Creates consumer group | 5.0+ |
| XGROUP CREATE | `XGroupCreateWithOptions(key, group, id, opts)` | Create with MKSTREAM | 5.0+ |
| XGROUP DESTROY | `XGroupDestroy(key, group)` | Destroys consumer group | 5.0+ |
| XGROUP SETID | `XGroupSetID(key, group, id)` | Sets group last delivered ID | 5.0+ |
| XREADGROUP | `XReadGroup(group, consumer, streams)` | Reads as consumer | 5.0+ |
| XREADGROUP | `XReadGroupWithOptions(group, consumer, streams, opts)` | Read with options | 5.0+ |
| XACK | `XAck(key, group, ids...)` | Acknowledges messages | 5.0+ |
| XCLAIM | `XClaim(key, group, consumer, minIdle, ids)` | Claims pending messages | 5.0+ |
| XCLAIM | `XClaimWithOptions(key, group, consumer, minIdle, ids, opts)` | Claim with options | 5.0+ |
| XPENDING | `XPending(key, group)` | Gets pending summary | 5.0+ |
| XPENDING | `XPendingWithOptions(key, group, opts)` | Gets detailed pending info | 5.0+ |

### Stream Information

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| XINFO STREAM | `XInfoStream(key)` | Gets stream information | 5.0+ |
| XINFO STREAM | `XInfoStreamFull(key, count)` | Gets detailed stream info | 5.0+ |
| XINFO GROUPS | `XInfoGroups(key)` | Gets consumer groups info | 5.0+ |
| XINFO CONSUMERS | `XInfoConsumers(key, group)` | Gets consumers info | 5.0+ |

### Stream Usage Examples

```go
// Basic stream operations
fields := map[string]string{"name": "Alice", "action": "login"}
id, _ := redis.XAdd("events", client.StreamIDAutoGenerate, fields)

// Reading from streams
streams := map[string]string{"events": "0-0"}
messages, _ := redis.XRead(streams)

// Consumer groups
redis.XGroupCreate("events", "processors", "$")
streams = map[string]string{"events": ">"}
messages, _ = redis.XReadGroup("processors", "worker1", streams)

// Acknowledge processing
redis.XAck("events", "processors", id)
```

---

## Geospatial Operations

Geospatial operations enable location-based applications with support for coordinates, distances, and area searches.

### Basic Geospatial Operations

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| GEOADD | `GeoAdd(key, members)` | Adds geospatial items | 3.2+ |
| GEOADD | `GeoAddWithOptions(key, members, opts)` | Add with NX/XX/CH options | 3.2+ |
| GEODIST | `GeoDist(key, member1, member2)` | Gets distance in meters | 3.2+ |
| GEODIST | `GeoDistWithUnit(key, member1, member2, unit)` | Distance with unit | 3.2+ |
| GEOHASH | `GeoHash(key, members...)` | Gets geohash strings | 3.2+ |
| GEOPOS | `GeoPos(key, members...)` | Gets coordinates | 3.2+ |

### Modern Geospatial Search (Redis 6.2+)

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| GEOSEARCH | `GeoSearch(key, opts)` | Search within area | 6.2+ |
| GEOSEARCHSTORE | `GeoSearchStore(dest, source, opts)` | Search and store results | 6.2+ |

### Legacy Geospatial Search (Deprecated)

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| GEORADIUS | `GeoRadius(key, lon, lat, radius, unit)` | Radius search from coords | 3.2+ |
| GEORADIUS | `GeoRadiusWithOptions(key, lon, lat, radius, unit, opts)` | Radius with options | 3.2+ |
| GEORADIUSBYMEMBER | `GeoRadiusByMember(key, member, radius, unit)` | Radius from member | 3.2+ |
| GEORADIUSBYMEMBER | `GeoRadiusByMemberWithOptions(key, member, radius, unit, opts)` | Member radius with options | 3.2+ |

### Geospatial Usage Examples

```go
// Adding locations
members := []client.GeoMember{
    {Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
    {Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
}
redis.GeoAdd("cities", members)

// Distance calculation
dist, _ := redis.GeoDist("cities", "San Francisco", "New York")

// Modern radius search
opts := client.GeoSearchOptions{
    FromLonLat: &client.GeoCoordinate{Longitude: -122.4194, Latitude: 37.7749},
    ByRadius:   &client.GeoRadius{Radius: 50, Unit: client.GeoUnitKilometers},
    WithCoord:  true,
    WithDist:   true,
}
locations, _ := redis.GeoSearch("cities", opts)
```

---

## Connection Operations

Connection operations manage client connections and authentication.

### Basic Connection Commands

| Command | Method | Description |
|---------|--------|-------------|
| PING | `Ping()` | Tests connection |
| ECHO | `Echo(message)` | Returns message |

### Authentication & Protocol

| Command | Method | Description | Version |
|---------|--------|-------------|---------|
| AUTH | `AuthWithUser(username, password)` | ACL authentication | 6.0+ |
| HELLO | `Hello(protocolVersion)` | Protocol handshake | 6.0+ |
| HELLO | `HelloWithOptions(opts)` | Hello with auth/name options | 6.0+ |
| RESET | `Reset()` | Resets connection state | 6.2+ |

---

## Server Operations

Server operations provide information about and control over the Redis server.

### Server Information

| Command | Method | Description |
|---------|--------|-------------|
| INFO | `Info()` | Gets server info as structured data |
| TIME | `Time()` | Gets server time |
| DBSIZE | `DBSize()` | Gets database size |
| COMMAND | `Command()` | Gets command information |

### Configuration

| Command | Method | Description |
|---------|--------|-------------|
| CONFIG GET | `ConfigGet(parameter)` | Gets configuration |
| CONFIG SET | `ConfigSet(parameter, value)` | Sets configuration |
| CONFIG REWRITE | `ConfigRewrite()` | Rewrites config file |
| CONFIG RESETSTAT | `ConfigResetStat()` | Resets statistics |

### Database Management

| Command | Method | Description |
|---------|--------|-------------|
| FLUSHDB | `FlushDB()` | Flushes current database |
| FLUSHALL | `FlushAll()` | Flushes all databases |
| SAVE | `Save()` | Synchronous save |
| BGSAVE | `BgSave()` | Background save |

---

## HyperLogLog Operations

HyperLogLog is a probabilistic data structure for cardinality estimation.

| Command | Method | Description |
|---------|--------|-------------|
| PFADD | `PFAdd(key, elements...)` | Adds elements to HyperLogLog |
| PFCOUNT | `PFCount(keys...)` | Gets cardinality estimate |
| PFMERGE | `PFMerge(destkey, sourcekeys...)` | Merges HyperLogLogs |

---

## Scripting Operations

Redis supports server-side Lua scripting for atomic operations.

### Script Management

| Command | Method | Description |
|---------|--------|-------------|
| SCRIPT LOAD | `ScriptLoad(script)` | Loads script, returns SHA1 |
| SCRIPT EXISTS | `ScriptExists(scripts...)` | Checks script existence |
| SCRIPT FLUSH | `ScriptFlush()` | Flushes script cache |
| SCRIPT KILL | `ScriptKill()` | Kills running script |

### Script Execution

| Command | Method | Description |
|---------|--------|-------------|
| EVAL | `Eval(script, keys, args)` | Evaluates Lua script |
| EVALSHA | `EvalSha(sha1, keys, args)` | Evaluates by SHA1 |

---

## Pub/Sub Operations

Publish/Subscribe messaging allows real-time message broadcasting.

### Publishing

| Command | Method | Description |
|---------|--------|-------------|
| PUBLISH | `Publish(channel, message)` | Publishes message to channel |

### Subscription Management

Subscriptions are managed through a dedicated `PubSub` object:

```go
pubsub, err := redis.PubSub()
pubsub.Subscribe("channel1", "channel2")
pubsub.PSubscribe("pattern*")
message := pubsub.Receive()
pubsub.Close()
```

---

## Transaction Operations

Transactions provide ACID properties for groups of commands.

### Transaction Management

Transactions are managed through a dedicated `Transaction` object:

```go
txn := redis.Transaction()
txn.Watch("key1", "key2")
txn.Command("SET", "key1", "value1")
txn.Command("INCR", "counter")
results, err := txn.Exec()
txn.Close()
```

---

## Sort Operations

The SORT command provides powerful sorting capabilities for lists, sets, and sorted sets.

### Sorting

Sorting is managed through a dedicated `SortCommand` object:

```go
sort := redis.Sort("mylist")
sort.By("weight_*").Limit(0, 10).Get("object_*").ASC()
results, err := sort.Run()
```

---

## Advanced Features

### Pipelining

Pipelining allows sending multiple commands without waiting for responses:

```go
pipeline := redis.Pipelining()
pipeline.Command("SET", "key1", "value1")
pipeline.Command("SET", "key2", "value2")
responses, err := pipeline.ReceiveAll()
pipeline.Close()
```

### Connection Pooling

LibRedis automatically manages connection pools with configurable parameters:

- Maximum idle connections
- Connection timeouts
- TCP keep-alive settings
- SSL/TLS support

### Error Handling

All commands return proper Go errors with detailed information about Redis-specific errors.

---

## Usage Examples

### Basic Operations

```go
// String operations
redis.Set("key", "value")
value, _ := redis.Get("key")

// List operations
redis.LPush("mylist", "item1", "item2")
items, _ := redis.LRange("mylist", 0, -1)

// Hash operations
redis.HSet("myhash", "field1", "value1")
value, _ := redis.HGet("myhash", "field1")
```

### Advanced Operations

```go
// Modern list operations (Redis 6.2+)
redis.LMove("source", "dest", client.ListDirectionRight, client.ListDirectionLeft)

// Bitmap operations (Redis 3.2+)
ops := []client.BitFieldOperation{
    {Type: "SET", Offset: 0, Value: 100},
    {Type: "GET", Offset: 0},
}
results, _ := redis.BitField("mykey", ops)

// Multi-member operations (Redis 6.2+)
results, _ := redis.SMIsMember("myset", "member1", "member2", "member3")

// Stream operations (Redis 5.0+)
fields := map[string]string{"event": "user_login", "user_id": "123"}
id, _ := redis.XAdd("events", client.StreamIDAutoGenerate, fields)

// Consumer group processing
redis.XGroupCreate("events", "processors", "$")
streams := map[string]string{"events": ">"}
messages, _ := redis.XReadGroup("processors", "worker1", streams)
redis.XAck("events", "processors", id)

// Geospatial operations (Redis 3.2+)
cities := []client.GeoMember{
    {Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
    {Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
}
redis.GeoAdd("cities", cities)

// Location-based search (Redis 6.2+)
searchOpts := client.GeoSearchOptions{
    FromLonLat: &client.GeoCoordinate{Longitude: -122.4194, Latitude: 37.7749},
    ByRadius:   &client.GeoRadius{Radius: 100, Unit: client.GeoUnitKilometers},
    WithCoord:  true,
}
nearbyLocations, _ := redis.GeoSearch("cities", searchOpts)
```

This comprehensive command reference covers all Redis operations supported by LibRedis, including the complete Redis Streams implementation and geospatial operations, bringing support for modern Redis 5.0+ through 7.0+ features.