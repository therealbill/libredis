# Advanced Features

This guide covers the advanced features of LibRedis including pipelining, transactions, pub/sub messaging, Lua scripting, and other sophisticated Redis operations.

## Table of Contents

- [Pipelining](#pipelining)
- [Transactions](#transactions)
- [Pub/Sub Messaging](#pubsub-messaging)
- [Lua Scripting](#lua-scripting)
- [HyperLogLog](#hyperloglog)
- [Advanced Bitmap Operations](#advanced-bitmap-operations)
- [Redis Streams](#redis-streams) **NEW**
- [Geospatial Operations](#geospatial-operations) **NEW**
- [Modern Redis Features](#modern-redis-features)
- [Monitoring and Debugging](#monitoring-and-debugging)
- [Performance Optimization](#performance-optimization)

## Pipelining

Pipelining allows you to send multiple commands to Redis without waiting for individual responses, significantly improving performance for bulk operations.

### Basic Pipelining

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/therealbill/libredis/client"
)

func main() {
    redis, err := client.DialWithConfig(&client.DialConfig{
        Network: "tcp",
        Address: "localhost:6379",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer redis.Close()
    
    // Create pipeline
    pipeline := redis.Pipelining()
    defer pipeline.Close()
    
    // Send multiple commands
    pipeline.Command("SET", "key1", "value1")
    pipeline.Command("SET", "key2", "value2")
    pipeline.Command("INCR", "counter")
    pipeline.Command("GET", "key1")
    
    // Receive all responses at once
    responses, err := pipeline.ReceiveAll()
    if err != nil {
        log.Fatal("Pipeline failed:", err)
    }
    
    // Process responses
    for i, response := range responses {
        fmt.Printf("Response %d: %v\n", i, response)
    }
}
```

### Advanced Pipelining Patterns

```go
// Bulk data insertion
func bulkInsert(redis *client.Redis, data map[string]string) error {
    pipeline := redis.Pipelining()
    defer pipeline.Close()
    
    // Send all SET commands
    for key, value := range data {
        pipeline.Command("SET", key, value)
    }
    
    // Execute all commands
    responses, err := pipeline.ReceiveAll()
    if err != nil {
        return fmt.Errorf("bulk insert failed: %w", err)
    }
    
    // Check for errors in responses
    for i, response := range responses {
        if response.Type == client.ErrorReply {
            return fmt.Errorf("command %d failed: %s", i, response.Error)
        }
    }
    
    return nil
}

// Batch operations with mixed commands
func batchOperations(redis *client.Redis) error {
    pipeline := redis.Pipelining()
    defer pipeline.Close()
    
    // Mix different types of operations
    pipeline.Command("LPUSH", "queue", "task1", "task2")
    pipeline.Command("SADD", "tags", "important", "urgent")
    pipeline.Command("HSET", "user:1", "last_seen", "2023-01-01")
    pipeline.Command("INCR", "visits")
    
    responses, err := pipeline.ReceiveAll()
    if err != nil {
        return err
    }
    
    fmt.Printf("Queue length: %v\n", responses[0])
    fmt.Printf("Tags added: %v\n", responses[1])
    fmt.Printf("Hash set: %v\n", responses[2])
    fmt.Printf("Visit count: %v\n", responses[3])
    
    return nil
}
```

## Transactions

Redis transactions provide ACID properties using MULTI/EXEC with optimistic locking via WATCH.

### Basic Transactions

```go
func basicTransaction(redis *client.Redis) error {
    // Start transaction
    txn := redis.Transaction()
    defer txn.Close()
    
    // Queue commands
    txn.Command("SET", "account:1:balance", "1000")
    txn.Command("SET", "account:2:balance", "500")
    txn.Command("INCR", "transaction_count")
    
    // Execute atomically
    results, err := txn.Exec()
    if err != nil {
        return fmt.Errorf("transaction failed: %w", err)
    }
    
    fmt.Printf("Transaction results: %v\n", results)
    return nil
}
```

### Optimistic Locking with WATCH

```go
func transferMoney(redis *client.Redis, fromAccount, toAccount string, amount int64) error {
    for {
        // Start transaction with WATCH
        txn := redis.Transaction()
        defer txn.Close()
        
        // Watch accounts for changes
        err := txn.Watch(fromAccount, toAccount)
        if err != nil {
            return err
        }
        
        // Get current balances
        fromBalance, err := redis.Get(fromAccount)
        if err != nil {
            txn.Discard()
            return err
        }
        
        toBalance, err := redis.Get(toAccount)
        if err != nil {
            txn.Discard()
            return err
        }
        
        // Convert to integers
        fromAmount, _ := strconv.ParseInt(string(fromBalance), 10, 64)
        toAmount, _ := strconv.ParseInt(string(toBalance), 10, 64)
        
        // Check sufficient funds
        if fromAmount < amount {
            txn.Discard()
            return fmt.Errorf("insufficient funds")
        }
        
        // Queue transfer commands
        txn.Command("SET", fromAccount, strconv.FormatInt(fromAmount-amount, 10))
        txn.Command("SET", toAccount, strconv.FormatInt(toAmount+amount, 10))
        txn.Command("INCR", "transfer_count")
        
        // Execute transaction
        results, err := txn.Exec()
        if err != nil {
            return err
        }
        
        // Check if transaction was executed (not discarded due to WATCH)
        if results != nil {
            fmt.Println("Transfer completed successfully")
            return nil
        }
        
        // Transaction was discarded, retry
        fmt.Println("Transaction conflict, retrying...")
        time.Sleep(time.Millisecond * 10)
    }
}
```

### Advanced Transaction Patterns

```go
// Conditional transaction execution
func conditionalUpdate(redis *client.Redis, key, expectedValue, newValue string) error {
    txn := redis.Transaction()
    defer txn.Close()
    
    // Watch the key
    err := txn.Watch(key)
    if err != nil {
        return err
    }
    
    // Check current value
    currentValue, err := redis.Get(key)
    if err != nil {
        txn.Discard()
        return err
    }
    
    // Only proceed if value matches expectation
    if string(currentValue) != expectedValue {
        txn.Discard()
        return fmt.Errorf("value mismatch: expected %s, got %s", expectedValue, string(currentValue))
    }
    
    // Queue update
    txn.Command("SET", key, newValue)
    txn.Command("INCR", "update_count")
    
    // Execute
    results, err := txn.Exec()
    if err != nil {
        return err
    }
    
    if results == nil {
        return fmt.Errorf("transaction aborted due to concurrent modification")
    }
    
    return nil
}
```

## Pub/Sub Messaging

Redis pub/sub provides real-time messaging capabilities for building event-driven applications.

### Basic Pub/Sub

```go
// Publisher
func publisher(redis *client.Redis) {
    for i := 0; i < 10; i++ {
        message := fmt.Sprintf("Message %d", i)
        subscribers, err := redis.Publish("news", message)
        if err != nil {
            log.Printf("Publish error: %v", err)
            continue
        }
        fmt.Printf("Published '%s' to %d subscribers\n", message, subscribers)
        time.Sleep(time.Second)
    }
}

// Subscriber
func subscriber(redis *client.Redis) {
    pubsub, err := redis.PubSub()
    if err != nil {
        log.Fatal("PubSub creation failed:", err)
    }
    defer pubsub.Close()
    
    // Subscribe to channels
    err = pubsub.Subscribe("news", "alerts")
    if err != nil {
        log.Fatal("Subscribe failed:", err)
    }
    
    fmt.Println("Listening for messages...")
    for {
        message := pubsub.Receive()
        if message == nil {
            break
        }
        
        fmt.Printf("Received: Channel=%s, Message=%s\n", 
                   message[1], message[2])
    }
}
```

### Pattern Subscriptions

```go
func patternSubscriber(redis *client.Redis) {
    pubsub, err := redis.PubSub()
    if err != nil {
        log.Fatal(err)
    }
    defer pubsub.Close()
    
    // Subscribe to patterns
    err = pubsub.PSubscribe("user:*", "system:*")
    if err != nil {
        log.Fatal("Pattern subscribe failed:", err)
    }
    
    for {
        message := pubsub.Receive()
        if message == nil {
            break
        }
        
        // Handle pattern messages
        if len(message) >= 4 {
            fmt.Printf("Pattern: %s, Channel: %s, Message: %s\n",
                       message[1], message[2], message[3])
        }
    }
}
```

### Advanced Pub/Sub Patterns

```go
// Event-driven architecture
type EventSystem struct {
    redis  *client.Redis
    pubsub *client.PubSubClient
}

func NewEventSystem(redis *client.Redis) *EventSystem {
    pubsub, _ := redis.PubSub()
    return &EventSystem{
        redis:  redis,
        pubsub: pubsub,
    }
}

func (es *EventSystem) PublishEvent(eventType, data string) error {
    channel := "events:" + eventType
    _, err := es.redis.Publish(channel, data)
    return err
}

func (es *EventSystem) SubscribeToEvents(eventTypes []string, handler func(string, string)) {
    channels := make([]string, len(eventTypes))
    for i, eventType := range eventTypes {
        channels[i] = "events:" + eventType
    }
    
    es.pubsub.Subscribe(channels...)
    
    go func() {
        for {
            message := es.pubsub.Receive()
            if message == nil {
                break
            }
            
            channel := message[1]
            data := message[2]
            eventType := strings.TrimPrefix(channel, "events:")
            
            handler(eventType, data)
        }
    }()
}

func (es *EventSystem) Close() {
    es.pubsub.Close()
}
```

## Lua Scripting

Lua scripting enables atomic server-side operations with custom logic.

### Basic Scripting

```go
// Load and execute Lua script
func basicLuaScript(redis *client.Redis) error {
    script := `
        local key = KEYS[1]
        local value = ARGV[1]
        
        redis.call('SET', key, value)
        local result = redis.call('GET', key)
        return result
    `
    
    // Load script and get SHA1
    sha1, err := redis.ScriptLoad(script)
    if err != nil {
        return err
    }
    
    // Execute by SHA1
    result, err := redis.EvalSha(sha1, []string{"mykey"}, []string{"myvalue"})
    if err != nil {
        return err
    }
    
    fmt.Printf("Script result: %v\n", result)
    return nil
}
```

### Advanced Lua Scripts

```go
// Atomic counter with maximum value
func atomicCounterWithMax(redis *client.Redis, key string, increment, maxValue int64) (int64, error) {
    script := `
        local key = KEYS[1]
        local increment = tonumber(ARGV[1])
        local maxValue = tonumber(ARGV[2])
        
        local current = redis.call('GET', key)
        if not current then
            current = 0
        else
            current = tonumber(current)
        end
        
        local newValue = current + increment
        if newValue > maxValue then
            return {current, false}
        end
        
        redis.call('SET', key, newValue)
        return {newValue, true}
    `
    
    result, err := redis.Eval(script, 
        []string{key}, 
        []string{strconv.FormatInt(increment, 10), strconv.FormatInt(maxValue, 10)})
    if err != nil {
        return 0, err
    }
    
    // Parse result
    if resultArray, ok := result.(*client.Reply); ok && resultArray.Type == client.MultiReply {
        value, _ := resultArray.Multi[0].IntegerValue()
        success, _ := resultArray.Multi[1].IntegerValue()
        if success == 1 {
            return value, nil
        }
        return value, fmt.Errorf("counter would exceed maximum")
    }
    
    return 0, fmt.Errorf("unexpected result format")
}

// Rate limiting with sliding window
func rateLimitSlidingWindow(redis *client.Redis, key string, windowSize, limit int64) (bool, error) {
    script := `
        local key = KEYS[1]
        local window = tonumber(ARGV[1])
        local limit = tonumber(ARGV[2])
        local now = tonumber(ARGV[3])
        
        -- Remove expired entries
        redis.call('ZREMRANGEBYSCORE', key, 0, now - window)
        
        -- Count current entries
        local current = redis.call('ZCARD', key)
        
        if current < limit then
            -- Add current request
            redis.call('ZADD', key, now, now)
            redis.call('EXPIRE', key, window)
            return 1
        else
            return 0
        end
    `
    
    now := time.Now().Unix()
    result, err := redis.Eval(script, 
        []string{key}, 
        []string{
            strconv.FormatInt(windowSize, 10),
            strconv.FormatInt(limit, 10),
            strconv.FormatInt(now, 10),
        })
    if err != nil {
        return false, err
    }
    
    if intResult, err := result.IntegerValue(); err == nil {
        return intResult == 1, nil
    }
    
    return false, fmt.Errorf("unexpected result type")
}
```

## HyperLogLog

HyperLogLog provides memory-efficient cardinality estimation for large datasets.

### Basic HyperLogLog Operations

```go
func hyperLogLogExample(redis *client.Redis) error {
    // Add elements to HyperLogLog
    added, err := redis.PFAdd("unique_visitors", "user1", "user2", "user3")
    if err != nil {
        return err
    }
    fmt.Printf("Added %d new unique elements\n", added)
    
    // Add more elements (duplicates won't increase count)
    redis.PFAdd("unique_visitors", "user1", "user4", "user5")
    
    // Get cardinality estimate
    count, err := redis.PFCount("unique_visitors")
    if err != nil {
        return err
    }
    fmt.Printf("Estimated unique visitors: %d\n", count)
    
    return nil
}
```

### HyperLogLog Merging

```go
func hyperLogLogMerging(redis *client.Redis) error {
    // Create separate HLLs for different pages
    redis.PFAdd("page1_visitors", "user1", "user2", "user3")
    redis.PFAdd("page2_visitors", "user2", "user4", "user5")
    redis.PFAdd("page3_visitors", "user1", "user5", "user6")
    
    // Merge all page visitors into total
    err := redis.PFMerge("total_visitors", "page1_visitors", "page2_visitors", "page3_visitors")
    if err != nil {
        return err
    }
    
    // Get total unique visitors across all pages
    totalCount, err := redis.PFCount("total_visitors")
    if err != nil {
        return err
    }
    
    fmt.Printf("Total unique visitors across all pages: %d\n", totalCount)
    return nil
}
```

## Advanced Bitmap Operations

Modern Redis bitmap operations for efficient bit manipulation.

### BitField Operations

```go
func bitFieldOperations(redis *client.Redis) error {
    key := "bitfield_example"
    
    // Define bit field operations
    operations := []client.BitFieldOperation{
        {Type: "SET", Offset: 0, Value: 100},     // Set 8-bit value at offset 0
        {Type: "SET", Offset: 8, Value: 200},     // Set 8-bit value at offset 8
        {Type: "INCRBY", Offset: 0, Value: 50},   // Increment value at offset 0
        {Type: "GET", Offset: 0},                 // Get value at offset 0
        {Type: "GET", Offset: 8},                 // Get value at offset 8
    }
    
    // Execute bit field operations
    results, err := redis.BitField(key, operations)
    if err != nil {
        return err
    }
    
    fmt.Printf("BitField results: %v\n", results)
    
    // Bit field with overflow control
    overflowOps := []client.BitFieldOperation{
        {Type: "INCRBY", Offset: 0, Value: 200}, // This might overflow
    }
    
    results, err = redis.BitFieldWithOverflow(key, client.BitFieldOverflowSat, overflowOps)
    if err != nil {
        return err
    }
    
    fmt.Printf("BitField with overflow control: %v\n", results)
    return nil
}
```

### Bit Position Operations

```go
func bitPositionOperations(redis *client.Redis) error {
    key := "bitpos_example"
    
    // Set some bits
    redis.SetBit(key, 10, 1)
    redis.SetBit(key, 20, 1)
    redis.SetBit(key, 30, 1)
    
    // Find first bit set to 1
    pos, err := redis.BitPos(key, 1)
    if err != nil {
        return err
    }
    fmt.Printf("First bit set to 1 at position: %d\n", pos)
    
    // Find first bit set to 1 in range
    start := int64(15)
    end := int64(25)
    opts := client.BitPosOptions{
        Start: &start,
        End:   &end,
    }
    
    pos, err = redis.BitPosWithRange(key, 1, opts)
    if err != nil {
        return err
    }
    fmt.Printf("First bit set to 1 in range [15,25]: %d\n", pos)
    
    return nil
}
```

## Redis Streams

Redis Streams provide a powerful log-like data structure for event streaming and message processing, with consumer group capabilities for building distributed applications.

### Basic Stream Operations

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/therealbill/libredis/client"
)

func basicStreamOperations(redis *client.Redis) error {
    // Add events to a stream
    eventData := map[string]string{
        "event":    "user_login",
        "user_id":  "12345",
        "ip":       "192.168.1.100",
        "timestamp": "2023-07-12T10:30:00Z",
    }
    
    // Add entry with auto-generated ID
    entryID, err := redis.XAdd("user_events", client.StreamIDAutoGenerate, eventData)
    if err != nil {
        return fmt.Errorf("failed to add stream entry: %w", err)
    }
    fmt.Printf("Added entry with ID: %s\n", entryID)
    
    // Add more events
    loginData := map[string]string{
        "event":   "user_logout",
        "user_id": "12345",
    }
    redis.XAdd("user_events", client.StreamIDAutoGenerate, loginData)
    
    // Read all events from the beginning
    streams := map[string]string{"user_events": "0-0"}
    messages, err := redis.XRead(streams)
    if err != nil {
        return fmt.Errorf("failed to read stream: %w", err)
    }
    
    // Process messages
    for _, streamMsg := range messages {
        fmt.Printf("Stream: %s\n", streamMsg.Stream)
        for _, entry := range streamMsg.Entries {
            fmt.Printf("  ID: %s\n", entry.ID)
            for field, value := range entry.Fields {
                fmt.Printf("    %s: %s\n", field, value)
            }
        }
    }
    
    return nil
}
```

### Consumer Groups for Distributed Processing

```go
func consumerGroupExample(redis *client.Redis) error {
    streamKey := "task_queue"
    groupName := "processors"
    
    // Create consumer group (start from latest messages)
    err := redis.XGroupCreate(streamKey, groupName, "$")
    if err != nil {
        return fmt.Errorf("failed to create consumer group: %w", err)
    }
    
    // Add some tasks to the stream
    tasks := []map[string]string{
        {"task": "process_image", "image_id": "img_001", "priority": "high"},
        {"task": "send_email", "user_id": "user_123", "template": "welcome"},
        {"task": "generate_report", "report_type": "daily", "date": "2023-07-12"},
    }
    
    for _, task := range tasks {
        redis.XAdd(streamKey, client.StreamIDAutoGenerate, task)
    }
    
    // Process messages as different consumers
    consumers := []string{"worker_1", "worker_2", "worker_3"}
    
    for _, consumer := range consumers {
        // Read new messages for this consumer
        streams := map[string]string{streamKey: ">"}
        messages, err := redis.XReadGroup(groupName, consumer, streams)
        if err != nil {
            continue // No new messages for this consumer
        }
        
        // Process each message
        for _, streamMsg := range messages {
            for _, entry := range streamMsg.Entries {
                fmt.Printf("Consumer %s processing task %s: %v\n", 
                    consumer, entry.ID, entry.Fields)
                
                // Simulate task processing
                if processTask(entry.Fields) {
                    // Acknowledge successful processing
                    acked, err := redis.XAck(streamKey, groupName, entry.ID)
                    if err == nil && acked > 0 {
                        fmt.Printf("Task %s acknowledged by %s\n", entry.ID, consumer)
                    }
                }
            }
        }
    }
    
    return nil
}

func processTask(task map[string]string) bool {
    // Simulate task processing logic
    fmt.Printf("Processing task: %s\n", task["task"])
    return true // Return true if task completed successfully
}
```

### Stream Management and Monitoring

```go
func streamManagement(redis *client.Redis) error {
    streamKey := "events"
    
    // Get stream length
    length, err := redis.XLen(streamKey)
    if err != nil {
        return err
    }
    fmt.Printf("Stream length: %d\n", length)
    
    // Get stream information
    info, err := redis.XInfoStream(streamKey)
    if err != nil {
        return err
    }
    fmt.Printf("Stream info: %+v\n", info)
    
    // Trim stream to keep only last 1000 entries
    trimmed, err := redis.XTrim(streamKey, "MAXLEN", "1000")
    if err != nil {
        return err
    }
    fmt.Printf("Trimmed %d entries\n", trimmed)
    
    // Get pending messages for consumer group
    groupName := "processors"
    pending, err := redis.XPending(streamKey, groupName)
    if err != nil {
        return err
    }
    fmt.Printf("Pending messages: %d\n", pending.Count)
    
    return nil
}
```

## Geospatial Operations

Geospatial operations enable location-based applications with support for coordinates, distance calculations, and proximity searches.

### Basic Geospatial Operations

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/therealbill/libredis/client"
)

func basicGeospatialOperations(redis *client.Redis) error {
    // Add locations to a geospatial index
    locations := []client.GeoMember{
        {Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
        {Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
        {Longitude: 2.3522, Latitude: 48.8566, Member: "Paris"},
        {Longitude: 139.6917, Latitude: 35.6895, Member: "Tokyo"},
        {Longitude: -0.1276, Latitude: 51.5074, Member: "London"},
    }
    
    added, err := redis.GeoAdd("cities", locations)
    if err != nil {
        return fmt.Errorf("failed to add locations: %w", err)
    }
    fmt.Printf("Added %d cities to geospatial index\n", added)
    
    // Calculate distance between cities
    distance, err := redis.GeoDist("cities", "San Francisco", "New York")
    if err != nil {
        return err
    }
    fmt.Printf("Distance SF to NYC: %.2f meters\n", distance)
    
    // Get distance in different units
    distanceKM, err := redis.GeoDistWithUnit("cities", "Paris", "London", client.GeoUnitKilometers)
    if err != nil {
        return err
    }
    fmt.Printf("Distance Paris to London: %.2f km\n", distanceKM)
    
    // Get coordinates for cities
    positions, err := redis.GeoPos("cities", "Tokyo", "Paris")
    if err != nil {
        return err
    }
    
    for i, city := range []string{"Tokyo", "Paris"} {
        if positions[i] != nil {
            fmt.Printf("%s coordinates: %.4f, %.4f\n", 
                city, positions[i].Longitude, positions[i].Latitude)
        }
    }
    
    return nil
}
```

### Proximity Search and Location-Based Queries

```go
func proximitySearch(redis *client.Redis) error {
    // Modern search using GEOSEARCH (Redis 6.2+)
    
    // Search for cities within 500km of San Francisco
    searchOpts := client.GeoSearchOptions{
        FromMember: stringPtr("San Francisco"),
        ByRadius:   &client.GeoRadius{Radius: 500, Unit: client.GeoUnitKilometers},
        WithCoord:  true,
        WithDist:   true,
        Order:      client.GeoOrderAsc, // Closest first
    }
    
    nearby, err := redis.GeoSearch("cities", searchOpts)
    if err != nil {
        return fmt.Errorf("geosearch failed: %w", err)
    }
    
    fmt.Println("Cities within 500km of San Francisco:")
    for _, location := range nearby {
        fmt.Printf("  %s", location.Member)
        if location.Distance != nil {
            fmt.Printf(" (%.2f km)", *location.Distance)
        }
        if location.Coordinates != nil {
            fmt.Printf(" at %.4f, %.4f", 
                location.Coordinates.Longitude, location.Coordinates.Latitude)
        }
        fmt.Println()
    }
    
    // Search within a rectangular area
    boxOpts := client.GeoSearchOptions{
        FromLonLat: &client.GeoCoordinate{Longitude: 0, Latitude: 50}, // Center of Europe
        ByBox:      &client.GeoBox{Width: 1000, Height: 1000, Unit: client.GeoUnitKilometers},
        WithCoord:  true,
        Count:      5, // Limit results
    }
    
    europeanCities, err := redis.GeoSearch("cities", boxOpts)
    if err != nil {
        return err
    }
    
    fmt.Println("\nCities in European region:")
    for _, location := range europeanCities {
        fmt.Printf("  %s at %.4f, %.4f\n", location.Member,
            location.Coordinates.Longitude, location.Coordinates.Latitude)
    }
    
    return nil
}
```

### Advanced Geospatial Features

```go
func advancedGeospatialFeatures(redis *client.Redis) error {
    // Store search results for later use
    storeOpts := client.GeoSearchStoreOptions{
        GeoSearchOptions: client.GeoSearchOptions{
            FromMember: stringPtr("New York"),
            ByRadius:   &client.GeoRadius{Radius: 6000, Unit: client.GeoUnitKilometers},
            WithDist:   true,
        },
        StoreDist: true, // Store distances as scores
    }
    
    stored, err := redis.GeoSearchStore("nearby_nyc", "cities", storeOpts)
    if err != nil {
        return err
    }
    fmt.Printf("Stored %d cities near NYC\n", stored)
    
    // Get geohashes for efficient indexing/storage
    hashes, err := redis.GeoHash("cities", "San Francisco", "Tokyo", "London")
    if err != nil {
        return err
    }
    
    fmt.Println("Geohashes:")
    cities := []string{"San Francisco", "Tokyo", "London"}
    for i, hash := range hashes {
        if hash != "" {
            fmt.Printf("  %s: %s\n", cities[i], hash)
        }
    }
    
    // Use legacy GEORADIUS for compatibility (deprecated but still supported)
    legacyResults, err := redis.GeoRadius("cities", -74.0060, 40.7128, 1000, client.GeoUnitKilometers)
    if err != nil {
        return err
    }
    fmt.Printf("Legacy radius search found %d cities\n", len(legacyResults))
    
    return nil
}

// Utility function for creating string pointers
func stringPtr(s string) *string {
    return &s
}
```

### Real-World Use Cases

```go
// Example: Delivery service application
func deliveryServiceExample(redis *client.Redis) error {
    // Add delivery drivers to geospatial index
    drivers := []client.GeoMember{
        {Longitude: -122.4194, Latitude: 37.7749, Member: "driver_001"},
        {Longitude: -122.4094, Latitude: 37.7849, Member: "driver_002"},
        {Longitude: -122.4294, Latitude: 37.7649, Member: "driver_003"},
    }
    
    redis.GeoAdd("active_drivers", drivers)
    
    // Customer requests delivery from specific location
    customerLoc := &client.GeoCoordinate{Longitude: -122.4150, Latitude: 37.7750}
    
    // Find nearest available drivers within 2km
    searchOpts := client.GeoSearchOptions{
        FromLonLat: customerLoc,
        ByRadius:   &client.GeoRadius{Radius: 2, Unit: client.GeoUnitKilometers},
        WithDist:   true,
        Order:      client.GeoOrderAsc, // Nearest first
        Count:      3, // Top 3 nearest drivers
    }
    
    nearbyDrivers, err := redis.GeoSearch("active_drivers", searchOpts)
    if err != nil {
        return err
    }
    
    fmt.Println("Available drivers:")
    for _, driver := range nearbyDrivers {
        fmt.Printf("  Driver %s: %.3f km away\n", driver.Member, *driver.Distance)
    }
    
    return nil
}
```

## Modern Redis Features

LibRedis supports the latest Redis features from versions 6.0+ and 7.0+.

### ACL Authentication (Redis 6.0+)

```go
func aclAuthentication(redis *client.Redis) error {
    // Authenticate with username and password
    err := redis.AuthWithUser("myuser", "mypassword")
    if err != nil {
        return fmt.Errorf("ACL authentication failed: %w", err)
    }
    
    // Use HELLO command for protocol negotiation
    info, err := redis.Hello(3)
    if err != nil {
        return err
    }
    
    fmt.Printf("Server info: %v\n", info)
    return nil
}
```

### Modern List Operations (Redis 6.2+/7.0+)

```go
func modernListOperations(redis *client.Redis) error {
    // LMOVE - atomic move between lists (Redis 6.2+)
    redis.LPush("source", "item1", "item2", "item3")
    
    value, err := redis.LMove("source", "dest", 
                             client.ListDirectionRight, 
                             client.ListDirectionLeft)
    if err != nil {
        return err
    }
    fmt.Printf("Moved value: %s\n", value)
    
    // LPOS - find position of element (Redis 6.0.6+)
    redis.LPush("mylist", "a", "b", "c", "b", "d")
    
    pos, err := redis.LPos("mylist", "b")
    if err != nil {
        return err
    }
    fmt.Printf("Position of 'b': %d\n", pos)
    
    // LPOS with options
    opts := client.LPosOptions{
        Rank:   2,  // Find 2nd occurrence
        Count:  2,  // Return up to 2 positions
        MaxLen: 10, // Search only first 10 elements
    }
    
    positions, err := redis.LPosWithOptions("mylist", "b", opts)
    if err != nil {
        return err
    }
    fmt.Printf("Positions with options: %v\n", positions)
    
    // LMPOP - pop from multiple lists (Redis 7.0+)
    redis.LPush("list1", "item1", "item2")
    redis.LPush("list2", "item3", "item4")
    
    result, err := redis.LMPop([]string{"list1", "list2"}, client.ListDirectionLeft)
    if err != nil {
        return err
    }
    fmt.Printf("Multi-pop result: %v\n", result)
    
    return nil
}
```

### Advanced Set/Hash Operations (Redis 6.2+)

```go
func modernSetHashOperations(redis *client.Redis) error {
    // SMISMEMBER - check multiple memberships (Redis 6.2+)
    redis.SAdd("myset", "member1", "member2", "member3")
    
    results, err := redis.SMIsMember("myset", "member1", "member4", "member2")
    if err != nil {
        return err
    }
    fmt.Printf("Membership results: %v\n", results)
    
    // HRANDFIELD - get random hash fields (Redis 6.2+)
    redis.HMSet("myhash", map[string]string{
        "field1": "value1",
        "field2": "value2",
        "field3": "value3",
    })
    
    field, err := redis.HRandField("myhash")
    if err != nil {
        return err
    }
    fmt.Printf("Random field: %s\n", field)
    
    return nil
}
```

## Monitoring and Debugging

Tools for monitoring Redis performance and debugging issues.

### Connection Monitoring

```go
func monitorRedis(redis *client.Redis) error {
    // Monitor Redis commands (use carefully in production)
    monitor, err := redis.Monitor()
    if err != nil {
        return err
    }
    
    go func() {
        for {
            command := monitor.Receive()
            if command == nil {
                break
            }
            fmt.Printf("Command: %s\n", strings.Join(command, " "))
        }
    }()
    
    // Let it run for a bit
    time.Sleep(10 * time.Second)
    monitor.Close()
    
    return nil
}
```

### Slow Query Analysis

```go
func analyzeSlowQueries(redis *client.Redis) error {
    // Get slow query log
    slowQueries, err := redis.SlowLogGet(10) // Get last 10 slow queries
    if err != nil {
        return err
    }
    
    fmt.Printf("Found %d slow queries:\n", len(slowQueries))
    for i, query := range slowQueries {
        fmt.Printf("Query %d: %v\n", i, query)
    }
    
    // Get slow log length
    length, err := redis.SlowLogLen()
    if err != nil {
        return err
    }
    fmt.Printf("Total slow queries: %d\n", length)
    
    return nil
}
```

## Performance Optimization

### Connection Pool Tuning

```go
func optimizeConnectionPool() *client.Redis {
    config := &client.DialConfig{
        Network:      "tcp",
        Address:      "localhost:6379",
        MaxIdle:      50,                    // Increase for high concurrency
        Timeout:      10 * time.Second,     // Reasonable timeout
        TCPKeepAlive: 30,                   // Keep connections alive
    }
    
    redis, err := client.DialWithConfig(config)
    if err != nil {
        log.Fatal(err)
    }
    
    return redis
}
```

### Efficient Bulk Operations

```go
func efficientBulkOperations(redis *client.Redis) error {
    // Use pipelining for bulk operations
    pipeline := redis.Pipelining()
    defer pipeline.Close()
    
    // Batch size optimization
    batchSize := 1000
    data := generateTestData(10000) // Your data source
    
    for i := 0; i < len(data); i += batchSize {
        end := i + batchSize
        if end > len(data) {
            end = len(data)
        }
        
        // Send batch
        for j := i; j < end; j++ {
            pipeline.Command("SET", data[j].Key, data[j].Value)
        }
        
        // Execute batch
        _, err := pipeline.ReceiveAll()
        if err != nil {
            return fmt.Errorf("batch %d failed: %w", i/batchSize, err)
        }
    }
    
    return nil
}
```

This comprehensive guide covers the advanced features of LibRedis. These patterns and techniques will help you build robust, high-performance applications with Redis.