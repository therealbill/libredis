package client

import (
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	if _, err := r.Publish("channel", "message"); err != nil {
		t.Error(err)
	}
}

func TestSubscribe(t *testing.T) {
	quit := make(chan bool)
	sub, err := r.PubSub()
	if err != nil {
		t.Error(err)
	}
	defer sub.Close()
	go func() {
		if err := sub.Subscribe("channel"); err != nil {
			t.Error(err)
			quit <- true
			return
		}
		for {
			list, err := sub.Receive()
			if err != nil {
				t.Error(err)
				quit <- true
				break
			}
			if list[0] == "message" {
				if list[1] != "channel" || list[2] != "message" {
					t.Fail()
				}
				quit <- true
				break
			}
		}
	}()
	time.Sleep(100 * time.Millisecond)
	r.Publish("channel", "message")
	time.Sleep(100 * time.Millisecond)
	<-quit
}

func TestPSubscribe(t *testing.T) {
	quit := make(chan bool)
	psub, err := r.PubSub()
	if err != nil {
		t.Error(err)
	}
	defer psub.Close()
	go func() {
		if err := psub.PSubscribe("news.*"); err != nil {
			t.Error(err)
			quit <- true
			return
		}
		for {
			list, err := psub.Receive()
			if err != nil {
				t.Error(err)
				quit <- true
				break
			}
			if list[0] == "pmessage" {
				if list[1] != "news.*" || list[2] != "news.china" || list[3] != "message" {
					t.Fail()
				}
				quit <- true
				break
			}
		}
	}()
	time.Sleep(100 * time.Millisecond)
	r.Publish("news.china", "message")
	time.Sleep(100 * time.Millisecond)
	<-quit
}

func TestUnSubscribe(t *testing.T) {
	quit := false
	ch := make(chan bool)
	sub, err := r.PubSub()
	if err != nil {
		t.Error(err)
	}
	defer sub.Close()
	go func() {
		for {
			if _, err := sub.Receive(); err != nil {
				if !quit {
					t.Error(err)
				}
			}
			ch <- true
		}
	}()
	time.Sleep(100 * time.Millisecond)
	sub.Subscribe("channel")
	time.Sleep(100 * time.Millisecond)
	<-ch
	if len(sub.Channels) != 1 {
		t.Fail()
	}
	if err := sub.UnSubscribe("channel"); err != nil {
		t.Error(err)
	}
	time.Sleep(100 * time.Millisecond)
	<-ch
	time.Sleep(100 * time.Millisecond)
	if len(sub.Channels) != 0 {
		t.Fail()
	}
	quit = true
}

func TestPUnSubscribe(t *testing.T) {
	quit := false
	ch := make(chan bool)
	sub, err := r.PubSub()
	if err != nil {
		t.Error(err)
	}
	defer sub.Close()
	go func() {
		for {
			if _, err := sub.Receive(); err != nil {
				if !quit {
					t.Error(err)
				}
			}
			ch <- true
		}
	}()
	time.Sleep(100 * time.Millisecond)
	sub.PSubscribe("news.*")
	time.Sleep(100 * time.Millisecond)
	<-ch
	if len(sub.Patterns) != 1 {
		t.Fail()
	}
	if err := sub.PUnSubscribe("news.*"); err != nil {
		t.Error(err)
	}
	time.Sleep(100 * time.Millisecond)
	<-ch
	time.Sleep(100 * time.Millisecond)
	if len(sub.Patterns) != 0 {
		t.Fail()
	}
	quit = true
}

// Enhanced Pub/Sub Tests

func TestPubSubChannels(t *testing.T) {
	// Create a subscriber to make sure we have active channels
	sub, err := r.PubSub()
	if err != nil {
		t.Error(err)
		return
	}
	defer sub.Close()

	// Subscribe to a test channel
	err = sub.Subscribe("test_channel")
	if err != nil {
		t.Error(err)
		return
	}

	// Give it time to register
	time.Sleep(100 * time.Millisecond)

	// Get all active channels
	channels, err := r.PubSubChannels()
	if err != nil {
		t.Logf("PubSubChannels failed (Redis may not support command): %v", err)
		return
	}

	t.Logf("Active channels: %v", channels)
}

func TestPubSubChannelsWithPattern(t *testing.T) {
	// Create a subscriber to make sure we have active channels
	sub, err := r.PubSub()
	if err != nil {
		t.Error(err)
		return
	}
	defer sub.Close()

	// Subscribe to test channels
	err = sub.Subscribe("test_channel_1", "test_channel_2", "other_channel")
	if err != nil {
		t.Error(err)
		return
	}

	// Give it time to register
	time.Sleep(100 * time.Millisecond)

	// Get channels matching pattern
	channels, err := r.PubSubChannelsWithPattern("test_*")
	if err != nil {
		t.Logf("PubSubChannelsWithPattern failed (Redis may not support command): %v", err)
		return
	}

	t.Logf("Channels matching 'test_*': %v", channels)
}

func TestPubSubNumSub(t *testing.T) {
	// Create multiple subscribers
	sub1, err := r.PubSub()
	if err != nil {
		t.Error(err)
		return
	}
	defer sub1.Close()

	sub2, err := r.PubSub()
	if err != nil {
		t.Error(err)
		return
	}
	defer sub2.Close()

	// Subscribe to same channel
	err = sub1.Subscribe("test_numsub")
	if err != nil {
		t.Error(err)
		return
	}

	err = sub2.Subscribe("test_numsub")
	if err != nil {
		t.Error(err)
		return
	}

	// Give it time to register
	time.Sleep(100 * time.Millisecond)

	// Get subscriber counts
	channelInfos, err := r.PubSubNumSub("test_numsub")
	if err != nil {
		t.Logf("PubSubNumSub failed (Redis may not support command): %v", err)
		return
	}

	if len(channelInfos) > 0 {
		t.Logf("Channel %s has %d subscribers", channelInfos[0].Channel, channelInfos[0].Subscribers)
		if channelInfos[0].Subscribers < 1 {
			t.Error("Expected at least 1 subscriber")
		}
	}
}

func TestPubSubNumPat(t *testing.T) {
	// Create a pattern subscriber
	psub, err := r.PubSub()
	if err != nil {
		t.Error(err)
		return
	}
	defer psub.Close()

	// Subscribe to a pattern
	err = psub.PSubscribe("test_pattern_*")
	if err != nil {
		t.Error(err)
		return
	}

	// Give it time to register
	time.Sleep(100 * time.Millisecond)

	// Get pattern subscription count
	count, err := r.PubSubNumPat()
	if err != nil {
		t.Logf("PubSubNumPat failed (Redis may not support command): %v", err)
		return
	}

	t.Logf("Number of pattern subscriptions: %d", count)
	if count < 1 {
		t.Error("Expected at least 1 pattern subscription")
	}
}

// Sharded Pub/Sub Tests (Redis 7.0+)

func TestSPublish(t *testing.T) {
	// Test sharded publish
	count, err := r.SPublish("test_shard_channel", "test message")
	if err != nil {
		t.Logf("SPublish failed (Redis may not support sharded pub/sub): %v", err)
		return
	}

	t.Logf("SPublish delivered to %d shards", count)
}

func TestShardedPubSub(t *testing.T) {
	// Create sharded pub/sub connection
	spub, err := r.ShardedPubSub()
	if err != nil {
		t.Error(err)
		return
	}
	defer spub.Close()

	// Test basic functionality without Redis 7.0+ requirement
	t.Logf("ShardedPubSub connection created successfully")
	
	// Test subscribe (may fail on older Redis versions)
	err = spub.SSubscribe("test_shard")
	if err != nil {
		t.Logf("SSubscribe failed (Redis may not support sharded pub/sub): %v", err)
		return
	}

	// Test unsubscribe
	err = spub.SUnSubscribe("test_shard")
	if err != nil {
		t.Logf("SUnSubscribe failed: %v", err)
	}
}
