# Phase 2 Implementation Plan - Major Feature Categories

## Overview
This document tracks the implementation of Phase 2: Major Feature Categories (Redis Streams and Geospatial) for the libredis library. All commands listed here can be implemented in a single Claude session.

**STATUS: ✅ COMPLETED** - All Phase 2 features have been successfully implemented and documented.

## Implementation Status Tracker

### Redis Streams (`client/streams.go`) - 14 commands
**New file creation required**

#### Basic Stream Operations
- [x] `XADD` - Add messages to stream (Redis 5.0+)
- [x] `XREAD` - Read messages from stream (Redis 5.0+)
- [x] `XRANGE` - Get range of messages (Redis 5.0+)
- [x] `XREVRANGE` - Get range in reverse (Redis 5.0+)
- [x] `XLEN` - Get stream length (Redis 5.0+)
- [x] `XDEL` - Delete messages (Redis 5.0+)
- [x] `XTRIM` - Trim stream length (Redis 5.0+)

#### Consumer Groups
- [x] `XGROUP CREATE` - Create consumer group (Redis 5.0+)
- [x] `XGROUP DESTROY` - Destroy consumer group (Redis 5.0+)
- [x] `XGROUP SETID` - Set group ID (Redis 5.0+)
- [x] `XREADGROUP` - Read as consumer group (Redis 5.0+)
- [x] `XACK` - Acknowledge processed messages (Redis 5.0+)
- [x] `XCLAIM` - Claim pending messages (Redis 5.0+)
- [x] `XPENDING` - Get pending messages info (Redis 5.0+)

#### Stream Information
- [x] `XINFO STREAM` - Get stream information (Redis 5.0+)
- [x] `XINFO GROUPS` - Get consumer groups info (Redis 5.0+)
- [x] `XINFO CONSUMERS` - Get consumers info (Redis 5.0+)

### Geospatial Commands (`client/geospatial.go`) - 8 commands
**New file creation required**

#### Basic Geospatial Operations
- [x] `GEOADD` - Add geospatial items (Redis 3.2+)
- [x] `GEODIST` - Get distance between points (Redis 3.2+)
- [x] `GEOHASH` - Get geohash strings (Redis 3.2+)
- [x] `GEOPOS` - Get coordinates (Redis 3.2+)

#### Geospatial Search (Modern)
- [x] `GEOSEARCH` - Search within area (Redis 6.2+)
- [x] `GEOSEARCHSTORE` - Store search results (Redis 6.2+)

#### Legacy Geospatial Search (Deprecated but still used)
- [x] `GEORADIUS` - Get items within radius (Redis 3.2+, deprecated 6.2+)
- [x] `GEORADIUSBYMEMBER` - Radius search by member (Redis 3.2+, deprecated 6.2+)

**Total Commands to Implement: 22**

## Implementation Details

### Method Signatures and Redis Command Mapping

#### Streams Commands (client/streams.go)

```go
// Stream constants and types
const (
	StreamIDAutoGenerate = "*"
	StreamIDLatest       = "$"
	StreamIDEarliest     = "0-0"
)

// XAddOptions represents options for XADD command
type XAddOptions struct {
	NoMkStream     bool   // NOMKSTREAM option
	MaxLen         int64  // MAXLEN option
	MinID          string // MINID option
	Approximate    bool   // ~ modifier for MAXLEN/MINID
	Limit          int64  // LIMIT option
}

// XReadOptions represents options for XREAD command  
type XReadOptions struct {
	Count int64 // COUNT option
	Block int64 // BLOCK option (milliseconds)
}

// XRangeOptions represents options for XRANGE/XREVRANGE commands
type XRangeOptions struct {
	Count int64 // COUNT option
}

// StreamEntry represents a single stream entry
type StreamEntry struct {
	ID     string
	Fields map[string]string
}

// StreamMessage represents messages from one stream
type StreamMessage struct {
	Stream  string
	Entries []StreamEntry
}

// XGroupCreateOptions represents options for XGROUP CREATE
type XGroupCreateOptions struct {
	MkStream     bool  // MKSTREAM option
	EntriesRead  int64 // ENTRIESREAD option
}

// XReadGroupOptions represents options for XREADGROUP
type XReadGroupOptions struct {
	Count int64 // COUNT option  
	Block int64 // BLOCK option (milliseconds)
	NoAck bool  // NOACK option
}

// XClaimOptions represents options for XCLAIM command
type XClaimOptions struct {
	Idle       int64  // IDLE option
	Time       int64  // TIME option  
	RetryCount int64  // RETRYCOUNT option
	Force      bool   // FORCE option
	JustID     bool   // JUSTID option
	LastID     string // LASTID option
}

// XPendingOptions represents options for XPENDING command
type XPendingOptions struct {
	Idle     int64  // IDLE option
	Start    string // Start ID
	End      string // End ID  
	Count    int64  // Count limit
	Consumer string // Specific consumer
}

// XPendingInfo represents summary of pending messages
type XPendingInfo struct {
	Count     int64
	Lower     string
	Higher    string
	Consumers map[string]int64
}

// XPendingMessage represents a pending message detail
type XPendingMessage struct {
	ID           string
	Consumer     string
	IdleTime     int64
	DeliveryCount int64
}

// Basic Stream Operations

// XADD key [NOMKSTREAM] [MAXLEN|MINID [=|~] threshold [LIMIT count]] *|ID field value [field value ...]
// XAdd appends a new entry to a stream.
func (r *Redis) XAdd(key, id string, fields map[string]string) (string, error)

// XAddWithOptions appends a new entry with additional options.
func (r *Redis) XAddWithOptions(key, id string, fields map[string]string, opts XAddOptions) (string, error)

// XREAD [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] id [id ...]
// XRead reads data from one or multiple streams.
func (r *Redis) XRead(streams map[string]string) ([]StreamMessage, error)

// XReadWithOptions reads from streams with additional options.
func (r *Redis) XReadWithOptions(streams map[string]string, opts XReadOptions) ([]StreamMessage, error)

// XRANGE key start end [COUNT count]
// XRange returns the stream entries matching a given range of IDs.
func (r *Redis) XRange(key, start, end string) ([]StreamEntry, error)

// XRangeWithOptions returns stream entries with count limit.
func (r *Redis) XRangeWithOptions(key, start, end string, opts XRangeOptions) ([]StreamEntry, error)

// XREVRANGE key end start [COUNT count]
// XRevRange returns stream entries in reverse order.
func (r *Redis) XRevRange(key, end, start string) ([]StreamEntry, error)

// XRevRangeWithOptions returns stream entries in reverse with count limit.
func (r *Redis) XRevRangeWithOptions(key, end, start string, opts XRangeOptions) ([]StreamEntry, error)

// XLEN key
// XLen returns the number of entries in a stream.
func (r *Redis) XLen(key string) (int64, error)

// XDEL key id [id ...]
// XDel removes specified entries from a stream.
func (r *Redis) XDel(key string, ids ...string) (int64, error)

// XTRIM key MAXLEN|MINID [=|~] threshold [LIMIT count]
// XTrim trims a stream to a given size or minimum ID.
func (r *Redis) XTrim(key string, strategy string, threshold string) (int64, error)

// XTrimWithOptions trims a stream with additional options.
func (r *Redis) XTrimWithOptions(key string, strategy string, threshold string, opts XAddOptions) (int64, error)

// Consumer Group Operations

// XGROUP CREATE key groupname id|$ [MKSTREAM] [ENTRIESREAD entries_read]
// XGroupCreate creates a new consumer group.
func (r *Redis) XGroupCreate(key, groupname, id string) error

// XGroupCreateWithOptions creates a consumer group with additional options.
func (r *Redis) XGroupCreateWithOptions(key, groupname, id string, opts XGroupCreateOptions) error

// XGROUP DESTROY key groupname
// XGroupDestroy destroys a consumer group.
func (r *Redis) XGroupDestroy(key, groupname string) (int64, error)

// XGROUP SETID key groupname id|$ [ENTRIESREAD entries_read]
// XGroupSetID sets the consumer group's last delivered ID.
func (r *Redis) XGroupSetID(key, groupname, id string) error

// XREADGROUP GROUP group consumer [COUNT count] [BLOCK milliseconds] [NOACK] STREAMS key [key ...] ID [ID ...]
// XReadGroup reads from streams as a consumer group member.
func (r *Redis) XReadGroup(group, consumer string, streams map[string]string) ([]StreamMessage, error)

// XReadGroupWithOptions reads from streams as consumer group with options.
func (r *Redis) XReadGroupWithOptions(group, consumer string, streams map[string]string, opts XReadGroupOptions) ([]StreamMessage, error)

// XACK key group id [id ...]
// XAck acknowledges processing of messages by a consumer.
func (r *Redis) XAck(key, group string, ids ...string) (int64, error)

// XCLAIM key group consumer min-idle-time id [id ...] [IDLE ms] [TIME unix-time] [RETRYCOUNT count] [FORCE] [JUSTID] [LASTID id]
// XClaim transfers ownership of pending messages to another consumer.
func (r *Redis) XClaim(key, group, consumer string, minIdleTime int64, ids []string) ([]StreamEntry, error)

// XClaimWithOptions claims messages with additional options.
func (r *Redis) XClaimWithOptions(key, group, consumer string, minIdleTime int64, ids []string, opts XClaimOptions) ([]StreamEntry, error)

// XPENDING key group [[IDLE min-idle-time] start end count [consumer]]
// XPending returns information about pending messages.
func (r *Redis) XPending(key, group string) (XPendingInfo, error)

// XPendingWithOptions returns detailed pending message information.
func (r *Redis) XPendingWithOptions(key, group string, opts XPendingOptions) ([]XPendingMessage, error)

// Stream Information

// XINFO STREAM key [FULL [COUNT count]]
// XInfoStream returns general information about a stream.
func (r *Redis) XInfoStream(key string) (map[string]interface{}, error)

// XInfoStreamFull returns detailed information about a stream.
func (r *Redis) XInfoStreamFull(key string, count int64) (map[string]interface{}, error)

// XINFO GROUPS key
// XInfoGroups returns information about consumer groups.
func (r *Redis) XInfoGroups(key string) ([]map[string]interface{}, error)

// XINFO CONSUMERS key groupname
// XInfoConsumers returns information about consumers in a group.
func (r *Redis) XInfoConsumers(key, groupname string) ([]map[string]interface{}, error)
```

#### Geospatial Commands (client/geospatial.go)

```go
// Geospatial constants
const (
	GeoUnitMeters     = "M"
	GeoUnitKilometers = "KM" 
	GeoUnitFeet       = "FT"
	GeoUnitMiles      = "MI"
	
	GeoOrderAsc  = "ASC"
	GeoOrderDesc = "DESC"
)

// GeoMember represents a geospatial member with coordinates
type GeoMember struct {
	Longitude float64
	Latitude  float64
	Member    string
}

// GeoAddOptions represents options for GEOADD command
type GeoAddOptions struct {
	NX bool // NX option - only add new elements
	XX bool // XX option - only update existing elements  
	CH bool // CH option - return count of changed elements
}

// GeoCoordinate represents longitude/latitude coordinates
type GeoCoordinate struct {
	Longitude float64
	Latitude  float64
}

// GeoSearchOptions represents options for GEOSEARCH command
type GeoSearchOptions struct {
	// Search center (exactly one must be specified)
	FromMember  *string        // FROMMEMBER option
	FromLonLat  *GeoCoordinate // FROMLONLAT option
	
	// Search area (exactly one must be specified)
	ByRadius *GeoRadius // BYRADIUS option
	ByBox    *GeoBox    // BYBOX option
	
	// Result options
	Order      string // ASC or DESC
	Count      int64  // COUNT option
	Any        bool   // ANY option
	WithCoord  bool   // WITHCOORD option
	WithDist   bool   // WITHDIST option
	WithHash   bool   // WITHHASH option
}

// GeoSearchStoreOptions represents options for GEOSEARCHSTORE command
type GeoSearchStoreOptions struct {
	GeoSearchOptions
	StoreDist bool // STOREDIST option
}

// GeoRadius represents a radius search parameter
type GeoRadius struct {
	Radius float64
	Unit   string
}

// GeoBox represents a box search parameter
type GeoBox struct {
	Width  float64
	Height float64
	Unit   string
}

// GeoRadiusOptions represents options for legacy GEORADIUS commands
type GeoRadiusOptions struct {
	WithCoord bool   // WITHCOORD option
	WithDist  bool   // WITHDIST option
	WithHash  bool   // WITHHASH option
	Count     int64  // COUNT option
	Any       bool   // ANY option
	Order     string // ASC or DESC
	Store     string // STORE option
	StoreDist string // STOREDIST option
}

// GeoLocation represents a geospatial search result
type GeoLocation struct {
	Member      string
	Coordinates *GeoCoordinate
	Distance    *float64
	Hash        *int64
}

// Basic Operations

// GEOADD key [NX|XX] [CH] longitude latitude member [longitude latitude member ...]
// GeoAdd adds geospatial items to a geospatial index.
func (r *Redis) GeoAdd(key string, members []GeoMember) (int64, error)

// GeoAddWithOptions adds geospatial items with additional options.
func (r *Redis) GeoAddWithOptions(key string, members []GeoMember, opts GeoAddOptions) (int64, error)

// GEODIST key member1 member2 [M|KM|FT|MI]
// GeoDist returns the distance between two geospatial members.
func (r *Redis) GeoDist(key, member1, member2 string) (float64, error)

// GeoDistWithUnit returns the distance with a specific unit.
func (r *Redis) GeoDistWithUnit(key, member1, member2, unit string) (float64, error)

// GEOHASH key member [member ...]
// GeoHash returns geohash strings for the specified members.
func (r *Redis) GeoHash(key string, members ...string) ([]string, error)

// GEOPOS key member [member ...]
// GeoPos returns coordinates for the specified members.
func (r *Redis) GeoPos(key string, members ...string) ([]*GeoCoordinate, error)

// Modern Search Commands

// GEOSEARCH key [FROMMEMBER member] [FROMLONLAT longitude latitude] [BYRADIUS radius M|KM|FT|MI] [BYBOX width height M|KM|FT|MI] [ASC|DESC] [COUNT count [ANY]] [WITHCOORD] [WITHDIST] [WITHHASH]
// GeoSearch queries a geospatial index for members within a specified area.
func (r *Redis) GeoSearch(key string, opts GeoSearchOptions) ([]GeoLocation, error)

// GEOSEARCHSTORE destination source [FROMMEMBER member] [FROMLONLAT longitude latitude] [BYRADIUS radius M|KM|FT|MI] [BYBOX width height M|KM|FT|MI] [ASC|DESC] [COUNT count [ANY]] [STOREDIST]
// GeoSearchStore executes a geospatial search and stores results in another key.
func (r *Redis) GeoSearchStore(destination, source string, opts GeoSearchStoreOptions) (int64, error)

// Legacy Search Commands (Deprecated but still supported)

// GEORADIUS key longitude latitude radius M|KM|FT|MI [WITHCOORD] [WITHDIST] [WITHHASH] [COUNT count [ANY]] [ASC|DESC] [STORE key] [STOREDIST key]
// GeoRadius returns members within a radius from coordinates (deprecated - use GeoSearch).
func (r *Redis) GeoRadius(key string, longitude, latitude, radius float64, unit string) ([]string, error)

// GeoRadiusWithOptions returns members with additional result information (deprecated).
func (r *Redis) GeoRadiusWithOptions(key string, longitude, latitude, radius float64, unit string, opts GeoRadiusOptions) ([]GeoLocation, error)

// GEORADIUSBYMEMBER key member radius M|KM|FT|MI [WITHCOORD] [WITHDIST] [WITHHASH] [COUNT count [ANY]] [ASC|DESC] [STORE key] [STOREDIST key]
// GeoRadiusByMember returns members within a radius from another member (deprecated).
func (r *Redis) GeoRadiusByMember(key, member string, radius float64, unit string) ([]string, error)

// GeoRadiusByMemberWithOptions returns members with additional information (deprecated).
func (r *Redis) GeoRadiusByMemberWithOptions(key, member string, radius float64, unit string, opts GeoRadiusOptions) ([]GeoLocation, error)
```

## Go Best Practices Applied

### 1. Type Safety and Constants
- Use typed constants for string literals (`StreamIDAutoGenerate`, `GeoUnitMeters`)
- Define unit constants to prevent typos and improve code clarity
- Use pointer types for optional/nullable fields (`*GeoCoordinate`, `*string`)

### 2. Structured Options Pattern
- Complex options grouped into structs (`XAddOptions`, `GeoSearchOptions`)
- Composition used where appropriate (`GeoSearchStoreOptions` embeds `GeoSearchOptions`)
- Clear field documentation with Redis option names

### 3. Return Type Consistency
- Return appropriate Go types: `int64` for counts, `[]StreamEntry` for lists
- Use structured types for complex responses (`StreamMessage`, `GeoLocation`)
- Pointer types for optional fields in response structures

### 4. Error Handling
- All methods return `(result, error)` tuples
- Descriptive error messages for validation failures
- Proper error propagation from Redis operations

### 5. Documentation Standards
- Include Redis command syntax in comments
- Clear descriptions of functionality and parameters
- Document Redis version requirements for new features
- Mark deprecated commands explicitly

### 6. Naming Conventions
- Go naming conventions consistently applied
- Stream operations prefixed with 'X' matching Redis commands
- Geospatial operations prefixed with 'Geo' for clarity
- Options structs follow `CommandNameOptions` pattern

### 7. Interface Design
- Simple methods for common use cases
- Complex methods with options for advanced scenarios
- Consistent parameter ordering across related methods
- Logical grouping of related functionality

### 8. Memory Efficiency
- Use slices instead of arrays for variable-length data
- Pointer types for optional/large structures
- Efficient field types (int64 for Redis numbers)

### 9. Null Safety
- Use pointer types for optional coordinates and values
- Clear distinction between zero values and missing values
- Proper nil handling in result structures

## Implementation Order

### Session 1: Basic Streams Operations (7 commands)
1. Create `client/streams.go` and supporting structures
2. Implement basic operations: `XADD`, `XREAD`, `XRANGE`, `XREVRANGE`
3. Implement management: `XLEN`, `XDEL`, `XTRIM`
4. Create `client/streams_test.go` with basic tests

### Session 2: Streams Consumer Groups (7 commands)
1. Implement consumer group management: `XGROUP CREATE`, `XGROUP DESTROY`, `XGROUP SETID`
2. Implement consumer operations: `XREADGROUP`, `XACK`, `XCLAIM`, `XPENDING`
3. Add consumer group tests to `client/streams_test.go`

### Session 3: Streams Info and Geospatial Basics (8 commands)
1. Implement stream info: `XINFO STREAM`, `XINFO GROUPS`, `XINFO CONSUMERS`
2. Create `client/geospatial.go` and supporting structures
3. Implement basic geo operations: `GEOADD`, `GEODIST`, `GEOHASH`, `GEOPOS`
4. Create `client/geospatial_test.go`

### Session 4: Geospatial Search and Finalization (4 commands + tests)
1. Implement modern search: `GEOSEARCH`, `GEOSEARCHSTORE`
2. Implement legacy search: `GEORADIUS`, `GEORADIUSBYMEMBER`
3. Complete all tests for geospatial commands
4. Update documentation

## Test Implementation Strategy

### Streams Tests (client/streams_test.go)
```go
func TestXAdd(t *testing.T) // Basic stream addition
func TestXRead(t *testing.T) // Reading from streams
func TestXRange(t *testing.T) // Range queries
func TestXLen(t *testing.T) // Stream length
func TestXTrim(t *testing.T) // Stream trimming
func TestXGroupOperations(t *testing.T) // Consumer groups
func TestXReadGroup(t *testing.T) // Group reading
func TestXAck(t *testing.T) // Message acknowledgment
func TestXPending(t *testing.T) // Pending messages
func TestXInfo(t *testing.T) // Stream information
```

### Geospatial Tests (client/geospatial_test.go)
```go
func TestGeoAdd(t *testing.T) // Adding geo data
func TestGeoDist(t *testing.T) // Distance calculations
func TestGeoHash(t *testing.T) // Geohash generation
func TestGeoPos(t *testing.T) // Position retrieval
func TestGeoSearch(t *testing.T) // Modern search
func TestGeoSearchStore(t *testing.T) // Search with storage
func TestGeoRadius(t *testing.T) // Legacy radius search
func TestGeoRadiusByMember(t *testing.T) // Legacy member radius
```

## Redis Version Compatibility

### Streams Commands:
- **Redis 5.0+**: All basic stream operations and consumer groups
- **Redis 6.0+**: Enhanced XINFO commands
- **Redis 6.2+**: Additional XTRIM options

### Geospatial Commands:
- **Redis 3.2+**: Basic geo operations and legacy search
- **Redis 6.2+**: Modern GEOSEARCH commands (GEORADIUS deprecated)

## Success Criteria

Phase 2 is complete when:
- [ ] All 22 commands are implemented with correct signatures
- [ ] Supporting data structures are properly defined
- [ ] All commands have comprehensive tests
- [ ] Integration tests pass with appropriate Redis versions
- [ ] Code follows existing libredis patterns
- [ ] Redis version requirements are documented
- [ ] Both new files (`streams.go`, `geospatial.go`) are complete
- [ ] CLAUDE.md is updated with new features

## Files to Create

### New Files:
- `client/streams.go` - Complete Redis Streams implementation
- `client/streams_test.go` - Comprehensive streams tests
- `client/geospatial.go` - Complete geospatial implementation  
- `client/geospatial_test.go` - Comprehensive geospatial tests

### Documentation to Update:
- `CLAUDE.md` - Add Streams and Geospatial information
- `README.md` - Update feature list with major new capabilities

## Special Implementation Notes

### Streams Considerations:
- Stream IDs can be auto-generated with "*" or explicitly provided
- Consumer groups require careful state management
- XREADGROUP supports blocking operations similar to BLPOP
- Message acknowledgment is crucial for reliable processing

### Geospatial Considerations:
- Coordinates must be valid longitude/latitude pairs
- Distance units: M (meters), KM (kilometers), FT (feet), MI (miles)
- GEOSEARCH replaces deprecated GEORADIUS commands
- Geohashes provide location approximation for indexing

This plan provides complete Redis Streams and Geospatial support, bringing libredis up to modern Redis standards for location-based and stream processing applications.