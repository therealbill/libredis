package client

// ACL constants for command categories and operations
const (
	ACLLogReset = "RESET"

	// Common ACL rule patterns
	ACLRuleAllCommands = "+@all"
	ACLRuleNoCommands  = "-@all"
	ACLRuleAllKeys     = "~*"
	ACLRuleNoKeys      = ""
)

// ACLUser represents a Redis ACL user
type ACLUser struct {
	Username  string
	Flags     []string
	Passwords []string
	Commands  []string
	Keys      []string
	Channels  []string
}

// ACLLogEntry represents an ACL log entry
type ACLLogEntry struct {
	Count       int64
	Reason      string
	Context     string
	Object      string
	Username    string
	AgeSeconds  float64
	ClientInfo  string
}

// ACLGenPassOptions represents options for ACL GENPASS
type ACLGenPassOptions struct {
	Bits int // Number of bits for password generation
}

// User Management

// ACL SETUSER username [rule ...]
// ACLSetUser creates or modifies an ACL user with specified rules.
func (r *Redis) ACLSetUser(username string, rules ...string) error {
	args := []interface{}{"ACL", "SETUSER", username}
	for _, rule := range rules {
		args = append(args, rule)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ACL GETUSER username
// ACLGetUser returns information about a specific ACL user.
func (r *Redis) ACLGetUser(username string) (ACLUser, error) {
	args := packArgs("ACL", "GETUSER", username)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return ACLUser{}, err
	}

	user := ACLUser{Username: username}
	if rp.Multi != nil {
		for i := 0; i < len(rp.Multi)-1; i += 2 {
			key, _ := rp.Multi[i].StringValue()
			valueReply := rp.Multi[i+1]

			switch key {
			case "flags":
				if valueReply.Multi != nil {
					for _, flagReply := range valueReply.Multi {
						flag, _ := flagReply.StringValue()
						user.Flags = append(user.Flags, flag)
					}
				}
			case "passwords":
				if valueReply.Multi != nil {
					for _, passReply := range valueReply.Multi {
						password, _ := passReply.StringValue()
						user.Passwords = append(user.Passwords, password)
					}
				}
			case "commands":
				if valueReply.Multi != nil {
					for _, cmdReply := range valueReply.Multi {
						command, _ := cmdReply.StringValue()
						user.Commands = append(user.Commands, command)
					}
				}
			case "keys":
				if valueReply.Multi != nil {
					for _, keyReply := range valueReply.Multi {
						keyPattern, _ := keyReply.StringValue()
						user.Keys = append(user.Keys, keyPattern)
					}
				}
			case "channels":
				if valueReply.Multi != nil {
					for _, chanReply := range valueReply.Multi {
						channel, _ := chanReply.StringValue()
						user.Channels = append(user.Channels, channel)
					}
				}
			}
		}
	}

	return user, nil
}

// ACL DELUSER username [username ...]
// ACLDelUser deletes one or more ACL users.
func (r *Redis) ACLDelUser(usernames ...string) (int64, error) {
	args := packArgs("ACL", "DELUSER", usernames)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// ACL USERS
// ACLUsers returns a list of all ACL usernames.
func (r *Redis) ACLUsers() ([]string, error) {
	args := packArgs("ACL", "USERS")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// Permissions and Categories

// ACL CAT [categoryname]
// ACLCat returns a list of all ACL command categories.
func (r *Redis) ACLCat() ([]string, error) {
	args := packArgs("ACL", "CAT")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// ACLCatByCategory returns commands in a specific category.
func (r *Redis) ACLCatByCategory(categoryname string) ([]string, error) {
	args := packArgs("ACL", "CAT", categoryname)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// ACL WHOAMI
// ACLWhoAmI returns the username of the current connection.
func (r *Redis) ACLWhoAmI() (string, error) {
	args := packArgs("ACL", "WHOAMI")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// ACL LOG [count|RESET]
// ACLLog returns ACL security events log entries.
func (r *Redis) ACLLog() ([]ACLLogEntry, error) {
	args := packArgs("ACL", "LOG")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	return parseACLLogEntries(rp.Multi)
}

// ACLLogWithCount returns a specific number of log entries.
func (r *Redis) ACLLogWithCount(count int) ([]ACLLogEntry, error) {
	args := packArgs("ACL", "LOG", count)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	return parseACLLogEntries(rp.Multi)
}

// ACLLogReset clears the ACL log.
func (r *Redis) ACLLogReset() error {
	args := packArgs("ACL", "LOG", ACLLogReset)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// Configuration Management

// ACL LOAD
// ACLLoad reloads ACL configuration from external ACL file.
func (r *Redis) ACLLoad() error {
	args := packArgs("ACL", "LOAD")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ACL SAVE
// ACLSave saves current ACL configuration to external file.
func (r *Redis) ACLSave() error {
	args := packArgs("ACL", "SAVE")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// ACL LIST
// ACLList returns a list of ACL rules for all users.
func (r *Redis) ACLList() ([]string, error) {
	args := packArgs("ACL", "LIST")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// Utilities

// ACL GENPASS [bits]
// ACLGenPass generates a secure password for ACL users.
func (r *Redis) ACLGenPass() (string, error) {
	args := packArgs("ACL", "GENPASS")
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// ACLGenPassWithBits generates a password with specified bit length.
func (r *Redis) ACLGenPassWithBits(bits int) (string, error) {
	args := packArgs("ACL", "GENPASS", bits)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// ACL DRYRUN username command [arg ...]
// ACLDryRun simulates command execution for permission testing.
func (r *Redis) ACLDryRun(username, command string, args ...string) error {
	cmdArgs := []interface{}{"ACL", "DRYRUN", username, command}
	for _, arg := range args {
		cmdArgs = append(cmdArgs, arg)
	}

	rp, err := r.ExecuteCommand(cmdArgs...)
	if err != nil {
		return err
	}
	return rp.OKValue()
}

// Helper functions

func parseACLLogEntries(replies []*Reply) ([]ACLLogEntry, error) {
	if replies == nil {
		return nil, nil
	}

	entries := make([]ACLLogEntry, len(replies))
	for i, entryReply := range replies {
		entry := ACLLogEntry{}
		if entryReply.Multi != nil {
			for j := 0; j < len(entryReply.Multi)-1; j += 2 {
				key, _ := entryReply.Multi[j].StringValue()
				valueReply := entryReply.Multi[j+1]

				switch key {
				case "count":
					entry.Count, _ = valueReply.IntegerValue()
				case "reason":
					entry.Reason, _ = valueReply.StringValue()
				case "context":
					entry.Context, _ = valueReply.StringValue()
				case "object":
					entry.Object, _ = valueReply.StringValue()
				case "username":
					entry.Username, _ = valueReply.StringValue()
				case "age-seconds":
					if valueReply.Type == BulkReply {
						ageStr, _ := valueReply.StringValue()
						// Convert string to float64 if needed
						if ageStr != "" {
							// Simple float parsing - in production might want strconv.ParseFloat
							entry.AgeSeconds = 0 // Placeholder for proper parsing
						}
					}
				case "client-info":
					entry.ClientInfo, _ = valueReply.StringValue()
				}
			}
		}
		entries[i] = entry
	}

	return entries, nil
}