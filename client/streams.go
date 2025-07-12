package client

// Stream constants and types
const (
	StreamIDAutoGenerate = "*"
	StreamIDLatest       = "$"
	StreamIDEarliest     = "0-0"
)

// XAddOptions represents options for XADD command
type XAddOptions struct {
	NoMkStream  bool   // NOMKSTREAM option
	MaxLen      int64  // MAXLEN option
	MinID       string // MINID option
	Approximate bool   // ~ modifier for MAXLEN/MINID
	Limit       int64  // LIMIT option
}

// XReadOptions represents options for XREAD command
type XReadOptions struct {
	Count int64 // COUNT option
	Block int64 // BLOCK option (milliseconds)
}

// XRangeOptions represents options for XRANGE/XREVRANGE commands
type XRangeOptions struct {
	Count int64 // COUNT option
}

// StreamEntry represents a single stream entry
type StreamEntry struct {
	ID     string
	Fields map[string]string
}

// StreamMessage represents messages from one stream
type StreamMessage struct {
	Stream  string
	Entries []StreamEntry
}

// XGroupCreateOptions represents options for XGROUP CREATE
type XGroupCreateOptions struct {
	MkStream    bool  // MKSTREAM option
	EntriesRead int64 // ENTRIESREAD option
}

// XReadGroupOptions represents options for XREADGROUP
type XReadGroupOptions struct {
	Count int64 // COUNT option
	Block int64 // BLOCK option (milliseconds)
	NoAck bool  // NOACK option
}

// XClaimOptions represents options for XCLAIM command
type XClaimOptions struct {
	Idle       int64  // IDLE option
	Time       int64  // TIME option
	RetryCount int64  // RETRYCOUNT option
	Force      bool   // FORCE option
	JustID     bool   // JUSTID option
	LastID     string // LASTID option
}

// XPendingOptions represents options for XPENDING command
type XPendingOptions struct {
	Idle     int64  // IDLE option
	Start    string // Start ID
	End      string // End ID
	Count    int64  // Count limit
	Consumer string // Specific consumer
}

// XPendingInfo represents summary of pending messages
type XPendingInfo struct {
	Count     int64
	Lower     string
	Higher    string
	Consumers map[string]int64
}

// XPendingMessage represents a pending message detail
type XPendingMessage struct {
	ID            string
	Consumer      string
	IdleTime      int64
	DeliveryCount int64
}

// Basic Stream Operations

// XADD key [NOMKSTREAM] [MAXLEN|MINID [=|~] threshold [LIMIT count]] *|ID field value [field value ...]
// XAdd appends a new entry to a stream.
func (r *Redis) XAdd(key, id string, fields map[string]string) (string, error) {
	args := []interface{}{"XADD", key, id}
	for field, value := range fields {
		args = append(args, field, value)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// XAddWithOptions appends a new entry with additional options.
func (r *Redis) XAddWithOptions(key, id string, fields map[string]string, opts XAddOptions) (string, error) {
	args := []interface{}{"XADD", key}

	if opts.NoMkStream {
		args = append(args, "NOMKSTREAM")
	}

	if opts.MaxLen > 0 {
		if opts.Approximate {
			args = append(args, "MAXLEN", "~", opts.MaxLen)
		} else {
			args = append(args, "MAXLEN", opts.MaxLen)
		}
		if opts.Limit > 0 {
			args = append(args, "LIMIT", opts.Limit)
		}
	}

	if opts.MinID != "" {
		if opts.Approximate {
			args = append(args, "MINID", "~", opts.MinID)
		} else {
			args = append(args, "MINID", opts.MinID)
		}
		if opts.Limit > 0 {
			args = append(args, "LIMIT", opts.Limit)
		}
	}

	args = append(args, id)
	for field, value := range fields {
		args = append(args, field, value)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// XREAD [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] id [id ...]
// XRead reads data from one or multiple streams.
func (r *Redis) XRead(streams map[string]string) ([]StreamMessage, error) {
	args := []interface{}{"XREAD", "STREAMS"}

	// Add stream keys
	for key := range streams {
		args = append(args, key)
	}

	// Add stream IDs
	for key := range streams {
		args = append(args, streams[key])
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	if rp.Multi == nil {
		return nil, nil
	}

	return parseStreamMessages(rp.Multi)
}

// XReadWithOptions reads from streams with additional options.
func (r *Redis) XReadWithOptions(streams map[string]string, opts XReadOptions) ([]StreamMessage, error) {
	args := []interface{}{"XREAD"}

	if opts.Count > 0 {
		args = append(args, "COUNT", opts.Count)
	}

	if opts.Block >= 0 {
		args = append(args, "BLOCK", opts.Block)
	}

	args = append(args, "STREAMS")

	// Add stream keys
	for key := range streams {
		args = append(args, key)
	}

	// Add stream IDs
	for key := range streams {
		args = append(args, streams[key])
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	if rp.Multi == nil {
		return nil, nil
	}

	return parseStreamMessages(rp.Multi)
}

// XRANGE key start end [COUNT count]
// XRange returns the stream entries matching a given range of IDs.
func (r *Redis) XRange(key, start, end string) ([]StreamEntry, error) {
	args := packArgs("XRANGE", key, start, end)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	return parseStreamEntries(rp.Multi)
}

// XRangeWithOptions returns stream entries with count limit.
func (r *Redis) XRangeWithOptions(key, start, end string, opts XRangeOptions) ([]StreamEntry, error) {
	args := []interface{}{"XRANGE", key, start, end}

	if opts.Count > 0 {
		args = append(args, "COUNT", opts.Count)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	return parseStreamEntries(rp.Multi)
}

// XREVRANGE key end start [COUNT count]
// XRevRange returns stream entries in reverse order.
func (r *Redis) XRevRange(key, end, start string) ([]StreamEntry, error) {
	args := packArgs("XREVRANGE", key, end, start)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	return parseStreamEntries(rp.Multi)
}

// XRevRangeWithOptions returns stream entries in reverse with count limit.
func (r *Redis) XRevRangeWithOptions(key, end, start string, opts XRangeOptions) ([]StreamEntry, error) {
	args := []interface{}{"XREVRANGE", key, end, start}

	if opts.Count > 0 {
		args = append(args, "COUNT", opts.Count)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	return parseStreamEntries(rp.Multi)
}

// XLEN key
// XLen returns the number of entries in a stream.
func (r *Redis) XLen(key string) (int64, error) {
	args := packArgs("XLEN", key)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// XDEL key id [id ...]
// XDel removes specified entries from a stream.
func (r *Redis) XDel(key string, ids ...string) (int64, error) {
	args := packArgs("XDEL", key, ids)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// XTRIM key MAXLEN|MINID [=|~] threshold [LIMIT count]
// XTrim trims a stream to a given size or minimum ID.
func (r *Redis) XTrim(key string, strategy string, threshold string) (int64, error) {
	args := packArgs("XTRIM", key, strategy, threshold)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// XTrimWithOptions trims a stream with additional options.
func (r *Redis) XTrimWithOptions(key string, strategy string, threshold string, opts XAddOptions) (int64, error) {
	args := []interface{}{"XTRIM", key, strategy}

	if opts.Approximate {
		args = append(args, "~")
	}

	args = append(args, threshold)

	if opts.Limit > 0 {
		args = append(args, "LIMIT", opts.Limit)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// Consumer Group Operations

// XGROUP CREATE key groupname id|$ [MKSTREAM] [ENTRIESREAD entries_read]
// XGroupCreate creates a new consumer group.
func (r *Redis) XGroupCreate(key, groupname, id string) error {
	args := packArgs("XGROUP", "CREATE", key, groupname, id)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// XGroupCreateWithOptions creates a consumer group with additional options.
func (r *Redis) XGroupCreateWithOptions(key, groupname, id string, opts XGroupCreateOptions) error {
	args := []interface{}{"XGROUP", "CREATE", key, groupname, id}

	if opts.MkStream {
		args = append(args, "MKSTREAM")
	}

	if opts.EntriesRead > 0 {
		args = append(args, "ENTRIESREAD", opts.EntriesRead)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// XGROUP DESTROY key groupname
// XGroupDestroy destroys a consumer group.
func (r *Redis) XGroupDestroy(key, groupname string) (int64, error) {
	args := packArgs("XGROUP", "DESTROY", key, groupname)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// XGROUP SETID key groupname id|$ [ENTRIESREAD entries_read]
// XGroupSetID sets the consumer group's last delivered ID.
func (r *Redis) XGroupSetID(key, groupname, id string) error {
	args := packArgs("XGROUP", "SETID", key, groupname, id)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// XREADGROUP GROUP group consumer [COUNT count] [BLOCK milliseconds] [NOACK] STREAMS key [key ...] ID [ID ...]
// XReadGroup reads from streams as a consumer group member.
func (r *Redis) XReadGroup(group, consumer string, streams map[string]string) ([]StreamMessage, error) {
	args := []interface{}{"XREADGROUP", "GROUP", group, consumer, "STREAMS"}

	// Add stream keys
	for key := range streams {
		args = append(args, key)
	}

	// Add stream IDs
	for key := range streams {
		args = append(args, streams[key])
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	if rp.Multi == nil {
		return nil, nil
	}

	return parseStreamMessages(rp.Multi)
}

// XReadGroupWithOptions reads from streams as consumer group with options.
func (r *Redis) XReadGroupWithOptions(group, consumer string, streams map[string]string, opts XReadGroupOptions) ([]StreamMessage, error) {
	args := []interface{}{"XREADGROUP", "GROUP", group, consumer}

	if opts.Count > 0 {
		args = append(args, "COUNT", opts.Count)
	}

	if opts.Block >= 0 {
		args = append(args, "BLOCK", opts.Block)
	}

	if opts.NoAck {
		args = append(args, "NOACK")
	}

	args = append(args, "STREAMS")

	// Add stream keys
	for key := range streams {
		args = append(args, key)
	}

	// Add stream IDs
	for key := range streams {
		args = append(args, streams[key])
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	if rp.Multi == nil {
		return nil, nil
	}

	return parseStreamMessages(rp.Multi)
}

// XACK key group id [id ...]
// XAck acknowledges processing of messages by a consumer.
func (r *Redis) XAck(key, group string, ids ...string) (int64, error) {
	args := packArgs("XACK", key, group, ids)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// XCLAIM key group consumer min-idle-time id [id ...] [IDLE ms] [TIME unix-time] [RETRYCOUNT count] [FORCE] [JUSTID] [LASTID id]
// XClaim transfers ownership of pending messages to another consumer.
func (r *Redis) XClaim(key, group, consumer string, minIdleTime int64, ids []string) ([]StreamEntry, error) {
	args := []interface{}{"XCLAIM", key, group, consumer, minIdleTime}
	for _, id := range ids {
		args = append(args, id)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	return parseStreamEntries(rp.Multi)
}

// XClaimWithOptions claims messages with additional options.
func (r *Redis) XClaimWithOptions(key, group, consumer string, minIdleTime int64, ids []string, opts XClaimOptions) ([]StreamEntry, error) {
	args := []interface{}{"XCLAIM", key, group, consumer, minIdleTime}
	for _, id := range ids {
		args = append(args, id)
	}

	if opts.Idle > 0 {
		args = append(args, "IDLE", opts.Idle)
	}

	if opts.Time > 0 {
		args = append(args, "TIME", opts.Time)
	}

	if opts.RetryCount > 0 {
		args = append(args, "RETRYCOUNT", opts.RetryCount)
	}

	if opts.Force {
		args = append(args, "FORCE")
	}

	if opts.JustID {
		args = append(args, "JUSTID")
	}

	if opts.LastID != "" {
		args = append(args, "LASTID", opts.LastID)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	if opts.JustID {
		// Return IDs only for JUSTID option
		entries := make([]StreamEntry, len(rp.Multi))
		for i, reply := range rp.Multi {
			id, _ := reply.StringValue()
			entries[i] = StreamEntry{ID: id, Fields: map[string]string{}}
		}
		return entries, nil
	}

	return parseStreamEntries(rp.Multi)
}

// XPENDING key group [[IDLE min-idle-time] start end count [consumer]]
// XPending returns information about pending messages.
func (r *Redis) XPending(key, group string) (XPendingInfo, error) {
	args := packArgs("XPENDING", key, group)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return XPendingInfo{}, err
	}

	if len(rp.Multi) < 4 {
		return XPendingInfo{}, err
	}

	count, _ := rp.Multi[0].IntegerValue()
	lower, _ := rp.Multi[1].StringValue()
	higher, _ := rp.Multi[2].StringValue()

	consumers := make(map[string]int64)
	if rp.Multi[3].Multi != nil {
		for _, consumerReply := range rp.Multi[3].Multi {
			if len(consumerReply.Multi) >= 2 {
				name, _ := consumerReply.Multi[0].StringValue()
				pending, _ := consumerReply.Multi[1].IntegerValue()
				consumers[name] = pending
			}
		}
	}

	return XPendingInfo{
		Count:     count,
		Lower:     lower,
		Higher:    higher,
		Consumers: consumers,
	}, nil
}

// XPendingWithOptions returns detailed pending message information.
func (r *Redis) XPendingWithOptions(key, group string, opts XPendingOptions) ([]XPendingMessage, error) {
	args := []interface{}{"XPENDING", key, group}

	if opts.Idle > 0 {
		args = append(args, "IDLE", opts.Idle)
	}

	args = append(args, opts.Start, opts.End, opts.Count)

	if opts.Consumer != "" {
		args = append(args, opts.Consumer)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	messages := make([]XPendingMessage, len(rp.Multi))
	for i, msgReply := range rp.Multi {
		if len(msgReply.Multi) >= 4 {
			id, _ := msgReply.Multi[0].StringValue()
			consumer, _ := msgReply.Multi[1].StringValue()
			idleTime, _ := msgReply.Multi[2].IntegerValue()
			deliveryCount, _ := msgReply.Multi[3].IntegerValue()

			messages[i] = XPendingMessage{
				ID:            id,
				Consumer:      consumer,
				IdleTime:      idleTime,
				DeliveryCount: deliveryCount,
			}
		}
	}

	return messages, nil
}

// Stream Information

// XINFO STREAM key [FULL [COUNT count]]
// XInfoStream returns general information about a stream.
func (r *Redis) XInfoStream(key string) (map[string]interface{}, error) {
	args := packArgs("XINFO", "STREAM", key)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if rp.Multi != nil {
		for i := 0; i < len(rp.Multi)-1; i += 2 {
			key, _ := rp.Multi[i].StringValue()
			value := parseInfoValue(rp.Multi[i+1])
			result[key] = value
		}
	}

	return result, nil
}

// XInfoStreamFull returns detailed information about a stream.
func (r *Redis) XInfoStreamFull(key string, count int64) (map[string]interface{}, error) {
	args := []interface{}{"XINFO", "STREAM", key, "FULL"}
	if count > 0 {
		args = append(args, "COUNT", count)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if rp.Multi != nil {
		for i := 0; i < len(rp.Multi)-1; i += 2 {
			key, _ := rp.Multi[i].StringValue()
			value := parseInfoValue(rp.Multi[i+1])
			result[key] = value
		}
	}

	return result, nil
}

// XINFO GROUPS key
// XInfoGroups returns information about consumer groups.
func (r *Redis) XInfoGroups(key string) ([]map[string]interface{}, error) {
	args := packArgs("XINFO", "GROUPS", key)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(rp.Multi))
	for i, groupReply := range rp.Multi {
		group := make(map[string]interface{})
		if groupReply.Multi != nil {
			for j := 0; j < len(groupReply.Multi)-1; j += 2 {
				key, _ := groupReply.Multi[j].StringValue()
				value := parseInfoValue(groupReply.Multi[j+1])
				group[key] = value
			}
		}
		result[i] = group
	}

	return result, nil
}

// XINFO CONSUMERS key groupname
// XInfoConsumers returns information about consumers in a group.
func (r *Redis) XInfoConsumers(key, groupname string) ([]map[string]interface{}, error) {
	args := packArgs("XINFO", "CONSUMERS", key, groupname)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(rp.Multi))
	for i, consumerReply := range rp.Multi {
		consumer := make(map[string]interface{})
		if consumerReply.Multi != nil {
			for j := 0; j < len(consumerReply.Multi)-1; j += 2 {
				key, _ := consumerReply.Multi[j].StringValue()
				value := parseInfoValue(consumerReply.Multi[j+1])
				consumer[key] = value
			}
		}
		result[i] = consumer
	}

	return result, nil
}

// Helper functions

func parseStreamMessages(replies []*Reply) ([]StreamMessage, error) {
	messages := make([]StreamMessage, len(replies))
	for i, streamReply := range replies {
		if len(streamReply.Multi) >= 2 {
			streamName, _ := streamReply.Multi[0].StringValue()
			entries, _ := parseStreamEntries(streamReply.Multi[1].Multi)
			messages[i] = StreamMessage{
				Stream:  streamName,
				Entries: entries,
			}
		}
	}
	return messages, nil
}

func parseStreamEntries(replies []*Reply) ([]StreamEntry, error) {
	entries := make([]StreamEntry, len(replies))
	for i, entryReply := range replies {
		if len(entryReply.Multi) >= 2 {
			id, _ := entryReply.Multi[0].StringValue()
			fields := make(map[string]string)

			if entryReply.Multi[1].Multi != nil {
				fieldReply := entryReply.Multi[1].Multi
				for j := 0; j < len(fieldReply)-1; j += 2 {
					field, _ := fieldReply[j].StringValue()
					value, _ := fieldReply[j+1].StringValue()
					fields[field] = value
				}
			}

			entries[i] = StreamEntry{
				ID:     id,
				Fields: fields,
			}
		}
	}
	return entries, nil
}

func parseInfoValue(reply *Reply) interface{} {
	if reply.Type == IntegerReply {
		return reply.Integer
	}
	if reply.Type == BulkReply {
		return string(reply.Bulk)
	}
	if reply.Type == MultiReply && reply.Multi != nil {
		result := make([]interface{}, len(reply.Multi))
		for i, r := range reply.Multi {
			result[i] = parseInfoValue(r)
		}
		return result
	}
	return nil
}
