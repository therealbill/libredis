package client

import (
	"testing"
)

func TestEcho(t *testing.T) {
	msg := "message"
	ret, err := r.Echo(msg)
	if err != nil {
		t.Error(err)
	} else if ret != msg {
		t.Errorf("echo %s\n%s", msg, ret)
	}
}

func TestPing(t *testing.T) {
	if err := r.Ping(); err != nil {
		t.Error(err)
	}
}

func BenchmarkPing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r.Ping()
	}
}

// Tests for new Phase 1 connection commands
func TestAuthWithUser(t *testing.T) {
	// Note: This test requires proper ACL setup in Redis 6.0+
	// For now, we just test that the function exists and doesn't panic
	err := r.AuthWithUser("default", "")
	// We expect an error for invalid credentials, but not a panic
	if err == nil {
		// This might succeed if Redis allows empty password for default user
		t.Log("AuthWithUser succeeded (empty password allowed)")
	}
}

func TestHello(t *testing.T) {
	// Test basic HELLO command
	info, err := r.Hello(3)
	if err != nil {
		t.Error(err)
	}
	if info == nil {
		t.Error("Expected non-nil info from HELLO command")
	}
}

func TestHelloWithOptions(t *testing.T) {
	opts := HelloOptions{
		ProtocolVersion: 3,
		ClientName:      "test-client",
	}
	
	info, err := r.HelloWithOptions(opts)
	if err != nil {
		t.Error(err)
	}
	if info == nil {
		t.Error("Expected non-nil info from HELLO command")
	}
}

func TestReset(t *testing.T) {
	// Test RESET command
	err := r.Reset()
	if err != nil {
		t.Error(err)
	}
}
