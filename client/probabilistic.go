package client

// BFReserveOptions represents options for BF.RESERVE command
type BFReserveOptions struct {
	Capacity    int64   // Initial capacity
	ErrorRate   float64 // Desired error rate
	Expansion   int     // Expansion factor when full
	NonScaling  bool    // Don't create additional filters
}

// CFReserveOptions represents options for CF.RESERVE command
type CFReserveOptions struct {
	Capacity    int64 // Initial capacity
	BucketSize  int   // Number of items per bucket
	MaxIterations int // Number of iterations before declaring filter full
	Expansion   int   // Expansion factor when full
}

// CMSInitByDimOptions represents options for CMS.INITBYDIM command
type CMSInitByDimOptions struct {
	Width int64 // Number of counters per array
	Depth int64 // Number of counter arrays
}

// CMSInitByProbOptions represents options for CMS.INITBYPROB command
type CMSInitByProbOptions struct {
	ErrorRate float64 // Overestimation error rate
	Probability float64 // Probability for accuracy
}

// Bloom Filter Commands

// BFReserve command:
// Create an empty Bloom filter with a given desired error ratio and initial capacity
// BF.RESERVE key error_rate capacity [EXPANSION expansion] [NONSCALING]
func (r *Redis) BFReserve(key string, errorRate float64, capacity int64, options ...*BFReserveOptions) (string, error) {
	args := []interface{}{"BF.RESERVE", key, errorRate, capacity}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.Expansion > 0 {
			args = append(args, "EXPANSION", opt.Expansion)
		}
		
		if opt.NonScaling {
			args = append(args, "NONSCALING")
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// BFAdd command:
// Add an item to a Bloom filter
// BF.ADD key item
func (r *Redis) BFAdd(key string, item interface{}) (bool, error) {
	rp, err := r.ExecuteCommand("BF.ADD", key, item)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// BFMAdd command:
// Add multiple items to a Bloom filter
// BF.MADD key item [item ...]
func (r *Redis) BFMAdd(key string, items ...interface{}) ([]bool, error) {
	args := []interface{}{"BF.MADD", key}
	for _, item := range items {
		args = append(args, item)
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make([]bool, len(multi))
	for i, reply := range multi {
		result[i], _ = reply.BoolValue()
	}
	
	return result, nil
}

// BFExists command:
// Check if an item exists in a Bloom filter
// BF.EXISTS key item
func (r *Redis) BFExists(key string, item interface{}) (bool, error) {
	rp, err := r.ExecuteCommand("BF.EXISTS", key, item)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// BFMExists command:
// Check if multiple items exist in a Bloom filter
// BF.MEXISTS key item [item ...]
func (r *Redis) BFMExists(key string, items ...interface{}) ([]bool, error) {
	args := []interface{}{"BF.MEXISTS", key}
	for _, item := range items {
		args = append(args, item)
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make([]bool, len(multi))
	for i, reply := range multi {
		result[i], _ = reply.BoolValue()
	}
	
	return result, nil
}

// Cuckoo Filter Commands

// CFReserve command:
// Create an empty Cuckoo filter with given capacity
// CF.RESERVE key capacity [BUCKETSIZE bucketsize] [MAXITERATIONS maxiterations] [EXPANSION expansion]
func (r *Redis) CFReserve(key string, capacity int64, options ...*CFReserveOptions) (string, error) {
	args := []interface{}{"CF.RESERVE", key, capacity}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.BucketSize > 0 {
			args = append(args, "BUCKETSIZE", opt.BucketSize)
		}
		
		if opt.MaxIterations > 0 {
			args = append(args, "MAXITERATIONS", opt.MaxIterations)
		}
		
		if opt.Expansion > 0 {
			args = append(args, "EXPANSION", opt.Expansion)
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// CFAdd command:
// Add an item to a Cuckoo filter
// CF.ADD key item
func (r *Redis) CFAdd(key string, item interface{}) (bool, error) {
	rp, err := r.ExecuteCommand("CF.ADD", key, item)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// CFExists command:
// Check if an item exists in a Cuckoo filter
// CF.EXISTS key item
func (r *Redis) CFExists(key string, item interface{}) (bool, error) {
	rp, err := r.ExecuteCommand("CF.EXISTS", key, item)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// CFDel command:
// Delete an item from a Cuckoo filter
// CF.DEL key item
func (r *Redis) CFDel(key string, item interface{}) (bool, error) {
	rp, err := r.ExecuteCommand("CF.DEL", key, item)
	if err != nil {
		return false, err
	}
	return rp.BoolValue()
}

// Count-Min Sketch Commands

// CMSInitByDim command:
// Initialize a Count-Min Sketch with specified dimensions
// CMS.INITBYDIM key width depth
func (r *Redis) CMSInitByDim(key string, width, depth int64) (string, error) {
	rp, err := r.ExecuteCommand("CMS.INITBYDIM", key, width, depth)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// CMSInitByProb command:
// Initialize a Count-Min Sketch with specified error rate and probability
// CMS.INITBYPROB key error_rate probability
func (r *Redis) CMSInitByProb(key string, errorRate, probability float64) (string, error) {
	rp, err := r.ExecuteCommand("CMS.INITBYPROB", key, errorRate, probability)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// CMSIncrBy command:
// Increment count of items by increment
// CMS.INCRBY key item increment [item increment ...]
func (r *Redis) CMSIncrBy(key string, itemIncrements ...interface{}) ([]int64, error) {
	if len(itemIncrements)%2 != 0 {
		return nil, nil // Must have even number of arguments (item/increment pairs)
	}
	
	args := []interface{}{"CMS.INCRBY", key}
	args = append(args, itemIncrements...)
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make([]int64, len(multi))
	for i, reply := range multi {
		result[i], _ = reply.IntegerValue()
	}
	
	return result, nil
}

// CMSQuery command:
// Get count of items
// CMS.QUERY key item [item ...]
func (r *Redis) CMSQuery(key string, items ...interface{}) ([]int64, error) {
	args := []interface{}{"CMS.QUERY", key}
	for _, item := range items {
		args = append(args, item)
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make([]int64, len(multi))
	for i, reply := range multi {
		result[i], _ = reply.IntegerValue()
	}
	
	return result, nil
}

// Additional Bloom Filter Commands

// BFInfo command:
// Get information about a Bloom filter
// BF.INFO key
func (r *Redis) BFInfo(key string) (map[string]interface{}, error) {
	rp, err := r.ExecuteCommand("BF.INFO", key)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]interface{})
	for i := 0; i < len(multi); i += 2 {
		if i+1 < len(multi) {
			key, _ := multi[i].StringValue()
			
			// Handle different value types
			if multi[i+1].Type == 2 { // IntegerReply
				result[key] = multi[i+1].Integer
			} else {
				value, _ := multi[i+1].StringValue()
				result[key] = value
			}
		}
	}
	
	return result, nil
}

// Additional Cuckoo Filter Commands

// CFInfo command:
// Get information about a Cuckoo filter
// CF.INFO key
func (r *Redis) CFInfo(key string) (map[string]interface{}, error) {
	rp, err := r.ExecuteCommand("CF.INFO", key)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]interface{})
	for i := 0; i < len(multi); i += 2 {
		if i+1 < len(multi) {
			key, _ := multi[i].StringValue()
			
			// Handle different value types
			if multi[i+1].Type == 2 { // IntegerReply
				result[key] = multi[i+1].Integer
			} else {
				value, _ := multi[i+1].StringValue()
				result[key] = value
			}
		}
	}
	
	return result, nil
}

// Additional Count-Min Sketch Commands

// CMSInfo command:
// Get information about a Count-Min Sketch
// CMS.INFO key
func (r *Redis) CMSInfo(key string) (map[string]interface{}, error) {
	rp, err := r.ExecuteCommand("CMS.INFO", key)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]interface{})
	for i := 0; i < len(multi); i += 2 {
		if i+1 < len(multi) {
			key, _ := multi[i].StringValue()
			
			// Handle different value types
			if multi[i+1].Type == 2 { // IntegerReply
				result[key] = multi[i+1].Integer
			} else {
				value, _ := multi[i+1].StringValue()
				result[key] = value
			}
		}
	}
	
	return result, nil
}

// CMSMerge command:
// Merge multiple Count-Min Sketches
// CMS.MERGE destkey numkeys source [source ...] [WEIGHTS weight [weight ...]]
func (r *Redis) CMSMerge(destKey string, sourceKeys []string, weights ...float64) (string, error) {
	args := []interface{}{"CMS.MERGE", destKey, len(sourceKeys)}
	
	for _, sourceKey := range sourceKeys {
		args = append(args, sourceKey)
	}
	
	if len(weights) > 0 {
		args = append(args, "WEIGHTS")
		for _, weight := range weights {
			args = append(args, weight)
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}