package client

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/therealbill/libredis/info"
	"github.com/therealbill/libredis/structures"
)

// BgRewriteAof Instruct Redis to start an Append Only File rewrite process.
// The rewrite will create a small optimized version of the current Append Only File.
func (r *Redis) BgRewriteAof() error {
	_, err := r.ExecuteCommand("BGREWRITEAOF")
	return err
}

// BgSave save the DB in background.
// The OK code is immediately returned.
// Redis forks, the parent continues to serve the clients, the child saves the DB on disk then exits.
// A client my be able to check if the operation succeeded using the LASTSAVE command.
func (r *Redis) BgSave() error {
	_, err := r.ExecuteCommand("BGSAVE")
	return err
}

// ClientKill closes a given client connection identified by ip:port.
// Due to the single-treaded nature of Redis,
// it is not possible to kill a client connection while it is executing a command.
// However, the client will notice the connection has been closed
// only when the next command is sent (and results in network error).
// Status code reply: OK if the connection exists and has been closed
func (r *Redis) ClientKill(ip string, port int) error {
	rp, err := r.ExecuteCommand("CLIENT", "KILL", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// Role() returns the current role name on the server. Requires Redis >= 2.8.12
func (r *Redis) Role() (sl []string, err error) {
	rp, err := r.ExecuteCommand("ROLE")
	if err != nil {
		return sl, err
	}
	return rp.Multi[1].ListValue()
}

// RoleName() returns the current role name on the server. Requires Redis >= 2.8.12
// This is a shorthand for those who only want the name of the role rather than
// all the details.
func (r *Redis) RoleName() (string, error) {
	rp, err := r.ExecuteCommand("ROLE")
	if err != nil {
		return "", err
	}
	return rp.Multi[0].StringValue()
}

// ClientList returns information and statistics
// about the client connections server in a mostly human readable format.
// Bulk reply: a unique string, formatted as follows:
// One client connection per line (separated by LF)
// Each line is composed of a succession of property=value fields separated by a space character.
func (r *Redis) ClientList() (string, error) {
	rp, err := r.ExecuteCommand("CLIENT", "LIST")
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// ClientGetName returns the name of the current connection as set by CLIENT SETNAME.
// Since every new connection starts without an associated name,
// if no name was assigned a null bulk reply is returned.
func (r *Redis) ClientGetName() ([]byte, error) {
	rp, err := r.ExecuteCommand("CLIENT", "GETNAME")
	if err != nil {
		return nil, err
	}
	return rp.BytesValue()
}

// ClientPause stops the server processing commands from clients for some time.
func (r *Redis) ClientPause(timeout uint64) error {
	rp, err := r.ExecuteCommand("CLIENT", "PAUSE", timeout)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ClientSetName assigns a name to the current connection.
func (r *Redis) ClientSetName(name string) error {
	rp, err := r.ExecuteCommand("CLIENT", "SETNAME", name)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ConfigGet is used to read the configuration parameters of a running Redis server.
// Not all the configuration parameters are supported in Redis 2.4,
// while Redis 2.6 can read the whole configuration of a server using this command.
// CONFIG GET takes a single argument, which is a glob-style pattern.
func (r *Redis) ConfigGet(parameter string) (map[string]string, error) {
	rp, err := r.ExecuteCommand("CONFIG", "GET", parameter)
	if err != nil {
		return nil, err
	}
	return rp.HashValue()
}

// ConfigRewrite rewrites the redis.conf file the server was started with,
// applying the minimal changes needed to make it reflecting the configuration currently used by the server,
// that may be different compared to the original one because of the use of the CONFIG SET command.
// Available since 2.8.0.
func (r *Redis) ConfigRewrite() error {
	rp, err := r.ExecuteCommand("CONFIG", "REWRITE")
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ConfigSet is used in order to reconfigure the server at run time without the need to restart Redis.
// You can change both trivial parameters or switch from one to another persistence option using this command.
func (r *Redis) ConfigSet(parameter, value string) error {
	rp, err := r.ExecuteCommand("CONFIG", "SET", parameter, value)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ConfigSetInt is a convenience wrapper for passing integers to ConfigSet
func (r *Redis) ConfigSetInt(parameter string, value int) error {
	sval := fmt.Sprintf("%d", value)
	rp, err := r.ExecuteCommand("CONFIG", "SET", parameter, sval)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ConfigResetStat resets the statistics reported by Redis using the INFO command.
// These are the counters that are reset:
// Keyspace hits
// Keyspace misses
// Number of commands processed
// Number of connections received
// Number of expired keys
// Number of rejected connections
// Latest fork(2) time
// The aof_delayed_fsync counter
func (r *Redis) ConfigResetStat() error {
	_, err := r.ExecuteCommand("CONFIG", "RESETSTAT")
	return err
}

// DBSize return the number of keys in the currently-selected database.
func (r *Redis) DBSize() (int64, error) {
	rp, err := r.ExecuteCommand("DBSIZE")
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// DebugObject is a debugging command that should not be used by clients.
func (r *Redis) DebugObject(key string) (string, error) {
	rp, err := r.ExecuteCommand("DEBUG", "OBJECT", key)
	if err != nil {
		return "", err
	}
	return rp.StatusValue()
}

// FlushAll delete all the keys of all the existing databases,
// not just the currently selected one.
// This command never fails.
func (r *Redis) FlushAll() error {
	_, err := r.ExecuteCommand("FLUSHALL")
	return err
}

// FlushDB delete all the keys of the currently selected DB.
// This command never fails.
func (r *Redis) FlushDB() error {
	_, err := r.ExecuteCommand("FLUSHDB")
	return err
}

// Info returns information and statistics about the server
// In RedisInfoAll struct see the github.com/therealbill/libredis/info package
// for details
func (r *Redis) Info() (sinfo structures.RedisInfoAll, err error) {
	rp, err := r.ExecuteCommand("info", "all")
	if err != nil {
		return
	}
	strval, _ := rp.StringValue()
	if err != nil {
		return
	}
	sinfo = info.GetAllInfo(strval)
	return
}

// SentinelInfo returns information and statistics for a sentinel instance
// In RedisInfoAll struct see the github.com/therealbill/libredis/info package
// for details
func (r *Redis) SentinelInfo() (sinfo structures.RedisInfoAll, err error) {
	rp, err := r.ExecuteCommand("info")
	if err != nil {
		return
	}
	strval, _ := rp.StringValue()
	if err != nil {
		return
	}
	sinfo = info.GetAllInfo(strval)
	return
}

// InfoString returns information and statistics about the server
// in a format that is simple to parse by computers and easy to read by humans.
// format document at http://redis.io/commands/info
func (r *Redis) InfoString(section string) (string, error) {
	args := packArgs("INFO", section)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// LastSave return the UNIX TIME of the last DB save executed with success.
// A client may check if a BGSAVE command succeeded reading the LASTSAVE value,
// then issuing a BGSAVE command and checking at regular intervals every N seconds if LASTSAVE changed.
// Integer reply: an UNIX time stamp.
func (r *Redis) LastSave() (int64, error) {
	rp, err := r.ExecuteCommand("LASTSAVE")
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// MonitorCommand is a debugging command that streams back every command processed by the Redis server.
type MonitorCommand struct {
	redis *Redis
	conn  *connection
}

// Monitor sned MONITOR command to redis server.
func (r *Redis) Monitor() (*MonitorCommand, error) {
	c, err := r.pool.Get()
	if err != nil {
		return nil, err
	}
	if err := c.SendCommand("MONITOR"); err != nil {
		return nil, err
	}
	rp, err := c.RecvReply()
	if err != nil {
		return nil, err
	}
	if err := rp.OKValue(); err != nil {
		return nil, err
	}
	return &MonitorCommand{r, c}, nil
}

// Receive read from redis server and return the reply.
func (m *MonitorCommand) Receive() (string, error) {
	rp, err := m.conn.RecvReply()
	if err != nil {
		return "", err
	}
	return rp.StatusValue()
}

// Close closes current monitor command.
func (m *MonitorCommand) Close() error {
	return m.conn.SendCommand("QUIT")
}

// Save performs a synchronous save of the dataset
// producing a point in time snapshot of all the data inside the Redis instance,
// in the form of an RDB file.
// You almost never want to call SAVE in production environments
// where it will block all the other clients. Instead usually BGSAVE is used.
func (r *Redis) Save() error {
	rp, err := r.ExecuteCommand("SAVE")
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// Shutdown behavior is the following:
// Stop all the clients.
// Perform a blocking SAVE if at least one save point is configured.
// Flush the Append Only File if AOF is enabled.
// Quit the server.
func (r *Redis) Shutdown(save bool) error {
	args := packArgs("SHUTDOWN")
	if save {
		args = append(args, "SAVE")
	} else {
		args = append(args, "NOSAVE")
	}
	rp, err := r.ExecuteCommand(args...)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.New(rp.Status)
}

// SlaveOf can change the replication settings of a slave on the fly.
// If a Redis server is already acting as slave, the command SLAVEOF NO ONE will turn off the replication,
// turning the Redis server into a MASTER.
// In the proper form SLAVEOF hostname port will make the server a slave of
// another server listening at the specified hostname and port.
//
// If a server is already a slave of some master,
// SLAVEOF hostname port will stop the replication against the old server
// and start the synchronization against the new one, discarding the old dataset.
// The form SLAVEOF NO ONE will stop replication, turning the server into a MASTER,
// but will not discard the replication.
// So, if the old master stops working,
// it is possible to turn the slave into a master and set the application to use this new master in read/write.
// Later when the other Redis server is fixed, it can be reconfigured to work as a slave.
func (r *Redis) SlaveOf(host, port string) error {
	rp, err := r.ExecuteCommand("SLAVEOF", host, port)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// SlowLog is used in order to read and reset the Redis slow queries log.
type SlowLog struct {
	ID           int64
	Timestamp    int64
	Microseconds int64
	Command      []string
}

// SlowLogGet returns slow logs.
func (r *Redis) SlowLogGet(n int64) ([]*SlowLog, error) {
	rp, err := r.ExecuteCommand("SLOWLOG", "GET", n)
	if err != nil {
		return nil, err
	}
	if rp.Type == ErrorReply {
		return nil, errors.New(rp.Error)
	}
	if rp.Type != MultiReply {
		return nil, errors.New("slowlog get protocol error")
	}
	var slow []*SlowLog
	for _, subrp := range rp.Multi {
		if subrp.Multi == nil || len(subrp.Multi) != 4 {
			return nil, errors.New("slowlog get protocol error")
		}
		id, err := subrp.Multi[0].IntegerValue()
		if err != nil {
			return nil, err
		}
		timestamp, err := subrp.Multi[1].IntegerValue()
		if err != nil {
			return nil, err
		}
		microseconds, err := subrp.Multi[2].IntegerValue()
		if err != nil {
			return nil, err
		}
		command, err := subrp.Multi[3].ListValue()
		if err != nil {
			return nil, err
		}
		slow = append(slow, &SlowLog{id, timestamp, microseconds, command})
	}
	return slow, nil
}

// SlowLogLen Obtaining the current length of the slow log
func (r *Redis) SlowLogLen() (int64, error) {
	rp, err := r.ExecuteCommand("SLOWLOG", "LEN")
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// SlowLogReset resetting the slow log.
// Once deleted the information is lost forever.
func (r *Redis) SlowLogReset() error {
	rp, err := r.ExecuteCommand("SLOWLOG", "RESET")
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// Time returns a multi bulk reply containing two elements:
// unix time in seconds,
// microseconds.
func (r *Redis) Time() ([]string, error) {
	rp, err := r.ExecuteCommand("TIME")
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// Command returns the Redis command info structures
func (r *Redis) Command() (comms []structures.CommandEntry, err error) {
	rp, err := r.ExecuteCommand("COMMAND")
	if err != nil {
		return nil, err
	}
	for _, subrp := range rp.Multi {
		name, _ := subrp.Multi[0].StringValue()
		arity, _ := subrp.Multi[1].IntegerValue()
		first, _ := subrp.Multi[3].IntegerValue()
		last, _ := subrp.Multi[4].IntegerValue()
		repeat, _ := subrp.Multi[5].IntegerValue()
		ce := structures.CommandEntry{Name: name, Arity: arity, FirstKey: first, LastKey: last, RepeatCount: repeat}
		flagmap := make(map[string]bool)
		for _, crp := range subrp.Multi[2].Multi {
			flag, _ := crp.StatusValue()
			flagmap[flag] = true
		}
		ce.Flags = flagmap
		comms = append(comms, ce)
	}
	return
}

// Memory Management (Redis 4.0+)

// Memory statistics and information types
type MemoryStats struct {
	PeakAllocated      int64
	TotalAllocated     int64
	StartupAllocated   int64
	ReplicationBacklog int64
	ClientsSlaves      int64
	ClientsNormal      int64
	AOFBuffer         int64
	LuaCaches         int64
	Overhead          MemoryOverhead
	Keys              MemoryKeys
	Dataset           MemoryDataset
}

type MemoryOverhead struct {
	Total     int64
	Hashtable int64
	Expires   int64
}

type MemoryKeys struct {
	Count               int64
	BucketsCount        int64
	ExpiringCount       int64
	ExpiringBucketsCount int64
}

type MemoryDataset struct {
	Bytes      int64
	Percentage float64
}

// MEMORY USAGE key [SAMPLES count]
// MemoryUsage returns memory usage information for a key.
func (r *Redis) MemoryUsage(key string) (int64, error) {
	args := packArgs("MEMORY", "USAGE", key)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// MemoryUsageWithSamples returns memory usage with specific sample count.
func (r *Redis) MemoryUsageWithSamples(key string, samples int) (int64, error) {
	args := packArgs("MEMORY", "USAGE", key, "SAMPLES", samples)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// MEMORY STATS
// MemoryStats returns detailed memory usage statistics.
func (r *Redis) MemoryStats() (MemoryStats, error) {
	args := packArgs("MEMORY", "STATS")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return MemoryStats{}, err
	}

	stats := MemoryStats{}
	if rp.Multi != nil {
		for i := 0; i < len(rp.Multi)-1; i += 2 {
			key, _ := rp.Multi[i].StringValue()
			value := rp.Multi[i+1]

			switch key {
			case "peak.allocated":
				stats.PeakAllocated, _ = value.IntegerValue()
			case "total.allocated":
				stats.TotalAllocated, _ = value.IntegerValue()
			case "startup.allocated":
				stats.StartupAllocated, _ = value.IntegerValue()
			case "replication.backlog":
				stats.ReplicationBacklog, _ = value.IntegerValue()
			case "clients.slaves":
				stats.ClientsSlaves, _ = value.IntegerValue()
			case "clients.normal":
				stats.ClientsNormal, _ = value.IntegerValue()
			case "aof.buffer":
				stats.AOFBuffer, _ = value.IntegerValue()
			case "lua.caches":
				stats.LuaCaches, _ = value.IntegerValue()
			case "overhead.hashtable.main":
				stats.Overhead.Hashtable, _ = value.IntegerValue()
			case "overhead.hashtable.expires":
				stats.Overhead.Expires, _ = value.IntegerValue()
			case "overhead.total":
				stats.Overhead.Total, _ = value.IntegerValue()
			case "keys.count":
				stats.Keys.Count, _ = value.IntegerValue()
			case "keys.bytes-per-key":
				// Skip calculated fields
			case "dataset.bytes":
				stats.Dataset.Bytes, _ = value.IntegerValue()
			case "dataset.percentage":
				if valueStr, err := value.StringValue(); err == nil {
					fmt.Sscanf(valueStr, "%f", &stats.Dataset.Percentage)
				}
			}
		}
	}

	return stats, nil
}

// MEMORY DOCTOR
// MemoryDoctor returns memory analysis and recommendations.
func (r *Redis) MemoryDoctor() (string, error) {
	args := packArgs("MEMORY", "DOCTOR")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// MEMORY PURGE
// MemoryPurge attempts to purge dirty pages for better memory reporting.
func (r *Redis) MemoryPurge() error {
	args := packArgs("MEMORY", "PURGE")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// Latency Monitoring (Redis 2.8.13+)

// Latency monitoring types
type LatencySample struct {
	Timestamp int64
	Latency   int64
}

type LatencyStats struct {
	Event    string
	Latest   int64
	AllTime  int64
	Samples  []LatencySample
}

// LATENCY LATEST
// LatencyLatest returns latest latency samples for all events.
func (r *Redis) LatencyLatest() ([]LatencyStats, error) {
	args := packArgs("LATENCY", "LATEST")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	if rp.Multi == nil {
		return nil, nil
	}

	var latencyStats []LatencyStats
	for _, eventReply := range rp.Multi {
		if len(eventReply.Multi) >= 4 {
			event, _ := eventReply.Multi[0].StringValue()
			latest, _ := eventReply.Multi[1].IntegerValue()
			allTime, _ := eventReply.Multi[2].IntegerValue()

			stats := LatencyStats{
				Event:   event,
				Latest:  latest,
				AllTime: allTime,
			}
			latencyStats = append(latencyStats, stats)
		}
	}

	return latencyStats, nil
}

// Note: LatencyHistory is already implemented in latency.go with different signature

// LATENCY RESET [event ...]
// LatencyReset resets latency data for all or specified events.
func (r *Redis) LatencyReset() (int64, error) {
	args := packArgs("LATENCY", "RESET")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// Note: LatencyResetEvents is already implemented in latency.go with different signature

// LATENCY GRAPH event
// LatencyGraph returns ASCII art latency graph for an event.
func (r *Redis) LatencyGraph(event string) (string, error) {
	args := packArgs("LATENCY", "GRAPH", event)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// Database Management (Redis 4.0+)

// SWAPDB index1 index2
// SwapDB swaps the contents of two Redis databases.
func (r *Redis) SwapDB(index1, index2 int) error {
	args := packArgs("SWAPDB", index1, index2)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// REPLICAOF host port / REPLICAOF NO ONE
// ReplicaOf configures Redis as a replica of another instance.
func (r *Redis) ReplicaOf(host, port string) error {
	args := packArgs("REPLICAOF", host, port)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ReplicaOfNoOne stops replication and promotes to master.
func (r *Redis) ReplicaOfNoOne() error {
	args := packArgs("REPLICAOF", "NO", "ONE")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// Module Management (Redis 4.0+)

// Module information
type ModuleInfo struct {
	Name    string
	Version int64
	Path    string
	Args    []string
}

// MODULE LIST
// ModuleList returns information about loaded Redis modules.
func (r *Redis) ModuleList() ([]ModuleInfo, error) {
	args := packArgs("MODULE", "LIST")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	if rp.Multi == nil {
		return nil, nil
	}

	var modules []ModuleInfo
	for _, moduleReply := range rp.Multi {
		if moduleReply.Multi != nil && len(moduleReply.Multi) >= 6 {
			module := ModuleInfo{}
			
			// Parse module information: [name, name_value, ver, version_value, path, path_value, args, args_array]
			for i := 0; i < len(moduleReply.Multi)-1; i += 2 {
				key, _ := moduleReply.Multi[i].StringValue()
				value := moduleReply.Multi[i+1]

				switch key {
				case "name":
					module.Name, _ = value.StringValue()
				case "ver":
					module.Version, _ = value.IntegerValue()
				case "path":
					module.Path, _ = value.StringValue()
				case "args":
					if value.Multi != nil {
						for _, argReply := range value.Multi {
							arg, _ := argReply.StringValue()
							module.Args = append(module.Args, arg)
						}
					}
				}
			}
			modules = append(modules, module)
		}
	}

	return modules, nil
}
