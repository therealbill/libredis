# Phase 4 Implementation Plan - Structured Data & Module Support

## Overview
This document tracks the implementation of Phase 4: Structured Data & Module Support (JSON, Search, Time Series, Probabilistic Data Structures) for the libredis library. All commands listed here can be implemented in a single Claude session.

**Note**: This phase implements Redis module commands that require Redis Stack or specific modules to be loaded.

## Implementation Status Tracker

### JSON Commands (`client/json.go`) - 15 commands
**New file creation required - Redis Stack module**

#### Basic JSON Operations
- [ ] `JSON.SET` - Set JSON value (RedisJSON 1.0+)
- [ ] `JSON.GET` - Get JSON value (RedisJSON 1.0+)
- [ ] `JSON.DEL` - Delete JSON path (RedisJSON 1.0+)
- [ ] `JSON.TYPE` - Get value type (RedisJSON 1.0+)

#### Numeric Operations
- [ ] `JSON.NUMINCRBY` - Increment number (RedisJSON 1.0+)
- [ ] `JSON.NUMMULTBY` - Multiply number (RedisJSON 1.0+)

#### String Operations
- [ ] `JSON.STRAPPEND` - Append to string (RedisJSON 1.0+)
- [ ] `JSON.STRLEN` - Get string length (RedisJSON 1.0+)

#### Array Operations
- [ ] `JSON.ARRAPPEND` - Append to array (RedisJSON 1.0+)
- [ ] `JSON.ARRINDEX` - Find array index (RedisJSON 1.0+)
- [ ] `JSON.ARRINSERT` - Insert into array (RedisJSON 1.0+)
- [ ] `JSON.ARRLEN` - Get array length (RedisJSON 1.0+)
- [ ] `JSON.ARRPOP` - Pop from array (RedisJSON 1.0+)
- [ ] `JSON.ARRTRIM` - Trim array (RedisJSON 1.0+)

#### Object Operations
- [ ] `JSON.OBJKEYS` - Get object keys (RedisJSON 1.0+)
- [ ] `JSON.OBJLEN` - Get object length (RedisJSON 1.0+)

### Search Commands (`client/search.go`) - 8 commands
**New file creation required - RediSearch module**

#### Index Management
- [ ] `FT.CREATE` - Create search index (RediSearch 1.0+)
- [ ] `FT.DROPINDEX` - Drop search index (RediSearch 1.0+)
- [ ] `FT.INFO` - Get index information (RediSearch 1.0+)

#### Search Operations
- [ ] `FT.SEARCH` - Search index (RediSearch 1.0+)
- [ ] `FT.AGGREGATE` - Aggregate search results (RediSearch 1.0+)
- [ ] `FT.EXPLAIN` - Explain query execution (RediSearch 1.0+)

#### Document Management
- [ ] `FT.ADD` - Add document to index (RediSearch 1.0+, deprecated 2.0+)
- [ ] `FT.DEL` - Delete document from index (RediSearch 1.0+, deprecated 2.0+)

### Time Series Commands (`client/timeseries.go`) - 10 commands
**New file creation required - RedisTimeSeries module**

#### Basic Operations
- [ ] `TS.CREATE` - Create time series (RedisTimeSeries 1.0+)
- [ ] `TS.ADD` - Add sample (RedisTimeSeries 1.0+)
- [ ] `TS.MADD` - Add multiple samples (RedisTimeSeries 1.0+)
- [ ] `TS.INCRBY` - Increment by value (RedisTimeSeries 1.0+)
- [ ] `TS.DECRBY` - Decrement by value (RedisTimeSeries 1.0+)

#### Query Operations
- [ ] `TS.RANGE` - Get time range (RedisTimeSeries 1.0+)
- [ ] `TS.REVRANGE` - Get reverse time range (RedisTimeSeries 1.0+)
- [ ] `TS.MRANGE` - Multi-series range query (RedisTimeSeries 1.0+)
- [ ] `TS.MREVRANGE` - Multi-series reverse range (RedisTimeSeries 1.0+)

#### Metadata Operations
- [ ] `TS.INFO` - Get series information (RedisTimeSeries 1.0+)

### Probabilistic Data Structures (`client/probabilistic.go`) - 12 commands
**New file creation required - RedisBloom module**

#### Bloom Filter Operations
- [ ] `BF.RESERVE` - Create Bloom filter (RedisBloom 1.0+)
- [ ] `BF.ADD` - Add item to Bloom filter (RedisBloom 1.0+)
- [ ] `BF.MADD` - Add multiple items (RedisBloom 1.0+)
- [ ] `BF.EXISTS` - Check if item exists (RedisBloom 1.0+)
- [ ] `BF.MEXISTS` - Check multiple items (RedisBloom 1.0+)

#### Cuckoo Filter Operations
- [ ] `CF.RESERVE` - Create Cuckoo filter (RedisBloom 2.0+)
- [ ] `CF.ADD` - Add item to Cuckoo filter (RedisBloom 2.0+)
- [ ] `CF.EXISTS` - Check if item exists (RedisBloom 2.0+)
- [ ] `CF.DEL` - Delete item from filter (RedisBloom 2.0+)

#### Count-Min Sketch Operations
- [ ] `CMS.INITBYDIM` - Initialize by dimensions (RedisBloom 2.0+)
- [ ] `CMS.INCRBY` - Increment counter (RedisBloom 2.0+)
- [ ] `CMS.QUERY` - Query counter value (RedisBloom 2.0+)

**Total Commands to Implement: 45**

## Implementation Details

### Method Signatures and Redis Command Mapping

#### JSON Commands (client/json.go)

```go
// JSON constants and types
const (
	JSONConditionNX = "NX" // Only set if key doesn't exist
	JSONConditionXX = "XX" // Only set if key exists
	
	JSONPathRoot = "$" // Root path
)

// JSONSetOptions represents options for JSON.SET command
type JSONSetOptions struct {
	Condition string // NX or XX
}

// JSONGetOptions represents options for JSON.GET command
type JSONGetOptions struct {
	Indent   string
	Newline  string
	Space    string
}

// JSONArrIndexOptions represents options for JSON.ARRINDEX command
type JSONArrIndexOptions struct {
	Start *int64 // Start index for search
	Stop  *int64 // Stop index for search
}

// Basic Operations

// JSON.SET key path value [NX|XX]
// JSONSet sets a JSON value at the specified path.
func (r *Redis) JSONSet(key, path string, value interface{}) error

// JSONSetWithOptions sets a JSON value with additional options.
func (r *Redis) JSONSetWithOptions(key, path string, value interface{}, opts JSONSetOptions) error

// JSON.GET key [INDENT indent] [NEWLINE newline] [SPACE space] [path ...]
// JSONGet retrieves JSON values from specified paths.
func (r *Redis) JSONGet(key string, paths ...string) (interface{}, error)

// JSONGetWithFormat retrieves JSON with formatting options.
func (r *Redis) JSONGetWithFormat(key string, opts JSONGetOptions, paths ...string) (string, error)

// JSON.DEL key [path]
// JSONDel deletes the entire JSON document or a specific path.
func (r *Redis) JSONDel(key string) (int64, error)

// JSONDelPath deletes a value at a specific JSON path.
func (r *Redis) JSONDelPath(key, path string) (int64, error)

// JSON.TYPE key [path]
// JSONType returns the type of JSON value at the specified path.
func (r *Redis) JSONType(key string) (string, error)

// JSONTypePath returns the type at a specific path.
func (r *Redis) JSONTypePath(key, path string) (string, error)

// Numeric Operations

// JSON.NUMINCRBY key path value
func (r *Redis) JSONNumIncrBy(key, path string, value float64) (*Reply, error)

// JSON.NUMMULTBY key path value
func (r *Redis) JSONNumMultBy(key, path string, value float64) (*Reply, error)

// String Operations

// JSON.STRAPPEND key [path] value
func (r *Redis) JSONStrAppend(key, value string) (*Reply, error)
func (r *Redis) JSONStrAppendPath(key, path, value string) (*Reply, error)

// JSON.STRLEN key [path]
func (r *Redis) JSONStrLen(key string) (*Reply, error)
func (r *Redis) JSONStrLenPath(key, path string) (*Reply, error)

// Array Operations

// JSON.ARRAPPEND key path value [value ...]
func (r *Redis) JSONArrAppend(key, path string, values ...interface{}) (*Reply, error)

// JSON.ARRINDEX key path value [start [stop]]
func (r *Redis) JSONArrIndex(key, path string, value interface{}) (*Reply, error)
func (r *Redis) JSONArrIndexWithRange(key, path string, value interface{}, start, stop int) (*Reply, error)

// JSON.ARRINSERT key path index value [value ...]
func (r *Redis) JSONArrInsert(key, path string, index int, values ...interface{}) (*Reply, error)

// JSON.ARRLEN key [path]
func (r *Redis) JSONArrLen(key string) (*Reply, error)
func (r *Redis) JSONArrLenPath(key, path string) (*Reply, error)

// JSON.ARRPOP key [path [index]]
func (r *Redis) JSONArrPop(key string) (*Reply, error)
func (r *Redis) JSONArrPopPath(key, path string, index int) (*Reply, error)

// JSON.ARRTRIM key path start stop
func (r *Redis) JSONArrTrim(key, path string, start, stop int) (*Reply, error)

// Object Operations

// JSON.OBJKEYS key [path]
func (r *Redis) JSONObjKeys(key string) (*Reply, error)
func (r *Redis) JSONObjKeysPath(key, path string) (*Reply, error)

// JSON.OBJLEN key [path]
func (r *Redis) JSONObjLen(key string) (*Reply, error)
func (r *Redis) JSONObjLenPath(key, path string) (*Reply, error)
```

#### Search Commands (client/search.go)

```go
// Index Management

// FT.CREATE index [ON HASH|JSON] [PREFIX count prefix [prefix ...]] SCHEMA field [field ...]
func (r *Redis) FTCreate(index string, schema FTSchema) (*Reply, error)
func (r *Redis) FTCreateWithOptions(index string, options FTCreateOptions, schema FTSchema) (*Reply, error)

// FT.DROPINDEX index [DD]
func (r *Redis) FTDropIndex(index string) (*Reply, error)
func (r *Redis) FTDropIndexWithDocs(index string) (*Reply, error)

// FT.INFO index
func (r *Redis) FTInfo(index string) (*Reply, error)

// Search Operations

// FT.SEARCH index query [NOCONTENT] [VERBATIM] [NOSTOPWORDS] [LIMIT offset num] [SORTBY sortby [ASC|DESC]]
func (r *Redis) FTSearch(index, query string) (*Reply, error)
func (r *Redis) FTSearchWithOptions(index, query string, options FTSearchOptions) (*Reply, error)

// FT.AGGREGATE index query [VERBATIM] [LOAD count field [field ...]] [GROUPBY nargs property [property ...]]
func (r *Redis) FTAggregate(index, query string) (*Reply, error)
func (r *Redis) FTAggregateWithOptions(index, query string, options FTAggregateOptions) (*Reply, error)

// FT.EXPLAIN index query [DIALECT dialect]
func (r *Redis) FTExplain(index, query string) (*Reply, error)

// Document Management (Deprecated in RediSearch 2.0+)

// FT.ADD index docId score [NOSAVE] [REPLACE] [PARTIAL] [LANGUAGE language] [PAYLOAD payload] FIELDS field content [field content ...]
func (r *Redis) FTAdd(index, docId string, score float64, fields map[string]string) (*Reply, error)
func (r *Redis) FTAddWithOptions(index, docId string, score float64, fields map[string]string, options FTAddOptions) (*Reply, error)

// FT.DEL index docId [DD]
func (r *Redis) FTDel(index, docId string) (*Reply, error)
func (r *Redis) FTDelWithDoc(index, docId string) (*Reply, error)
```

#### Time Series Commands (client/timeseries.go)

```go
// Basic Operations

// TS.CREATE key [RETENTION retentionTime] [ENCODING [COMPRESSED|UNCOMPRESSED]] [CHUNK_SIZE size] [DUPLICATE_POLICY policy] [LABELS label value [label value ...]]
func (r *Redis) TSCreate(key string) (*Reply, error)
func (r *Redis) TSCreateWithOptions(key string, options TSCreateOptions) (*Reply, error)

// TS.ADD key timestamp value [RETENTION retentionTime] [ENCODING [COMPRESSED|UNCOMPRESSED]] [CHUNK_SIZE size] [ON_DUPLICATE policy] [LABELS label value [label value ...]]
func (r *Redis) TSAdd(key string, timestamp int64, value float64) (*Reply, error)
func (r *Redis) TSAddWithOptions(key string, timestamp int64, value float64, options TSAddOptions) (*Reply, error)

// TS.MADD key timestamp value [key timestamp value ...]
func (r *Redis) TSMAdd(samples []TSSample) (*Reply, error)

// TS.INCRBY key value [TIMESTAMP timestamp] [RETENTION retentionTime] [ENCODING [COMPRESSED|UNCOMPRESSED]] [CHUNK_SIZE size] [LABELS label value [label value ...]]
func (r *Redis) TSIncrBy(key string, value float64) (*Reply, error)
func (r *Redis) TSIncrByWithOptions(key string, value float64, options TSIncrOptions) (*Reply, error)

// TS.DECRBY key value [TIMESTAMP timestamp] [RETENTION retentionTime] [ENCODING [COMPRESSED|UNCOMPRESSED]] [CHUNK_SIZE size] [LABELS label value [label value ...]]
func (r *Redis) TSDecrBy(key string, value float64) (*Reply, error)
func (r *Redis) TSDecrByWithOptions(key string, value float64, options TSDecrOptions) (*Reply, error)

// Query Operations

// TS.RANGE key fromTimestamp toTimestamp [LATEST] [FILTER_BY_TS ts ...] [FILTER_BY_VALUE min max] [COUNT count] [ALIGN value] [AGGREGATION aggregationType bucketDuration [BUCKETTIMESTAMP bt] [EMPTY]]
func (r *Redis) TSRange(key string, fromTimestamp, toTimestamp int64) (*Reply, error)
func (r *Redis) TSRangeWithOptions(key string, fromTimestamp, toTimestamp int64, options TSRangeOptions) (*Reply, error)

// TS.REVRANGE key fromTimestamp toTimestamp [LATEST] [FILTER_BY_TS ts ...] [FILTER_BY_VALUE min max] [COUNT count] [ALIGN value] [AGGREGATION aggregationType bucketDuration [BUCKETTIMESTAMP bt] [EMPTY]]
func (r *Redis) TSRevRange(key string, fromTimestamp, toTimestamp int64) (*Reply, error)
func (r *Redis) TSRevRangeWithOptions(key string, fromTimestamp, toTimestamp int64, options TSRangeOptions) (*Reply, error)

// TS.MRANGE fromTimestamp toTimestamp [LATEST] [FILTER_BY_TS ts ...] [FILTER_BY_VALUE min max] [WITHLABELS | SELECTED_LABELS label ...] [COUNT count] [ALIGN value] [AGGREGATION aggregationType bucketDuration [BUCKETTIMESTAMP bt] [EMPTY]] FILTER filter ...
func (r *Redis) TSMRange(fromTimestamp, toTimestamp int64, filters []string) (*Reply, error)
func (r *Redis) TSMRangeWithOptions(fromTimestamp, toTimestamp int64, filters []string, options TSMRangeOptions) (*Reply, error)

// TS.MREVRANGE fromTimestamp toTimestamp [LATEST] [FILTER_BY_TS ts ...] [FILTER_BY_VALUE min max] [WITHLABELS | SELECTED_LABELS label ...] [COUNT count] [ALIGN value] [AGGREGATION aggregationType bucketDuration [BUCKETTIMESTAMP bt] [EMPTY]] FILTER filter ...
func (r *Redis) TSMRevRange(fromTimestamp, toTimestamp int64, filters []string) (*Reply, error)
func (r *Redis) TSMRevRangeWithOptions(fromTimestamp, toTimestamp int64, filters []string, options TSMRangeOptions) (*Reply, error)

// Metadata Operations

// TS.INFO key [DEBUG]
func (r *Redis) TSInfo(key string) (*Reply, error)
func (r *Redis) TSInfoDebug(key string) (*Reply, error)
```

#### Probabilistic Data Structures (client/probabilistic.go)

```go
// Bloom Filter Operations

// BF.RESERVE key error_rate capacity [EXPANSION expansion] [NONSCALING]
func (r *Redis) BFReserve(key string, errorRate float64, capacity int64) (*Reply, error)
func (r *Redis) BFReserveWithOptions(key string, errorRate float64, capacity int64, expansion int, nonscaling bool) (*Reply, error)

// BF.ADD key item
func (r *Redis) BFAdd(key, item string) (*Reply, error)

// BF.MADD key item [item ...]
func (r *Redis) BFMAdd(key string, items ...string) (*Reply, error)

// BF.EXISTS key item
func (r *Redis) BFExists(key, item string) (*Reply, error)

// BF.MEXISTS key item [item ...]
func (r *Redis) BFMExists(key string, items ...string) (*Reply, error)

// Cuckoo Filter Operations

// CF.RESERVE key capacity [BUCKETSIZE bucketsize] [MAXITERATIONS maxiterations] [EXPANSION expansion]
func (r *Redis) CFReserve(key string, capacity int64) (*Reply, error)
func (r *Redis) CFReserveWithOptions(key string, capacity int64, options CFReserveOptions) (*Reply, error)

// CF.ADD key item
func (r *Redis) CFAdd(key, item string) (*Reply, error)

// CF.EXISTS key item
func (r *Redis) CFExists(key, item string) (*Reply, error)

// CF.DEL key item
func (r *Redis) CFDel(key, item string) (*Reply, error)

// Count-Min Sketch Operations

// CMS.INITBYDIM key width depth
func (r *Redis) CMSInitByDim(key string, width, depth int64) (*Reply, error)

// CMS.INCRBY key item increment [item increment ...]
func (r *Redis) CMSIncrBy(key string, items map[string]int64) (*Reply, error)

// CMS.QUERY key item [item ...]
func (r *Redis) CMSQuery(key string, items ...string) (*Reply, error)
```

## Go Best Practices Applied

### 1. Type Safety and Constants
- Use typed constants for JSON conditions (`JSONConditionNX`, `JSONConditionXX`)
- Define search field type constants (`FTFieldTypeText`, `FTFieldTypeNumeric`)
- Use constants for time series encoding and aggregation types
- Prevent string typos with well-defined constant sets

### 2. Structured Options Pattern
- Complex options grouped into clear structs (`JSONSetOptions`, `FTCreateOptions`, `TSCreateOptions`)
- Optional parameters use pointer types (`*int64`, `*TSValueFilter`)
- Composition used for related option sets
- Clear field documentation with Redis option names

### 3. Return Type Consistency
- Return appropriate Go types matching Redis module semantics
- Use `int64` for counts and numeric results
- Use structured types for complex module responses
- Consistent error handling across all module commands

### 4. Error Handling and Module Safety
- All methods return `(result, error)` tuples
- Module availability checks before command execution
- Descriptive error messages for module-specific failures
- Graceful degradation when modules are not loaded

### 5. Documentation Standards
- Include Redis module command syntax in comments
- Clear descriptions of module-specific functionality
- Document module version requirements
- Include examples for complex operations

### 6. Naming Conventions
- Module commands clearly prefixed (`JSON`, `FT`, `TS`, `BF`, `CF`, `CMS`)
- Method names follow Redis command structure
- Options structs follow `CommandNameOptions` pattern
- Consistent naming across related module operations

### 7. Interface Design for Modules
- Simple methods for common module operations
- Complex methods with options for advanced features
- Logical grouping by module functionality
- Consistent parameter ordering within modules

### 8. Memory and Performance Efficiency
- Efficient data structures for time series and search results
- Appropriate use of slices for variable-length module data
- Structured types optimized for JSON marshaling/unmarshaling
- Minimal allocations for frequently used operations

### 9. Module Integration Patterns
- Clear separation between core Redis and module commands
- Proper abstraction of module-specific concepts
- Consistent error handling across different modules
- Module availability detection and graceful fallbacks

### 10. Advanced Data Structure Support
- Professional API design for JSON document operations
- Full-featured search capabilities with proper typing
- Time series operations with appropriate temporal types
- Probabilistic data structures with statistical accuracy considerations

## Supporting Data Structures

All supporting data structures are defined inline with their respective command sections above, following Go best practices for:
- Clear field naming and documentation
- Appropriate use of pointer types for optional fields
- Logical grouping of related fields
- Efficient memory layout and JSON serialization compatibility

## Implementation Order

### Session 1: JSON Basic Operations (8 commands)
1. Create `client/json.go` and supporting structures
2. Implement basic operations: `JSON.SET`, `JSON.GET`, `JSON.DEL`, `JSON.TYPE`
3. Implement numeric operations: `JSON.NUMINCRBY`, `JSON.NUMMULTBY`
4. Implement string operations: `JSON.STRAPPEND`, `JSON.STRLEN`
5. Create `client/json_test.go`

### Session 2: JSON Arrays/Objects and Search Basics (11 commands)
1. Complete JSON implementation: Array operations (6 commands), Object operations (2 commands)
2. Create `client/search.go` and supporting structures
3. Implement search index management: `FT.CREATE`, `FT.DROPINDEX`, `FT.INFO`
4. Create `client/search_test.go`

### Session 3: Search Operations and Time Series (15 commands)
1. Complete search implementation: `FT.SEARCH`, `FT.AGGREGATE`, `FT.EXPLAIN`, `FT.ADD`, `FT.DEL`
2. Create `client/timeseries.go` and supporting structures
3. Implement basic time series: `TS.CREATE`, `TS.ADD`, `TS.MADD`, `TS.INCRBY`, `TS.DECRBY`
4. Create `client/timeseries_test.go`

### Session 4: Time Series Queries and Probabilistic Data (11 commands + tests)
1. Complete time series: `TS.RANGE`, `TS.REVRANGE`, `TS.MRANGE`, `TS.MREVRANGE`, `TS.INFO`
2. Create `client/probabilistic.go` and supporting structures
3. Implement probabilistic data structures: Bloom filters, Cuckoo filters, Count-Min Sketch
4. Create `client/probabilistic_test.go`
5. Finalize all tests and documentation

## Test Implementation Strategy

### Module Availability Checks
All tests must check for module availability:
```go
func TestJSONCommands(t *testing.T) {
    if !isModuleLoaded("ReJSON") {
        t.Skip("RedisJSON module not loaded")
    }
    // Test implementation
}
```

### JSON Tests (client/json_test.go)
```go
func TestJSONSet(t *testing.T) // Basic JSON setting
func TestJSONGet(t *testing.T) // JSON retrieval
func TestJSONNumeric(t *testing.T) // Numeric operations
func TestJSONString(t *testing.T) // String operations
func TestJSONArray(t *testing.T) // Array operations
func TestJSONObject(t *testing.T) // Object operations
```

### Search Tests (client/search_test.go)
```go
func TestFTCreate(t *testing.T) // Index creation
func TestFTSearch(t *testing.T) // Search operations
func TestFTAggregate(t *testing.T) // Aggregation queries
func TestFTDocumentOps(t *testing.T) // Document management
```

### Time Series Tests (client/timeseries_test.go)
```go
func TestTSCreate(t *testing.T) // Series creation
func TestTSAdd(t *testing.T) // Adding samples
func TestTSRange(t *testing.T) // Range queries
func TestTSAggregation(t *testing.T) // Aggregated queries
```

### Probabilistic Tests (client/probabilistic_test.go)
```go
func TestBloomFilter(t *testing.T) // Bloom filter operations
func TestCuckooFilter(t *testing.T) // Cuckoo filter operations
func TestCountMinSketch(t *testing.T) // Count-Min Sketch operations
```

## Module Dependencies

### Required Redis Modules:
- **RedisJSON**: For JSON commands
- **RediSearch**: For search commands
- **RedisTimeSeries**: For time series commands
- **RedisBloom**: For probabilistic data structures

### Installation Notes:
- Redis Stack includes all required modules
- Individual modules can be compiled and loaded separately
- Commands will return module-specific errors if modules aren't loaded

## Success Criteria

Phase 4 is complete when:
- [ ] All 45 module commands are implemented with correct signatures
- [ ] Supporting data structures for all modules are defined
- [ ] Module availability checks are in place
- [ ] All commands have comprehensive tests with module checks
- [ ] Tests pass when appropriate modules are loaded
- [ ] Tests skip gracefully when modules are unavailable
- [ ] Code follows existing libredis patterns
- [ ] Module version requirements are documented
- [ ] Advanced data structure capabilities are functional
- [ ] CLAUDE.md is updated with module information

## Files to Create

### New Files:
- `client/json.go` - RedisJSON implementation
- `client/json_test.go` - JSON command tests
- `client/search.go` - RediSearch implementation
- `client/search_test.go` - Search command tests
- `client/timeseries.go` - RedisTimeSeries implementation
- `client/timeseries_test.go` - Time series tests
- `client/probabilistic.go` - RedisBloom implementation
- `client/probabilistic_test.go` - Probabilistic data structure tests

### Documentation to Update:
- `CLAUDE.md` - Add module support information
- `README.md` - Update with advanced data structures and search capabilities

## Special Implementation Notes

### JSON Considerations:
- JSON paths use JSONPath syntax
- Values can be any JSON-serializable Go type
- Path operations return arrays when multiple paths match
- Encoding/decoding requires proper JSON marshaling

### Search Considerations:
- Index schemas define searchable fields and their types
- Query syntax follows RediSearch query language
- Document management is deprecated in RediSearch 2.0+
- Full-text search requires proper tokenization

### Time Series Considerations:
- Timestamps are Unix milliseconds
- Retention policies automatically expire old data
- Aggregations provide downsampling capabilities
- Labels enable time series filtering and grouping

### Probabilistic Data Structure Considerations:
- Bloom filters have false positives but no false negatives
- Cuckoo filters support deletions unlike Bloom filters
- Count-Min Sketch provides frequency estimation
- Error rates and capacities must be set appropriately

This plan adds advanced data structure support to libredis, enabling modern Redis Stack capabilities including JSON document storage, full-text search, time series analytics, and probabilistic data structures for large-scale applications.