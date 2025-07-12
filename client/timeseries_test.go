// +build integration

package client

import (
	"testing"
	"time"
)

// Helper function to check if RedisTimeSeries module is available
func isTimeSeriesModuleAvailable(t *testing.T) bool {
	// Try to create a simple time series - if it fails, module is not available
	_, err := r.TSCreate("test_ts_module")
	if err != nil {
		t.Skip("RedisTimeSeries module not available, skipping TimeSeries tests")
		return false
	}
	r.Del("test_ts_module") // Clean up
	return true
}

func TestTSCreate(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Test basic time series creation
	result, err := r.TSCreate("test_ts")
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Test time series creation with options
	labels := map[string]string{
		"sensor":   "temperature",
		"location": "living_room",
	}
	result, err = r.TSCreate("test_ts_with_options", &TSCreateOptions{
		RetentionMsecs:  3600000, // 1 hour
		ChunkSize:       4096,
		DuplicatePolicy: "LAST",
		Labels:          labels,
	})
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Clean up
	r.Del("test_ts", "test_ts_with_options")
}

func TestTSAdd(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Create time series
	r.TSCreate("test_ts_add")

	// Test adding a sample
	now := time.Now().UnixMilli()
	timestamp, err := r.TSAdd("test_ts_add", now, 25.5)
	if err != nil {
		t.Error(err)
	}
	if timestamp != now {
		t.Errorf("Expected timestamp %d, got %d", now, timestamp)
	}

	// Test adding with auto timestamp (0)
	timestamp, err = r.TSAdd("test_ts_add", 0, 26.0)
	if err != nil {
		t.Error(err)
	}
	if timestamp == 0 {
		t.Error("Expected non-zero auto-generated timestamp")
	}

	// Clean up
	r.Del("test_ts_add")
}

func TestTSMAdd(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Create multiple time series
	r.TSCreate("test_ts_madd_1")
	r.TSCreate("test_ts_madd_2")

	// Test adding multiple samples
	now := time.Now().UnixMilli()
	samples := []TSMAddSample{
		{Key: "test_ts_madd_1", Timestamp: now, Value: 10.0},
		{Key: "test_ts_madd_2", Timestamp: now, Value: 20.0},
		{Key: "test_ts_madd_1", Timestamp: now + 1000, Value: 11.0},
		{Key: "test_ts_madd_2", Timestamp: now + 1000, Value: 21.0},
	}

	timestamps, err := r.TSMAdd(samples...)
	if err != nil {
		t.Error(err)
	}
	if len(timestamps) != 4 {
		t.Errorf("Expected 4 timestamps, got %d", len(timestamps))
	}

	// Clean up
	r.Del("test_ts_madd_1", "test_ts_madd_2")
}

func TestTSIncrBy(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Create time series
	r.TSCreate("test_ts_incr")

	// Test incrementing value
	timestamp, err := r.TSIncrBy("test_ts_incr", 5.0)
	if err != nil {
		t.Error(err)
	}
	if timestamp == 0 {
		t.Error("Expected non-zero timestamp")
	}

	// Test incrementing with explicit timestamp
	now := time.Now().UnixMilli()
	timestamp, err = r.TSIncrBy("test_ts_incr", 3.0, &TSIncrByOptions{
		Timestamp: now,
	})
	if err != nil {
		t.Error(err)
	}
	if timestamp != now {
		t.Errorf("Expected timestamp %d, got %d", now, timestamp)
	}

	// Clean up
	r.Del("test_ts_incr")
}

func TestTSDecrBy(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Create time series
	r.TSCreate("test_ts_decr")

	// Test decrementing value
	timestamp, err := r.TSDecrBy("test_ts_decr", 2.0)
	if err != nil {
		t.Error(err)
	}
	if timestamp == 0 {
		t.Error("Expected non-zero timestamp")
	}

	// Clean up
	r.Del("test_ts_decr")
}

func TestTSRange(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Create time series and add samples
	r.TSCreate("test_ts_range")
	now := time.Now().UnixMilli()
	
	// Add several samples
	r.TSAdd("test_ts_range", now, 10.0)
	r.TSAdd("test_ts_range", now+1000, 20.0)
	r.TSAdd("test_ts_range", now+2000, 30.0)
	r.TSAdd("test_ts_range", now+3000, 40.0)
	r.TSAdd("test_ts_range", now+4000, 50.0)

	// Test range query
	samples, err := r.TSRange("test_ts_range", now, now+4000)
	if err != nil {
		t.Error(err)
	}
	if len(samples) != 5 {
		t.Errorf("Expected 5 samples, got %d", len(samples))
	}

	// Verify first sample
	if samples[0].Timestamp != now || samples[0].Value != 10.0 {
		t.Errorf("Expected first sample: {%d, 10.0}, got {%d, %f}", now, samples[0].Timestamp, samples[0].Value)
	}

	// Test range query with count limit
	samples, err = r.TSRange("test_ts_range", now, now+4000, &TSRangeOptions{
		Count: 3,
	})
	if err != nil {
		t.Error(err)
	}
	if len(samples) != 3 {
		t.Errorf("Expected 3 samples with count limit, got %d", len(samples))
	}

	// Test range query with aggregation
	samples, err = r.TSRange("test_ts_range", now, now+4000, &TSRangeOptions{
		Aggregation: &TSAggregation{
			Type:       "avg",
			TimeBucket: 2000,
		},
	})
	if err != nil {
		t.Error(err)
	}
	if len(samples) == 0 {
		t.Error("Expected aggregated samples")
	}

	// Clean up
	r.Del("test_ts_range")
}

func TestTSRevRange(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Create time series and add samples
	r.TSCreate("test_ts_revrange")
	now := time.Now().UnixMilli()
	
	// Add several samples
	r.TSAdd("test_ts_revrange", now, 10.0)
	r.TSAdd("test_ts_revrange", now+1000, 20.0)
	r.TSAdd("test_ts_revrange", now+2000, 30.0)

	// Test reverse range query
	samples, err := r.TSRevRange("test_ts_revrange", now, now+2000)
	if err != nil {
		t.Error(err)
	}
	if len(samples) != 3 {
		t.Errorf("Expected 3 samples, got %d", len(samples))
	}

	// Verify samples are in reverse order (latest first)
	if samples[0].Value != 30.0 {
		t.Errorf("Expected first sample value 30.0, got %f", samples[0].Value)
	}
	if samples[2].Value != 10.0 {
		t.Errorf("Expected last sample value 10.0, got %f", samples[2].Value)
	}

	// Clean up
	r.Del("test_ts_revrange")
}

func TestTSMRange(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Create multiple time series with labels
	labels1 := map[string]string{"sensor": "temperature", "location": "room1"}
	labels2 := map[string]string{"sensor": "temperature", "location": "room2"}
	
	r.TSCreate("test_ts_mrange_1", &TSCreateOptions{Labels: labels1})
	r.TSCreate("test_ts_mrange_2", &TSCreateOptions{Labels: labels2})

	now := time.Now().UnixMilli()
	
	// Add samples to both series
	r.TSAdd("test_ts_mrange_1", now, 25.0)
	r.TSAdd("test_ts_mrange_1", now+1000, 26.0)
	r.TSAdd("test_ts_mrange_2", now, 22.0)
	r.TSAdd("test_ts_mrange_2", now+1000, 23.0)

	// Test multi-range query
	filters := []string{"sensor=temperature"}
	results, err := r.TSMRange(now, now+1000, filters)
	if err != nil {
		t.Error(err)
	}

	if len(results) == 0 {
		t.Error("Expected results from multi-range query")
	}

	// Verify we got data for both series
	expectedKeys := []string{"test_ts_mrange_1", "test_ts_mrange_2"}
	for _, key := range expectedKeys {
		if samples, exists := results[key]; !exists || len(samples) == 0 {
			t.Errorf("Expected samples for key %s", key)
		}
	}

	// Clean up
	r.Del("test_ts_mrange_1", "test_ts_mrange_2")
}

func TestTSMRevRange(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Create multiple time series with labels
	labels := map[string]string{"sensor": "humidity"}
	
	r.TSCreate("test_ts_mrevrange_1", &TSCreateOptions{Labels: labels})
	r.TSCreate("test_ts_mrevrange_2", &TSCreateOptions{Labels: labels})

	now := time.Now().UnixMilli()
	
	// Add samples
	r.TSAdd("test_ts_mrevrange_1", now, 60.0)
	r.TSAdd("test_ts_mrevrange_1", now+1000, 65.0)
	r.TSAdd("test_ts_mrevrange_2", now, 55.0)
	r.TSAdd("test_ts_mrevrange_2", now+1000, 58.0)

	// Test multi-reverse range query
	filters := []string{"sensor=humidity"}
	results, err := r.TSMRevRange(now, now+1000, filters)
	if err != nil {
		t.Error(err)
	}

	if len(results) == 0 {
		t.Error("Expected results from multi-reverse range query")
	}

	// Verify samples are in reverse order
	for key, samples := range results {
		if len(samples) >= 2 {
			if samples[0].Timestamp <= samples[1].Timestamp {
				t.Errorf("Expected reverse order for key %s", key)
			}
		}
	}

	// Clean up
	r.Del("test_ts_mrevrange_1", "test_ts_mrevrange_2")
}

func TestTSInfo(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Create time series with labels
	labels := map[string]string{
		"sensor":   "temperature",
		"location": "kitchen",
	}
	r.TSCreate("test_ts_info", &TSCreateOptions{
		RetentionMsecs:  3600000,
		ChunkSize:       4096,
		DuplicatePolicy: "LAST",
		Labels:          labels,
	})

	// Add some samples
	now := time.Now().UnixMilli()
	r.TSAdd("test_ts_info", now, 23.5)
	r.TSAdd("test_ts_info", now+1000, 24.0)

	// Test getting info
	info, err := r.TSInfo("test_ts_info")
	if err != nil {
		t.Error(err)
	}

	// Verify some basic info
	if info.TotalSamples == 0 {
		t.Error("Expected non-zero total samples")
	}
	if info.RetentionTime != 3600000 {
		t.Errorf("Expected retention time 3600000, got %d", info.RetentionTime)
	}
	if info.ChunkSize != 4096 {
		t.Errorf("Expected chunk size 4096, got %d", info.ChunkSize)
	}
	if info.DuplicatePolicy != "LAST" {
		t.Errorf("Expected duplicate policy LAST, got %s", info.DuplicatePolicy)
	}

	// Verify labels
	if len(info.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(info.Labels))
	}
	if info.Labels["sensor"] != "temperature" {
		t.Errorf("Expected sensor=temperature, got %s", info.Labels["sensor"])
	}
	if info.Labels["location"] != "kitchen" {
		t.Errorf("Expected location=kitchen, got %s", info.Labels["location"])
	}

	// Clean up
	r.Del("test_ts_info")
}

func TestTSWithFilters(t *testing.T) {
	if !isTimeSeriesModuleAvailable(t) {
		return
	}

	// Create time series and add samples with different values
	r.TSCreate("test_ts_filter")
	now := time.Now().UnixMilli()
	
	// Add samples with various values
	r.TSAdd("test_ts_filter", now, 5.0)
	r.TSAdd("test_ts_filter", now+1000, 15.0)
	r.TSAdd("test_ts_filter", now+2000, 25.0)
	r.TSAdd("test_ts_filter", now+3000, 35.0)
	r.TSAdd("test_ts_filter", now+4000, 45.0)

	// Test range query with value filter
	samples, err := r.TSRange("test_ts_filter", now, now+4000, &TSRangeOptions{
		FilterBy: &TSFilterBy{
			Min: 20.0,
			Max: 40.0,
		},
	})
	if err != nil {
		t.Error(err)
	}

	// Should only return samples with values between 20.0 and 40.0
	expectedCount := 2 // 25.0 and 35.0
	if len(samples) != expectedCount {
		t.Errorf("Expected %d filtered samples, got %d", expectedCount, len(samples))
	}

	// Verify filtered values
	for _, sample := range samples {
		if sample.Value < 20.0 || sample.Value > 40.0 {
			t.Errorf("Sample value %f outside filter range [20.0, 40.0]", sample.Value)
		}
	}

	// Clean up
	r.Del("test_ts_filter")
}