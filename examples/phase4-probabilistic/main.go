package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/therealbill/libredis/client"
)

func main() {
	// Connect to Redis
	redis, err := client.DialURL("tcp://localhost:6379/0")
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redis.Close()

	fmt.Println("=== Phase 4 Probabilistic Data Structures Example ===")

	// Bloom Filter Examples
	fmt.Println("\n=== Bloom Filter Examples ===")

	// Example 1: Create a Bloom filter for tracking visited users
	result, err := redis.BFReserve("users:visited", 0.001, 10000) // 0.1% error rate, 10k capacity
	if err != nil {
		fmt.Printf("Error creating Bloom filter: %v\n", err)
		return
	}
	fmt.Printf("Created Bloom filter: %s\n", result)

	// Example 2: Add users to the filter
	users := []string{"user:123", "user:456", "user:789", "user:101", "user:202"}
	fmt.Println("\n--- Adding users to Bloom filter ---")
	for _, user := range users {
		added, err := redis.BFAdd("users:visited", user)
		if err != nil {
			fmt.Printf("Error adding %s: %v\n", user, err)
		} else {
			fmt.Printf("Added %s: %t (true = new, false = might exist)\n", user, added)
		}
	}

	// Example 3: Add multiple users at once
	fmt.Println("\n--- Adding multiple users at once ---")
	newUsers := []interface{}{"user:303", "user:404", "user:123"} // user:123 already exists
	results, err := redis.BFMAdd("users:visited", newUsers...)
	if err != nil {
		fmt.Printf("Error with bulk add: %v\n", err)
	} else {
		for i, added := range results {
			fmt.Printf("User %v: %t\n", newUsers[i], added)
		}
	}

	// Example 4: Check if users exist
	fmt.Println("\n--- Checking user existence ---")
	testUsers := []interface{}{"user:123", "user:999", "user:456", "user:888"}
	existResults, err := redis.BFMExists("users:visited", testUsers...)
	if err != nil {
		fmt.Printf("Error checking existence: %v\n", err)
	} else {
		for i, exists := range existResults {
			status := "might exist"
			if !exists {
				status = "definitely not in set"
			}
			fmt.Printf("User %v: %s\n", testUsers[i], status)
		}
	}

	// Cuckoo Filter Examples
	fmt.Println("\n=== Cuckoo Filter Examples ===")

	// Example 5: Create a Cuckoo filter (supports deletions)
	result, err = redis.CFReserve("active:sessions", 5000)
	if err != nil {
		fmt.Printf("Error creating Cuckoo filter: %v\n", err)
	} else {
		fmt.Printf("Created Cuckoo filter: %s\n", result)
	}

	// Example 6: Add and manage active sessions
	sessions := []string{"session:abc123", "session:def456", "session:ghi789"}
	fmt.Println("\n--- Managing active sessions ---")
	
	for _, session := range sessions {
		added, err := redis.CFAdd("active:sessions", session)
		if err != nil {
			fmt.Printf("Error adding session %s: %v\n", session, err)
		} else {
			fmt.Printf("Added session %s: %t\n", session, added)
		}
	}

	// Example 7: Check session existence
	exists, err := redis.CFExists("active:sessions", "session:abc123")
	if err != nil {
		fmt.Printf("Error checking session: %v\n", err)
	} else {
		fmt.Printf("Session abc123 active: %t\n", exists)
	}

	// Example 8: Delete a session (unique to Cuckoo filters)
	deleted, err := redis.CFDel("active:sessions", "session:def456")
	if err != nil {
		fmt.Printf("Error deleting session: %v\n", err)
	} else {
		fmt.Printf("Deleted session def456: %t\n", deleted)
	}

	// Verify deletion
	exists, err = redis.CFExists("active:sessions", "session:def456")
	if err != nil {
		fmt.Printf("Error checking deleted session: %v\n", err)
	} else {
		fmt.Printf("Deleted session still exists: %t (should be false)\n", exists)
	}

	// Count-Min Sketch Examples
	fmt.Println("\n=== Count-Min Sketch Examples ===")

	// Example 9: Create a Count-Min Sketch for page view counting
	result, err = redis.CMSInitByDim("page:views", 2000, 5) // 2000 width, 5 depth
	if err != nil {
		fmt.Printf("Error creating Count-Min Sketch: %v\n", err)
	} else {
		fmt.Printf("Created Count-Min Sketch: %s\n", result)
	}

	// Example 10: Simulate page views
	fmt.Println("\n--- Simulating page views ---")
	pages := []string{"/home", "/about", "/contact", "/products", "/blog"}
	rand.Seed(time.Now().UnixNano())

	// Generate random page views
	viewCounts := make(map[string]int)
	for i := 0; i < 100; i++ {
		page := pages[rand.Intn(len(pages))]
		viewCounts[page]++
		
		_, err := redis.CMSIncrBy("page:views", page, 1)
		if err != nil {
			fmt.Printf("Error incrementing %s: %v\n", page, err)
		}
	}

	fmt.Println("Generated 100 random page views")
	fmt.Println("Actual counts:")
	for page, count := range viewCounts {
		fmt.Printf("  %s: %d views\n", page, count)
	}

	// Example 11: Query page view counts
	fmt.Println("\n--- Querying estimated counts ---")
	queryPages := []interface{}{"/home", "/about", "/contact", "/products", "/blog", "/nonexistent"}
	counts, err := redis.CMSQuery("page:views", queryPages...)
	if err != nil {
		fmt.Printf("Error querying counts: %v\n", err)
	} else {
		fmt.Println("Estimated counts from Count-Min Sketch:")
		for i, count := range counts {
			page := queryPages[i].(string)
			actual := viewCounts[page]
			fmt.Printf("  %s: %d (actual: %d, error: %d)\n", page, count, actual, int64(actual)-count)
		}
	}

	// Example 12: Bulk increment operation
	fmt.Println("\n--- Bulk increment operations ---")
	increments := []interface{}{"/api/users", 5, "/api/products", 3, "/api/orders", 8}
	results, err = redis.CMSIncrBy("page:views", increments...)
	if err != nil {
		fmt.Printf("Error with bulk increment: %v\n", err)
	} else {
		fmt.Printf("Bulk increment results: %v\n", results)
	}

	// Example 13: Create another CMS and merge
	fmt.Println("\n--- CMS merge example ---")
	result, err = redis.CMSInitByDim("page:views:backup", 2000, 5)
	if err != nil {
		fmt.Printf("Error creating backup CMS: %v\n", err)
	} else {
		// Add some data to backup
		redis.CMSIncrBy("page:views:backup", "/admin", 10, "/login", 25)
		
		// Merge backup into main CMS
		result, err = redis.CMSMerge("page:views:merged", []string{"page:views", "page:views:backup"})
		if err != nil {
			fmt.Printf("Error merging CMS: %v\n", err)
		} else {
			fmt.Printf("Merged CMS: %s\n", result)
			
			// Query merged data
			mergedCounts, err := redis.CMSQuery("page:views:merged", "/home", "/admin", "/login")
			if err != nil {
				fmt.Printf("Error querying merged data: %v\n", err)
			} else {
				fmt.Printf("Merged counts - /home: %d, /admin: %d, /login: %d\n", 
					mergedCounts[0], mergedCounts[1], mergedCounts[2])
			}
		}
	}

	// Example 14: Get information about data structures
	fmt.Println("\n--- Data structure information ---")
	
	// Bloom Filter info
	bfInfo, err := redis.BFInfo("users:visited")
	if err != nil {
		fmt.Printf("Error getting BF info: %v\n", err)
	} else {
		fmt.Printf("Bloom Filter info:\n")
		for key, value := range bfInfo {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	// Cuckoo Filter info
	cfInfo, err := redis.CFInfo("active:sessions")
	if err != nil {
		fmt.Printf("Error getting CF info: %v\n", err)
	} else {
		fmt.Printf("Cuckoo Filter info:\n")
		for key, value := range cfInfo {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	// Count-Min Sketch info
	cmsInfo, err := redis.CMSInfo("page:views")
	if err != nil {
		fmt.Printf("Error getting CMS info: %v\n", err)
	} else {
		fmt.Printf("Count-Min Sketch info:\n")
		for key, value := range cmsInfo {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	// Example 15: Clean up
	fmt.Println("\n--- Cleanup ---")
	cleanupKeys := []string{
		"users:visited", 
		"active:sessions", 
		"page:views", 
		"page:views:backup", 
		"page:views:merged",
	}
	
	for _, key := range cleanupKeys {
		deleted, err := redis.Del(key)
		if err != nil {
			fmt.Printf("Error deleting %s: %v\n", key, err)
		} else {
			fmt.Printf("Deleted %s: %d\n", key, deleted)
		}
	}

	fmt.Println("\n=== Probabilistic Data Structures Example Complete ===")
	fmt.Println("\nKey takeaways:")
	fmt.Println("- Bloom Filters: Memory-efficient set membership (no deletions)")
	fmt.Println("- Cuckoo Filters: Set membership with deletion support")
	fmt.Println("- Count-Min Sketch: Frequency counting with bounded error")
	fmt.Println("- All structures trade memory for accuracy but provide huge space savings")
}