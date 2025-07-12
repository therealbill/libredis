# API Reference

This document provides detailed API documentation for LibRedis, including method signatures, parameters, return values, and usage examples.

## Table of Contents

- [Connection Management](#connection-management)
- [Core Data Types](#core-data-types)
- [String Operations](#string-operations)
- [List Operations](#list-operations)
- [Set Operations](#set-operations)
- [Sorted Set Operations](#sorted-set-operations)
- [Hash Operations](#hash-operations)
- [Key Operations](#key-operations)
- [Bitmap Operations](#bitmap-operations)
- [Redis Streams Operations](#redis-streams-operations) **NEW**
- [Geospatial Operations](#geospatial-operations) **NEW**
- [JSON Operations](#json-operations) **PHASE 4**
- [Search Operations](#search-operations) **PHASE 4**
- [Time Series Operations](#time-series-operations) **PHASE 4**
- [Probabilistic Data Structures](#probabilistic-data-structures) **PHASE 4**
- [Connection Commands](#connection-commands)
- [Server Operations](#server-operations)
- [Advanced Features](#advanced-features)
- [Error Handling](#error-handling)

## Connection Management

### DialConfig

Configuration structure for Redis connections:

```go
type DialConfig struct {
    Network       string        // "tcp" or "unix"
    Address       string        // "host:port" or "/path/to/socket"
    Database      int           // Redis database number (0-15)
    Password      string        // Redis password
    Timeout       time.Duration // Connection timeout
    MaxIdle       int           // Maximum idle connections in pool
    SSL           bool          // Enable SSL/TLS
    SSLSkipVerify bool          // Skip SSL certificate verification
    SSLCert       string        // SSL certificate file path
    SSLKey        string        // SSL private key file path
    SSLCA         string        // SSL CA certificate file path
    TCPKeepAlive  int           // TCP keep-alive interval (seconds)
}
```

### Connection Functions

#### DialWithConfig

```go
func DialWithConfig(config *DialConfig) (*Redis, error)
```

Creates a new Redis connection with specified configuration.

**Parameters:**
- `config`: Pointer to DialConfig struct

**Returns:**
- `*Redis`: Redis client instance
- `error`: Connection error, if any

**Example:**
```go
config := &client.DialConfig{
    Network:  "tcp",
    Address:  "localhost:6379",
    Database: 0,
    Timeout:  5 * time.Second,
    MaxIdle:  10,
}

redis, err := client.DialWithConfig(config)
if err != nil {
    log.Fatal(err)
}
defer redis.Close()
```

#### DialURL

```go
func DialURL(url string) (*Redis, error)
```

Creates a Redis connection from a URL string.

**Parameters:**
- `url`: Redis URL in format `tcp://[auth:password@]host:port[/db][?param=value]`

**Returns:**
- `*Redis`: Redis client instance
- `error`: Connection error, if any

**Example:**
```go
redis, err := client.DialURL("tcp://auth:password@localhost:6379/0?timeout=5s&maxidle=10")
```

## Core Data Types

### Reply

The Reply struct represents responses from Redis:

```go
type Reply struct {
    Type   int       // Reply type constant
    Error  string    // Error message (for ErrorReply)
    Multi  []*Reply  // Multi-bulk reply data
    Bulk   []byte    // Bulk reply data
    Status string    // Status reply data
    Int    int64     // Integer reply data
}
```

**Reply Type Constants:**
- `ErrorReply`: Redis error response
- `StatusReply`: Simple string response (e.g., "OK")
- `IntegerReply`: Numeric response
- `BulkReply`: Binary data response
- `MultiReply`: Array of replies

### Helper Methods

#### IntegerValue

```go
func (r *Reply) IntegerValue() (int64, error)
```

Extracts integer value from reply.

#### StringValue

```go
func (r *Reply) StringValue() (string, error)
```

Extracts string value from reply.

#### BoolValue

```go
func (r *Reply) BoolValue() (bool, error)
```

Extracts boolean value from reply (1 = true, 0 = false).

#### ListValue

```go
func (r *Reply) ListValue() ([]string, error)
```

Extracts string array from multi-bulk reply.

## String Operations

### Basic String Commands

#### Set

```go
func (r *Redis) Set(key, value string) error
```

Sets the string value of a key.

**Parameters:**
- `key`: The key name
- `value`: The value to set

**Returns:**
- `error`: Command error, if any

**Example:**
```go
err := redis.Set("username", "john_doe")
```

#### Get

```go
func (r *Redis) Get(key string) ([]byte, error)
```

Gets the value of a key as bytes.

**Parameters:**
- `key`: The key name

**Returns:**
- `[]byte`: The value as bytes (nil if key doesn't exist)
- `error`: Command error, if any

**Example:**
```go
value, err := redis.Get("username")
if err != nil {
    return err
}
if value != nil {
    fmt.Println("Username:", string(value))
}
```

#### GetInt

```go
func (r *Redis) GetInt(key string) (int64, error)
```

Gets the value of a key as an integer.

**Example:**
```go
count, err := redis.GetInt("counter")
```

### String Operations with Expiration

#### Setex

```go
func (r *Redis) Setex(key string, seconds int, value string) error
```

Sets a key with expiration in seconds.

**Parameters:**
- `key`: The key name
- `seconds`: Expiration time in seconds
- `value`: The value to set

**Example:**
```go
err := redis.Setex("session_token", 3600, "abc123")
```

#### PSetex

```go
func (r *Redis) PSetex(key string, milliseconds int, value string) error
```

Sets a key with expiration in milliseconds.

### Multiple String Operations

#### MGet

```go
func (r *Redis) MGet(keys ...string) ([][]byte, error)
```

Gets multiple key values in a single command.

**Parameters:**
- `keys`: Variable number of key names

**Returns:**
- `[][]byte`: Array of values (nil elements for non-existent keys)
- `error`: Command error, if any

**Example:**
```go
values, err := redis.MGet("key1", "key2", "key3")
```

#### MSet

```go
func (r *Redis) MSet(pairs map[string]string) error
```

Sets multiple key-value pairs atomically.

**Parameters:**
- `pairs`: Map of key-value pairs

**Example:**
```go
pairs := map[string]string{
    "key1": "value1",
    "key2": "value2",
}
err := redis.MSet(pairs)
```

### Numeric Operations

#### Incr

```go
func (r *Redis) Incr(key string) (int64, error)
```

Increments the integer value of a key by 1.

**Returns:**
- `int64`: The new value after increment

**Example:**
```go
newValue, err := redis.Incr("counter")
```

#### IncrBy

```go
func (r *Redis) IncrBy(key string, increment int64) (int64, error)
```

Increments the integer value of a key by the given amount.

#### IncrByFloat

```go
func (r *Redis) IncrByFloat(key string, increment float64) (float64, error)
```

Increments the float value of a key by the given amount.

## List Operations

### Basic List Operations

#### LPush

```go
func (r *Redis) LPush(key string, values ...string) (int64, error)
```

Prepends values to a list.

**Parameters:**
- `key`: The list key
- `values`: Variable number of values to prepend

**Returns:**
- `int64`: Length of the list after the operation
- `error`: Command error, if any

**Example:**
```go
length, err := redis.LPush("tasks", "task1", "task2", "task3")
```

#### RPush

```go
func (r *Redis) RPush(key string, values ...string) (int64, error)
```

Appends values to a list.

#### LPop

```go
func (r *Redis) LPop(key string) ([]byte, error)
```

Removes and returns the first element of a list.

**Returns:**
- `[]byte`: The popped element (nil if list is empty)
- `error`: Command error, if any

#### LRange

```go
func (r *Redis) LRange(key string, start, end int64) ([]string, error)
```

Gets a range of elements from a list.

**Parameters:**
- `key`: The list key
- `start`: Start index (0-based)
- `end`: End index (-1 for last element)

**Returns:**
- `[]string`: Array of elements in the range

**Example:**
```go
elements, err := redis.LRange("mylist", 0, -1) // Get all elements
```

### Modern List Operations (Redis 6.2+)

#### LMove

```go
func (r *Redis) LMove(source, destination, wherefrom, whereto string) (string, error)
```

Atomically moves an element from one list to another.

**Parameters:**
- `source`: Source list key
- `destination`: Destination list key
- `wherefrom`: Direction to pop from source ("LEFT" or "RIGHT")
- `whereto`: Direction to push to destination ("LEFT" or "RIGHT")

**Returns:**
- `string`: The moved element

**Example:**
```go
element, err := redis.LMove("source", "dest", 
                           client.ListDirectionRight, 
                           client.ListDirectionLeft)
```

#### LPos

```go
func (r *Redis) LPos(key, element string) (int64, error)
```

Returns the index of the first matching element.

**Parameters:**
- `key`: The list key
- `element`: The element to find

**Returns:**
- `int64`: Index of the element (-1 if not found)

#### LPosWithOptions

```go
func (r *Redis) LPosWithOptions(key, element string, opts LPosOptions) ([]int64, error)
```

Returns positions with additional options.

**LPosOptions:**
```go
type LPosOptions struct {
    Rank   int  // Which occurrence to find (1st, 2nd, etc.)
    Count  int  // Maximum number of positions to return
    MaxLen int  // Maximum number of elements to examine
}
```

### List Constants

```go
const (
    ListDirectionLeft  = "LEFT"
    ListDirectionRight = "RIGHT"
)
```

## Set Operations

### Basic Set Operations

#### SAdd

```go
func (r *Redis) SAdd(key string, members ...string) (int64, error)
```

Adds members to a set.

**Returns:**
- `int64`: Number of new members added

**Example:**
```go
added, err := redis.SAdd("tags", "redis", "database", "nosql")
```

#### SMembers

```go
func (r *Redis) SMembers(key string) ([]string, error)
```

Gets all members of a set.

#### SIsMember

```go
func (r *Redis) SIsMember(key, member string) (bool, error)
```

Checks if a member exists in a set.

### Modern Set Operations (Redis 6.2+)

#### SMIsMember

```go
func (r *Redis) SMIsMember(key string, members ...string) ([]bool, error)
```

Checks if multiple members exist in a set.

**Parameters:**
- `key`: The set key
- `members`: Variable number of members to check

**Returns:**
- `[]bool`: Array of boolean values indicating membership

**Example:**
```go
results, err := redis.SMIsMember("myset", "member1", "member2", "member3")
// results[0] = true if member1 exists, false otherwise
```

## Sorted Set Operations

### Basic Sorted Set Operations

#### ZAdd

```go
func (r *Redis) ZAdd(key string, score float64, member string) (int64, error)
```

Adds a member with score to a sorted set.

#### ZAddVariadic

```go
func (r *Redis) ZAddVariadic(key string, pairs map[string]float64) (int64, error)
```

Adds multiple members with scores.

**Example:**
```go
scores := map[string]float64{
    "player1": 100.5,
    "player2": 87.2,
    "player3": 95.0,
}
added, err := redis.ZAddVariadic("leaderboard", scores)
```

#### ZRange

```go
func (r *Redis) ZRange(key string, start, stop int64, withscores bool) ([]string, error)
```

Gets members in a score range by rank.

**Parameters:**
- `withscores`: If true, includes scores in the result

### Modern Sorted Set Operations

#### ZPopMax

```go
func (r *Redis) ZPopMax(key string) (ZMember, error)
```

Pops the member with the highest score.

**ZMember struct:**
```go
type ZMember struct {
    Member string
    Score  float64
}
```

#### ZMScore

```go
func (r *Redis) ZMScore(key string, members ...string) ([]float64, error)
```

Gets scores for multiple members.

**Example:**
```go
scores, err := redis.ZMScore("leaderboard", "player1", "player2")
// scores[0] = score for player1, scores[1] = score for player2
```

## Hash Operations

### Basic Hash Operations

#### HSet

```go
func (r *Redis) HSet(key, field, value string) (bool, error)
```

Sets a field in a hash.

**Returns:**
- `bool`: true if field is new, false if updated

#### HGet

```go
func (r *Redis) HGet(key, field string) ([]byte, error)
```

Gets a field value from a hash.

#### HMSet

```go
func (r *Redis) HMSet(key string, pairs map[string]string) error
```

Sets multiple fields in a hash.

#### HGetAll

```go
func (r *Redis) HGetAll(key string) (map[string]string, error)
```

Gets all field-value pairs from a hash.

### Modern Hash Operations

#### HStrLen

```go
func (r *Redis) HStrLen(key, field string) (int64, error)
```

Gets the string length of a field value.

#### HRandField

```go
func (r *Redis) HRandField(key string) (string, error)
```

Returns a random field from the hash.

## Key Operations

### Basic Key Operations

#### Exists

```go
func (r *Redis) Exists(key string) (bool, error)
```

Checks if a key exists.

#### Del

```go
func (r *Redis) Del(keys ...string) (int64, error)
```

Deletes one or more keys.

**Returns:**
- `int64`: Number of keys that were deleted

#### Type

```go
func (r *Redis) Type(key string) (string, error)
```

Gets the type of a key.

**Possible return values:**
- "string", "list", "set", "zset", "hash", "stream", "none"

### Key Expiration

#### Expire

```go
func (r *Redis) Expire(key string, seconds int) (bool, error)
```

Sets expiration time in seconds.

#### TTL

```go
func (r *Redis) TTL(key string) (int64, error)
```

Gets time to live in seconds.

**Returns:**
- `int64`: TTL in seconds (-1 if no expiration, -2 if key doesn't exist)

### Modern Key Operations

#### Copy

```go
func (r *Redis) Copy(source, destination string) (bool, error)
```

Copies a key to another key.

#### CopyWithOptions

```go
func (r *Redis) CopyWithOptions(source, destination string, opts CopyOptions) (bool, error)
```

Copies a key with additional options.

**CopyOptions:**
```go
type CopyOptions struct {
    DestinationDB int  // Target database
    Replace       bool // Replace if destination exists
}
```

#### Touch

```go
func (r *Redis) Touch(keys ...string) (int64, error)
```

Updates the last access time of keys.

#### Unlink

```go
func (r *Redis) Unlink(keys ...string) (int64, error)
```

Non-blocking deletion of keys.

## Bitmap Operations

### BitField Operations

#### BitField

```go
func (r *Redis) BitField(key string, operations []BitFieldOperation) ([]int64, error)
```

Performs arbitrary bit field operations.

**BitFieldOperation:**
```go
type BitFieldOperation struct {
    Type   string      // "GET", "SET", "INCRBY"
    Offset int64       // Bit offset
    Value  interface{} // Value for SET/INCRBY operations
}
```

**Example:**
```go
ops := []client.BitFieldOperation{
    {Type: "SET", Offset: 0, Value: 100},
    {Type: "GET", Offset: 0},
    {Type: "INCRBY", Offset: 8, Value: 50},
}
results, err := redis.BitField("mykey", ops)
```

#### BitPos

```go
func (r *Redis) BitPos(key string, bit int) (int64, error)
```

Finds the position of the first bit set to the specified value.

**Parameters:**
- `bit`: 0 or 1

## Redis Streams Operations

Redis Streams provide a log-like data structure for event streaming and message processing.

### Stream Data Types

#### StreamEntry

```go
type StreamEntry struct {
    ID     string
    Fields map[string]string
}
```

Represents a single entry in a Redis stream.

#### StreamMessage

```go
type StreamMessage struct {
    Stream  string
    Entries []StreamEntry
}
```

Represents messages from one stream.

#### XAddOptions

```go
type XAddOptions struct {
    NoMkStream  bool   // NOMKSTREAM option
    MaxLen      int64  // MAXLEN option
    MinID       string // MINID option
    Approximate bool   // ~ modifier for MAXLEN/MINID
    Limit       int64  // LIMIT option
}
```

Options for XADD command.

### Basic Stream Operations

#### XAdd

```go
func (r *Redis) XAdd(key, id string, fields map[string]string) (string, error)
func (r *Redis) XAddWithOptions(key, id string, fields map[string]string, opts XAddOptions) (string, error)
```

Adds a new entry to a stream.

**Parameters:**
- `key`: Stream key
- `id`: Entry ID ("*" for auto-generation)
- `fields`: Field-value pairs for the entry
- `opts`: Additional options (for WithOptions variant)

**Returns:**
- `string`: Generated or specified entry ID
- `error`: Error, if any

**Example:**
```go
fields := map[string]string{"name": "Alice", "action": "login"}
id, err := redis.XAdd("events", client.StreamIDAutoGenerate, fields)
```

#### XRead

```go
func (r *Redis) XRead(streams map[string]string) ([]StreamMessage, error)
func (r *Redis) XReadWithOptions(streams map[string]string, opts XReadOptions) ([]StreamMessage, error)
```

Reads entries from one or more streams.

**Parameters:**
- `streams`: Map of stream names to starting IDs
- `opts`: Reading options (for WithOptions variant)

**Returns:**
- `[]StreamMessage`: Array of stream messages
- `error`: Error, if any

#### XLen

```go
func (r *Redis) XLen(key string) (int64, error)
```

Returns the number of entries in a stream.

### Consumer Group Operations

#### XGroupCreate

```go
func (r *Redis) XGroupCreate(key, groupname, id string) error
func (r *Redis) XGroupCreateWithOptions(key, groupname, id string, opts XGroupCreateOptions) error
```

Creates a new consumer group.

**Parameters:**
- `key`: Stream key
- `groupname`: Consumer group name
- `id`: Starting ID for the group ("$" for latest)
- `opts`: Additional options (for WithOptions variant)

#### XReadGroup

```go
func (r *Redis) XReadGroup(group, consumer string, streams map[string]string) ([]StreamMessage, error)
func (r *Redis) XReadGroupWithOptions(group, consumer string, streams map[string]string, opts XReadGroupOptions) ([]StreamMessage, error)
```

Reads entries as a consumer group member.

**Parameters:**
- `group`: Consumer group name
- `consumer`: Consumer name
- `streams`: Map of stream names to IDs (">" for new messages)
- `opts`: Reading options (for WithOptions variant)

#### XAck

```go
func (r *Redis) XAck(key, group string, ids ...string) (int64, error)
```

Acknowledges processing of messages.

**Parameters:**
- `key`: Stream key
- `group`: Consumer group name
- `ids`: Message IDs to acknowledge

**Returns:**
- `int64`: Number of acknowledged messages
- `error`: Error, if any

## Geospatial Operations

Geospatial operations enable location-based applications with coordinate storage and proximity searches.

### Geospatial Data Types

#### GeoMember

```go
type GeoMember struct {
    Longitude float64
    Latitude  float64
    Member    string
}
```

Represents a geospatial member with coordinates.

#### GeoCoordinate

```go
type GeoCoordinate struct {
    Longitude float64
    Latitude  float64
}
```

Represents longitude/latitude coordinates.

#### GeoLocation

```go
type GeoLocation struct {
    Member      string
    Coordinates *GeoCoordinate
    Distance    *float64
    Hash        *int64
}
```

Represents a geospatial search result.

#### GeoSearchOptions

```go
type GeoSearchOptions struct {
    FromMember *string        // FROMMEMBER option
    FromLonLat *GeoCoordinate // FROMLONLAT option
    ByRadius   *GeoRadius     // BYRADIUS option
    ByBox      *GeoBox        // BYBOX option
    Order      string         // ASC or DESC
    Count      int64          // COUNT option
    Any        bool           // ANY option
    WithCoord  bool           // WITHCOORD option
    WithDist   bool           // WITHDIST option
    WithHash   bool           // WITHHASH option
}
```

Options for geospatial search operations.

### Basic Geospatial Operations

#### GeoAdd

```go
func (r *Redis) GeoAdd(key string, members []GeoMember) (int64, error)
func (r *Redis) GeoAddWithOptions(key string, members []GeoMember, opts GeoAddOptions) (int64, error)
```

Adds geospatial items to a geospatial index.

**Parameters:**
- `key`: Geospatial index key
- `members`: Array of geospatial members to add
- `opts`: Additional options (for WithOptions variant)

**Returns:**
- `int64`: Number of added elements
- `error`: Error, if any

**Example:**
```go
members := []client.GeoMember{
    {Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
    {Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
}
count, err := redis.GeoAdd("cities", members)
```

#### GeoDist

```go
func (r *Redis) GeoDist(key, member1, member2 string) (float64, error)
func (r *Redis) GeoDistWithUnit(key, member1, member2, unit string) (float64, error)
```

Returns the distance between two geospatial members.

**Parameters:**
- `key`: Geospatial index key
- `member1`, `member2`: Member names
- `unit`: Distance unit (M, KM, FT, MI)

**Returns:**
- `float64`: Distance in specified unit
- `error`: Error, if any

#### GeoSearch

```go
func (r *Redis) GeoSearch(key string, opts GeoSearchOptions) ([]GeoLocation, error)
```

Searches for members within a specified area.

**Parameters:**
- `key`: Geospatial index key
- `opts`: Search options specifying center point and search area

**Returns:**
- `[]GeoLocation`: Array of matching locations
- `error`: Error, if any

**Example:**
```go
opts := client.GeoSearchOptions{
    FromLonLat: &client.GeoCoordinate{Longitude: -122.4194, Latitude: 37.7749},
    ByRadius:   &client.GeoRadius{Radius: 50, Unit: client.GeoUnitKilometers},
    WithCoord:  true,
    WithDist:   true,
}
locations, err := redis.GeoSearch("cities", opts)
```

## Connection Commands

### Basic Connection

#### Ping

```go
func (r *Redis) Ping() error
```

Tests the connection to Redis.

#### Echo

```go
func (r *Redis) Echo(message string) (string, error)
```

Returns the given message.

### Authentication (Redis 6.0+)

#### AuthWithUser

```go
func (r *Redis) AuthWithUser(username, password string) error
```

Authenticates with username and password (ACL).

#### Hello

```go
func (r *Redis) Hello(protocolVersion int) (map[string]interface{}, error)
```

Performs protocol handshake.

#### HelloWithOptions

```go
func (r *Redis) HelloWithOptions(opts HelloOptions) (map[string]interface{}, error)
```

Hello with additional options.

**HelloOptions:**
```go
type HelloOptions struct {
    ProtocolVersion int
    Username        string
    Password        string
    ClientName      string
}
```

## Advanced Features

### Pipelining

#### Pipelining

```go
func (r *Redis) Pipelining() *Pipelined
```

Creates a new pipeline for batch operations.

**Pipelined methods:**
- `Command(args ...interface{}) error`
- `Receive() (*Reply, error)`
- `ReceiveAll() ([]*Reply, error)`
- `Close() error`

### Transactions

#### Transaction

```go
func (r *Redis) Transaction() *Transaction
```

Creates a new transaction.

**Transaction methods:**
- `Watch(keys ...string) error`
- `UnWatch() error`
- `Command(args ...interface{}) error`
- `Exec() ([]*Reply, error)`
- `Discard() error`
- `Close() error`

### Pub/Sub

#### PubSub

```go
func (r *Redis) PubSub() (*PubSubClient, error)
```

Creates a pub/sub client.

**PubSubClient methods:**
- `Subscribe(channels ...string) error`
- `PSubscribe(patterns ...string) error`
- `UnSubscribe(channels ...string) error`
- `Receive() []string`
- `Close() error`

## JSON Operations

**Phase 4 Feature - Requires RedisJSON Module**

LibRedis supports Redis JSON operations for storing, retrieving, and manipulating JSON documents.

### JSONSet

```go
func (r *Redis) JSONSet(key, path string, value interface{}, options ...*JSONSetOptions) (string, error)
```

Set JSON value at path in key.

**Parameters:**
- `key`: Redis key
- `path`: JSON path (e.g., ".", ".field", ".array[0]")
- `value`: JSON value to set
- `options`: Optional settings (NX, XX)

**Example:**
```go
// Set entire document
result, err := redis.JSONSet("user:1", ".", `{"name": "John", "age": 30}`)

// Set specific field
result, err := redis.JSONSet("user:1", ".email", `"john@example.com"`)

// Conditional set (only if path doesn't exist)
result, err := redis.JSONSet("user:1", ".phone", `"555-0123"`, &client.JSONSetOptions{NX: true})
```

### JSONGet

```go
func (r *Redis) JSONGet(key string, options ...*JSONGetOptions) ([]byte, error)
```

Get JSON value from key.

**Example:**
```go
// Get entire document
data, err := redis.JSONGet("user:1")

// Get specific paths with formatting
data, err := redis.JSONGet("user:1", &client.JSONGetOptions{
    Paths:   []string{".name", ".age"},
    Indent:  "  ",
    NewLine: "\n",
})
```

### Other JSON Commands

- `JSONDel(key, path)` - Delete JSON value at path
- `JSONType(key, path)` - Get type of JSON value
- `JSONNumIncrBy(key, path, number)` - Increment numeric value
- `JSONNumMultBy(key, path, number)` - Multiply numeric value
- `JSONStrAppend(key, path, jsonString)` - Append to string value
- `JSONStrLen(key, path)` - Get string length
- `JSONArrAppend(key, path, values...)` - Append to array
- `JSONArrIndex(key, path, value, startStop...)` - Find array index
- `JSONArrInsert(key, path, index, values...)` - Insert into array
- `JSONArrLen(key, path)` - Get array length
- `JSONArrPop(key, path, index...)` - Pop from array
- `JSONArrTrim(key, path, start, stop)` - Trim array
- `JSONObjKeys(key, path)` - Get object keys
- `JSONObjLen(key, path)` - Get object length

## Search Operations

**Phase 4 Feature - Requires RediSearch Module**

LibRedis supports full-text search capabilities with index management and advanced query features.

### FTCreate

```go
func (r *Redis) FTCreate(index string, schema []FTFieldSchema, options ...*FTCreateOptions) (string, error)
```

Create a search index.

**Example:**
```go
schema := []client.FTFieldSchema{
    {Name: "title", Type: "TEXT", Weight: 1.0},
    {Name: "price", Type: "NUMERIC", Sortable: true},
    {Name: "location", Type: "GEO"},
    {Name: "tags", Type: "TAG", Separator: ","},
}

result, err := redis.FTCreate("products", schema, &client.FTCreateOptions{
    OnHash: true,
    Prefix: []string{"product:"},
})
```

### FTSearch

```go
func (r *Redis) FTSearch(index, query string, options ...*FTSearchOptions) ([]interface{}, error)
```

Search the index with a textual query.

**Example:**
```go
// Basic search
results, err := redis.FTSearch("products", "smartphone")

// Advanced search with filters and options
results, err := redis.FTSearch("products", "smartphone", &client.FTSearchOptions{
    Filter: []client.FTNumericFilter{
        {Field: "price", Min: 100, Max: 500},
    },
    WithScores: true,
    Limit: &client.FTLimit{Offset: 0, Num: 10},
})
```

### Other Search Commands

- `FTDropIndex(index, deleteDocuments...)` - Delete search index
- `FTInfo(index)` - Get index information
- `FTAggregate(index, query, options...)` - Aggregate search results
- `FTExplain(index, query, dialect...)` - Explain query execution
- `FTAdd(index, docID, score, fields, options...)` - Add document (deprecated)
- `FTDel(index, docID, deleteDocument...)` - Delete document (deprecated)

## Time Series Operations

**Phase 4 Feature - Requires RedisTimeSeries Module**

LibRedis supports time series data with automatic compression and aggregation.

### TSCreate

```go
func (r *Redis) TSCreate(key string, options ...*TSCreateOptions) (string, error)
```

Create a new time series.

**Example:**
```go
labels := map[string]string{
    "sensor":   "temperature",
    "location": "living_room",
}

result, err := redis.TSCreate("temp:living_room", &client.TSCreateOptions{
    RetentionMsecs:  3600000, // 1 hour
    Labels:          labels,
    DuplicatePolicy: "LAST",
})
```

### TSAdd

```go
func (r *Redis) TSAdd(key string, timestamp int64, value float64, options ...*TSAddOptions) (int64, error)
```

Add a sample to time series.

**Example:**
```go
// Add with current timestamp
timestamp, err := redis.TSAdd("temp:living_room", 0, 23.5)

// Add with specific timestamp
now := time.Now().UnixMilli()
timestamp, err := redis.TSAdd("temp:living_room", now, 24.0)
```

### TSRange

```go
func (r *Redis) TSRange(key string, fromTimestamp, toTimestamp int64, options ...*TSRangeOptions) ([]TSSample, error)
```

Query time series data in a range.

**Example:**
```go
// Get all samples in last hour
end := time.Now().UnixMilli()
start := end - (60 * 60 * 1000)
samples, err := redis.TSRange("temp:living_room", start, end)

// Get aggregated data
samples, err := redis.TSRange("temp:living_room", start, end, &client.TSRangeOptions{
    Aggregation: &client.TSAggregation{
        Type:       "avg",
        TimeBucket: 300000, // 5 minutes
    },
})
```

### Other Time Series Commands

- `TSMAdd(samples...)` - Add multiple samples
- `TSIncrBy(key, value, options...)` - Increment sample value
- `TSDecrBy(key, value, options...)` - Decrement sample value
- `TSRevRange(key, fromTimestamp, toTimestamp, options...)` - Query in reverse order
- `TSMRange(fromTimestamp, toTimestamp, filters, options...)` - Query multiple series
- `TSMRevRange(fromTimestamp, toTimestamp, filters, options...)` - Query multiple series in reverse
- `TSInfo(key)` - Get time series information

## Probabilistic Data Structures

**Phase 4 Feature - Requires RedisBloom Module**

LibRedis supports probabilistic data structures for memory-efficient approximate operations.

### Bloom Filters

#### BFReserve

```go
func (r *Redis) BFReserve(key string, errorRate float64, capacity int64, options ...*BFReserveOptions) (string, error)
```

Create a Bloom filter.

**Example:**
```go
// Create Bloom filter with 0.1% error rate and 10,000 capacity
result, err := redis.BFReserve("users:visited", 0.001, 10000)
```

#### BFAdd / BFExists

```go
func (r *Redis) BFAdd(key string, item interface{}) (bool, error)
func (r *Redis) BFExists(key string, item interface{}) (bool, error)
```

Add items and check existence.

**Example:**
```go
// Add user to filter
added, err := redis.BFAdd("users:visited", "user:123")

// Check if user exists
exists, err := redis.BFExists("users:visited", "user:123")
```

### Cuckoo Filters

#### CFReserve

```go
func (r *Redis) CFReserve(key string, capacity int64, options ...*CFReserveOptions) (string, error)
```

Create a Cuckoo filter (supports deletions).

**Example:**
```go
result, err := redis.CFReserve("active:users", 10000)
```

#### CF Operations

```go
func (r *Redis) CFAdd(key string, item interface{}) (bool, error)
func (r *Redis) CFExists(key string, item interface{}) (bool, error)
func (r *Redis) CFDel(key string, item interface{}) (bool, error)
```

### Count-Min Sketch

#### CMSInitByDim

```go
func (r *Redis) CMSInitByDim(key string, width, depth int64) (string, error)
```

Create a Count-Min Sketch for frequency counting.

**Example:**
```go
// Create CMS with 2000 width and 5 depth
result, err := redis.CMSInitByDim("page:views", 2000, 5)

// Increment page view count
counts, err := redis.CMSIncrBy("page:views", "/home", 1, "/about", 1)

// Query page view counts
counts, err := redis.CMSQuery("page:views", "/home", "/about", "/contact")
```

### Other Probabilistic Commands

- `BFMAdd(key, items...)` - Add multiple items to Bloom filter
- `BFMExists(key, items...)` - Check multiple items in Bloom filter
- `BFInfo(key)` - Get Bloom filter information
- `CFInfo(key)` - Get Cuckoo filter information
- `CMSInitByProb(key, errorRate, probability)` - Create CMS by probability
- `CMSInfo(key)` - Get Count-Min Sketch information
- `CMSMerge(destKey, sourceKeys, weights...)` - Merge multiple CMS

## Error Handling

### Common Error Patterns

```go
// Check for connection errors
if err != nil {
    if strings.Contains(err.Error(), "connection") {
        // Handle connection error
        return fmt.Errorf("Redis connection lost: %w", err)
    }
    return err
}

// Check for nil values (key doesn't exist)
value, err := redis.Get("key")
if err != nil {
    return err
}
if value == nil {
    // Key doesn't exist
    return nil
}

// Handle Redis errors vs nil values
reply, err := redis.ExecuteCommand("GET", "key")
if err != nil {
    return err
}
if reply.Type == client.ErrorReply {
    return fmt.Errorf("Redis error: %s", reply.Error)
}
```

### Error Types

LibRedis returns standard Go errors with context:

```go
// Connection errors
var ErrConnectionLost = errors.New("connection lost")

// Command errors  
var ErrInvalidCommand = errors.New("invalid command")

// Type conversion errors
var ErrInvalidType = errors.New("invalid reply type")
```

This API reference covers the core functionality of LibRedis. For complete examples and advanced usage patterns, see the other documentation files.