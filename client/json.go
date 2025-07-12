package client

import (
	"strconv"
)

// JSONOptions represents options for JSON commands
type JSONOptions struct {
	Path     string
	Indent   string
	NewLine  string
	Space    string
	NX       bool // Only set if the path does not exist
	XX       bool // Only set if the path exists
}

// JSONSetOptions represents options for JSON.SET command
type JSONSetOptions struct {
	NX bool // Only set if the path does not exist
	XX bool // Only set if the path exists
}

// JSONGetOptions represents options for JSON.GET command
type JSONGetOptions struct {
	Indent  string
	NewLine string
	Space   string
	Paths   []string
}

// JSONArrInsertOptions represents options for JSON.ARRINSERT command
type JSONArrInsertOptions struct {
	Values []interface{}
}

// Basic JSON Operations

// JSONSet command:
// Set the JSON value at path in key
// JSON.SET key path value [NX|XX]
func (r *Redis) JSONSet(key, path string, value interface{}, options ...*JSONSetOptions) (string, error) {
	args := []interface{}{"JSON.SET", key, path, value}
	
	if len(options) > 0 && options[0] != nil {
		if options[0].NX {
			args = append(args, "NX")
		} else if options[0].XX {
			args = append(args, "XX")
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// JSONGet command:
// Return the value at path in JSON serialized form
// JSON.GET key [INDENT indent] [NEWLINE newline] [SPACE space] [path ...]
func (r *Redis) JSONGet(key string, options ...*JSONGetOptions) ([]byte, error) {
	args := []interface{}{"JSON.GET", key}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		if opt.Indent != "" {
			args = append(args, "INDENT", opt.Indent)
		}
		if opt.NewLine != "" {
			args = append(args, "NEWLINE", opt.NewLine)
		}
		if opt.Space != "" {
			args = append(args, "SPACE", opt.Space)
		}
		if len(opt.Paths) > 0 {
			for _, path := range opt.Paths {
				args = append(args, path)
			}
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.BytesValue()
}

// JSONDel command:
// Delete a value
// JSON.DEL key [path]
func (r *Redis) JSONDel(key string, path ...string) (int64, error) {
	args := []interface{}{"JSON.DEL", key}
	if len(path) > 0 {
		args = append(args, path[0])
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// JSONType command:
// Report the type of JSON value at path
// JSON.TYPE key [path]
func (r *Redis) JSONType(key string, path ...string) (string, error) {
	args := []interface{}{"JSON.TYPE", key}
	if len(path) > 0 {
		args = append(args, path[0])
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// Numeric Operations

// JSONNumIncrBy command:
// Increment the number value stored at path by number
// JSON.NUMINCRBY key path number
func (r *Redis) JSONNumIncrBy(key, path string, number float64) (float64, error) {
	rp, err := r.ExecuteCommand("JSON.NUMINCRBY", key, path, number)
	if err != nil {
		return 0.0, err
	}
	s, err := rp.StringValue()
	if err != nil {
		return 0.0, err
	}
	return strconv.ParseFloat(s, 64)
}

// JSONNumMultBy command:
// Multiply the number value stored at path by number
// JSON.NUMMULTBY key path number
func (r *Redis) JSONNumMultBy(key, path string, number float64) (float64, error) {
	rp, err := r.ExecuteCommand("JSON.NUMMULTBY", key, path, number)
	if err != nil {
		return 0.0, err
	}
	s, err := rp.StringValue()
	if err != nil {
		return 0.0, err
	}
	return strconv.ParseFloat(s, 64)
}

// String Operations

// JSONStrAppend command:
// Append the json-string value(s) the string at path
// JSON.STRAPPEND key [path] json-string
func (r *Redis) JSONStrAppend(key string, path string, jsonString string) (int64, error) {
	rp, err := r.ExecuteCommand("JSON.STRAPPEND", key, path, jsonString)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// JSONStrLen command:
// Report the length of the JSON String at path in key
// JSON.STRLEN key [path]
func (r *Redis) JSONStrLen(key string, path ...string) (int64, error) {
	args := []interface{}{"JSON.STRLEN", key}
	if len(path) > 0 {
		args = append(args, path[0])
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// Array Operations

// JSONArrAppend command:
// Append the json value(s) into the array at path after the last element in it
// JSON.ARRAPPEND key path value [value ...]
func (r *Redis) JSONArrAppend(key, path string, values ...interface{}) (int64, error) {
	args := []interface{}{"JSON.ARRAPPEND", key, path}
	for _, value := range values {
		args = append(args, value)
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// JSONArrIndex command:
// Search for the first occurrence of a scalar JSON value in an array
// JSON.ARRINDEX key path value [start [stop]]
func (r *Redis) JSONArrIndex(key, path string, value interface{}, startStop ...int) (int64, error) {
	args := []interface{}{"JSON.ARRINDEX", key, path, value}
	if len(startStop) > 0 {
		args = append(args, startStop[0])
		if len(startStop) > 1 {
			args = append(args, startStop[1])
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// JSONArrInsert command:
// Insert the json value(s) into the array at path before the index (shifts to the right)
// JSON.ARRINSERT key path index value [value ...]
func (r *Redis) JSONArrInsert(key, path string, index int, values ...interface{}) (int64, error) {
	args := []interface{}{"JSON.ARRINSERT", key, path, index}
	for _, value := range values {
		args = append(args, value)
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// JSONArrLen command:
// Report the length of the JSON Array at path in key
// JSON.ARRLEN key [path]
func (r *Redis) JSONArrLen(key string, path ...string) (int64, error) {
	args := []interface{}{"JSON.ARRLEN", key}
	if len(path) > 0 {
		args = append(args, path[0])
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// JSONArrPop command:
// Remove and return element from the index in the array
// JSON.ARRPOP key [path [index]]
func (r *Redis) JSONArrPop(key string, path string, index ...int) ([]byte, error) {
	args := []interface{}{"JSON.ARRPOP", key, path}
	if len(index) > 0 {
		args = append(args, index[0])
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.BytesValue()
}

// JSONArrTrim command:
// Trim an array so that it contains only the specified inclusive range of elements
// JSON.ARRTRIM key path start stop
func (r *Redis) JSONArrTrim(key, path string, start, stop int) (int64, error) {
	rp, err := r.ExecuteCommand("JSON.ARRTRIM", key, path, start, stop)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// Object Operations

// JSONObjKeys command:
// Return the keys in the object that's referenced by path
// JSON.OBJKEYS key [path]
func (r *Redis) JSONObjKeys(key string, path ...string) ([]string, error) {
	args := []interface{}{"JSON.OBJKEYS", key}
	if len(path) > 0 {
		args = append(args, path[0])
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	return rp.ListValue()
}

// JSONObjLen command:
// Report the number of keys in the JSON Object at path in key
// JSON.OBJLEN key [path]
func (r *Redis) JSONObjLen(key string, path ...string) (int64, error) {
	args := []interface{}{"JSON.OBJLEN", key}
	if len(path) > 0 {
		args = append(args, path[0])
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}