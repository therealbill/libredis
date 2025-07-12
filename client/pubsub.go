package client

import (
	"errors"
	"strconv"
	"strings"
)

// Publish posts a message to the given channel.
// Integer reply: the number of clients that received the message.
func (r *Redis) Publish(channel, message string) (int64, error) {
	rp, err := r.ExecuteCommand("PUBLISH", channel, message)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// PubSub doc: http://redis.io/topics/pubsub
type PubSub struct {
	redis *Redis
	conn  *connection

	Patterns map[string]bool
	Channels map[string]bool
}

// GetName returns the address/name of the sentinel we are connected to
func (p *PubSub) GetName() string {
	return p.redis.GetName()
}

// PubSub new a PubSub from *redis.
func (r *Redis) PubSub() (*PubSub, error) {
	c, err := r.pool.Get()
	if err != nil {
		return nil, err
	}
	return &PubSub{
		redis:    r,
		conn:     c,
		Patterns: make(map[string]bool),
		Channels: make(map[string]bool),
	}, nil
}

// Close closes current pubsub command.
func (p *PubSub) Close() error {
	return p.conn.Conn.Close()
}

// Receive returns the reply of pubsub command.
// A message is a Multi-bulk reply with three elements.
// The first element is the kind of message:
// 1) subscribe: means that we successfully subscribed to the channel given as the second element in the reply.
// The third argument represents the number of channels we are currently subscribed to.
// 2) unsubscribe: means that we successfully unsubscribed from the channel given as second element in the reply.
// third argument represents the number of channels we are currently subscribed to.
// When the last argument is zero, we are no longer subscribed to any channel,
// and the client can issue any kind of Redis command as we are outside the Pub/Sub state.
// 3) message: it is a message received as result of a PUBLISH command issued by another client.
// The second element is the name of the originating channel, and the third argument is the actual message payload.
func (p *PubSub) Receive() ([]string, error) {
	rp, err := p.conn.RecvReply()
	if err != nil {
		return nil, err
	}
	command, err := rp.Multi[0].StringValue()
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(command) {
	case "psubscribe", "punsubscribe":
		pattern, err := rp.Multi[1].StringValue()
		if err != nil {
			return nil, err
		}
		if command == "psubscribe" {
			p.Patterns[pattern] = true
		} else {
			delete(p.Patterns, pattern)
		}
		number, err := rp.Multi[2].IntegerValue()
		if err != nil {
			return nil, err
		}
		return []string{command, pattern, strconv.FormatInt(number, 10)}, nil
	case "subscribe", "unsubscribe":
		channel, err := rp.Multi[1].StringValue()
		if err != nil {
			return nil, err
		}
		if command == "subscribe" {
			p.Channels[channel] = true
		} else {
			delete(p.Channels, channel)
		}
		number, err := rp.Multi[2].IntegerValue()
		if err != nil {
			return nil, err
		}
		return []string{command, channel, strconv.FormatInt(number, 10)}, nil
	case "pmessage":
		pattern, err := rp.Multi[1].StringValue()
		if err != nil {
			return nil, err
		}
		channel, err := rp.Multi[2].StringValue()
		if err != nil {
			return nil, err
		}
		message, err := rp.Multi[3].StringValue()
		if err != nil {
			return nil, err
		}
		return []string{command, pattern, channel, message}, nil
	case "message":
		channel, err := rp.Multi[1].StringValue()
		if err != nil {
			return nil, err
		}
		message, err := rp.Multi[2].StringValue()
		if err != nil {
			return nil, err
		}
		return []string{command, channel, message}, nil
	}
	return nil, errors.New("pubsub protocol error")
}

// Subscribe channel [channel ...]
func (p *PubSub) Subscribe(channels ...string) error {
	args := packArgs("SUBSCRIBE", channels)
	return p.conn.SendCommand(args...)
}

// PSubscribe pattern [pattern ...]
func (p *PubSub) PSubscribe(patterns ...string) error {
	args := packArgs("PSUBSCRIBE", patterns)
	return p.conn.SendCommand(args...)
}

// UnSubscribe [channel [channel ...]]
func (p *PubSub) UnSubscribe(channels ...string) error {
	args := packArgs("UNSUBSCRIBE", channels)
	return p.conn.SendCommand(args...)
}

// PUnSubscribe [pattern [pattern ...]]
func (p *PubSub) PUnSubscribe(patterns ...string) error {
	args := packArgs("PUNSUBSCRIBE", patterns)
	return p.conn.SendCommand(args...)
}

// Enhanced Pub/Sub Information Commands

// PubSubChannelInfo represents channel subscription information
type PubSubChannelInfo struct {
	Channel     string
	Subscribers int64
}

// PUBSUB CHANNELS [pattern]
// PubSubChannels returns active channels.
func (r *Redis) PubSubChannels() ([]string, error) {
	args := packArgs("PUBSUB", "CHANNELS")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// PubSubChannelsWithPattern returns active channels matching pattern.
func (r *Redis) PubSubChannelsWithPattern(pattern string) ([]string, error) {
	args := packArgs("PUBSUB", "CHANNELS", pattern)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// PUBSUB NUMSUB [channel ...]
// PubSubNumSub returns subscriber counts for specified channels.
func (r *Redis) PubSubNumSub(channels ...string) ([]PubSubChannelInfo, error) {
	args := packArgs("PUBSUB", "NUMSUB", channels)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	if rp.Multi == nil {
		return nil, nil
	}

	// Response format is: [channel1, count1, channel2, count2, ...]
	var channelInfos []PubSubChannelInfo
	for i := 0; i < len(rp.Multi)-1; i += 2 {
		channel, err := rp.Multi[i].StringValue()
		if err != nil {
			continue
		}
		count, err := rp.Multi[i+1].IntegerValue()
		if err != nil {
			continue
		}
		channelInfos = append(channelInfos, PubSubChannelInfo{
			Channel:     channel,
			Subscribers: count,
		})
	}

	return channelInfos, nil
}

// PUBSUB NUMPAT
// PubSubNumPat returns the number of pattern subscriptions.
func (r *Redis) PubSubNumPat() (int64, error) {
	args := packArgs("PUBSUB", "NUMPAT")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// Sharded Pub/Sub (Redis 7.0+)

// ShardedPubSubMessage represents a sharded pub/sub message
type ShardedPubSubMessage struct {
	Type         string // subscribe, unsubscribe, message
	ShardChannel string
	Message      string
	Count        int64 // subscription count
}

// ShardedPubSub represents a sharded pub/sub connection
type ShardedPubSub struct {
	redis *Redis
	conn  *connection

	ShardChannels map[string]bool
}

// GetName returns the address/name of the redis instance we are connected to
func (sp *ShardedPubSub) GetName() string {
	return sp.redis.GetName()
}

// SPUBLISH shardchannel message
// SPublish publishes a message to a sharded channel.
func (r *Redis) SPublish(shardchannel, message string) (int64, error) {
	args := packArgs("SPUBLISH", shardchannel, message)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// ShardedPubSub creates a new sharded pub/sub connection.
func (r *Redis) ShardedPubSub() (*ShardedPubSub, error) {
	c, err := r.pool.Get()
	if err != nil {
		return nil, err
	}
	return &ShardedPubSub{
		redis:         r,
		conn:          c,
		ShardChannels: make(map[string]bool),
	}, nil
}

// Close closes the sharded pub/sub connection.
func (sp *ShardedPubSub) Close() error {
	return sp.conn.Conn.Close()
}

// SSUBSCRIBE shardchannel [shardchannel ...]
// SSubscribe subscribes to one or more sharded channels.
func (sp *ShardedPubSub) SSubscribe(shardchannels ...string) error {
	args := packArgs("SSUBSCRIBE", shardchannels)
	return sp.conn.SendCommand(args...)
}

// SUNSUBSCRIBE [shardchannel ...]
// SUnSubscribe unsubscribes from sharded channels.
func (sp *ShardedPubSub) SUnSubscribe(shardchannels ...string) error {
	args := packArgs("SUNSUBSCRIBE", shardchannels)
	return sp.conn.SendCommand(args...)
}

// Receive receives messages from sharded subscriptions.
func (sp *ShardedPubSub) Receive() (ShardedPubSubMessage, error) {
	rp, err := sp.conn.RecvReply()
	if err != nil {
		return ShardedPubSubMessage{}, err
	}

	if len(rp.Multi) < 3 {
		return ShardedPubSubMessage{}, errors.New("invalid sharded pubsub message format")
	}

	command, err := rp.Multi[0].StringValue()
	if err != nil {
		return ShardedPubSubMessage{}, err
	}

	msg := ShardedPubSubMessage{Type: strings.ToLower(command)}

	switch strings.ToLower(command) {
	case "ssubscribe", "sunsubscribe":
		shardchannel, err := rp.Multi[1].StringValue()
		if err != nil {
			return ShardedPubSubMessage{}, err
		}
		count, err := rp.Multi[2].IntegerValue()
		if err != nil {
			return ShardedPubSubMessage{}, err
		}

		if command == "ssubscribe" {
			sp.ShardChannels[shardchannel] = true
		} else {
			delete(sp.ShardChannels, shardchannel)
		}

		msg.ShardChannel = shardchannel
		msg.Count = count

	case "smessage":
		if len(rp.Multi) < 3 {
			return ShardedPubSubMessage{}, errors.New("invalid smessage format")
		}

		shardchannel, err := rp.Multi[1].StringValue()
		if err != nil {
			return ShardedPubSubMessage{}, err
		}
		message, err := rp.Multi[2].StringValue()
		if err != nil {
			return ShardedPubSubMessage{}, err
		}

		msg.ShardChannel = shardchannel
		msg.Message = message

	default:
		return ShardedPubSubMessage{}, errors.New("unknown sharded pubsub message type: " + command)
	}

	return msg, nil
}
