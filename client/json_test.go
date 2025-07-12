// +build integration

package client

import (
	"encoding/json"
	"testing"
)

// Helper function to check if RedisJSON module is available
func isJSONModuleAvailable(t *testing.T) bool {
	// Try a basic JSON command - if it fails, module is not available
	_, err := r.JSONSet("test_json_module", ".", `"test"`)
	if err != nil {
		t.Skip("RedisJSON module not available, skipping JSON tests")
		return false
	}
	r.Del("test_json_module") // Clean up
	return true
}

func TestJSONSet(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Test basic JSON.SET
	result, err := r.JSONSet("json_key", ".", `{"name": "John", "age": 30}`)
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Test JSON.SET with NX option (should fail since key exists)
	_, err = r.JSONSet("json_key", ".", `{"name": "Jane"}`, &JSONSetOptions{NX: true})
	if err == nil {
		t.Error("Expected error with NX option on existing key")
	}

	// Test JSON.SET with XX option (should succeed since key exists)
	result, err = r.JSONSet("json_key", ".", `{"name": "Jane"}`, &JSONSetOptions{XX: true})
	if err != nil {
		t.Error(err)
	}
	if result != "OK" {
		t.Errorf("Expected OK, got %s", result)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONGet(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"name": "John", "age": 30, "city": "New York"}`)

	// Test basic JSON.GET
	data, err := r.JSONGet("json_key")
	if err != nil {
		t.Error(err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Error(err)
	}
	if result["name"] != "John" {
		t.Errorf("Expected name to be John, got %v", result["name"])
	}

	// Test JSON.GET with specific paths
	data, err = r.JSONGet("json_key", &JSONGetOptions{Paths: []string{".name"}})
	if err != nil {
		t.Error(err)
	}
	if string(data) != `"John"` {
		t.Errorf("Expected \"John\", got %s", string(data))
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONDel(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"name": "John", "age": 30, "city": "New York"}`)

	// Test JSON.DEL with specific path
	deleted, err := r.JSONDel("json_key", ".age")
	if err != nil {
		t.Error(err)
	}
	if deleted != 1 {
		t.Errorf("Expected 1 deleted, got %d", deleted)
	}

	// Verify age was deleted
	data, _ := r.JSONGet("json_key")
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	if _, exists := result["age"]; exists {
		t.Error("Expected age to be deleted")
	}

	// Test JSON.DEL entire key
	deleted, err = r.JSONDel("json_key")
	if err != nil {
		t.Error(err)
	}
	if deleted != 1 {
		t.Errorf("Expected 1 deleted, got %d", deleted)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONType(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Test different JSON types
	r.JSONSet("json_key", ".string", `"hello"`)
	r.JSONSet("json_key", ".number", `42`)
	r.JSONSet("json_key", ".array", `[1,2,3]`)
	r.JSONSet("json_key", ".object", `{"key": "value"}`)
	r.JSONSet("json_key", ".bool", `true`)
	r.JSONSet("json_key", ".null", `null`)

	// Test string type
	typeResult, err := r.JSONType("json_key", ".string")
	if err != nil {
		t.Error(err)
	}
	if typeResult != "string" {
		t.Errorf("Expected string, got %s", typeResult)
	}

	// Test number type
	typeResult, err = r.JSONType("json_key", ".number")
	if err != nil {
		t.Error(err)
	}
	if typeResult != "number" {
		t.Errorf("Expected number, got %s", typeResult)
	}

	// Test array type
	typeResult, err = r.JSONType("json_key", ".array")
	if err != nil {
		t.Error(err)
	}
	if typeResult != "array" {
		t.Errorf("Expected array, got %s", typeResult)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONNumIncrBy(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"count": 10}`)

	// Test JSON.NUMINCRBY
	result, err := r.JSONNumIncrBy("json_key", ".count", 5.5)
	if err != nil {
		t.Error(err)
	}
	if result != 15.5 {
		t.Errorf("Expected 15.5, got %f", result)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONNumMultBy(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"count": 10}`)

	// Test JSON.NUMMULTBY
	result, err := r.JSONNumMultBy("json_key", ".count", 2.5)
	if err != nil {
		t.Error(err)
	}
	if result != 25.0 {
		t.Errorf("Expected 25.0, got %f", result)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONStrAppend(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"greeting": "hello"}`)

	// Test JSON.STRAPPEND
	length, err := r.JSONStrAppend("json_key", ".greeting", `" world"`)
	if err != nil {
		t.Error(err)
	}
	if length != 11 { // "hello world" length
		t.Errorf("Expected 11, got %d", length)
	}

	// Verify the result
	data, _ := r.JSONGet("json_key", &JSONGetOptions{Paths: []string{".greeting"}})
	if string(data) != `"hello world"` {
		t.Errorf("Expected \"hello world\", got %s", string(data))
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONStrLen(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"greeting": "hello"}`)

	// Test JSON.STRLEN
	length, err := r.JSONStrLen("json_key", ".greeting")
	if err != nil {
		t.Error(err)
	}
	if length != 5 { // "hello" length
		t.Errorf("Expected 5, got %d", length)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONArrAppend(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"numbers": [1, 2, 3]}`)

	// Test JSON.ARRAPPEND
	length, err := r.JSONArrAppend("json_key", ".numbers", 4, 5)
	if err != nil {
		t.Error(err)
	}
	if length != 5 { // [1, 2, 3, 4, 5] length
		t.Errorf("Expected 5, got %d", length)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONArrIndex(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"numbers": [1, 2, 3, 2, 4]}`)

	// Test JSON.ARRINDEX
	index, err := r.JSONArrIndex("json_key", ".numbers", 2)
	if err != nil {
		t.Error(err)
	}
	if index != 1 { // First occurrence of 2 is at index 1
		t.Errorf("Expected 1, got %d", index)
	}

	// Test with start position
	index, err = r.JSONArrIndex("json_key", ".numbers", 2, 2)
	if err != nil {
		t.Error(err)
	}
	if index != 3 { // Second occurrence of 2 is at index 3
		t.Errorf("Expected 3, got %d", index)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONArrInsert(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"numbers": [1, 3]}`)

	// Test JSON.ARRINSERT
	length, err := r.JSONArrInsert("json_key", ".numbers", 1, 2)
	if err != nil {
		t.Error(err)
	}
	if length != 3 { // [1, 2, 3] length
		t.Errorf("Expected 3, got %d", length)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONArrLen(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"numbers": [1, 2, 3, 4, 5]}`)

	// Test JSON.ARRLEN
	length, err := r.JSONArrLen("json_key", ".numbers")
	if err != nil {
		t.Error(err)
	}
	if length != 5 {
		t.Errorf("Expected 5, got %d", length)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONArrPop(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"numbers": [1, 2, 3]}`)

	// Test JSON.ARRPOP (default last element)
	popped, err := r.JSONArrPop("json_key", ".numbers")
	if err != nil {
		t.Error(err)
	}
	if string(popped) != "3" {
		t.Errorf("Expected \"3\", got %s", string(popped))
	}

	// Test JSON.ARRPOP with specific index
	popped, err = r.JSONArrPop("json_key", ".numbers", 0)
	if err != nil {
		t.Error(err)
	}
	if string(popped) != "1" {
		t.Errorf("Expected \"1\", got %s", string(popped))
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONArrTrim(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"numbers": [1, 2, 3, 4, 5]}`)

	// Test JSON.ARRTRIM
	length, err := r.JSONArrTrim("json_key", ".numbers", 1, 3)
	if err != nil {
		t.Error(err)
	}
	if length != 3 { // [2, 3, 4] length
		t.Errorf("Expected 3, got %d", length)
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONObjKeys(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"name": "John", "age": 30, "city": "New York"}`)

	// Test JSON.OBJKEYS
	keys, err := r.JSONObjKeys("json_key")
	if err != nil {
		t.Error(err)
	}
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Check if all expected keys are present
	expectedKeys := map[string]bool{"name": false, "age": false, "city": false}
	for _, key := range keys {
		if _, exists := expectedKeys[key]; exists {
			expectedKeys[key] = true
		}
	}
	for key, found := range expectedKeys {
		if !found {
			t.Errorf("Expected key %s not found", key)
		}
	}

	// Clean up
	r.Del("json_key")
}

func TestJSONObjLen(t *testing.T) {
	if !isJSONModuleAvailable(t) {
		return
	}

	// Set up test data
	r.JSONSet("json_key", ".", `{"name": "John", "age": 30, "city": "New York"}`)

	// Test JSON.OBJLEN
	length, err := r.JSONObjLen("json_key")
	if err != nil {
		t.Error(err)
	}
	if length != 3 {
		t.Errorf("Expected 3, got %d", length)
	}

	// Clean up
	r.Del("json_key")
}