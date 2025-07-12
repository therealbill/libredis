# Phase 1 Implementation Plan - Core Missing Commands

## ✅ COMPLETED
Phase 1 implementation has been successfully completed! All 24 core missing Redis commands have been implemented with comprehensive tests.

## Overview
This document tracked the implementation of Phase 1: Core Missing Commands for the libredis library. All commands were successfully implemented in a single Claude session.

## Implementation Status Tracker

### Lists Commands (`client/lists.go`) - 5 commands ✅
- [x] `LMOVE` - Move element between lists (Redis 6.2+)
- [x] `BLMOVE` - Blocking version of LMOVE (Redis 6.2+)
- [x] `LPOS` - Find element position (Redis 6.0.6+)
- [x] `LMPOP` - Pop from multiple lists (Redis 7.0+)
- [x] `BLMPOP` - Blocking version of LMPOP (Redis 7.0+)

### Sets Commands (`client/sets.go`) - 1 command ✅
- [x] `SMISMEMBER` - Check multiple members (Redis 6.2+)

### Sorted Sets Commands (`client/sorted_sets.go`) - 6 commands ✅
- [x] `ZPOPMAX` - Pop maximum elements (Redis 5.0+)
- [x] `ZPOPMIN` - Pop minimum elements (Redis 5.0+)
- [x] `BZPOPMAX` - Blocking ZPOPMAX (Redis 5.0+)
- [x] `BZPOPMIN` - Blocking ZPOPMIN (Redis 5.0+)
- [x] `ZRANDMEMBER` - Get random members (Redis 6.2+)
- [x] `ZMSCORE` - Get multiple member scores (Redis 6.2+)
### Hashes Commands (`client/hashes.go`) - 2 commands ✅
- [x] `HSTRLEN` - Get hash field string length (Redis 3.2+)
- [x] `HRANDFIELD` - Get random hash fields (Redis 6.2+)

### Keys Commands (`client/keys.go`) - 4 commands ✅
- [x] `COPY` - Copy key to another key (Redis 6.2+)
- [x] `TOUCH` - Update last access time (Redis 3.2.1+)
- [x] `UNLINK` - Non-blocking delete (Redis 4.0+)
- [x] `WAIT` - Wait for write replication (Redis 3.0+)

### Bitmap Commands (new file: `client/bitmaps.go`) - 3 commands ✅
- [x] `BITFIELD` - Perform bitfield operations (Redis 3.2+)
- [x] `BITFIELD_RO` - Read-only bitfield operations (Redis 6.0+)
- [x] `BITPOS` - Find bit position (Redis 2.8.7+)

### Connection Commands (`client/connection.go`) - 3 commands ✅
- [x] `AUTH` - Enhanced for ACL support (Redis 6.0+)
- [x] `HELLO` - Handshake with authentication (Redis 6.0+)
- [x] `RESET` - Reset connection state (Redis 6.2+)

**Total Commands Implemented: 24/24** ✅

## Implementation Details

### Method Signatures and Redis Command Mapping

#### Lists Commands (client/lists.go)
```go
// Direction constants for list operations
const (
	ListDirectionLeft  = "LEFT"
	ListDirectionRight = "RIGHT"
)

// LPosOptions represents options for LPOS command
type LPosOptions struct {
	Rank   int  // RANK option
	Count  int  // COUNT option  
	MaxLen int  // MAXLEN option
}

// LMPopOptions represents options for LMPOP command
type LMPopOptions struct {
	Count int // COUNT option
}

// LMOVE source destination LEFT|RIGHT LEFT|RIGHT
// LMove atomically returns and removes the first/last element from source list,
// and pushes the element as the first/last element of destination list.
func (r *Redis) LMove(source, destination, wherefrom, whereto string) (string, error)

// BLMOVE source destination LEFT|RIGHT LEFT|RIGHT timeout
// BLMove is the blocking variant of LMOVE.
func (r *Redis) BLMove(source, destination, wherefrom, whereto string, timeout int) (string, error)

// LPOS key element [RANK rank] [COUNT num-matches] [MAXLEN len]
// LPos returns the index of matching elements inside a Redis list.
func (r *Redis) LPos(key, element string) (int64, error)

// LPosWithOptions returns the index of matching elements with additional options.
func (r *Redis) LPosWithOptions(key, element string, opts LPosOptions) ([]int64, error)

// LMPOP numkeys key [key ...] LEFT|RIGHT [COUNT count]
// LMPop pops one or more elements from the first non-empty list key.
func (r *Redis) LMPop(keys []string, direction string) (map[string][]string, error)

// LMPopWithCount pops count elements from the first non-empty list key.
func (r *Redis) LMPopWithCount(keys []string, direction string, count int) (map[string][]string, error)

// BLMPOP timeout numkeys key [key ...] LEFT|RIGHT [COUNT count]
// BLMPop is the blocking variant of LMPOP.
func (r *Redis) BLMPop(timeout int, keys []string, direction string) (map[string][]string, error)

// BLMPopWithCount is the blocking variant of LMPOP with count.
func (r *Redis) BLMPopWithCount(timeout int, keys []string, direction string, count int) (map[string][]string, error)
```

#### Sets Commands (client/sets.go)
```go
// SMISMEMBER key member [member ...]
// SMIsMember returns whether each member is a member of the set stored at key.
func (r *Redis) SMIsMember(key string, members ...string) ([]bool, error)
```

#### Sorted Sets Commands (client/sorted_sets.go)
```go
// ZMember represents a sorted set member with score
type ZMember struct {
	Member string
	Score  float64
}

// ZPopResult represents the result of a ZPOP operation
type ZPopResult struct {
	Key     string
	Members []ZMember
}

// ZRandMemberOptions represents options for ZRANDMEMBER command
type ZRandMemberOptions struct {
	Count      int
	WithScores bool
}

// ZPOPMAX key [count]
// ZPopMax removes and returns up to count members with the highest scores.
func (r *Redis) ZPopMax(key string) (ZMember, error)

// ZPopMaxCount removes and returns up to count members with the highest scores.
func (r *Redis) ZPopMaxCount(key string, count int) ([]ZMember, error)

// ZPOPMIN key [count]
// ZPopMin removes and returns up to count members with the lowest scores.
func (r *Redis) ZPopMin(key string) (ZMember, error)

// ZPopMinCount removes and returns up to count members with the lowest scores.
func (r *Redis) ZPopMinCount(key string, count int) ([]ZMember, error)

// BZPOPMAX key [key ...] timeout
// BZPopMax is the blocking variant of ZPOPMAX.
func (r *Redis) BZPopMax(keys []string, timeout int) (ZPopResult, error)

// BZPOPMIN key [key ...] timeout
// BZPopMin is the blocking variant of ZPOPMIN.
func (r *Redis) BZPopMin(keys []string, timeout int) (ZPopResult, error)

// ZRANDMEMBER key [count [WITHSCORES]]
// ZRandMember returns a random member from the sorted set.
func (r *Redis) ZRandMember(key string) (string, error)

// ZRandMemberWithOptions returns random members with additional options.
func (r *Redis) ZRandMemberWithOptions(key string, opts ZRandMemberOptions) ([]ZMember, error)

// ZMSCORE key member [member ...]
// ZMScore returns the scores associated with the specified members.
func (r *Redis) ZMScore(key string, members ...string) ([]float64, error)
```

#### Hashes Commands (client/hashes.go)
```go
// HField represents a hash field-value pair
type HField struct {
	Field string
	Value string
}

// HRandFieldOptions represents options for HRANDFIELD command
type HRandFieldOptions struct {
	Count      int
	WithValues bool
}

// HSTRLEN key field
// HStrLen returns the string length of the value associated with field.
func (r *Redis) HStrLen(key, field string) (int64, error)

// HRANDFIELD key [count [WITHVALUES]]
// HRandField returns a random field from the hash stored at key.
func (r *Redis) HRandField(key string) (string, error)

// HRandFieldWithOptions returns random fields with additional options.
func (r *Redis) HRandFieldWithOptions(key string, opts HRandFieldOptions) ([]HField, error)
```

#### Keys Commands (client/keys.go)
```go
// CopyOptions represents options for COPY command
type CopyOptions struct {
	DestinationDB int  // DB option
	Replace       bool // REPLACE option
}

// COPY source destination [DB destination-db] [REPLACE]
// Copy makes a copy of the value stored at the source key to the destination key.
func (r *Redis) Copy(source, destination string) (bool, error)

// CopyWithOptions copies a key with additional options.
func (r *Redis) CopyWithOptions(source, destination string, opts CopyOptions) (bool, error)

// TOUCH key [key ...]
// Touch alters the last access time of one or more keys.
func (r *Redis) Touch(keys ...string) (int64, error)

// UNLINK key [key ...]
// Unlink is similar to DEL but performs non-blocking deletion.
func (r *Redis) Unlink(keys ...string) (int64, error)

// WAIT numreplicas timeout
// Wait blocks until all write commands are successfully synced to replicas.
func (r *Redis) Wait(numreplicas, timeout int) (int64, error)
```

#### Bitmap Commands (new file: client/bitmaps.go)
```go
// BitFieldOperation represents a single bitfield operation
type BitFieldOperation struct {
	Type   string      // GET, SET, INCRBY
	Offset int64       // Bit offset
	Value  interface{} // Value for SET/INCRBY operations
}

// BitFieldOverflow represents overflow behavior
type BitFieldOverflow string

const (
	BitFieldOverflowWrap BitFieldOverflow = "WRAP"
	BitFieldOverflowSat  BitFieldOverflow = "SAT" 
	BitFieldOverflowFail BitFieldOverflow = "FAIL"
)

// BitPosOptions represents options for BITPOS command
type BitPosOptions struct {
	Start *int64 // Start position
	End   *int64 // End position
}

// BITFIELD key [GET type offset] [SET type offset value] [INCRBY type offset increment] [OVERFLOW WRAP|SAT|FAIL]
// BitField performs arbitrary bit field integer operations on strings.
func (r *Redis) BitField(key string, operations []BitFieldOperation) ([]int64, error)

// BitFieldWithOverflow performs bitfield operations with overflow control.
func (r *Redis) BitFieldWithOverflow(key string, overflow BitFieldOverflow, operations []BitFieldOperation) ([]int64, error)

// BITFIELD_RO key [GET type offset] [GET type offset ...]
// BitFieldRO is the read-only variant of BITFIELD.
func (r *Redis) BitFieldRO(key string, getOps []BitFieldOperation) ([]int64, error)

// BITPOS key bit [start] [end]
// BitPos returns the position of the first bit set to 1 or 0.
func (r *Redis) BitPos(key string, bit int) (int64, error)

// BitPosWithRange returns the bit position within a specified range.
func (r *Redis) BitPosWithRange(key string, bit int, opts BitPosOptions) (int64, error)
```

#### Connection Commands (client/connection.go)
```go
// HelloOptions represents options for HELLO command
type HelloOptions struct {
	ProtocolVersion int
	Username        string
	Password        string
	ClientName      string
}

// AUTH [username] password
// AuthWithUser authenticates using username and password (ACL).
func (r *Redis) AuthWithUser(username, password string) error

// HELLO [protover [AUTH username password] [SETNAME clientname]]
// Hello switches to a different protocol version and authenticates.
func (r *Redis) Hello(protocolVersion int) (map[string]interface{}, error)

// HelloWithOptions performs handshake with additional options.
func (r *Redis) HelloWithOptions(opts HelloOptions) (map[string]interface{}, error)

// RESET
// Reset resets the connection state.
func (r *Redis) Reset() error
```

## Go Best Practices Applied

### 1. Type Safety and Constants
- Use typed constants for string literals (e.g., `ListDirectionLeft`, `BitFieldOverflowWrap`)
- Define custom types for better type safety (`BitFieldOverflow`, `HRandFieldOptions`)
- Use pointer types for optional parameters (`*int64` in `BitPosOptions`)

### 2. Structured Options Pattern
- Use option structs instead of many parameters (`LPosOptions`, `CopyOptions`, `HelloOptions`)
- Provide both simple and complex variants of methods
- Use descriptive field names with clear types

### 3. Return Type Consistency
- Return appropriate Go types instead of generic `*Reply`
- Use `int64` for Redis integer responses
- Use `bool` for Redis boolean responses
- Use `[]bool` for multi-member checks
- Use structured types for complex responses (`ZMember`, `ZPopResult`)

### 4. Error Handling
- All methods return `(result, error)` tuples
- Use descriptive error messages
- Validate input parameters before Redis calls

### 5. Documentation Standards
- Include Redis command syntax in comments
- Provide clear descriptions of what each method does
- Document parameter meanings and constraints
- Include Redis version requirements

### 6. Naming Conventions
- Use Go conventions (PascalCase for exported, camelCase for unexported)
- Use descriptive names that clearly indicate purpose
- Avoid abbreviations unless commonly understood
- Use consistent naming patterns across similar methods

### 7. Interface Design
- Provide simple methods for common cases
- Provide complex methods with options for advanced cases
- Use meaningful return types that match the operation semantics
- Group related functionality logically

## Implementation Order

### Session 1: Lists and Sets (6 commands)
1. `LMOVE` and `BLMOVE` (related commands)
2. `LPOS` (simpler command)
3. `LMPOP` and `BLMPOP` (related commands)
4. `SMISMEMBER` (simple set command)

### Session 2: Sorted Sets (6 commands)
1. `ZPOPMAX` and `ZPOPMIN` (related commands)
2. `BZPOPMAX` and `BZPOPMIN` (related commands)
3. `ZRANDMEMBER` (with variants)
4. `ZMSCORE`

### Session 3: Hashes, Keys, and Bitmaps (9 commands)
1. `HSTRLEN` and `HRANDFIELD` (hash commands)
2. `COPY`, `TOUCH`, `UNLINK`, `WAIT` (key commands)
3. Create `client/bitmaps.go` with `BITFIELD`, `BITFIELD_RO`, `BITPOS`

### Session 4: Connection Commands and Testing (3 commands + tests)
1. Enhanced `AUTH`, `HELLO`, `RESET` in connection.go
2. Add comprehensive tests for all new commands
3. Update documentation

## Test Implementation Strategy

### Test Structure and Best Practices

#### Test Organization
```go
func TestLMove(t *testing.T) {
    tests := []struct {
        name         string
        source       string
        destination  string
        wherefrom    string
        whereto      string
        setup        func(*Redis) error
        wantResult   string
        wantErr      bool
    }{
        {
            name:        "move from left to right",
            source:      "src",
            destination: "dst", 
            wherefrom:   ListDirectionLeft,
            whereto:     ListDirectionRight,
            setup:       func(r *Redis) error { return r.LPush("src", "value1", "value2") },
            wantResult:  "value2",
            wantErr:     false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

#### Input Validation Tests
```go
func TestLMoveValidation(t *testing.T) {
    r := setupRedis(t)
    
    // Test invalid direction
    _, err := r.LMove("src", "dst", "INVALID", ListDirectionRight)
    if err == nil {
        t.Error("expected error for invalid direction")
    }
    
    // Test empty key names
    _, err = r.LMove("", "dst", ListDirectionLeft, ListDirectionRight)
    if err == nil {
        t.Error("expected error for empty source key")
    }
}
```

#### Benchmark Tests
```go
func BenchmarkLMove(b *testing.B) {
    r := setupRedis(b)
    r.LPush("bench_src", "value")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        r.LMove("bench_src", "bench_dst", ListDirectionLeft, ListDirectionRight)
        r.LMove("bench_dst", "bench_src", ListDirectionLeft, ListDirectionRight)
    }
}
```

#### Integration Test Helpers
```go
func setupRedis(t testing.TB) *Redis {
    t.Helper()
    
    r, err := Dial(&DialConfig{Address: "127.0.0.1:6379"})
    if err != nil {
        t.Skip("Redis server not available:", err)
    }
    
    // Clean up function
    t.Cleanup(func() {
        r.FlushDB()
        r.Close()
    })
    
    return r
}

func requireRedisVersion(t *testing.T, r *Redis, minVersion string) {
    t.Helper()
    
    info, err := r.Info()
    if err != nil {
        t.Skip("Cannot get Redis version:", err)
    }
    
    // Version check logic...
    if !isVersionAtLeast(info.RedisVersion, minVersion) {
        t.Skipf("Redis version %s required, got %s", minVersion, info.RedisVersion)
    }
}
```

### Test Coverage Requirements

For each command, implement:
1. **Basic functionality test** - Happy path scenarios
2. **Parameter validation test** - Invalid inputs and edge cases
3. **Error handling test** - Redis error conditions
4. **Integration test** - Real Redis server interaction
5. **Benchmark test** - Performance measurement
6. **Version compatibility test** - Redis version requirements

### Test File Locations
- `client/lists_test.go` - Add new list command tests
- `client/sets_test.go` - Add SMISMEMBER test
- `client/sorted_sets_test.go` - Add new sorted set command tests
- `client/hashes_test.go` - Add new hash command tests
- `client/keys_test.go` - Add new key command tests
- `client/bitmaps_test.go` - New file for bitmap command tests
- `client/connection_test.go` - Add new connection command tests

## Redis Version Compatibility

Commands will include version requirements in comments:
- Redis 3.0+: `WAIT`
- Redis 3.2+: `HSTRLEN`, `BITFIELD`
- Redis 4.0+: `UNLINK`
- Redis 5.0+: `ZPOPMAX`, `ZPOPMIN`, `BZPOPMAX`, `BZPOPMIN`
- Redis 6.0+: `HELLO`, `BITFIELD_RO`, `AUTH` (ACL)
- Redis 6.2+: `LMOVE`, `BLMOVE`, `COPY`, `RESET`, `SMISMEMBER`, `ZRANDMEMBER`, `ZMSCORE`, `HRANDFIELD`
- Redis 7.0+: `LMPOP`, `BLMPOP`

## Success Criteria

Phase 1 is complete when:
- [ ] All 24 commands are implemented with correct signatures
- [ ] All commands have corresponding tests
- [ ] All tests pass with Redis server running
- [ ] Code follows existing libredis patterns and conventions
- [ ] Redis version requirements are documented
- [ ] CLAUDE.md is updated with new command information

## Files to Modify/Create

### Existing Files to Modify:
- `client/lists.go` - Add 5 new list commands
- `client/sets.go` - Add 1 new set command
- `client/sorted_sets.go` - Add 6 new sorted set commands
- `client/hashes.go` - Add 2 new hash commands
- `client/keys.go` - Add 4 new key commands
- `client/connection.go` - Add/enhance 3 connection commands

### New Files to Create:
- `client/bitmaps.go` - New file for bitmap operations
- `client/bitmaps_test.go` - Tests for bitmap operations

### Test Files to Update:
- `client/lists_test.go`
- `client/sets_test.go`
- `client/sorted_sets_test.go`
- `client/hashes_test.go`
- `client/keys_test.go`
- `client/connection_test.go`

### Documentation to Update:
- `CLAUDE.md` - Add Phase 1 command information
- `README.md` - Update supported commands list (optional)

This plan ensures all Phase 1 commands can be implemented systematically in a single Claude session while maintaining code quality and test coverage.