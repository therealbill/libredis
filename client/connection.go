package client

// Echo command returns message.
func (r *Redis) Echo(message string) (string, error) {
	rp, err := r.ExecuteCommand("ECHO", message)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// Ping command returns PONG.
// This command is often used to test if a connection is still alive, or to measure latency.
func (r *Redis) Ping() error {
	_, err := r.ExecuteCommand("PING")
	return err
}

func (r *Redis) Address() string {
	return r.address
}

// HelloOptions represents options for HELLO command
type HelloOptions struct {
	ProtocolVersion int
	Username        string
	Password        string
	ClientName      string
}

// AUTH [username] password
// AuthWithUser authenticates using username and password (ACL).
// Redis 6.0+
func (r *Redis) AuthWithUser(username, password string) error {
	_, err := r.ExecuteCommand("AUTH", username, password)
	return err
}

// HELLO [protover [AUTH username password] [SETNAME clientname]]
// Hello switches to a different protocol version and authenticates.
// Redis 6.0+
func (r *Redis) Hello(protocolVersion int) (map[string]interface{}, error) {
	rp, err := r.ExecuteCommand("HELLO", protocolVersion)
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]interface{})
	if rp.Type == MultiReply {
		for i := 0; i < len(rp.Multi); i += 2 {
			if i+1 < len(rp.Multi) {
				key, _ := rp.Multi[i].StringValue()
				value, _ := rp.Multi[i+1].StringValue()
				result[key] = value
			}
		}
	}
	
	return result, nil
}

// HelloWithOptions performs handshake with additional options.
// Redis 6.0+
func (r *Redis) HelloWithOptions(opts HelloOptions) (map[string]interface{}, error) {
	args := []interface{}{"HELLO"}
	
	if opts.ProtocolVersion != 0 {
		args = append(args, opts.ProtocolVersion)
	}
	
	if opts.Username != "" && opts.Password != "" {
		args = append(args, "AUTH", opts.Username, opts.Password)
	}
	
	if opts.ClientName != "" {
		args = append(args, "SETNAME", opts.ClientName)
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]interface{})
	if rp.Type == MultiReply {
		for i := 0; i < len(rp.Multi); i += 2 {
			if i+1 < len(rp.Multi) {
				key, _ := rp.Multi[i].StringValue()
				value, _ := rp.Multi[i+1].StringValue()
				result[key] = value
			}
		}
	}
	
	return result, nil
}

// RESET
// Reset resets the connection state.
// Redis 6.2+
func (r *Redis) Reset() error {
	_, err := r.ExecuteCommand("RESET")
	return err
}
