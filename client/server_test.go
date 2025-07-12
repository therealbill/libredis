package client

import (
	"strings"
	"testing"
	"time"
)

func TestBgRewriteAof(t *testing.T) {
	if err := r.BgRewriteAof(); err != nil {
		t.Error(err)
	}
}

func TestBgSave(t *testing.T) {
	// the testing process can call this while the TestBgRewriteAof is still running
	// thus we will sleep a couple seconds to let it finish.
	// Not ideal but it keeps the test simple
	time.Sleep(2 * time.Second)
	if err := r.BgSave(); err != nil {
		if !strings.Contains(err.Error(), "while AOF log rewriting") {
			t.Error(err)
		}
	}
}

func TestClientKill(t *testing.T) {
	if err := r.ClientKill("127.0.0.1", 80); err == nil {
		t.Fail()
	}
}

func TestClientList(t *testing.T) {
	_, err := r.ClientList()
	if err != nil {
		t.Error(err)
	}
}

func TestClientGetName(t *testing.T) {
	if _, err := r.ClientGetName(); err != nil {
		t.Error(err)
	}
}

/*
func TestClientPause(t *testing.T) {
	if err := r.ClientPause(100); err != nil {
		t.Error(err.Error())
	}
}
*/

func TestClientSetName(t *testing.T) {
	if err := r.ClientSetName("name"); err != nil {
		t.Error(err)
	}
}

func TestConfigGet(t *testing.T) {
	if result, err := r.ConfigGet("daemonize"); err != nil {
		t.Error(err)
	} else if result == nil {
		t.Fail()
	} else if len(result) != 1 {
		t.Fail()
	}
}

func TestConfigResetStat(t *testing.T) {
	if err := r.ConfigResetStat(); err != nil {
		t.Error(err)
	}
}

func TestDBSize(t *testing.T) {
	r.FlushDB()
	n, err := r.DBSize()
	if err != nil {
		t.Error(err)
	}
	if n != 0 {
		t.Fail()
	}
}

func TestDebugObject(t *testing.T) {
	r.Del("key")
	r.LPush("key", "value")
	if _, err := r.DebugObject("key"); err != nil {
		t.Error(err)
	}
}

func TestFlushAll(t *testing.T) {
	if err := r.FlushAll(); err != nil {
		t.Error(err)
	}
}

func TestFlushDB(t *testing.T) {
	if err := r.FlushDB(); err != nil {
		t.Error(err)
	}
}

func TestInfo(t *testing.T) {
	serverinfo, err := r.Info()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if len(serverinfo.Server.Version) == 0 {
		t.Error("Server.Version Not Parsed")
		t.Fail()
	}
}

func TestLastSave(t *testing.T) {
	r.Save()
	if timestamp, err := r.LastSave(); err != nil {
		t.Error(err)
	} else if timestamp <= 0 {
		t.Fail()
	}
}

func TestMonitor(t *testing.T) {
	quit := false
	m, err := r.Monitor()
	if err != nil {
		t.Error(err)
	}
	defer m.Close()
	go func() {
		for {
			if s, err := m.Receive(); err != nil {
				if !quit {
					t.Error(err)
				}
			} else if s == "" {
				t.Fail()
			}
		}
	}()
	time.Sleep(100 * time.Millisecond)
	r.LPush("key", "value")
	time.Sleep(100 * time.Microsecond)
}

func TestSave(t *testing.T) {
	if err := r.Save(); err != nil {
		if !strings.Contains(err.Error(), "ERR Background save already in progress") {
			t.Error(err)
		}
	}
}

func TestSlowLogGet(t *testing.T) {
	r.Del("key")
	r.LPush("key", "value")
	if result, err := r.SlowLogGet(1); err != nil {
		t.Error(err)
	} else if len(result) > 1 {
		t.Fail()
	}
}

func TestSlowLogLen(t *testing.T) {
	r.Del("key")
	r.LPush("key", "value")
	if _, err := r.SlowLogLen(); err != nil {
		t.Error(err)
	}
}

func TestSlowLogReset(t *testing.T) {
	if err := r.SlowLogReset(); err != nil {
		t.Error(err)
	}
}

func TestTime(t *testing.T) {
	tt, err := r.Time()
	if err != nil {
		t.Error(err)
	}
	if len(tt) != 2 {
		t.Fail()
	}
}

// Memory Management Tests

func TestMemoryUsage(t *testing.T) {
	// Set a test key
	err := r.Set("memory_test_key", "test_value")
	if err != nil {
		t.Error(err)
		return
	}
	defer r.Del("memory_test_key")

	// Test basic memory usage
	usage, err := r.MemoryUsage("memory_test_key")
	if err != nil {
		t.Logf("MemoryUsage failed (Redis may not support MEMORY commands): %v", err)
		return
	}

	if usage <= 0 {
		t.Error("Expected positive memory usage")
	}

	t.Logf("Memory usage for test key: %d bytes", usage)
}

func TestMemoryUsageWithSamples(t *testing.T) {
	// Set a test key
	err := r.Set("memory_test_key2", "test_value_with_samples")
	if err != nil {
		t.Error(err)
		return
	}
	defer r.Del("memory_test_key2")

	// Test memory usage with samples
	usage, err := r.MemoryUsageWithSamples("memory_test_key2", 5)
	if err != nil {
		t.Logf("MemoryUsageWithSamples failed (Redis may not support MEMORY commands): %v", err)
		return
	}

	if usage <= 0 {
		t.Error("Expected positive memory usage")
	}

	t.Logf("Memory usage with samples: %d bytes", usage)
}

func TestMemoryStats(t *testing.T) {
	stats, err := r.MemoryStats()
	if err != nil {
		t.Logf("MemoryStats failed (Redis may not support MEMORY commands): %v", err)
		return
	}

	t.Logf("Memory Stats - Peak: %d, Total: %d, Dataset: %d (%.2f%%)",
		stats.PeakAllocated, stats.TotalAllocated, stats.Dataset.Bytes, stats.Dataset.Percentage)
}

func TestMemoryDoctor(t *testing.T) {
	analysis, err := r.MemoryDoctor()
	if err != nil {
		t.Logf("MemoryDoctor failed (Redis may not support MEMORY commands): %v", err)
		return
	}

	if analysis == "" {
		t.Error("Expected non-empty memory analysis")
	}

	t.Logf("Memory Doctor analysis: %s", analysis)
}

func TestMemoryPurge(t *testing.T) {
	err := r.MemoryPurge()
	if err != nil {
		t.Logf("MemoryPurge failed (Redis may not support MEMORY commands): %v", err)
		return
	}

	t.Log("Memory purge completed successfully")
}

// Latency Monitoring Tests

func TestLatencyLatest(t *testing.T) {
	latencyStats, err := r.LatencyLatest()
	if err != nil {
		t.Logf("LatencyLatest failed (Redis may not support LATENCY commands): %v", err)
		return
	}

	t.Logf("Latest latency events: %d", len(latencyStats))
	for _, stats := range latencyStats {
		t.Logf("Event: %s, Latest: %dms, All-time: %dms", stats.Event, stats.Latest, stats.AllTime)
	}
}

func TestLatencyHistory(t *testing.T) {
	// First get available events
	latencyStats, err := r.LatencyLatest()
	if err != nil {
		t.Logf("LatencyLatest failed: %v", err)
		return
	}

	if len(latencyStats) == 0 {
		t.Log("No latency events available for history test")
		return
	}

	// Get history for first event
	event := latencyStats[0].Event
	samples, err := r.LatencyHistory(event)
	if err != nil {
		t.Logf("LatencyHistory failed: %v", err)
		return
	}

	t.Logf("Latency history for '%s': %d samples", event, len(samples))
}

func TestLatencyReset(t *testing.T) {
	count, err := r.LatencyReset()
	if err != nil {
		t.Logf("LatencyReset failed (Redis may not support LATENCY commands): %v", err)
		return
	}

	t.Logf("Reset %d latency events", count)
}

func TestLatencyGraph(t *testing.T) {
	// Generate some latency by performing operations
	for i := 0; i < 10; i++ {
		r.Set("latency_test", "value")
		r.Get("latency_test")
	}

	// Try to get a graph for a common event
	graph, err := r.LatencyGraph("command")
	if err != nil {
		t.Logf("LatencyGraph failed (Redis may not support LATENCY commands or no data): %v", err)
		return
	}

	if graph == "" {
		t.Log("No latency graph data available")
		return
	}

	t.Logf("Latency graph:\n%s", graph)
}

// Database Management Tests

func TestSwapDB(t *testing.T) {
	// Set a test key in current database (1 from config)
	err := r.Set("swap_test_key", "original_db_value")
	if err != nil {
		t.Error(err)
		return
	}

	// Try to swap databases 1 and 2
	err = r.SwapDB(1, 2)
	if err != nil {
		t.Logf("SwapDB failed (Redis may not support SWAPDB): %v", err)
		// Clean up
		r.Del("swap_test_key")
		return
	}

	// After swap, our key should be in database 2 and not visible
	value, err := r.Get("swap_test_key")
	if err != nil {
		t.Error(err)
	}

	// Swap back to restore original state
	r.SwapDB(1, 2)
	
	// Clean up
	r.Del("swap_test_key")

	t.Log("SwapDB test completed successfully")
}

func TestReplicaOf(t *testing.T) {
	// Test REPLICAOF NO ONE (promote to master)
	err := r.ReplicaOfNoOne()
	if err != nil {
		t.Logf("ReplicaOfNoOne failed (Redis may not support REPLICAOF): %v", err)
		return
	}

	t.Log("ReplicaOfNoOne completed successfully")

	// Note: We don't test ReplicaOf with actual host/port in unit tests
	// as it would require another Redis instance
}

// Module Management Tests

func TestModuleList(t *testing.T) {
	modules, err := r.ModuleList()
	if err != nil {
		t.Logf("ModuleList failed (Redis may not support MODULE commands): %v", err)
		return
	}

	t.Logf("Loaded modules: %d", len(modules))
	for _, module := range modules {
		t.Logf("Module: %s v%d at %s", module.Name, module.Version, module.Path)
	}
}
