// +build integration

package client

import (
	"testing"
)

// Helper function to check if RediSearch module is available
func isSearchModuleAvailable(t *testing.T) bool {
	// Try to create a simple index - if it fails, module is not available
	_, err := r.FTCreate("test_search_module", []FTFieldSchema{
		{Name: "title", Type: "TEXT"},
	})
	if err != nil {
		t.Skip("RediSearch module not available, skipping Search tests")
		return false
	}
	r.FTDropIndex("test_search_module", true) // Clean up
	return true
}

func TestFTCreate(t *testing.T) {
	if !isSearchModuleAvailable(t) {
		return
	}

	// Test basic index creation
	schema := []FTFieldSchema{
		{Name: "title", Type: "TEXT", Weight: 1.0},
		{Name: "body", Type: "TEXT"},
		{Name: "price", Type: "NUMERIC", Sortable: true},
		{Name: "location", Type: "GEO"},
		{Name: "tags", Type: "TAG", Separator: ","},
	}

	result, err := r.FTCreate("test_index", schema, &FTCreateOptions{
		OnHash: true,
		Prefix: []string{"product:"},
	})
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Clean up
	r.FTDropIndex("test_index", true)
}

func TestFTDropIndex(t *testing.T) {
	if !isSearchModuleAvailable(t) {
		return
	}

	// Create an index first
	schema := []FTFieldSchema{
		{Name: "title", Type: "TEXT"},
	}
	r.FTCreate("test_drop_index", schema)

	// Test dropping the index
	result, err := r.FTDropIndex("test_drop_index")
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}
}

func TestFTInfo(t *testing.T) {
	if !isSearchModuleAvailable(t) {
		return
	}

	// Create an index first
	schema := []FTFieldSchema{
		{Name: "title", Type: "TEXT"},
		{Name: "price", Type: "NUMERIC"},
	}
	r.FTCreate("test_info_index", schema)

	// Test getting index info
	info, err := r.FTInfo("test_info_index")
	if err != nil {
		t.Error(err)
	}
	if len(info) == 0 {
		t.Error("Expected info to contain data")
	}

	// Check if some expected keys are present
	if _, exists := info["index_name"]; !exists {
		t.Error("Expected index_name in info")
	}

	// Clean up
	r.FTDropIndex("test_info_index", true)
}

func TestFTSearch(t *testing.T) {
	if !isSearchModuleAvailable(t) {
		return
	}

	// Create an index
	schema := []FTFieldSchema{
		{Name: "title", Type: "TEXT"},
		{Name: "body", Type: "TEXT"},
		{Name: "price", Type: "NUMERIC"},
	}
	r.FTCreate("test_search_index", schema, &FTCreateOptions{
		OnHash: true,
		Prefix: []string{"doc:"},
	})

	// Add some test documents
	r.HSet("doc:1", "title", "First Product", "body", "This is a great product", "price", "100")
	r.HSet("doc:2", "title", "Second Product", "body", "This is another product", "price", "200")
	r.HSet("doc:3", "title", "Third Item", "body", "This is an item", "price", "50")

	// Test basic search
	results, err := r.FTSearch("test_search_index", "product")
	if err != nil {
		t.Error(err)
	}
	if len(results) < 2 { // Should return at least count + first result
		t.Errorf("Expected at least 2 results, got %d", len(results))
	}

	// Test search with options
	results, err = r.FTSearch("test_search_index", "product", &FTSearchOptions{
		WithScores: true,
		Limit:      &FTLimit{Offset: 0, Num: 10},
	})
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Error("Expected results with scores")
	}

	// Clean up
	r.Del("doc:1", "doc:2", "doc:3")
	r.FTDropIndex("test_search_index", true)
}

func TestFTAggregate(t *testing.T) {
	if !isSearchModuleAvailable(t) {
		return
	}

	// Create an index
	schema := []FTFieldSchema{
		{Name: "category", Type: "TAG"},
		{Name: "price", Type: "NUMERIC"},
		{Name: "quantity", Type: "NUMERIC"},
	}
	r.FTCreate("test_agg_index", schema, &FTCreateOptions{
		OnHash: true,
		Prefix: []string{"item:"},
	})

	// Add some test documents
	r.HSet("item:1", "category", "electronics", "price", "100", "quantity", "5")
	r.HSet("item:2", "category", "electronics", "price", "200", "quantity", "3")
	r.HSet("item:3", "category", "books", "price", "20", "quantity", "10")
	r.HSet("item:4", "category", "books", "price", "30", "quantity", "7")

	// Test aggregation
	results, err := r.FTAggregate("test_agg_index", "*", &FTAggregateOptions{
		GroupBy: &FTGroupBy{
			Fields: []string{"@category"},
			Reduce: []FTReduce{
				{Function: "COUNT", Args: []string{}, As: "count"},
				{Function: "AVG", Args: []string{"@price"}, As: "avg_price"},
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Error("Expected aggregation results")
	}

	// Clean up
	r.Del("item:1", "item:2", "item:3", "item:4")
	r.FTDropIndex("test_agg_index", true)
}

func TestFTExplain(t *testing.T) {
	if !isSearchModuleAvailable(t) {
		return
	}

	// Create an index
	schema := []FTFieldSchema{
		{Name: "title", Type: "TEXT"},
		{Name: "body", Type: "TEXT"},
	}
	r.FTCreate("test_explain_index", schema)

	// Test query explanation
	explanation, err := r.FTExplain("test_explain_index", "hello world")
	if err != nil {
		t.Error(err)
	}
	if explanation == "" {
		t.Error("Expected explanation to contain data")
	}

	// Clean up
	r.FTDropIndex("test_explain_index", true)
}

func TestFTAdd(t *testing.T) {
	if !isSearchModuleAvailable(t) {
		return
	}

	// Create an index
	schema := []FTFieldSchema{
		{Name: "title", Type: "TEXT"},
		{Name: "body", Type: "TEXT"},
	}
	r.FTCreate("test_add_index", schema)

	// Test adding a document (deprecated command)
	fields := map[string]interface{}{
		"title": "Test Document",
		"body":  "This is a test document",
	}
	result, err := r.FTAdd("test_add_index", "doc1", 1.0, fields)
	if err != nil {
		// FT.ADD might not be available in newer versions
		t.Skip("FT.ADD command not available (deprecated in RediSearch 2.0+)")
		return
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Clean up
	r.FTDropIndex("test_add_index", true)
}

func TestFTDel(t *testing.T) {
	if !isSearchModuleAvailable(t) {
		return
	}

	// Create an index
	schema := []FTFieldSchema{
		{Name: "title", Type: "TEXT"},
	}
	r.FTCreate("test_del_index", schema)

	// Add a document using FT.ADD if available, otherwise skip
	fields := map[string]interface{}{
		"title": "Test Document",
	}
	_, err := r.FTAdd("test_del_index", "doc1", 1.0, fields)
	if err != nil {
		t.Skip("FT.ADD command not available, skipping FT.DEL test")
		return
	}

	// Test deleting the document
	deleted, err := r.FTDel("test_del_index", "doc1")
	if err != nil {
		t.Error(err)
	}
	if deleted != 1 {
		t.Errorf("Expected 1 deleted, got %d", deleted)
	}

	// Clean up
	r.FTDropIndex("test_del_index", true)
}

func TestFTSearchWithFilters(t *testing.T) {
	if !isSearchModuleAvailable(t) {
		return
	}

	// Create an index with numeric and geo fields
	schema := []FTFieldSchema{
		{Name: "title", Type: "TEXT"},
		{Name: "price", Type: "NUMERIC"},
		{Name: "location", Type: "GEO"},
	}
	r.FTCreate("test_filter_index", schema, &FTCreateOptions{
		OnHash: true,
		Prefix: []string{"prod:"},
	})

	// Add test documents
	r.HSet("prod:1", "title", "Cheap Product", "price", "50", "location", "-122.4194,37.7749")
	r.HSet("prod:2", "title", "Expensive Product", "price", "500", "location", "-122.4094,37.7849")
	r.HSet("prod:3", "title", "Medium Product", "price", "100", "location", "-122.4294,37.7649")

	// Test search with numeric filter
	results, err := r.FTSearch("test_filter_index", "*", &FTSearchOptions{
		Filter: []FTNumericFilter{
			{Field: "price", Min: 40, Max: 150},
		},
		Limit: &FTLimit{Offset: 0, Num: 10},
	})
	if err != nil {
		t.Error(err)
	}
	if len(results) < 2 { // Should return count + results
		t.Errorf("Expected at least 2 results with numeric filter, got %d", len(results))
	}

	// Test search with geo filter
	results, err = r.FTSearch("test_filter_index", "*", &FTSearchOptions{
		GeoFilter: &FTGeoFilter{
			Field:     "location",
			Longitude: -122.4194,
			Latitude:  37.7749,
			Radius:    10,
			Unit:      "km",
		},
		Limit: &FTLimit{Offset: 0, Num: 10},
	})
	if err != nil {
		t.Error(err)
	}
	if len(results) < 1 {
		t.Error("Expected at least 1 result with geo filter")
	}

	// Clean up
	r.Del("prod:1", "prod:2", "prod:3")
	r.FTDropIndex("test_filter_index", true)
}

func TestFTSearchWithHighlight(t *testing.T) {
	if !isSearchModuleAvailable(t) {
		return
	}

	// Create an index
	schema := []FTFieldSchema{
		{Name: "title", Type: "TEXT"},
		{Name: "body", Type: "TEXT"},
	}
	r.FTCreate("test_highlight_index", schema, &FTCreateOptions{
		OnHash: true,
		Prefix: []string{"article:"},
	})

	// Add test document
	r.HSet("article:1", "title", "Redis Search Tutorial", "body", "This tutorial explains Redis search capabilities")

	// Test search with highlighting
	results, err := r.FTSearch("test_highlight_index", "Redis search", &FTSearchOptions{
		Highlight: &FTHighlightOptions{
			Fields: []string{"title", "body"},
			Tags: &FTHighlightTags{
				Open:  "<b>",
				Close: "</b>",
			},
		},
		Return: []string{"title", "body"},
	})
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Error("Expected highlighted results")
	}

	// Clean up
	r.Del("article:1")
	r.FTDropIndex("test_highlight_index", true)
}