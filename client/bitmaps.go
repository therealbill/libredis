package client

// BitFieldOperation represents a single bitfield operation
type BitFieldOperation struct {
	Type   string      // GET, SET, INCRBY
	Offset int64       // Bit offset
	Value  interface{} // Value for SET/INCRBY operations
}

// BitFieldOverflow represents overflow behavior
type BitFieldOverflow string

const (
	BitFieldOverflowWrap BitFieldOverflow = "WRAP"
	BitFieldOverflowSat  BitFieldOverflow = "SAT" 
	BitFieldOverflowFail BitFieldOverflow = "FAIL"
)

// BitPosOptions represents options for BITPOS command
type BitPosOptions struct {
	Start *int64 // Start position
	End   *int64 // End position
}

// BITFIELD key [GET type offset] [SET type offset value] [INCRBY type offset increment] [OVERFLOW WRAP|SAT|FAIL]
// BitField performs arbitrary bit field integer operations on strings.
// Redis 3.2+
func (r *Redis) BitField(key string, operations []BitFieldOperation) ([]int64, error) {
	args := []interface{}{"BITFIELD", key}
	
	for _, op := range operations {
		switch op.Type {
		case "GET":
			args = append(args, "GET", op.Offset)
		case "SET":
			args = append(args, "SET", op.Offset, op.Value)
		case "INCRBY":
			args = append(args, "INCRBY", op.Offset, op.Value)
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	if rp.Type == MultiReply {
		result := make([]int64, len(rp.Multi))
		for i, item := range rp.Multi {
			if val, err := item.IntegerValue(); err == nil {
				result[i] = val
			}
		}
		return result, nil
	}
	
	return nil, nil
}

// BitFieldWithOverflow performs bitfield operations with overflow control.
// Redis 3.2+
func (r *Redis) BitFieldWithOverflow(key string, overflow BitFieldOverflow, operations []BitFieldOperation) ([]int64, error) {
	args := []interface{}{"BITFIELD", key, "OVERFLOW", string(overflow)}
	
	for _, op := range operations {
		switch op.Type {
		case "GET":
			args = append(args, "GET", op.Offset)
		case "SET":
			args = append(args, "SET", op.Offset, op.Value)
		case "INCRBY":
			args = append(args, "INCRBY", op.Offset, op.Value)
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	if rp.Type == MultiReply {
		result := make([]int64, len(rp.Multi))
		for i, item := range rp.Multi {
			if val, err := item.IntegerValue(); err == nil {
				result[i] = val
			}
		}
		return result, nil
	}
	
	return nil, nil
}

// BITFIELD_RO key [GET type offset] [GET type offset ...]
// BitFieldRO is the read-only variant of BITFIELD.
// Redis 6.0+
func (r *Redis) BitFieldRO(key string, getOps []BitFieldOperation) ([]int64, error) {
	args := []interface{}{"BITFIELD_RO", key}
	
	for _, op := range getOps {
		if op.Type == "GET" {
			args = append(args, "GET", op.Offset)
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	if rp.Type == MultiReply {
		result := make([]int64, len(rp.Multi))
		for i, item := range rp.Multi {
			if val, err := item.IntegerValue(); err == nil {
				result[i] = val
			}
		}
		return result, nil
	}
	
	return nil, nil
}

// BITPOS key bit [start] [end]
// BitPos returns the position of the first bit set to 1 or 0.
// Redis 2.8.7+
func (r *Redis) BitPos(key string, bit int) (int64, error) {
	rp, err := r.ExecuteCommand("BITPOS", key, bit)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// BitPosWithRange returns the bit position within a specified range.
// Redis 2.8.7+
func (r *Redis) BitPosWithRange(key string, bit int, opts BitPosOptions) (int64, error) {
	args := []interface{}{"BITPOS", key, bit}
	
	if opts.Start != nil {
		args = append(args, *opts.Start)
		if opts.End != nil {
			args = append(args, *opts.End)
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}