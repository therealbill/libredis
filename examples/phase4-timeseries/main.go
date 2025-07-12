package main

import (
	"fmt"
	"log"
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

	fmt.Println("=== Phase 4 Time Series Operations Example ===")

	// Example 1: Create time series for temperature sensors
	labels := map[string]string{
		"sensor_type": "temperature",
		"location":    "living_room",
		"unit":        "celsius",
	}

	result, err := redis.TSCreate("temp:living_room", &client.TSCreateOptions{
		RetentionMsecs:  3600000, // 1 hour retention
		ChunkSize:       4096,
		DuplicatePolicy: "LAST",
		Labels:          labels,
	})
	if err != nil {
		fmt.Printf("Error creating time series: %v\n", err)
		return
	}
	fmt.Printf("Created time series: %s\n", result)

	// Create another time series for humidity
	humidityLabels := map[string]string{
		"sensor_type": "humidity",
		"location":    "living_room",
		"unit":        "percent",
	}

	result, err = redis.TSCreate("humidity:living_room", &client.TSCreateOptions{
		RetentionMsecs:  3600000,
		ChunkSize:       4096,
		DuplicatePolicy: "LAST",
		Labels:          humidityLabels,
	})
	if err != nil {
		fmt.Printf("Error creating humidity series: %v\n", err)
		return
	}
	fmt.Printf("Created humidity series: %s\n", result)

	// Example 2: Add current temperature reading
	now := time.Now().UnixMilli()
	timestamp, err := redis.TSAdd("temp:living_room", now, 23.5)
	if err != nil {
		fmt.Printf("Error adding temperature sample: %v\n", err)
		return
	}
	fmt.Printf("Added temperature sample at timestamp: %d\n", timestamp)

	// Example 3: Add multiple samples with auto timestamp
	fmt.Println("\n--- Adding multiple samples ---")
	temperatures := []float64{23.8, 24.1, 24.0, 23.9, 24.2}
	for i, temp := range temperatures {
		timestamp, err := redis.TSAdd("temp:living_room", 0, temp) // 0 = auto timestamp
		if err != nil {
			fmt.Printf("Error adding sample %d: %v\n", i, err)
		} else {
			fmt.Printf("Added temperature %.1f°C at %d\n", temp, timestamp)
		}
		time.Sleep(100 * time.Millisecond) // Small delay between samples
	}

	// Example 4: Add multiple samples at once
	fmt.Println("\n--- Adding multiple samples at once ---")
	baseTime := time.Now().UnixMilli()
	samples := []client.TSMAddSample{
		{Key: "temp:living_room", Timestamp: baseTime + 1000, Value: 24.5},
		{Key: "humidity:living_room", Timestamp: baseTime + 1000, Value: 65.0},
		{Key: "temp:living_room", Timestamp: baseTime + 2000, Value: 24.3},
		{Key: "humidity:living_room", Timestamp: baseTime + 2000, Value: 66.5},
	}

	timestamps, err := redis.TSMAdd(samples...)
	if err != nil {
		fmt.Printf("Error adding multiple samples: %v\n", err)
	} else {
		fmt.Printf("Added %d samples successfully\n", len(timestamps))
		for i, ts := range timestamps {
			fmt.Printf("  Sample %d timestamp: %d\n", i+1, ts)
		}
	}

	// Example 5: Increment/decrement operations
	fmt.Println("\n--- Increment/decrement operations ---")
	newTimestamp, err := redis.TSIncrBy("temp:living_room", 0.5)
	if err != nil {
		fmt.Printf("Error incrementing: %v\n", err)
	} else {
		fmt.Printf("Incremented temperature by 0.5°C at %d\n", newTimestamp)
	}

	newTimestamp, err = redis.TSDecrBy("humidity:living_room", 2.0)
	if err != nil {
		fmt.Printf("Error decrementing: %v\n", err)
	} else {
		fmt.Printf("Decremented humidity by 2%% at %d\n", newTimestamp)
	}

	// Example 6: Query time series data
	fmt.Println("\n--- Querying time series data ---")
	endTime := time.Now().UnixMilli()
	startTime := endTime - (10 * 60 * 1000) // Last 10 minutes

	samples, err := redis.TSRange("temp:living_room", startTime, endTime)
	if err != nil {
		fmt.Printf("Error querying temperature data: %v\n", err)
	} else {
		fmt.Printf("Temperature samples in last 10 minutes: %d\n", len(samples))
		for i, sample := range samples {
			if i < 5 { // Show first 5 samples
				t := time.UnixMilli(sample.Timestamp)
				fmt.Printf("  %s: %.1f°C\n", t.Format("15:04:05"), sample.Value)
			}
		}
		if len(samples) > 5 {
			fmt.Printf("  ... and %d more samples\n", len(samples)-5)
		}
	}

	// Example 7: Query with aggregation
	fmt.Println("\n--- Querying with aggregation ---")
	samples, err = redis.TSRange("temp:living_room", startTime, endTime, &client.TSRangeOptions{
		Aggregation: &client.TSAggregation{
			Type:       "avg",
			TimeBucket: 60000, // 1 minute buckets
		},
	})
	if err != nil {
		fmt.Printf("Error querying aggregated data: %v\n", err)
	} else {
		fmt.Printf("Average temperature per minute: %d data points\n", len(samples))
		for _, sample := range samples {
			t := time.UnixMilli(sample.Timestamp)
			fmt.Printf("  %s: %.2f°C (avg)\n", t.Format("15:04"), sample.Value)
		}
	}

	// Example 8: Query multiple time series
	fmt.Println("\n--- Querying multiple time series ---")
	filters := []string{"location=living_room"}
	results, err := redis.TSMRange(startTime, endTime, filters, &client.TSMRangeOptions{
		WithLabels: true,
		Count:      5, // Limit to 5 samples per series
	})
	if err != nil {
		fmt.Printf("Error querying multiple series: %v\n", err)
	} else {
		fmt.Printf("Multiple series query results: %d series\n", len(results))
		for key, samples := range results {
			fmt.Printf("  %s: %d samples\n", key, len(samples))
			if len(samples) > 0 {
				latest := samples[len(samples)-1]
				t := time.UnixMilli(latest.Timestamp)
				fmt.Printf("    Latest: %s = %.2f\n", t.Format("15:04:05"), latest.Value)
			}
		}
	}

	// Example 9: Reverse range query
	fmt.Println("\n--- Reverse range query ---")
	samples, err = redis.TSRevRange("temp:living_room", startTime, endTime, &client.TSRangeOptions{
		Count: 3, // Get last 3 samples
	})
	if err != nil {
		fmt.Printf("Error with reverse range: %v\n", err)
	} else {
		fmt.Printf("Last 3 temperature readings (reverse order):\n")
		for _, sample := range samples {
			t := time.UnixMilli(sample.Timestamp)
			fmt.Printf("  %s: %.1f°C\n", t.Format("15:04:05"), sample.Value)
		}
	}

	// Example 10: Get time series information
	fmt.Println("\n--- Time series information ---")
	info, err := redis.TSInfo("temp:living_room")
	if err != nil {
		fmt.Printf("Error getting time series info: %v\n", err)
	} else {
		fmt.Printf("Temperature series info:\n")
		fmt.Printf("  Total samples: %d\n", info.TotalSamples)
		fmt.Printf("  Memory usage: %d bytes\n", info.MemoryUsage)
		fmt.Printf("  Retention time: %d ms\n", info.RetentionTime)
		fmt.Printf("  Chunk count: %d\n", info.ChunkCount)
		fmt.Printf("  Duplicate policy: %s\n", info.DuplicatePolicy)
		fmt.Printf("  Labels:\n")
		for label, value := range info.Labels {
			fmt.Printf("    %s: %s\n", label, value)
		}
		if info.FirstTimestamp > 0 {
			firstTime := time.UnixMilli(info.FirstTimestamp)
			lastTime := time.UnixMilli(info.LastTimestamp)
			fmt.Printf("  Time range: %s to %s\n", 
				firstTime.Format("15:04:05"), lastTime.Format("15:04:05"))
		}
	}

	// Example 11: Query with value filtering
	fmt.Println("\n--- Query with value filtering ---")
	samples, err = redis.TSRange("temp:living_room", startTime, endTime, &client.TSRangeOptions{
		FilterBy: &client.TSFilterBy{
			Min: 24.0,
			Max: 25.0,
		},
	})
	if err != nil {
		fmt.Printf("Error with filtered query: %v\n", err)
	} else {
		fmt.Printf("Temperature readings between 24-25°C: %d samples\n", len(samples))
		for _, sample := range samples {
			t := time.UnixMilli(sample.Timestamp)
			fmt.Printf("  %s: %.1f°C\n", t.Format("15:04:05"), sample.Value)
		}
	}

	// Example 12: Clean up
	fmt.Println("\n--- Cleanup ---")
	deleted, err := redis.Del("temp:living_room")
	if err != nil {
		fmt.Printf("Error deleting temperature series: %v\n", err)
	} else {
		fmt.Printf("Deleted temperature series: %d\n", deleted)
	}

	deleted, err = redis.Del("humidity:living_room")
	if err != nil {
		fmt.Printf("Error deleting humidity series: %v\n", err)
	} else {
		fmt.Printf("Deleted humidity series: %d\n", deleted)
	}

	fmt.Println("\n=== Time Series Operations Example Complete ===")
}