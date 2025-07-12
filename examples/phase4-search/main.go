package main

import (
	"fmt"
	"log"

	"github.com/therealbill/libredis/client"
)

func main() {
	// Connect to Redis
	redis, err := client.DialURL("tcp://localhost:6379/0")
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redis.Close()

	fmt.Println("=== Phase 4 Search Operations Example ===")

	// Example 1: Create a search index
	schema := []client.FTFieldSchema{
		{Name: "title", Type: "TEXT", Weight: 1.0, Sortable: true},
		{Name: "description", Type: "TEXT"},
		{Name: "price", Type: "NUMERIC", Sortable: true},
		{Name: "category", Type: "TAG", Separator: ","},
		{Name: "location", Type: "GEO"},
		{Name: "rating", Type: "NUMERIC"},
	}

	result, err := redis.FTCreate("products", schema, &client.FTCreateOptions{
		OnHash: true,
		Prefix: []string{"product:"},
	})
	if err != nil {
		fmt.Printf("Error creating index (may already exist): %v\n", err)
	} else {
		fmt.Printf("Created search index: %s\n", result)
	}

	// Example 2: Add some sample products using Redis hashes
	products := map[string]map[string]interface{}{
		"product:1": {
			"title":       "Wireless Bluetooth Headphones",
			"description": "High-quality wireless headphones with noise cancellation",
			"price":       "199.99",
			"category":    "electronics,audio",
			"location":    "-122.4194,37.7749", // San Francisco
			"rating":      "4.5",
		},
		"product:2": {
			"title":       "Smartphone with 5G",
			"description": "Latest smartphone with 5G connectivity and advanced camera",
			"price":       "899.99",
			"category":    "electronics,mobile",
			"location":    "-74.0060,40.7128", // New York
			"rating":      "4.8",
		},
		"product:3": {
			"title":       "Coffee Maker",
			"description": "Automatic coffee maker with programmable settings",
			"price":       "79.99",
			"category":    "appliances,kitchen",
			"location":    "-87.6298,41.8781", // Chicago
			"rating":      "4.2",
		},
		"product:4": {
			"title":       "Gaming Laptop",
			"description": "High-performance laptop for gaming and professional work",
			"price":       "1299.99",
			"category":    "electronics,computers",
			"location":    "-118.2437,34.0522", // Los Angeles
			"rating":      "4.7",
		},
	}

	for key, fields := range products {
		for field, value := range fields {
			_, err := redis.HSet(key, field, value)
			if err != nil {
				fmt.Printf("Error setting field %s for %s: %v\n", field, key, err)
			}
		}
		fmt.Printf("Added product: %s\n", key)
	}

	// Example 3: Basic text search
	fmt.Println("\n--- Basic Text Search ---")
	results, err := redis.FTSearch("products", "wireless")
	if err != nil {
		fmt.Printf("Error searching: %v\n", err)
	} else {
		fmt.Printf("Search results for 'wireless': %d total\n", len(results))
		if len(results) > 1 {
			fmt.Printf("First result: %v\n", results[1])
		}
	}

	// Example 4: Search with numeric filter
	fmt.Println("\n--- Search with Price Filter ---")
	results, err = redis.FTSearch("products", "*", &client.FTSearchOptions{
		Filter: []client.FTNumericFilter{
			{Field: "price", Min: 100, Max: 500},
		},
		WithScores: true,
		Limit:      &client.FTLimit{Offset: 0, Num: 10},
	})
	if err != nil {
		fmt.Printf("Error searching with filter: %v\n", err)
	} else {
		fmt.Printf("Products priced $100-500: %d results\n", len(results))
		for i, result := range results {
			if i == 0 {
				continue // Skip count
			}
			fmt.Printf("Result %d: %v\n", i, result)
		}
	}

	// Example 5: Search by category (tag)
	fmt.Println("\n--- Search by Category ---")
	results, err = redis.FTSearch("products", "@category:{electronics}")
	if err != nil {
		fmt.Printf("Error searching by category: %v\n", err)
	} else {
		fmt.Printf("Electronics products: %d results\n", len(results))
	}

	// Example 6: Geographic search
	fmt.Println("\n--- Geographic Search ---")
	results, err = redis.FTSearch("products", "*", &client.FTSearchOptions{
		GeoFilter: &client.FTGeoFilter{
			Field:     "location",
			Longitude: -122.4194, // San Francisco coordinates
			Latitude:  37.7749,
			Radius:    500,       // 500 km radius
			Unit:      "km",
		},
		Limit: &client.FTLimit{Offset: 0, Num: 10},
	})
	if err != nil {
		fmt.Printf("Error with geographic search: %v\n", err)
	} else {
		fmt.Printf("Products within 500km of San Francisco: %d results\n", len(results))
	}

	// Example 7: Complex search query
	fmt.Println("\n--- Complex Search Query ---")
	results, err = redis.FTSearch("products", "smartphone OR laptop", &client.FTSearchOptions{
		Filter: []client.FTNumericFilter{
			{Field: "rating", Min: 4.0, Max: 5.0},
		},
		SortBy:    "price",
		SortOrder: "ASC",
		Limit:     &client.FTLimit{Offset: 0, Num: 5},
	})
	if err != nil {
		fmt.Printf("Error with complex search: %v\n", err)
	} else {
		fmt.Printf("High-rated smartphones or laptops (sorted by price): %d results\n", len(results))
	}

	// Example 8: Aggregation query
	fmt.Println("\n--- Aggregation Query ---")
	results, err = redis.FTAggregate("products", "*", &client.FTAggregateOptions{
		GroupBy: &client.FTGroupBy{
			Fields: []string{"@category"},
			Reduce: []client.FTReduce{
				{Function: "COUNT", Args: []string{}, As: "count"},
				{Function: "AVG", Args: []string{"@price"}, As: "avg_price"},
			},
		},
	})
	if err != nil {
		fmt.Printf("Error with aggregation: %v\n", err)
	} else {
		fmt.Printf("Category aggregation results: %d groups\n", len(results))
		for i, result := range results {
			fmt.Printf("Group %d: %v\n", i, result)
		}
	}

	// Example 9: Get index information
	fmt.Println("\n--- Index Information ---")
	info, err := redis.FTInfo("products")
	if err != nil {
		fmt.Printf("Error getting index info: %v\n", err)
	} else {
		fmt.Printf("Index information:\n")
		for key, value := range info {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	// Example 10: Search with highlighting
	fmt.Println("\n--- Search with Highlighting ---")
	results, err = redis.FTSearch("products", "gaming", &client.FTSearchOptions{
		Highlight: &client.FTHighlightOptions{
			Fields: []string{"title", "description"},
			Tags: &client.FTHighlightTags{
				Open:  "<mark>",
				Close: "</mark>",
			},
		},
		Return: []string{"title", "description"},
	})
	if err != nil {
		fmt.Printf("Error with highlighting: %v\n", err)
	} else {
		fmt.Printf("Search with highlighting: %d results\n", len(results))
		for i, result := range results {
			if i == 0 {
				continue // Skip count
			}
			fmt.Printf("Highlighted result: %v\n", result)
		}
	}

	// Example 11: Explain query
	fmt.Println("\n--- Query Explanation ---")
	explanation, err := redis.FTExplain("products", "wireless OR bluetooth")
	if err != nil {
		fmt.Printf("Error explaining query: %v\n", err)
	} else {
		fmt.Printf("Query explanation:\n%s\n", explanation)
	}

	// Example 12: Clean up
	fmt.Println("\n--- Cleanup ---")
	
	// Delete the documents
	for key := range products {
		deleted, err := redis.Del(key)
		if err != nil {
			fmt.Printf("Error deleting %s: %v\n", key, err)
		} else {
			fmt.Printf("Deleted %s: %d\n", key, deleted)
		}
	}

	// Drop the index
	result, err = redis.FTDropIndex("products", true)
	if err != nil {
		fmt.Printf("Error dropping index: %v\n", err)
	} else {
		fmt.Printf("Dropped index: %s\n", result)
	}

	fmt.Println("\n=== Search Operations Example Complete ===")
}