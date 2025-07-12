package client

import (
	"testing"
	"time"
)

func TestDel(t *testing.T) {
	r.Set("key", "value")
	if n, err := r.Del("key"); err != nil {
		t.Error(err)
	} else if n != 1 {
		t.Fail()
	}
}

func TestDump(t *testing.T) {
	r.Set("key", "value")
	data, err := r.Dump("key")
	if err != nil {
		t.Error(err)
	}
	if data == nil || len(data) == 0 {
		t.Fail()
	}
}

func TestExists(t *testing.T) {
	r.Del("key")
	b, err := r.Exists("key")
	if err != nil {
		t.Error(err)
	}
	if b {
		t.Fail()
	}
}

func TestExpire(t *testing.T) {
	r.Set("key", "value")
	if b, err := r.Expire("key", 10); err != nil {
		t.Error(err)
	} else if !b {
		t.Fail()
	}
	if n, err := r.TTL("key"); err != nil {
		t.Error(err)
	} else if n != 10 {
		t.Fail()
	}
}

func TestExpireAt(t *testing.T) {
	r.Set("key", "value")
	if b, err := r.ExpireAt("key", time.Now().Add(10*time.Second).Unix()); err != nil {
		t.Error(err)
	} else if !b {
		t.Fail()
	}
	if n, err := r.TTL("key"); err != nil {
		t.Error(err)
	} else if n < 0 {
		t.Fail()
	}
}

func TestKeys(t *testing.T) {
	r.FlushDB()
	keys, err := r.Keys("*")
	if err != nil {
		t.Error(err)
	}
	if len(keys) != 0 {
		t.Fail()
	}
	r.Set("key", "value")
	keys, err = r.Keys("*")
	if err != nil {
		t.Error(err)
	}
	if len(keys) != 1 || keys[0] != "key" {
		t.Fail()
	}
}

func TestMove(t *testing.T) {
	r.Set("key", "value")
	if _, err := r.Move("key", db+1); err != nil {
		t.Error(err)
	}
}

func TestObject(t *testing.T) {
	r.Del("key")
	r.LPush("key", "hello world")
	if rp, err := r.Object("refcount", "key"); err != nil {
		t.Error(err)
	} else if rp.Type != IntegerReply {
		t.Fail()
	}
	if rp, err := r.Object("encoding", "key"); err != nil {
		t.Error(err)
	} else if rp.Type != BulkReply {
		t.Fail()
	}
	if rp, err := r.Object("idletime", "key"); err != nil {
		t.Error(err)
	} else if rp.Type != IntegerReply {
		t.Fail()
	}
}

func TestPersist(t *testing.T) {
	r.Set("key", "value")
	r.Expire("key", 500)
	if n, _ := r.TTL("key"); n < 0 {
		t.Fail()
	}
	if b, err := r.Persist("key"); err != nil {
		t.Error(err)
	} else if !b {
		t.Fail()
	}
	if n, _ := r.TTL("key"); n > 0 {
		t.Fail()
	}
}

func TestPExpire(t *testing.T) {
	r.Set("key", "value")
	if b, err := r.PExpire("key", 100); err != nil {
		t.Error(err)
	} else if !b {
		t.Fail()
	}
}

func TestPExpireAt(t *testing.T) {
	r.Set("key", "value")
	if b, err := r.PExpireAt("key", time.Now().Add(500*time.Second).Unix()*1000); err != nil {
		t.Error(err)
	} else if !b {
		t.Fail()
	}
}

func TestPTTL(t *testing.T) {
	r.Set("key", "value")
	r.PExpire("key", 1000)
	if n, err := r.PTTL("key"); err != nil {
		t.Error(err)
	} else if n < 0 {
		t.Fail()
	}
}

func TestRandomKey(t *testing.T) {
	r.FlushDB()
	key, err := r.RandomKey()
	if err != nil {
		t.Error(err)
	}
	if key != nil {
		t.Fail()
	}
	r.Set("key", "value")
	key, _ = r.RandomKey()
	if string(key) != "key" {
		t.Fail()
	}
}

func TestRename(t *testing.T) {
	r.Set("key", "value")
	if err := r.Rename("key", "newkey"); err != nil {
		t.Error(err)
	}
	b, _ := r.Exists("key")
	if b {
		t.Fail()
	}
	v, _ := r.Get("newkey")
	if string(v) != "value" {
		t.Fail()
	}
}

func TestRenamenx(t *testing.T) {
	r.Set("key", "value")
	r.Set("newkey", "value")
	if b, err := r.Renamenx("key", "newkey"); err != nil {
		t.Error(err)
	} else if b {
		t.Fail()
	}
}

func TestRestore(t *testing.T) {
	r.Set("key", "value")
	data, _ := r.Dump("key")
	r.Del("key")
	if err := r.Restore("key", 0, string(data)); err != nil {
		t.Error(err)
	}
}

func TestTTL(t *testing.T) {
	r.Set("key", "value")
	r.Expire("key", 100)
	n, err := r.TTL("key")
	if err != nil {
		t.Error(err)
	}
	if n < 0 {
		t.Fail()
	}
	r.Persist("key")
	n, _ = r.TTL("key")
	if n > 0 {
		t.Fail()
	}
}

func TestType(t *testing.T) {
	r.Set("key", "value")
	ty, err := r.Type("key")
	if err != nil {
		t.Error(err)
	}
	if ty != "string" {
		t.Fail()
	}
}

func TestScan(t *testing.T) {
	r.FlushDB()
	cursor, list, err := r.Scan(0, "", 0)
	if err != nil {
		t.Error(err)
	} else if len(list) != 0 {
		t.Fail()
	} else if cursor != 0 {
		t.Fail()
	}
}

// Tests for new Phase 1 key commands
func TestCopy(t *testing.T) {
	r.Del("source", "destination")
	r.Set("source", "hello")
	
	copied, err := r.Copy("source", "destination")
	if err != nil {
		t.Error(err)
	}
	if !copied {
		t.Error("Expected copy to succeed")
	}
	
	// Verify destination has the value
	value, err := r.Get("destination")
	if err != nil {
		t.Error(err)
	}
	if string(value) != "hello" {
		t.Error("Expected 'hello', got", string(value))
	}
}

func TestCopyWithOptions(t *testing.T) {
	r.Del("source", "destination")
	r.Set("source", "hello")
	r.Set("destination", "existing")
	
	opts := CopyOptions{
		DestinationDB: 0,
		Replace:       true,
	}
	
	copied, err := r.CopyWithOptions("source", "destination", opts)
	if err != nil {
		t.Error(err)
	}
	if !copied {
		t.Error("Expected copy with replace to succeed")
	}
}

func TestTouch(t *testing.T) {
	r.Del("key1", "key2", "key3")
	r.Set("key1", "value1")
	r.Set("key2", "value2")
	
	touched, err := r.Touch("key1", "key2", "key3")
	if err != nil {
		t.Error(err)
	}
	if touched != 2 {
		t.Error("Expected 2 keys touched, got", touched)
	}
}

func TestUnlink(t *testing.T) {
	r.Del("key1", "key2", "key3")
	r.Set("key1", "value1")
	r.Set("key2", "value2")
	
	unlinked, err := r.Unlink("key1", "key2", "key3")
	if err != nil {
		t.Error(err)
	}
	if unlinked != 2 {
		t.Error("Expected 2 keys unlinked, got", unlinked)
	}
}

func TestWait(t *testing.T) {
	// Test WAIT command with 0 replicas and short timeout
	synced, err := r.Wait(0, 100)
	if err != nil {
		t.Error(err)
	}
	if synced < 0 {
		t.Error("Expected non-negative number of synced replicas")
	}
}
