package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/therealbill/libredis/client"
)

// User represents a user document
type User struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Age   int      `json:"age"`
	Tags  []string `json:"tags"`
}

func main() {
	// Connect to Redis
	redis, err := client.DialURL("tcp://localhost:6379/0")
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redis.Close()

	fmt.Println("=== Phase 4 JSON Operations Example ===")

	// Example 1: Store a JSON document
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		Tags:  []string{"developer", "golang", "redis"},
	}

	userJSON, _ := json.Marshal(user)
	result, err := redis.JSONSet("user:1", ".", string(userJSON))
	if err != nil {
		fmt.Printf("Error setting JSON: %v\n", err)
		return
	}
	fmt.Printf("Stored user document: %s\n", result)

	// Example 2: Retrieve the full document
	data, err := redis.JSONGet("user:1")
	if err != nil {
		fmt.Printf("Error getting JSON: %v\n", err)
		return
	}
	fmt.Printf("Retrieved document: %s\n", string(data))

	// Example 3: Get specific fields
	data, err = redis.JSONGet("user:1", &client.JSONGetOptions{
		Paths: []string{".name", ".email"},
	})
	if err != nil {
		fmt.Printf("Error getting specific fields: %v\n", err)
		return
	}
	fmt.Printf("Name and email: %s\n", string(data))

	// Example 4: Update a specific field
	result, err = redis.JSONSet("user:1", ".age", "31")
	if err != nil {
		fmt.Printf("Error updating age: %v\n", err)
		return
	}
	fmt.Printf("Updated age: %s\n", result)

	// Example 5: Increment a numeric field
	newAge, err := redis.JSONNumIncrBy("user:1", ".age", 1)
	if err != nil {
		fmt.Printf("Error incrementing age: %v\n", err)
		return
	}
	fmt.Printf("Age after increment: %.0f\n", newAge)

	// Example 6: Work with arrays
	arrayLength, err := redis.JSONArrAppend("user:1", ".tags", "json", "redis-stack")
	if err != nil {
		fmt.Printf("Error appending to tags: %v\n", err)
		return
	}
	fmt.Printf("Tags array length after append: %d\n", arrayLength)

	// Example 7: Find array index
	index, err := redis.JSONArrIndex("user:1", ".tags", "golang")
	if err != nil {
		fmt.Printf("Error finding index: %v\n", err)
		return
	}
	fmt.Printf("Index of 'golang' in tags: %d\n", index)

	// Example 8: Get array length
	length, err := redis.JSONArrLen("user:1", ".tags")
	if err != nil {
		fmt.Printf("Error getting array length: %v\n", err)
		return
	}
	fmt.Printf("Tags array length: %d\n", length)

	// Example 9: Work with string operations
	result, err = redis.JSONSet("user:1", ".bio", `"Software developer"`)
	if err != nil {
		fmt.Printf("Error setting bio: %v\n", err)
		return
	}

	newLength, err := redis.JSONStrAppend("user:1", ".bio", `" with 10+ years experience"`)
	if err != nil {
		fmt.Printf("Error appending to bio: %v\n", err)
		return
	}
	fmt.Printf("Bio length after append: %d\n", newLength)

	// Example 10: Get type of JSON value
	jsonType, err := redis.JSONType("user:1", ".tags")
	if err != nil {
		fmt.Printf("Error getting type: %v\n", err)
		return
	}
	fmt.Printf("Type of .tags field: %s\n", jsonType)

	// Example 11: Conditional operations
	result, err = redis.JSONSet("user:1", ".last_login", `"2023-12-07T10:30:00Z"`, &client.JSONSetOptions{NX: true})
	if err != nil {
		fmt.Printf("Error setting last_login: %v\n", err)
		return
	}
	fmt.Printf("Set last_login (NX): %s\n", result)

	// Try to set the same field again with NX (should fail)
	result, err = redis.JSONSet("user:1", ".last_login", `"2023-12-08T09:15:00Z"`, &client.JSONSetOptions{NX: true})
	if err != nil {
		fmt.Printf("Expected error for NX on existing field: %v\n", err)
	} else {
		fmt.Printf("Unexpected success with NX on existing field: %s\n", result)
	}

	// Example 12: Get final document
	data, err = redis.JSONGet("user:1", &client.JSONGetOptions{
		Indent:  "  ",
		NewLine: "\n",
	})
	if err != nil {
		fmt.Printf("Error getting final document: %v\n", err)
		return
	}
	fmt.Printf("Final document (formatted):\n%s\n", string(data))

	// Example 13: Clean up
	deleted, err := redis.JSONDel("user:1")
	if err != nil {
		fmt.Printf("Error deleting document: %v\n", err)
		return
	}
	fmt.Printf("Deleted document: %d items removed\n", deleted)

	fmt.Println("\n=== JSON Operations Example Complete ===")
}