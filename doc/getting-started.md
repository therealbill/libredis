# Getting Started with LibRedis

LibRedis is a comprehensive Redis client library for Go that provides full support for Redis commands, connection pooling, and advanced features like pipelining, transactions, and pub/sub messaging.

## Table of Contents

- [Installation](#installation)
- [Basic Connection](#basic-connection)
- [Connection Configuration](#connection-configuration)
- [Basic Operations](#basic-operations)
- [Redis Streams](#redis-streams-operations-redis-50)
- [Geospatial Operations](#geospatial-operations-redis-32)
- [Error Handling](#error-handling)
- [Connection Pooling](#connection-pooling)
- [SSL/TLS Support](#ssltls-support)
- [Best Practices](#best-practices)

## Installation

Install LibRedis using Go modules:

```bash
go get github.com/therealbill/libredis
```

Import in your Go code:

```go
import "github.com/therealbill/libredis/client"
```

## Basic Connection

### Simple Connection

Connect to Redis running on localhost with default settings:

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/therealbill/libredis/client"
)

func main() {
    // Connect to Redis
    redis, err := client.DialWithConfig(&client.DialConfig{
        Network:  "tcp",
        Address:  "127.0.0.1:6379",
        Database: 0,
        Password: "",
        Timeout:  5 * time.Second,
        MaxIdle:  10,
    })
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }
    defer redis.Close()
    
    // Test connection
    err = redis.Ping()
    if err != nil {
        log.Fatal("Failed to ping Redis:", err)
    }
    
    fmt.Println("Connected to Redis successfully!")
}
```

### URL-based Connection

Connect using a Redis URL:

```go
redis, err := client.DialURL("tcp://auth:password@127.0.0.1:6379/0?timeout=5s&maxidle=10")
if err != nil {
    log.Fatal("Failed to connect:", err)
}
defer redis.Close()
```

## Connection Configuration

### DialConfig Structure

```go
type DialConfig struct {
    Network       string        // "tcp" or "unix"
    Address       string        // "host:port" or "/path/to/socket"
    Database      int           // Redis database number
    Password      string        // Redis password
    Timeout       time.Duration // Connection timeout
    MaxIdle       int           // Max idle connections in pool
    SSL           bool          // Enable SSL/TLS
    SSLSkipVerify bool          // Skip SSL certificate verification
    SSLCert       string        // SSL certificate file path
    SSLKey        string        // SSL private key file path
    SSLCA         string        // SSL CA certificate file path
    TCPKeepAlive  int           // TCP keep-alive interval (seconds)
}
```

### Configuration Examples

**Basic Configuration:**
```go
config := &client.DialConfig{
    Network:      "tcp",
    Address:      "localhost:6379",
    Database:     0,
    Password:     "",
    Timeout:      5 * time.Second,
    MaxIdle:      10,
    TCPKeepAlive: 60,
}
```

**SSL Configuration:**
```go
config := &client.DialConfig{
    Network:       "tcp",
    Address:       "redis.example.com:6380",
    Database:      0,
    Password:      "secret",
    Timeout:       10 * time.Second,
    MaxIdle:       20,
    SSL:           true,
    SSLSkipVerify: false,
    SSLCert:       "/path/to/client.crt",
    SSLKey:        "/path/to/client.key",
    SSLCA:         "/path/to/ca.crt",
}
```

## Basic Operations

### String Operations

```go
// Set and get string values
err := redis.Set("username", "john_doe")
if err != nil {
    log.Fatal("Set failed:", err)
}

value, err := redis.Get("username")
if err != nil {
    log.Fatal("Get failed:", err)
}
fmt.Println("Username:", string(value))

// Set with expiration
err = redis.Setex("session_token", 3600, "abc123")
if err != nil {
    log.Fatal("Setex failed:", err)
}

// Increment counters
count, err := redis.Incr("page_views")
if err != nil {
    log.Fatal("Incr failed:", err)
}
fmt.Println("Page views:", count)
```

### List Operations

```go
// Push items to list
length, err := redis.LPush("tasks", "task1", "task2", "task3")
if err != nil {
    log.Fatal("LPush failed:", err)
}
fmt.Println("List length:", length)

// Get list items
items, err := redis.LRange("tasks", 0, -1)
if err != nil {
    log.Fatal("LRange failed:", err)
}
fmt.Println("Tasks:", items)

// Pop items from list
item, err := redis.LPop("tasks")
if err != nil {
    log.Fatal("LPop failed:", err)
}
if item != nil {
    fmt.Println("Popped task:", string(item))
}
```

### Hash Operations

```go
// Set hash fields
success, err := redis.HSet("user:1", "name", "John Doe")
if err != nil {
    log.Fatal("HSet failed:", err)
}

// Get hash field
name, err := redis.HGet("user:1", "name")
if err != nil {
    log.Fatal("HGet failed:", err)
}
if name != nil {
    fmt.Println("User name:", string(name))
}

// Set multiple fields
fields := map[string]string{
    "email": "john@example.com",
    "age":   "30",
}
err = redis.HMSet("user:1", fields)
if err != nil {
    log.Fatal("HMSet failed:", err)
}

// Get all fields
userdata, err := redis.HGetAll("user:1")
if err != nil {
    log.Fatal("HGetAll failed:", err)
}
fmt.Println("User data:", userdata)
```

### Set Operations

```go
// Add members to set
count, err := redis.SAdd("tags", "redis", "database", "nosql")
if err != nil {
    log.Fatal("SAdd failed:", err)
}
fmt.Println("Added tags:", count)

// Check membership
isMember, err := redis.SIsMember("tags", "redis")
if err != nil {
    log.Fatal("SIsMember failed:", err)
}
fmt.Println("Is 'redis' a tag?", isMember)

// Get all members
members, err := redis.SMembers("tags")
if err != nil {
    log.Fatal("SMembers failed:", err)
}
fmt.Println("All tags:", members)
```

### Sorted Set Operations

```go
// Add scored members
count, err := redis.ZAdd("leaderboard", 100, "player1")
if err != nil {
    log.Fatal("ZAdd failed:", err)
}

// Add multiple members
scores := map[string]float64{
    "player2": 150,
    "player3": 75,
}
count, err = redis.ZAddVariadic("leaderboard", scores)
if err != nil {
    log.Fatal("ZAddVariadic failed:", err)
}

// Get top players
topPlayers, err := redis.ZRevRange("leaderboard", 0, 2, true)
if err != nil {
    log.Fatal("ZRevRange failed:", err)
}
fmt.Println("Top players:", topPlayers)
```

### Redis Streams Operations (Redis 5.0+)

Redis Streams provide powerful event streaming capabilities for building real-time applications:

```go
// Add events to a stream
fields := map[string]string{
    "user_id": "123",
    "action":  "login",
    "ip":      "192.168.1.100",
}
entryID, err := redis.XAdd("events", client.StreamIDAutoGenerate, fields)
if err != nil {
    log.Fatal("XAdd failed:", err)
}
fmt.Println("Added event with ID:", entryID)

// Read from streams
streams := map[string]string{"events": "0-0"} // Read from beginning
messages, err := redis.XRead(streams)
if err != nil {
    log.Fatal("XRead failed:", err)
}

for _, msg := range messages {
    fmt.Printf("Stream: %s\n", msg.Stream)
    for _, entry := range msg.Entries {
        fmt.Printf("  ID: %s, Fields: %v\n", entry.ID, entry.Fields)
    }
}

// Consumer groups for reliable processing
err = redis.XGroupCreate("events", "processors", client.StreamIDLatest)
if err != nil {
    log.Printf("Group creation failed (may already exist): %v", err)
}

// Read as a consumer
consumerStreams := map[string]string{"events": ">"}
groupMessages, err := redis.XReadGroup("processors", "worker1", consumerStreams)
if err != nil {
    log.Fatal("XReadGroup failed:", err)
}

// Process and acknowledge messages
for _, msg := range groupMessages {
    for _, entry := range msg.Entries {
        fmt.Printf("Processing: %v\n", entry.Fields)
        
        // Acknowledge message after processing
        ackCount, err := redis.XAck("events", "processors", entry.ID)
        if err != nil {
            log.Printf("XAck failed: %v", err)
        } else {
            fmt.Printf("Acknowledged %d messages\n", ackCount)
        }
    }
}
```

### Geospatial Operations (Redis 3.2+)

Geospatial operations enable location-based applications with coordinate storage and proximity searches:

```go
// Add locations to a geospatial index
locations := []client.GeoMember{
    {Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
    {Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
    {Longitude: -87.6298, Latitude: 41.8781, Member: "Chicago"},
    {Longitude: 2.3522, Latitude: 48.8566, Member: "Paris"},
}

count, err := redis.GeoAdd("cities", locations)
if err != nil {
    log.Fatal("GeoAdd failed:", err)
}
fmt.Printf("Added %d cities\n", count)

// Calculate distance between cities
distance, err := redis.GeoDist("cities", "San Francisco", "New York")
if err != nil {
    log.Fatal("GeoDist failed:", err)
}
fmt.Printf("Distance SF to NYC: %.2f meters\n", distance)

// Distance with specific unit
distanceKM, err := redis.GeoDistWithUnit("cities", "San Francisco", "New York", "km")
if err != nil {
    log.Fatal("GeoDistWithUnit failed:", err)
}
fmt.Printf("Distance SF to NYC: %.2f km\n", distanceKM)

// Find nearby locations (modern search - Redis 6.2+)
searchOpts := client.GeoSearchOptions{
    FromLonLat: &client.GeoCoordinate{
        Longitude: -122.4194, 
        Latitude:  37.7749,
    },
    ByRadius: &client.GeoRadius{
        Radius: 1000, 
        Unit:   "km",
    },
    WithCoord: true,
    WithDist:  true,
    Count:     5,
}

nearbyLocations, err := redis.GeoSearch("cities", searchOpts)
if err != nil {
    log.Fatal("GeoSearch failed:", err)
}

fmt.Println("Nearby locations:")
for _, location := range nearbyLocations {
    fmt.Printf("  %s", location.Member)
    if location.Distance != nil {
        fmt.Printf(" (%.2f km)", *location.Distance)
    }
    if location.Coordinates != nil {
        fmt.Printf(" [%.4f, %.4f]", location.Coordinates.Longitude, location.Coordinates.Latitude)
    }
    fmt.Println()
}

// Get coordinates of specific locations
coords, err := redis.GeoPos("cities", "San Francisco", "Paris")
if err != nil {
    log.Fatal("GeoPos failed:", err)
}

for i, coord := range coords {
    if coord != nil {
        cities := []string{"San Francisco", "Paris"}
        fmt.Printf("%s: [%.4f, %.4f]\n", cities[i], coord.Longitude, coord.Latitude)
    }
}
```

## Error Handling

LibRedis provides detailed error information for different scenarios:

```go
// Handle different types of errors
value, err := redis.Get("nonexistent_key")
if err != nil {
    log.Printf("Error: %v", err)
    return
}

// Check for nil values (key doesn't exist)
if value == nil {
    fmt.Println("Key doesn't exist")
} else {
    fmt.Println("Value:", string(value))
}

// Connection errors
err = redis.Ping()
if err != nil {
    log.Fatal("Connection lost:", err)
}
```

### Common Error Patterns

```go
// Graceful error handling
func getUser(redis *client.Redis, userID string) (map[string]string, error) {
    userData, err := redis.HGetAll("user:" + userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %s: %w", userID, err)
    }
    
    if len(userData) == 0 {
        return nil, fmt.Errorf("user %s not found", userID)
    }
    
    return userData, nil
}
```

## Connection Pooling

LibRedis automatically manages connection pools to optimize performance:

### Pool Configuration

```go
config := &client.DialConfig{
    Network:      "tcp",
    Address:      "localhost:6379",
    MaxIdle:      50,                // Maximum idle connections
    Timeout:      10 * time.Second,  // Connection timeout
    TCPKeepAlive: 30,                // Keep-alive interval
}

redis, err := client.DialWithConfig(config)
if err != nil {
    log.Fatal("Connection failed:", err)
}
```

### Pool Best Practices

1. **Set appropriate MaxIdle**: Based on your application's concurrency needs
2. **Configure timeouts**: Prevent hanging connections
3. **Monitor pool usage**: Check for connection leaks
4. **Graceful shutdown**: Always close connections

```go
// Graceful shutdown
func gracefulShutdown(redis *client.Redis) {
    redis.Close() // Closes all pooled connections
}
```

## SSL/TLS Support

### Basic SSL Configuration

```go
config := &client.DialConfig{
    Network:  "tcp",
    Address:  "secure-redis.example.com:6380",
    Password: "secure_password",
    SSL:      true,
    SSLSkipVerify: false, // Verify certificates in production
}

redis, err := client.DialWithConfig(config)
if err != nil {
    log.Fatal("SSL connection failed:", err)
}
```

### Client Certificate Authentication

```go
config := &client.DialConfig{
    Network:  "tcp",
    Address:  "redis.example.com:6380",
    SSL:      true,
    SSLCert:  "/path/to/client.crt",
    SSLKey:   "/path/to/client.key",
    SSLCA:    "/path/to/ca.crt",
}
```

## Best Practices

### 1. Connection Management

```go
// ✅ Good: Single connection instance
var redisClient *client.Redis

func init() {
    var err error
    redisClient, err = client.DialWithConfig(&client.DialConfig{
        Network: "tcp",
        Address: "localhost:6379",
        MaxIdle: 10,
        Timeout: 5 * time.Second,
    })
    if err != nil {
        log.Fatal("Redis connection failed:", err)
    }
}

// ❌ Bad: Creating new connections for each operation
func badExample() {
    redis, _ := client.DialWithConfig(config) // Don't do this repeatedly
    redis.Set("key", "value")
    redis.Close()
}
```

### 2. Error Handling

```go
// ✅ Good: Proper error handling
func setUserData(userID string, data map[string]string) error {
    err := redisClient.HMSet("user:"+userID, data)
    if err != nil {
        return fmt.Errorf("failed to save user data: %w", err)
    }
    return nil
}

// ❌ Bad: Ignoring errors
func badExample() {
    redisClient.Set("key", "value") // Don't ignore errors
}
```

### 3. Key Naming

```go
// ✅ Good: Consistent key naming
const (
    UserPrefix    = "user:"
    SessionPrefix = "session:"
    CachePrefix   = "cache:"
)

func getUserKey(userID string) string {
    return UserPrefix + userID
}

// ❌ Bad: Inconsistent naming
func badExample() {
    redisClient.Set("user_123", "data")
    redisClient.Set("USER:456", "data")
    redisClient.Set("users/789", "data")
}
```

### 4. Resource Cleanup

```go
// ✅ Good: Proper cleanup
func main() {
    redis, err := client.DialWithConfig(config)
    if err != nil {
        log.Fatal(err)
    }
    defer redis.Close() // Always close connections
    
    // Your application logic here
}
```

### 5. Configuration Management

```go
// ✅ Good: Environment-based configuration
func createRedisConfig() *client.DialConfig {
    return &client.DialConfig{
        Network:  getEnv("REDIS_NETWORK", "tcp"),
        Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
        Password: getEnv("REDIS_PASSWORD", ""),
        Database: getEnvInt("REDIS_DB", 0),
        Timeout:  time.Duration(getEnvInt("REDIS_TIMEOUT", 5)) * time.Second,
        MaxIdle:  getEnvInt("REDIS_MAX_IDLE", 10),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

## Next Steps

Now that you have the basics covered, explore more advanced features:

- **[Advanced Features](advanced-features.md)** - Pipelining, transactions, pub/sub
- **[Commands Reference](commands.md)** - Complete command documentation
- **[API Reference](api-reference.md)** - Detailed API documentation
- **[Sentinel Guide](sentinel.md)** - High availability with Redis Sentinel

## Common Use Cases

### Caching

```go
func getFromCache(key string) ([]byte, error) {
    data, err := redisClient.Get("cache:" + key)
    if err != nil {
        return nil, err
    }
    if data == nil {
        return nil, nil // Cache miss
    }
    return data, nil
}

func setCache(key string, data []byte, ttl int) error {
    return redisClient.Setex("cache:"+key, ttl, string(data))
}
```

### Session Management

```go
func saveSession(sessionID string, userData map[string]string) error {
    key := "session:" + sessionID
    err := redisClient.HMSet(key, userData)
    if err != nil {
        return err
    }
    return redisClient.Expire(key, 3600) // 1 hour expiration
}

func getSession(sessionID string) (map[string]string, error) {
    return redisClient.HGetAll("session:" + sessionID)
}
```

### Rate Limiting

```go
func checkRateLimit(userID string, limit int64, window int) (bool, error) {
    key := "rate_limit:" + userID
    
    current, err := redisClient.Incr(key)
    if err != nil {
        return false, err
    }
    
    if current == 1 {
        // First request, set expiration
        err = redisClient.Expire(key, window)
        if err != nil {
            return false, err
        }
    }
    
    return current <= limit, nil
}
```

### Event Processing with Streams

```go
// Event producer
func publishEvent(eventType, userID string, data map[string]string) error {
    fields := make(map[string]string)
    fields["event_type"] = eventType
    fields["user_id"] = userID
    fields["timestamp"] = time.Now().Format(time.RFC3339)
    
    // Copy additional data
    for k, v := range data {
        fields[k] = v
    }
    
    _, err := redisClient.XAdd("events", client.StreamIDAutoGenerate, fields)
    return err
}

// Event consumer worker
func processEvents(consumerGroup, consumerName string) {
    // Create consumer group if it doesn't exist
    redisClient.XGroupCreate("events", consumerGroup, client.StreamIDLatest)
    
    for {
        streams := map[string]string{"events": ">"}
        opts := client.XReadGroupOptions{
            Count: 10,
            Block: 5000, // 5 second timeout
        }
        
        messages, err := redisClient.XReadGroupWithOptions(consumerGroup, consumerName, streams, opts)
        if err != nil {
            log.Printf("Error reading from stream: %v", err)
            continue
        }
        
        for _, msg := range messages {
            for _, entry := range msg.Entries {
                // Process the event
                if err := handleEvent(entry.Fields); err != nil {
                    log.Printf("Error processing event %s: %v", entry.ID, err)
                    continue
                }
                
                // Acknowledge successful processing
                redisClient.XAck("events", consumerGroup, entry.ID)
            }
        }
    }
}

func handleEvent(fields map[string]string) error {
    eventType := fields["event_type"]
    userID := fields["user_id"]
    
    switch eventType {
    case "user_login":
        log.Printf("User %s logged in", userID)
    case "purchase":
        log.Printf("User %s made a purchase: %s", userID, fields["product"])
    default:
        log.Printf("Unknown event type: %s", eventType)
    }
    
    return nil
}
```

### Location-Based Services

```go
// Store and find nearby locations
func findNearbyStores(userLat, userLon float64, radiusKM float64) ([]string, error) {
    // Add stores to geospatial index (typically done during setup)
    stores := []client.GeoMember{
        {Longitude: -122.4194, Latitude: 37.7749, Member: "Downtown Store"},
        {Longitude: -122.4094, Latitude: 37.7849, Member: "Marina Store"},
        {Longitude: -122.4294, Latitude: 37.7649, Member: "Mission Store"},
    }
    redisClient.GeoAdd("stores", stores)
    
    // Search for nearby stores
    searchOpts := client.GeoSearchOptions{
        FromLonLat: &client.GeoCoordinate{
            Longitude: userLon,
            Latitude:  userLat,
        },
        ByRadius: &client.GeoRadius{
            Radius: radiusKM,
            Unit:   "km",
        },
        WithDist: true,
        Count:    10,
    }
    
    locations, err := redisClient.GeoSearch("stores", searchOpts)
    if err != nil {
        return nil, err
    }
    
    var nearbyStores []string
    for _, location := range locations {
        distance := *location.Distance
        nearbyStores = append(nearbyStores, 
            fmt.Sprintf("%s (%.2f km)", location.Member, distance))
    }
    
    return nearbyStores, nil
}

// Delivery tracking
func trackDelivery(driverID string, lat, lon float64) error {
    driver := client.GeoMember{
        Longitude: lon,
        Latitude:  lat,
        Member:    driverID,
    }
    
    _, err := redisClient.GeoAdd("drivers", []client.GeoMember{driver})
    return err
}

func findNearestDriver(customerLat, customerLon float64) (string, float64, error) {
    searchOpts := client.GeoSearchOptions{
        FromLonLat: &client.GeoCoordinate{
            Longitude: customerLon,
            Latitude:  customerLat,
        },
        ByRadius: &client.GeoRadius{
            Radius: 50, // 50 km search radius
            Unit:   "km",
        },
        WithDist: true,
        Count:    1,
    }
    
    drivers, err := redisClient.GeoSearch("drivers", searchOpts)
    if err != nil {
        return "", 0, err
    }
    
    if len(drivers) == 0 {
        return "", 0, fmt.Errorf("no drivers found nearby")
    }
    
    nearest := drivers[0]
    return nearest.Member, *nearest.Distance, nil
}
```

This getting started guide covers the fundamentals of using LibRedis, including the powerful new Redis Streams and Geospatial features. For more advanced usage patterns and features, continue with the other documentation files.