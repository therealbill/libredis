package structures

// CommandEntry defines the entries for a Redis command as shown in the
// results of the COMMAND command.
type CommandEntry struct {
	Name        string
	Arity       int64
	Flags       map[string]bool
	FirstKey    int64
	LastKey     int64
	RepeatCount int64
}

// ReadOnly returns true if the command has the "readonly" flag set.
func (c *CommandEntry) ReadOnly() bool {
	_, set := c.Flags["readonly"]
	return set
}

// Writes returns true if the command has the "write" flag set.
// This indicates modifications are made to the datastore
func (c *CommandEntry) Writes() bool {
	_, set := c.Flags["write"]
	return set
}

// Admin returns true if the command has the "admin" flag set.
func (c *CommandEntry) Admin() bool {
	_, set := c.Flags["admin"]
	return set
}

// Pubsub returns true if the command is pubsub related
func (c *CommandEntry) Pubsub() bool {
	_, set := c.Flags["pubsub"]
	return set
}
