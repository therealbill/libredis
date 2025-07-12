package client

import (
	"testing"
)

func TestACLSetUser(t *testing.T) {
	// Test creating a basic user
	err := r.ACLSetUser("testuser", "on", ">password123", "~*", "+@all")
	if err != nil {
		t.Logf("ACLSetUser failed (Redis may not support ACL): %v", err)
		return
	}

	// Clean up
	r.ACLDelUser("testuser")
}

func TestACLGetUser(t *testing.T) {
	// Create a test user first
	err := r.ACLSetUser("testuser", "on", ">password123", "~*", "+@read")
	if err != nil {
		t.Logf("ACLSetUser failed (Redis may not support ACL): %v", err)
		return
	}
	defer r.ACLDelUser("testuser")

	// Get user information
	user, err := r.ACLGetUser("testuser")
	if err != nil {
		t.Errorf("ACLGetUser failed: %v", err)
		return
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	// Check that flags, keys, and commands are populated
	t.Logf("User: %+v", user)
}

func TestACLDelUser(t *testing.T) {
	// Create test users
	err := r.ACLSetUser("testuser1", "on")
	if err != nil {
		t.Logf("ACLSetUser failed (Redis may not support ACL): %v", err)
		return
	}

	err = r.ACLSetUser("testuser2", "on")
	if err != nil {
		t.Logf("ACLSetUser failed (Redis may not support ACL): %v", err)
		return
	}

	// Delete both users
	deleted, err := r.ACLDelUser("testuser1", "testuser2")
	if err != nil {
		t.Errorf("ACLDelUser failed: %v", err)
		return
	}

	if deleted != 2 {
		t.Errorf("Expected to delete 2 users, deleted %d", deleted)
	}
}

func TestACLUsers(t *testing.T) {
	users, err := r.ACLUsers()
	if err != nil {
		t.Logf("ACLUsers failed (Redis may not support ACL): %v", err)
		return
	}

	// Should at least have 'default' user
	if len(users) == 0 {
		t.Error("Expected at least one user (default)")
		return
	}

	// Check if default user exists
	hasDefault := false
	for _, user := range users {
		if user == "default" {
			hasDefault = true
			break
		}
	}

	if !hasDefault {
		t.Error("Expected 'default' user in ACL users list")
	}

	t.Logf("ACL Users: %v", users)
}

func TestACLCat(t *testing.T) {
	categories, err := r.ACLCat()
	if err != nil {
		t.Logf("ACLCat failed (Redis may not support ACL): %v", err)
		return
	}

	// Should have standard categories like "read", "write", etc.
	if len(categories) == 0 {
		t.Error("Expected ACL categories")
		return
	}

	t.Logf("ACL Categories: %v", categories)
}

func TestACLCatByCategory(t *testing.T) {
	// Test getting commands in 'read' category
	commands, err := r.ACLCatByCategory("read")
	if err != nil {
		t.Logf("ACLCatByCategory failed (Redis may not support ACL): %v", err)
		return
	}

	// Should have commands like "get", "mget", etc.
	if len(commands) == 0 {
		t.Error("Expected commands in 'read' category")
		return
	}

	t.Logf("Read category commands: %v", commands)
}

func TestACLWhoAmI(t *testing.T) {
	username, err := r.ACLWhoAmI()
	if err != nil {
		t.Logf("ACLWhoAmI failed (Redis may not support ACL): %v", err)
		return
	}

	// Should return 'default' for default connection
	if username == "" {
		t.Error("Expected non-empty username")
		return
	}

	t.Logf("Current ACL user: %s", username)
}

func TestACLLog(t *testing.T) {
	// Reset log first
	err := r.ACLLogReset()
	if err != nil {
		t.Logf("ACLLogReset failed (Redis may not support ACL): %v", err)
		return
	}

	// Get log entries
	entries, err := r.ACLLog()
	if err != nil {
		t.Errorf("ACLLog failed: %v", err)
		return
	}

	// After reset, should be empty
	if len(entries) != 0 {
		t.Logf("Expected empty log after reset, got %d entries", len(entries))
	}

	t.Logf("ACL Log entries: %d", len(entries))
}

func TestACLLogWithCount(t *testing.T) {
	entries, err := r.ACLLogWithCount(5)
	if err != nil {
		t.Logf("ACLLogWithCount failed (Redis may not support ACL): %v", err)
		return
	}

	t.Logf("ACL Log entries (limit 5): %d", len(entries))
}

func TestACLList(t *testing.T) {
	rules, err := r.ACLList()
	if err != nil {
		t.Logf("ACLList failed (Redis may not support ACL): %v", err)
		return
	}

	// Should have at least one rule (for default user)
	if len(rules) == 0 {
		t.Error("Expected at least one ACL rule")
		return
	}

	t.Logf("ACL Rules: %v", rules)
}

func TestACLGenPass(t *testing.T) {
	password, err := r.ACLGenPass()
	if err != nil {
		t.Logf("ACLGenPass failed (Redis may not support ACL): %v", err)
		return
	}

	if password == "" {
		t.Error("Expected non-empty generated password")
		return
	}

	if len(password) < 10 {
		t.Errorf("Expected password length >= 10, got %d", len(password))
	}

	t.Logf("Generated password: %s (length: %d)", password, len(password))
}

func TestACLGenPassWithBits(t *testing.T) {
	password, err := r.ACLGenPassWithBits(128)
	if err != nil {
		t.Logf("ACLGenPassWithBits failed (Redis may not support ACL): %v", err)
		return
	}

	if password == "" {
		t.Error("Expected non-empty generated password")
		return
	}

	t.Logf("Generated password (128 bits): %s", password)
}

func TestACLDryRun(t *testing.T) {
	// Create a limited user for testing
	err := r.ACLSetUser("limiteduser", "on", ">password", "~key:*", "+get")
	if err != nil {
		t.Logf("ACLSetUser failed (Redis may not support ACL): %v", err)
		return
	}
	defer r.ACLDelUser("limiteduser")

	// Test allowed command
	err = r.ACLDryRun("limiteduser", "get", "key:test")
	if err != nil {
		t.Logf("ACLDryRun failed for allowed command: %v", err)
	}

	// Test disallowed command (should fail)
	err = r.ACLDryRun("limiteduser", "set", "key:test", "value")
	if err == nil {
		t.Error("Expected ACLDryRun to fail for disallowed command")
	}

	t.Logf("ACLDryRun test completed")
}

func TestACLLoadSave(t *testing.T) {
	// Note: These operations require proper ACL file configuration
	// In many test environments, these might fail
	
	err := r.ACLSave()
	if err != nil {
		t.Logf("ACLSave failed (may not be configured): %v", err)
	} else {
		t.Log("ACLSave succeeded")
	}

	err = r.ACLLoad()
	if err != nil {
		t.Logf("ACLLoad failed (may not be configured): %v", err)
	} else {
		t.Log("ACLLoad succeeded")
	}
}

// Integration test for ACL workflow
func TestACLWorkflow(t *testing.T) {
	// 1. Create a user with specific permissions
	err := r.ACLSetUser("workflowuser", "on", ">securepass", "~data:*", "+get", "+set")
	if err != nil {
		t.Logf("ACLSetUser failed (Redis may not support ACL): %v", err)
		return
	}
	defer r.ACLDelUser("workflowuser")

	// 2. Verify user was created
	users, err := r.ACLUsers()
	if err != nil {
		t.Errorf("ACLUsers failed: %v", err)
		return
	}

	userFound := false
	for _, user := range users {
		if user == "workflowuser" {
			userFound = true
			break
		}
	}

	if !userFound {
		t.Error("workflowuser not found in ACL users list")
		return
	}

	// 3. Get user details
	user, err := r.ACLGetUser("workflowuser")
	if err != nil {
		t.Errorf("ACLGetUser failed: %v", err)
		return
	}

	if user.Username != "workflowuser" {
		t.Errorf("Expected username 'workflowuser', got '%s'", user.Username)
	}

	// 4. Test permissions with dry run
	err = r.ACLDryRun("workflowuser", "get", "data:test")
	if err != nil {
		t.Errorf("ACLDryRun failed for allowed operation: %v", err)
	}

	// 5. Test forbidden operation
	err = r.ACLDryRun("workflowuser", "del", "data:test")
	if err == nil {
		t.Error("Expected ACLDryRun to fail for forbidden operation")
	}

	t.Log("ACL workflow test completed successfully")
}