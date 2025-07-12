package client

import (
	"testing"
	"time"
)

// Test stream constants
func TestStreamConstants(t *testing.T) {
	if StreamIDAutoGenerate != "*" {
		t.Error("StreamIDAutoGenerate constant incorrect")
	}

	if StreamIDLatest != "$" {
		t.Error("StreamIDLatest constant incorrect")
	}

	if StreamIDEarliest != "0-0" {
		t.Error("StreamIDEarliest constant incorrect")
	}
}

// Test data structures
func TestStreamStructures(t *testing.T) {
	entry := StreamEntry{
		ID:     "1234567890123-0",
		Fields: map[string]string{"field1": "value1", "field2": "value2"},
	}

	if entry.ID != "1234567890123-0" || len(entry.Fields) != 2 {
		t.Error("StreamEntry struct not working correctly")
	}

	message := StreamMessage{
		Stream:  "mystream",
		Entries: []StreamEntry{entry},
	}

	if message.Stream != "mystream" || len(message.Entries) != 1 {
		t.Error("StreamMessage struct not working correctly")
	}
}

// Basic Stream Operations Tests

func TestXAdd(t *testing.T) {
	r.Del("mystream")

	fields := map[string]string{
		"name":    "Alice",
		"surname": "Jones",
		"age":     "30",
	}

	// Test auto-generated ID
	id, err := r.XAdd("mystream", StreamIDAutoGenerate, fields)
	if err != nil {
		t.Error(err)
	}
	if id == "" {
		t.Error("Expected auto-generated ID")
	}

	// Test explicit ID
	explicitID := "1234567890123-0"
	id2, err := r.XAdd("mystream", explicitID, fields)
	if err != nil {
		t.Error(err)
	}
	if id2 != explicitID {
		t.Errorf("Expected ID %s, got %s", explicitID, id2)
	}
}

func TestXAddWithOptions(t *testing.T) {
	r.Del("mystream")

	fields := map[string]string{"test": "value"}

	// Test with MAXLEN option
	opts := XAddOptions{
		MaxLen: 5,
	}

	// Add several entries
	for i := 0; i < 10; i++ {
		_, err := r.XAddWithOptions("mystream", StreamIDAutoGenerate, fields, opts)
		if err != nil {
			t.Error(err)
		}
	}

	// Check stream length
	length, err := r.XLen("mystream")
	if err != nil {
		t.Error(err)
	}
	if length > 5 {
		t.Errorf("Expected stream length <= 5, got %d", length)
	}
}

func TestXLen(t *testing.T) {
	r.Del("mystream")

	// Empty stream
	length, err := r.XLen("mystream")
	if err != nil {
		t.Error(err)
	}
	if length != 0 {
		t.Errorf("Expected length 0, got %d", length)
	}

	// Add entries
	fields := map[string]string{"test": "value"}
	r.XAdd("mystream", StreamIDAutoGenerate, fields)
	r.XAdd("mystream", StreamIDAutoGenerate, fields)

	length, err = r.XLen("mystream")
	if err != nil {
		t.Error(err)
	}
	if length != 2 {
		t.Errorf("Expected length 2, got %d", length)
	}
}

func TestXRange(t *testing.T) {
	r.Del("mystream")

	// Add test entries
	fields1 := map[string]string{"name": "Alice", "age": "30"}
	fields2 := map[string]string{"name": "Bob", "age": "25"}

	id1, _ := r.XAdd("mystream", StreamIDAutoGenerate, fields1)
	id2, _ := r.XAdd("mystream", StreamIDAutoGenerate, fields2)

	// Test range query
	entries, err := r.XRange("mystream", "-", "+")
	if err != nil {
		t.Error(err)
	}
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}

	// Verify entry content
	if entries[0].ID != id1 || entries[0].Fields["name"] != "Alice" {
		t.Error("First entry content incorrect")
	}
	if entries[1].ID != id2 || entries[1].Fields["name"] != "Bob" {
		t.Error("Second entry content incorrect")
	}
}

func TestXRangeWithOptions(t *testing.T) {
	r.Del("mystream")

	// Add multiple entries
	fields := map[string]string{"test": "value"}
	for i := 0; i < 5; i++ {
		r.XAdd("mystream", StreamIDAutoGenerate, fields)
	}

	// Test with COUNT option
	opts := XRangeOptions{Count: 2}
	entries, err := r.XRangeWithOptions("mystream", "-", "+", opts)
	if err != nil {
		t.Error(err)
	}
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
}

func TestXRevRange(t *testing.T) {
	r.Del("mystream")

	// Add test entries
	fields1 := map[string]string{"seq": "1"}
	fields2 := map[string]string{"seq": "2"}

	id1, _ := r.XAdd("mystream", StreamIDAutoGenerate, fields1)
	id2, _ := r.XAdd("mystream", StreamIDAutoGenerate, fields2)

	// Test reverse range
	entries, err := r.XRevRange("mystream", "+", "-")
	if err != nil {
		t.Error(err)
	}
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}

	// Should be in reverse order
	if entries[0].ID != id2 || entries[1].ID != id1 {
		t.Error("Entries not in reverse order")
	}
}

func TestXRead(t *testing.T) {
	r.Del("stream1", "stream2")

	// Add entries to streams
	fields1 := map[string]string{"msg": "hello"}
	fields2 := map[string]string{"msg": "world"}

	r.XAdd("stream1", StreamIDAutoGenerate, fields1)
	r.XAdd("stream2", StreamIDAutoGenerate, fields2)

	// Read from both streams
	streams := map[string]string{
		"stream1": "0-0",
		"stream2": "0-0",
	}

	messages, err := r.XRead(streams)
	if err != nil {
		t.Error(err)
	}
	if len(messages) != 2 {
		t.Errorf("Expected 2 stream messages, got %d", len(messages))
	}

	// Verify stream names
	streamNames := make(map[string]bool)
	for _, msg := range messages {
		streamNames[msg.Stream] = true
	}
	if !streamNames["stream1"] || !streamNames["stream2"] {
		t.Error("Missing expected streams in response")
	}
}

func TestXReadWithOptions(t *testing.T) {
	r.Del("mystream")

	// Add multiple entries
	fields := map[string]string{"data": "test"}
	for i := 0; i < 5; i++ {
		r.XAdd("mystream", StreamIDAutoGenerate, fields)
	}

	// Test with COUNT option
	opts := XReadOptions{Count: 2}
	streams := map[string]string{"mystream": "0-0"}

	messages, err := r.XReadWithOptions(streams, opts)
	if err != nil {
		t.Error(err)
	}
	if len(messages) != 1 || len(messages[0].Entries) != 2 {
		t.Error("COUNT option not working correctly")
	}
}

func TestXDel(t *testing.T) {
	r.Del("mystream")

	// Add entries
	fields := map[string]string{"test": "value"}
	id1, _ := r.XAdd("mystream", StreamIDAutoGenerate, fields)
	id2, _ := r.XAdd("mystream", StreamIDAutoGenerate, fields)

	// Delete one entry
	deleted, err := r.XDel("mystream", id1)
	if err != nil {
		t.Error(err)
	}
	if deleted != 1 {
		t.Errorf("Expected 1 deleted entry, got %d", deleted)
	}

	// Verify deletion
	entries, _ := r.XRange("mystream", "-", "+")
	if len(entries) != 1 || entries[0].ID != id2 {
		t.Error("Entry not deleted correctly")
	}
}

func TestXTrim(t *testing.T) {
	r.Del("mystream")

	// Add multiple entries
	fields := map[string]string{"test": "value"}
	for i := 0; i < 10; i++ {
		r.XAdd("mystream", StreamIDAutoGenerate, fields)
	}

	// Trim to 5 entries
	trimmed, err := r.XTrim("mystream", "MAXLEN", "5")
	if err != nil {
		t.Error(err)
	}
	if trimmed != 5 {
		t.Errorf("Expected 5 trimmed entries, got %d", trimmed)
	}

	// Verify trim
	length, _ := r.XLen("mystream")
	if length != 5 {
		t.Errorf("Expected length 5 after trim, got %d", length)
	}
}

// Consumer Group Tests

func TestXGroupCreate(t *testing.T) {
	r.Del("mystream")

	// Add an entry first
	fields := map[string]string{"test": "value"}
	r.XAdd("mystream", StreamIDAutoGenerate, fields)

	// Create consumer group
	err := r.XGroupCreate("mystream", "mygroup", "0")
	if err != nil {
		t.Error(err)
	}

	// Creating the same group again should fail
	err = r.XGroupCreate("mystream", "mygroup", "0")
	if err == nil {
		t.Error("Expected error when creating duplicate group")
	}
}

func TestXGroupCreateWithOptions(t *testing.T) {
	r.Del("newstream")

	// Test with MKSTREAM option
	opts := XGroupCreateOptions{MkStream: true}
	err := r.XGroupCreateWithOptions("newstream", "mygroup", "$", opts)
	if err != nil {
		t.Error(err)
	}

	// Stream should exist now
	length, err := r.XLen("newstream")
	if err != nil {
		t.Error(err)
	}
	if length != 0 {
		t.Errorf("Expected empty stream, got length %d", length)
	}
}

func TestXGroupDestroy(t *testing.T) {
	r.Del("mystream")

	// Create group
	fields := map[string]string{"test": "value"}
	r.XAdd("mystream", StreamIDAutoGenerate, fields)
	r.XGroupCreate("mystream", "testgroup", "0")

	// Destroy group
	destroyed, err := r.XGroupDestroy("mystream", "testgroup")
	if err != nil {
		t.Error(err)
	}
	if destroyed != 1 {
		t.Errorf("Expected 1 destroyed group, got %d", destroyed)
	}

	// Destroying again should return 0
	destroyed, err = r.XGroupDestroy("mystream", "testgroup")
	if err != nil {
		t.Error(err)
	}
	if destroyed != 0 {
		t.Errorf("Expected 0 destroyed groups, got %d", destroyed)
	}
}

func TestXReadGroup(t *testing.T) {
	r.Del("mystream")

	// Setup stream and group
	fields := map[string]string{"msg": "hello"}
	r.XAdd("mystream", StreamIDAutoGenerate, fields)
	r.XGroupCreate("mystream", "mygroup", "0")

	// Read as consumer
	streams := map[string]string{"mystream": ">"}
	messages, err := r.XReadGroup("mygroup", "consumer1", streams)
	if err != nil {
		t.Error(err)
	}
	if len(messages) != 1 || len(messages[0].Entries) != 1 {
		t.Error("XReadGroup not working correctly")
	}

	// Reading again should return nothing (messages already delivered)
	messages2, err := r.XReadGroup("mygroup", "consumer1", streams)
	if err != nil {
		t.Error(err)
	}
	if messages2 != nil && len(messages2) > 0 {
		t.Error("Expected no new messages on second read")
	}
}

func TestXAck(t *testing.T) {
	r.Del("mystream")

	// Setup stream and group
	fields := map[string]string{"msg": "hello"}
	id, _ := r.XAdd("mystream", StreamIDAutoGenerate, fields)
	r.XGroupCreate("mystream", "mygroup", "0")

	// Read message
	streams := map[string]string{"mystream": ">"}
	r.XReadGroup("mygroup", "consumer1", streams)

	// Acknowledge message
	acked, err := r.XAck("mystream", "mygroup", id)
	if err != nil {
		t.Error(err)
	}
	if acked != 1 {
		t.Errorf("Expected 1 acknowledged message, got %d", acked)
	}
}

func TestXPending(t *testing.T) {
	r.Del("mystream")

	// Setup stream and group
	fields := map[string]string{"msg": "hello"}
	r.XAdd("mystream", StreamIDAutoGenerate, fields)
	r.XGroupCreate("mystream", "mygroup", "0")

	// Read message (creates pending entry)
	streams := map[string]string{"mystream": ">"}
	r.XReadGroup("mygroup", "consumer1", streams)

	// Check pending info
	pending, err := r.XPending("mystream", "mygroup")
	if err != nil {
		t.Error(err)
	}
	if pending.Count != 1 {
		t.Errorf("Expected 1 pending message, got %d", pending.Count)
	}
	if pending.Consumers["consumer1"] != 1 {
		t.Error("Consumer pending count incorrect")
	}
}

func TestXInfoStream(t *testing.T) {
	r.Del("mystream")

	// Add entry
	fields := map[string]string{"test": "value"}
	r.XAdd("mystream", StreamIDAutoGenerate, fields)

	// Get stream info
	info, err := r.XInfoStream("mystream")
	if err != nil {
		t.Error(err)
	}
	if info == nil {
		t.Error("Expected stream info")
	}

	// Should contain basic stream information
	if length, ok := info["length"]; !ok || length.(int64) != 1 {
		t.Error("Stream length not reported correctly")
	}
}

func TestXInfoGroups(t *testing.T) {
	r.Del("mystream")

	// Setup stream and group
	fields := map[string]string{"test": "value"}
	r.XAdd("mystream", StreamIDAutoGenerate, fields)
	r.XGroupCreate("mystream", "testgroup", "0")

	// Get groups info
	groups, err := r.XInfoGroups("mystream")
	if err != nil {
		t.Error(err)
	}
	if len(groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(groups))
	}

	// Verify group name
	if name, ok := groups[0]["name"]; !ok || name.(string) != "testgroup" {
		t.Error("Group name not reported correctly")
	}
}

func TestXInfoConsumers(t *testing.T) {
	r.Del("mystream")

	// Setup stream and group
	fields := map[string]string{"test": "value"}
	r.XAdd("mystream", StreamIDAutoGenerate, fields)
	r.XGroupCreate("mystream", "testgroup", "0")

	// Create consumer by reading
	streams := map[string]string{"mystream": ">"}
	r.XReadGroup("testgroup", "consumer1", streams)

	// Get consumers info
	consumers, err := r.XInfoConsumers("mystream", "testgroup")
	if err != nil {
		t.Error(err)
	}
	if len(consumers) != 1 {
		t.Errorf("Expected 1 consumer, got %d", len(consumers))
	}

	// Verify consumer name
	if name, ok := consumers[0]["name"]; !ok || name.(string) != "consumer1" {
		t.Error("Consumer name not reported correctly")
	}
}

// Advanced Tests

func TestStreamBlocking(t *testing.T) {
	r.Del("blockstream")

	// Test with very short timeout to avoid hanging the test
	opts := XReadOptions{Block: 1} // 1ms timeout
	streams := map[string]string{"blockstream": "$"}

	start := time.Now()
	messages, err := r.XReadWithOptions(streams, opts)
	elapsed := time.Since(start)

	if err != nil {
		t.Error(err)
	}
	if messages != nil && len(messages) > 0 {
		t.Error("Expected no messages on blocking read with timeout")
	}

	// Should have blocked for approximately the timeout duration
	if elapsed < time.Millisecond {
		t.Error("Blocking read returned too quickly")
	}
}

func TestXClaimBasic(t *testing.T) {
	r.Del("mystream")

	// Setup stream and group
	fields := map[string]string{"msg": "test"}
	id, _ := r.XAdd("mystream", StreamIDAutoGenerate, fields)
	r.XGroupCreate("mystream", "mygroup", "0")

	// Read with one consumer
	streams := map[string]string{"mystream": ">"}
	r.XReadGroup("mygroup", "consumer1", streams)

	// Sleep briefly to ensure idle time
	time.Sleep(10 * time.Millisecond)

	// Claim with another consumer
	entries, err := r.XClaim("mystream", "mygroup", "consumer2", 1, []string{id})
	if err != nil {
		t.Error(err)
	}
	if len(entries) != 1 || entries[0].ID != id {
		t.Error("XClaim not working correctly")
	}
}
