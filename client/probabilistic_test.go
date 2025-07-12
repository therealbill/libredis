// +build integration

package client

import (
	"testing"
)

// Helper function to check if RedisBloom module is available
func isBloomModuleAvailable(t *testing.T) bool {
	// Try to create a simple bloom filter - if it fails, module is not available
	_, err := r.BFReserve("test_bloom_module", 0.01, 1000)
	if err != nil {
		t.Skip("RedisBloom module not available, skipping Probabilistic tests")
		return false
	}
	r.Del("test_bloom_module") // Clean up
	return true
}

// Bloom Filter Tests

func TestBFReserve(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Test basic Bloom filter creation
	result, err := r.BFReserve("test_bf", 0.01, 1000)
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Test Bloom filter creation with options
	result, err = r.BFReserve("test_bf_options", 0.01, 1000, &BFReserveOptions{
		Expansion:  2,
		NonScaling: false,
	})
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Clean up
	r.Del("test_bf", "test_bf_options")
}

func TestBFAdd(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Bloom filter
	r.BFReserve("test_bf_add", 0.01, 1000)

	// Test adding an item
	added, err := r.BFAdd("test_bf_add", "item1")
	if err != nil {
		t.Error(err)
	}
	if !added {
		t.Error("Expected item to be added (true)")
	}

	// Test adding the same item again (should return false)
	added, err = r.BFAdd("test_bf_add", "item1")
	if err != nil {
		t.Error(err)
	}
	if added {
		t.Error("Expected item not to be added again (false)")
	}

	// Clean up
	r.Del("test_bf_add")
}

func TestBFMAdd(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Bloom filter
	r.BFReserve("test_bf_madd", 0.01, 1000)

	// Test adding multiple items
	results, err := r.BFMAdd("test_bf_madd", "item1", "item2", "item3")
	if err != nil {
		t.Error(err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// All items should be new (true)
	for i, added := range results {
		if !added {
			t.Errorf("Expected item %d to be added", i)
		}
	}

	// Test adding same items again
	results, err = r.BFMAdd("test_bf_madd", "item1", "item2", "item4")
	if err != nil {
		t.Error(err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// First two should be false (already exist), third should be true (new)
	if results[0] || results[1] {
		t.Error("Expected first two items to not be added (false)")
	}
	if !results[2] {
		t.Error("Expected third item to be added (true)")
	}

	// Clean up
	r.Del("test_bf_madd")
}

func TestBFExists(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Bloom filter and add items
	r.BFReserve("test_bf_exists", 0.01, 1000)
	r.BFAdd("test_bf_exists", "existing_item")

	// Test checking existing item
	exists, err := r.BFExists("test_bf_exists", "existing_item")
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error("Expected existing item to be found")
	}

	// Test checking non-existing item
	exists, err = r.BFExists("test_bf_exists", "non_existing_item")
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Error("Expected non-existing item not to be found")
	}

	// Clean up
	r.Del("test_bf_exists")
}

func TestBFMExists(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Bloom filter and add items
	r.BFReserve("test_bf_mexists", 0.01, 1000)
	r.BFMAdd("test_bf_mexists", "item1", "item2")

	// Test checking multiple items
	results, err := r.BFMExists("test_bf_mexists", "item1", "item2", "item3")
	if err != nil {
		t.Error(err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// First two should exist, third should not
	if !results[0] || !results[1] {
		t.Error("Expected first two items to exist")
	}
	if results[2] {
		t.Error("Expected third item not to exist")
	}

	// Clean up
	r.Del("test_bf_mexists")
}

// Cuckoo Filter Tests

func TestCFReserve(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Test basic Cuckoo filter creation
	result, err := r.CFReserve("test_cf", 1000)
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Test Cuckoo filter creation with options
	result, err = r.CFReserve("test_cf_options", 1000, &CFReserveOptions{
		BucketSize:    4,
		MaxIterations: 20,
		Expansion:     1,
	})
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Clean up
	r.Del("test_cf", "test_cf_options")
}

func TestCFAdd(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Cuckoo filter
	r.CFReserve("test_cf_add", 1000)

	// Test adding an item
	added, err := r.CFAdd("test_cf_add", "item1")
	if err != nil {
		t.Error(err)
	}
	if !added {
		t.Error("Expected item to be added (true)")
	}

	// Test adding the same item again (should return false)
	added, err = r.CFAdd("test_cf_add", "item1")
	if err != nil {
		t.Error(err)
	}
	if added {
		t.Error("Expected item not to be added again (false)")
	}

	// Clean up
	r.Del("test_cf_add")
}

func TestCFExists(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Cuckoo filter and add item
	r.CFReserve("test_cf_exists", 1000)
	r.CFAdd("test_cf_exists", "existing_item")

	// Test checking existing item
	exists, err := r.CFExists("test_cf_exists", "existing_item")
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Error("Expected existing item to be found")
	}

	// Test checking non-existing item
	exists, err = r.CFExists("test_cf_exists", "non_existing_item")
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Error("Expected non-existing item not to be found")
	}

	// Clean up
	r.Del("test_cf_exists")
}

func TestCFDel(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Cuckoo filter and add item
	r.CFReserve("test_cf_del", 1000)
	r.CFAdd("test_cf_del", "item_to_delete")

	// Verify item exists
	exists, _ := r.CFExists("test_cf_del", "item_to_delete")
	if !exists {
		t.Error("Item should exist before deletion")
	}

	// Test deleting the item
	deleted, err := r.CFDel("test_cf_del", "item_to_delete")
	if err != nil {
		t.Error(err)
	}
	if !deleted {
		t.Error("Expected item to be deleted (true)")
	}

	// Verify item no longer exists
	exists, _ = r.CFExists("test_cf_del", "item_to_delete")
	if exists {
		t.Error("Item should not exist after deletion")
	}

	// Test deleting non-existing item
	deleted, err = r.CFDel("test_cf_del", "non_existing_item")
	if err != nil {
		t.Error(err)
	}
	if deleted {
		t.Error("Expected non-existing item not to be deleted (false)")
	}

	// Clean up
	r.Del("test_cf_del")
}

// Count-Min Sketch Tests

func TestCMSInitByDim(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Test creating Count-Min Sketch with dimensions
	result, err := r.CMSInitByDim("test_cms_dim", 2000, 5)
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Clean up
	r.Del("test_cms_dim")
}

func TestCMSInitByProb(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Test creating Count-Min Sketch with error rate and probability
	result, err := r.CMSInitByProb("test_cms_prob", 0.001, 0.99)
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Clean up
	r.Del("test_cms_prob")
}

func TestCMSIncrBy(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Count-Min Sketch
	r.CMSInitByDim("test_cms_incr", 2000, 5)

	// Test incrementing items
	results, err := r.CMSIncrBy("test_cms_incr", "item1", 5, "item2", 3, "item1", 2)
	if err != nil {
		t.Error(err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// First increment of item1 should be 5
	if results[0] != 5 {
		t.Errorf("Expected first result to be 5, got %d", results[0])
	}
	
	// First increment of item2 should be 3
	if results[1] != 3 {
		t.Errorf("Expected second result to be 3, got %d", results[1])
	}
	
	// Second increment of item1 should be 7 (5 + 2)
	if results[2] != 7 {
		t.Errorf("Expected third result to be 7, got %d", results[2])
	}

	// Clean up
	r.Del("test_cms_incr")
}

func TestCMSQuery(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Count-Min Sketch and add items
	r.CMSInitByDim("test_cms_query", 2000, 5)
	r.CMSIncrBy("test_cms_query", "item1", 10, "item2", 5)

	// Test querying items
	results, err := r.CMSQuery("test_cms_query", "item1", "item2", "item3")
	if err != nil {
		t.Error(err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// item1 should have count 10
	if results[0] != 10 {
		t.Errorf("Expected item1 count to be 10, got %d", results[0])
	}
	
	// item2 should have count 5
	if results[1] != 5 {
		t.Errorf("Expected item2 count to be 5, got %d", results[1])
	}
	
	// item3 should have count 0 (never incremented)
	if results[2] != 0 {
		t.Errorf("Expected item3 count to be 0, got %d", results[2])
	}

	// Clean up
	r.Del("test_cms_query")
}

func TestBFInfo(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Bloom filter and add some items
	r.BFReserve("test_bf_info", 0.01, 1000)
	r.BFMAdd("test_bf_info", "item1", "item2", "item3")

	// Test getting Bloom filter info
	info, err := r.BFInfo("test_bf_info")
	if err != nil {
		t.Error(err)
	}
	if len(info) == 0 {
		t.Error("Expected info to contain data")
	}

	// Check for some expected fields
	if _, exists := info["Capacity"]; !exists {
		t.Error("Expected Capacity in info")
	}
	if _, exists := info["Size"]; !exists {
		t.Error("Expected Size in info")
	}

	// Clean up
	r.Del("test_bf_info")
}

func TestCFInfo(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Cuckoo filter and add some items
	r.CFReserve("test_cf_info", 1000)
	r.CFAdd("test_cf_info", "item1")
	r.CFAdd("test_cf_info", "item2")

	// Test getting Cuckoo filter info
	info, err := r.CFInfo("test_cf_info")
	if err != nil {
		t.Error(err)
	}
	if len(info) == 0 {
		t.Error("Expected info to contain data")
	}

	// Check for some expected fields
	if _, exists := info["Size"]; !exists {
		t.Error("Expected Size in info")
	}

	// Clean up
	r.Del("test_cf_info")
}

func TestCMSInfo(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create Count-Min Sketch and add some items
	r.CMSInitByDim("test_cms_info", 2000, 5)
	r.CMSIncrBy("test_cms_info", "item1", 10, "item2", 5)

	// Test getting Count-Min Sketch info
	info, err := r.CMSInfo("test_cms_info")
	if err != nil {
		t.Error(err)
	}
	if len(info) == 0 {
		t.Error("Expected info to contain data")
	}

	// Check for some expected fields
	if _, exists := info["width"]; !exists {
		t.Error("Expected width in info")
	}
	if _, exists := info["depth"]; !exists {
		t.Error("Expected depth in info")
	}

	// Clean up
	r.Del("test_cms_info")
}

func TestCMSMerge(t *testing.T) {
	if !isBloomModuleAvailable(t) {
		return
	}

	// Create multiple Count-Min Sketches
	r.CMSInitByDim("test_cms_merge_1", 2000, 5)
	r.CMSInitByDim("test_cms_merge_2", 2000, 5)
	r.CMSInitByDim("test_cms_merge_dest", 2000, 5)

	// Add different items to each sketch
	r.CMSIncrBy("test_cms_merge_1", "item1", 10, "item2", 5)
	r.CMSIncrBy("test_cms_merge_2", "item1", 3, "item3", 7)

	// Test merging sketches
	result, err := r.CMSMerge("test_cms_merge_dest", []string{"test_cms_merge_1", "test_cms_merge_2"})
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Verify merged results
	results, err := r.CMSQuery("test_cms_merge_dest", "item1", "item2", "item3")
	if err != nil {
		t.Error(err)
	}

	// item1 should be 13 (10 + 3), item2 should be 5, item3 should be 7
	if results[0] != 13 {
		t.Errorf("Expected item1 count to be 13, got %d", results[0])
	}
	if results[1] != 5 {
		t.Errorf("Expected item2 count to be 5, got %d", results[1])
	}
	if results[2] != 7 {
		t.Errorf("Expected item3 count to be 7, got %d", results[2])
	}

	// Clean up
	r.Del("test_cms_merge_1", "test_cms_merge_2", "test_cms_merge_dest")
}