package client

import (
	"testing"
)

// Test new constants and data structures
func TestListDirectionConstants(t *testing.T) {
	if ListDirectionLeft != "LEFT" {
		t.Error("ListDirectionLeft constant incorrect")
	}
	
	if ListDirectionRight != "RIGHT" {
		t.Error("ListDirectionRight constant incorrect")
	}
}

func TestLPosOptions(t *testing.T) {
	opts := LPosOptions{
		Rank:   1,
		Count:  5,
		MaxLen: 100,
	}
	
	if opts.Rank != 1 || opts.Count != 5 || opts.MaxLen != 100 {
		t.Error("LPosOptions struct not working correctly")
	}
}

func TestBLPop(t *testing.T) {
	r.Del("key")
	result, err := r.BLPop([]string{"key"}, 1)
	if err != nil {
		t.Error(err)
	}
	if len(result) != 0 {
		t.Fail()
	}
	r.LPush("key", "value")
	result, err = r.BLPop([]string{"key"}, 0)
	if err != nil {
		t.Error(err)
	}
	if len(result) == 0 {
		t.Fail()
	}
	if result[0] != "key" || result[1] != "value" {
		t.Fail()
	}
}

func TestBRPop(t *testing.T) {
	r.Del("key")
	result, err := r.BRPop([]string{"key"}, 1)
	if err != nil {
		t.Error(err)
	}
	if len(result) != 0 {
		t.Fail()
	}
	r.RPush("key", "value")
	result, _ = r.BRPop([]string{"key"}, 1)
	if result == nil {
		t.Fail()
	}
	if result[0] != "key" || result[1] != "value" {
		t.Fail()
	}
}

func TestBRPopLPush(t *testing.T) {
	r.Del("key", "key1")
	result, err := r.BRPopLPush("key", "key1", 1)
	if err != nil {
		t.Error(err)
	} else if result != nil {
		t.Fail()
	}
	r.RPush("key", "value")
	result, _ = r.BRPopLPush("key", "key1", 1)
	if result == nil {
		t.Fail()
	}
}

func TestLIndex(t *testing.T) {
	r.Del("key")
	r.LPush("key", "world", "hello")
	if value, err := r.LIndex("key", 0); err != nil {
		t.Error(err)
	} else if string(value) != "hello" {
		t.Fail()
	}
	if value, err := r.LIndex("key", -1); err != nil {
		t.Error(err)
	} else if string(value) != "world" {
		t.Fail()
	}
	if value, err := r.LIndex("key", 3); err != nil {
		t.Error(err)
	} else if value != nil {
		t.Fail()
	}
}

func TestLInsert(t *testing.T) {
	r.Del("key")
	r.RPush("key", "hello", "world")
	if n, err := r.LInsert("key", "before", "world", "three"); err != nil {
		t.Error(err)
	} else if n != 3 {
		t.Fail()
	}
}

func TestLLen(t *testing.T) {
	r.Del("key")
	r.RPush("key", "hello", "world")
	if n, err := r.LLen("key"); err != nil {
		t.Error(err)
	} else if n != 2 {
		t.Fail()
	}
}

func TestLPop(t *testing.T) {
	r.Del("key")
	r.RPush("key", "one", "two", "three")
	if value, err := r.LPop("key"); err != nil {
		t.Error(err)
	} else if string(value) != "one" {
		t.Fail()
	}
	r.Del("key")
	if value, _ := r.LPop("key"); value != nil {
		t.Fail()
	}
}

func TestLPush(t *testing.T) {
	r.Del("key")
	if n, err := r.LPush("key", "value"); err != nil {
		t.Error(err)
	} else if n != 1 {
		t.Fail()
	}
}

func BenchmarkLPush(b *testing.B) {
	r.Del("key")
	for i := 0; i < b.N; i++ {
		r.LPush("key", "value")
	}
}

func TestLPushx(t *testing.T) {
	r.Del("key")
	if n, err := r.LPushx("key", "value"); err != nil {
		t.Error(err)
	} else if n != 0 {
		t.Fail()
	}
	r.LPush("key", "value")
	if n, _ := r.LPushx("key", "value"); n != 2 {
		t.Fail()
	}
}

func TestLRange(t *testing.T) {
	r.Del("key")
	r.RPush("key", "one", "two", "three")
	if data, err := r.LRange("key", 0, 0); err != nil {
		t.Error(err)
	} else if len(data) != 1 {
		t.Fail()
	} else if data[0] != "one" {
		t.Fail()
	}
	if data, _ := r.LRange("key", 5, 10); len(data) != 0 {
		t.Fail()
	}
}

func BenchmarkLRange(b *testing.B) {
	r.Del("key")
	r.RPush("key", "one", "two", "three")
	for i := 0; i < b.N; i++ {
		r.LRange("key", 0, 10)
	}
}

func TestLRem(t *testing.T) {
	r.Del("key")
	r.RPush("key", "hello", "hello", "foo", "hello")
	if n, err := r.LRem("key", -2, "hello"); err != nil {
		t.Error(err)
	} else if n != 2 {
		t.Fail()
	}
}

func TestLSet(t *testing.T) {
	r.Del("key")
	r.RPush("key", "value")
	if err := r.LSet("key", 0, "value2"); err != nil {
		t.Error(err)
	}
	if err := r.LSet("key", 1, "value"); err == nil {
		t.Fail()
	}
}

func TestLTrim(t *testing.T) {
	r.Del("key")
	r.RPush("key", "one", "two", "three")
	if err := r.LTrim("key", 1, -1); err != nil {
		t.Error(err)
	}
}

func TestRPop(t *testing.T) {
	r.Del("key")
	r.RPush("key", "one", "two", "three")
	if value, err := r.RPop("key"); err != nil {
		t.Error(err)
	} else if string(value) != "three" {
		t.Fail()
	}
	r.Del("key")
	if value, _ := r.RPop("key"); value != nil {
		t.Fail()
	}
}

func TestRPopLPush(t *testing.T) {
	r.Del("key")
	if value, err := r.RPopLPush("key", "key"); err != nil {
		t.Error(err)
	} else if value != nil {
		t.Fail()
	}
}

func TestRPush(t *testing.T) {
	r.Del("key")
	if n, err := r.RPush("key", "one", "two"); err != nil {
		t.Error(err)
	} else if n != 2 {
		t.Fail()
	}
}

func TestRPushx(t *testing.T) {
	r.Del("key")
	if n, err := r.RPushx("key", "value"); err != nil {
		t.Error(err)
	} else if n != 0 {
		t.Fail()
	}
}

// Tests for new Phase 1 commands
func TestLMove(t *testing.T) {
	r.Del("source", "dest")
	r.RPush("source", "one", "two", "three")
	
	value, err := r.LMove("source", "dest", ListDirectionRight, ListDirectionLeft)
	if err != nil {
		t.Error(err)
	}
	if value != "three" {
		t.Error("Expected 'three', got", value)
	}
}

func TestBLMove(t *testing.T) {
	r.Del("source", "dest")
	r.RPush("source", "one", "two")
	
	value, err := r.BLMove("source", "dest", ListDirectionLeft, ListDirectionRight, 1)
	if err != nil {
		t.Error(err)
	}
	if value != "one" {
		t.Error("Expected 'one', got", value)
	}
}

func TestLPos(t *testing.T) {
	r.Del("key")
	r.RPush("key", "a", "b", "c", "b")
	
	pos, err := r.LPos("key", "b")
	if err != nil {
		t.Error(err)
	}
	if pos != 1 {
		t.Error("Expected position 1, got", pos)
	}
}

func TestLPosWithOptions(t *testing.T) {
	r.Del("key")
	r.RPush("key", "a", "b", "c", "b", "b")
	
	opts := LPosOptions{
		Rank:   2,
		Count:  2,
		MaxLen: 10,
	}
	
	positions, err := r.LPosWithOptions("key", "b", opts)
	if err != nil {
		t.Error(err)
	}
	if len(positions) == 0 {
		t.Error("Expected at least one position")
	}
}

func TestLMPop(t *testing.T) {
	r.Del("list1", "list2")
	r.RPush("list1", "one", "two")
	r.RPush("list2", "three", "four")
	
	result, err := r.LMPop([]string{"list1", "list2"}, ListDirectionLeft)
	if err != nil {
		t.Error(err)
	}
	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}
}
