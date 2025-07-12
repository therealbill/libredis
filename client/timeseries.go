package client

import (
	"strconv"
)

// TSCreateOptions represents options for TS.CREATE command
type TSCreateOptions struct {
	RetentionMsecs   int64             // Retention period in milliseconds
	ChunkSize        int               // Chunk size for compressed data
	DuplicatePolicy  string            // Policy for handling duplicates
	Labels           map[string]string // Labels for the time series
}

// TSAddOptions represents options for TS.ADD command
type TSAddOptions struct {
	RetentionMsecs  int64             // Retention period in milliseconds
	ChunkSize       int               // Chunk size for compressed data
	DuplicatePolicy string            // Policy for handling duplicates
	Labels          map[string]string // Labels for the time series
	OnDuplicate     string            // Action on duplicate timestamp
}

// TSIncrByOptions represents options for TS.INCRBY command
type TSIncrByOptions struct {
	Timestamp       int64             // Explicit timestamp
	RetentionMsecs  int64             // Retention period in milliseconds
	ChunkSize       int               // Chunk size for compressed data
	Labels          map[string]string // Labels for the time series
}

// TSDecrByOptions represents options for TS.DECRBY command
type TSDecrByOptions struct {
	Timestamp       int64             // Explicit timestamp
	RetentionMsecs  int64             // Retention period in milliseconds
	ChunkSize       int               // Chunk size for compressed data
	Labels          map[string]string // Labels for the time series
}

// TSMAddSample represents a sample for TS.MADD command
type TSMAddSample struct {
	Key       string
	Timestamp int64
	Value     float64
}

// TSRangeOptions represents options for TS.RANGE command
type TSRangeOptions struct {
	Count       int              // Maximum number of samples
	Aggregation *TSAggregation   // Aggregation function
	FilterBy    *TSFilterBy      // Filter by value
}

// TSMRangeOptions represents options for TS.MRANGE command
type TSMRangeOptions struct {
	Count       int              // Maximum number of samples
	Aggregation *TSAggregation   // Aggregation function
	FilterBy    *TSFilterBy      // Filter by value
	WithLabels  bool             // Include labels in response
	SelectedLabels []string      // Specific labels to include
	GroupBy     *TSGroupBy       // Group by labels
}

// TSAggregation represents aggregation options
type TSAggregation struct {
	Type       string // avg, sum, min, max, range, count, std.p, std.s, var.p, var.s, first, last
	TimeBucket int64  // Time bucket for aggregation
}

// TSFilterBy represents filter options
type TSFilterBy struct {
	Min float64
	Max float64
}

// TSGroupBy represents grouping options
type TSGroupBy struct {
	Label  string
	Reduce string // sum, min, max, avg, std.p, std.s, var.p, var.s, count, range
}

// TSSample represents a time series sample
type TSSample struct {
	Timestamp int64
	Value     float64
}

// TSInfo represents time series information
type TSInfo struct {
	TotalSamples     int64
	MemoryUsage      int64
	FirstTimestamp   int64
	LastTimestamp    int64
	RetentionTime    int64
	ChunkCount       int64
	ChunkSize        int64
	DuplicatePolicy  string
	Labels           map[string]string
	SourceKey        string
	Rules            []TSRule
}

// TSRule represents a downsampling rule
type TSRule struct {
	DestKey     string
	TimeBucket  int64
	Aggregation string
}

// Basic Time Series Operations

// TSCreate command:
// Create a new time series
// TS.CREATE key [RETENTION retentionTime] [CHUNK_SIZE size] [DUPLICATE_POLICY policy] [LABELS label value ...]
func (r *Redis) TSCreate(key string, options ...*TSCreateOptions) (string, error) {
	args := []interface{}{"TS.CREATE", key}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.RetentionMsecs > 0 {
			args = append(args, "RETENTION", opt.RetentionMsecs)
		}
		
		if opt.ChunkSize > 0 {
			args = append(args, "CHUNK_SIZE", opt.ChunkSize)
		}
		
		if opt.DuplicatePolicy != "" {
			args = append(args, "DUPLICATE_POLICY", opt.DuplicatePolicy)
		}
		
		if len(opt.Labels) > 0 {
			args = append(args, "LABELS")
			for label, value := range opt.Labels {
				args = append(args, label, value)
			}
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// TSAdd command:
// Add a sample to a time series
// TS.ADD key timestamp value [RETENTION retentionTime] [CHUNK_SIZE size] [ON_DUPLICATE policy] [LABELS label value ...]
func (r *Redis) TSAdd(key string, timestamp int64, value float64, options ...*TSAddOptions) (int64, error) {
	args := []interface{}{"TS.ADD", key, timestamp, value}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.RetentionMsecs > 0 {
			args = append(args, "RETENTION", opt.RetentionMsecs)
		}
		
		if opt.ChunkSize > 0 {
			args = append(args, "CHUNK_SIZE", opt.ChunkSize)
		}
		
		if opt.OnDuplicate != "" {
			args = append(args, "ON_DUPLICATE", opt.OnDuplicate)
		}
		
		if len(opt.Labels) > 0 {
			args = append(args, "LABELS")
			for label, value := range opt.Labels {
				args = append(args, label, value)
			}
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// TSMAdd command:
// Add multiple samples to multiple time series
// TS.MADD key timestamp value [key timestamp value ...]
func (r *Redis) TSMAdd(samples ...TSMAddSample) ([]int64, error) {
	if len(samples) == 0 {
		return nil, nil
	}
	
	args := []interface{}{"TS.MADD"}
	for _, sample := range samples {
		args = append(args, sample.Key, sample.Timestamp, sample.Value)
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

// TSIncrBy command:
// Increment the value of a sample
// TS.INCRBY key value [TIMESTAMP timestamp] [RETENTION retentionTime] [CHUNK_SIZE size] [LABELS label value ...]
func (r *Redis) TSIncrBy(key string, value float64, options ...*TSIncrByOptions) (int64, error) {
	args := []interface{}{"TS.INCRBY", key, value}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.Timestamp > 0 {
			args = append(args, "TIMESTAMP", opt.Timestamp)
		}
		
		if opt.RetentionMsecs > 0 {
			args = append(args, "RETENTION", opt.RetentionMsecs)
		}
		
		if opt.ChunkSize > 0 {
			args = append(args, "CHUNK_SIZE", opt.ChunkSize)
		}
		
		if len(opt.Labels) > 0 {
			args = append(args, "LABELS")
			for label, value := range opt.Labels {
				args = append(args, label, value)
			}
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// TSDecrBy command:
// Decrement the value of a sample
// TS.DECRBY key value [TIMESTAMP timestamp] [RETENTION retentionTime] [CHUNK_SIZE size] [LABELS label value ...]
func (r *Redis) TSDecrBy(key string, value float64, options ...*TSDecrByOptions) (int64, error) {
	args := []interface{}{"TS.DECRBY", key, value}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.Timestamp > 0 {
			args = append(args, "TIMESTAMP", opt.Timestamp)
		}
		
		if opt.RetentionMsecs > 0 {
			args = append(args, "RETENTION", opt.RetentionMsecs)
		}
		
		if opt.ChunkSize > 0 {
			args = append(args, "CHUNK_SIZE", opt.ChunkSize)
		}
		
		if len(opt.Labels) > 0 {
			args = append(args, "LABELS")
			for label, value := range opt.Labels {
				args = append(args, label, value)
			}
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// Query Operations

// TSRange command:
// Query a range of samples from a time series
// TS.RANGE key fromTimestamp toTimestamp [COUNT count] [AGGREGATION aggregationType timeBucket] [FILTER_BY_TS timestamp ...]
func (r *Redis) TSRange(key string, fromTimestamp, toTimestamp int64, options ...*TSRangeOptions) ([]TSSample, error) {
	args := []interface{}{"TS.RANGE", key, fromTimestamp, toTimestamp}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.Count > 0 {
			args = append(args, "COUNT", opt.Count)
		}
		
		if opt.Aggregation != nil {
			args = append(args, "AGGREGATION", opt.Aggregation.Type, opt.Aggregation.TimeBucket)
		}
		
		if opt.FilterBy != nil {
			args = append(args, "FILTER_BY_VALUE", opt.FilterBy.Min, opt.FilterBy.Max)
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make([]TSSample, len(multi))
	for i, reply := range multi {
		sampleMulti, _ := reply.MultiValue()
		if len(sampleMulti) >= 2 {
			timestamp, _ := sampleMulti[0].IntegerValue()
			valueStr, _ := sampleMulti[1].StringValue()
			value, _ := strconv.ParseFloat(valueStr, 64)
			
			result[i] = TSSample{
				Timestamp: timestamp,
				Value:     value,
			}
		}
	}
	
	return result, nil
}

// TSRevRange command:
// Query a range of samples from a time series in reverse order
// TS.REVRANGE key fromTimestamp toTimestamp [COUNT count] [AGGREGATION aggregationType timeBucket] [FILTER_BY_TS timestamp ...]
func (r *Redis) TSRevRange(key string, fromTimestamp, toTimestamp int64, options ...*TSRangeOptions) ([]TSSample, error) {
	args := []interface{}{"TS.REVRANGE", key, fromTimestamp, toTimestamp}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.Count > 0 {
			args = append(args, "COUNT", opt.Count)
		}
		
		if opt.Aggregation != nil {
			args = append(args, "AGGREGATION", opt.Aggregation.Type, opt.Aggregation.TimeBucket)
		}
		
		if opt.FilterBy != nil {
			args = append(args, "FILTER_BY_VALUE", opt.FilterBy.Min, opt.FilterBy.Max)
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make([]TSSample, len(multi))
	for i, reply := range multi {
		sampleMulti, _ := reply.MultiValue()
		if len(sampleMulti) >= 2 {
			timestamp, _ := sampleMulti[0].IntegerValue()
			valueStr, _ := sampleMulti[1].StringValue()
			value, _ := strconv.ParseFloat(valueStr, 64)
			
			result[i] = TSSample{
				Timestamp: timestamp,
				Value:     value,
			}
		}
	}
	
	return result, nil
}

// TSMRange command:
// Query a range of samples from multiple time series
// TS.MRANGE fromTimestamp toTimestamp [COUNT count] [AGGREGATION aggregationType timeBucket] [WITHLABELS | SELECTED_LABELS label ...] [GROUPBY label REDUCE reducer] [FILTER_BY_TS timestamp ...] FILTER filter ...
func (r *Redis) TSMRange(fromTimestamp, toTimestamp int64, filters []string, options ...*TSMRangeOptions) (map[string][]TSSample, error) {
	args := []interface{}{"TS.MRANGE", fromTimestamp, toTimestamp}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.Count > 0 {
			args = append(args, "COUNT", opt.Count)
		}
		
		if opt.Aggregation != nil {
			args = append(args, "AGGREGATION", opt.Aggregation.Type, opt.Aggregation.TimeBucket)
		}
		
		if opt.WithLabels {
			args = append(args, "WITHLABELS")
		} else if len(opt.SelectedLabels) > 0 {
			args = append(args, "SELECTED_LABELS")
			for _, label := range opt.SelectedLabels {
				args = append(args, label)
			}
		}
		
		if opt.GroupBy != nil {
			args = append(args, "GROUPBY", opt.GroupBy.Label, "REDUCE", opt.GroupBy.Reduce)
		}
		
		if opt.FilterBy != nil {
			args = append(args, "FILTER_BY_VALUE", opt.FilterBy.Min, opt.FilterBy.Max)
		}
	}
	
	// Add filters
	args = append(args, "FILTER")
	for _, filter := range filters {
		args = append(args, filter)
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make(map[string][]TSSample)
	for _, reply := range multi {
		seriesMulti, _ := reply.MultiValue()
		if len(seriesMulti) >= 2 {
			key, _ := seriesMulti[0].StringValue()
			samplesMulti, _ := seriesMulti[len(seriesMulti)-1].MultiValue()
			
			samples := make([]TSSample, len(samplesMulti))
			for i, sampleReply := range samplesMulti {
				sampleMulti, _ := sampleReply.MultiValue()
				if len(sampleMulti) >= 2 {
					timestamp, _ := sampleMulti[0].IntegerValue()
					valueStr, _ := sampleMulti[1].StringValue()
					value, _ := strconv.ParseFloat(valueStr, 64)
					
					samples[i] = TSSample{
						Timestamp: timestamp,
						Value:     value,
					}
				}
			}
			
			result[key] = samples
		}
	}
	
	return result, nil
}

// TSMRevRange command:
// Query a range of samples from multiple time series in reverse order
// TS.MREVRANGE fromTimestamp toTimestamp [COUNT count] [AGGREGATION aggregationType timeBucket] [WITHLABELS | SELECTED_LABELS label ...] [GROUPBY label REDUCE reducer] [FILTER_BY_TS timestamp ...] FILTER filter ...
func (r *Redis) TSMRevRange(fromTimestamp, toTimestamp int64, filters []string, options ...*TSMRangeOptions) (map[string][]TSSample, error) {
	args := []interface{}{"TS.MREVRANGE", fromTimestamp, toTimestamp}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.Count > 0 {
			args = append(args, "COUNT", opt.Count)
		}
		
		if opt.Aggregation != nil {
			args = append(args, "AGGREGATION", opt.Aggregation.Type, opt.Aggregation.TimeBucket)
		}
		
		if opt.WithLabels {
			args = append(args, "WITHLABELS")
		} else if len(opt.SelectedLabels) > 0 {
			args = append(args, "SELECTED_LABELS")
			for _, label := range opt.SelectedLabels {
				args = append(args, label)
			}
		}
		
		if opt.GroupBy != nil {
			args = append(args, "GROUPBY", opt.GroupBy.Label, "REDUCE", opt.GroupBy.Reduce)
		}
		
		if opt.FilterBy != nil {
			args = append(args, "FILTER_BY_VALUE", opt.FilterBy.Min, opt.FilterBy.Max)
		}
	}
	
	// Add filters
	args = append(args, "FILTER")
	for _, filter := range filters {
		args = append(args, filter)
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	result := make(map[string][]TSSample)
	for _, reply := range multi {
		seriesMulti, _ := reply.MultiValue()
		if len(seriesMulti) >= 2 {
			key, _ := seriesMulti[0].StringValue()
			samplesMulti, _ := seriesMulti[len(seriesMulti)-1].MultiValue()
			
			samples := make([]TSSample, len(samplesMulti))
			for i, sampleReply := range samplesMulti {
				sampleMulti, _ := sampleReply.MultiValue()
				if len(sampleMulti) >= 2 {
					timestamp, _ := sampleMulti[0].IntegerValue()
					valueStr, _ := sampleMulti[1].StringValue()
					value, _ := strconv.ParseFloat(valueStr, 64)
					
					samples[i] = TSSample{
						Timestamp: timestamp,
						Value:     value,
					}
				}
			}
			
			result[key] = samples
		}
	}
	
	return result, nil
}

// Metadata Operations

// TSInfo command:
// Get information and statistics for a time series
// TS.INFO key
func (r *Redis) TSInfo(key string) (*TSInfo, error) {
	rp, err := r.ExecuteCommand("TS.INFO", key)
	if err != nil {
		return nil, err
	}
	
	multi, err := rp.MultiValue()
	if err != nil {
		return nil, err
	}
	
	info := &TSInfo{
		Labels: make(map[string]string),
		Rules:  make([]TSRule, 0),
	}
	
	for i := 0; i < len(multi); i += 2 {
		if i+1 < len(multi) {
			key, _ := multi[i].StringValue()
			
			switch key {
			case "totalSamples":
				info.TotalSamples, _ = multi[i+1].IntegerValue()
			case "memoryUsage":
				info.MemoryUsage, _ = multi[i+1].IntegerValue()
			case "firstTimestamp":
				info.FirstTimestamp, _ = multi[i+1].IntegerValue()
			case "lastTimestamp":
				info.LastTimestamp, _ = multi[i+1].IntegerValue()
			case "retentionTime":
				info.RetentionTime, _ = multi[i+1].IntegerValue()
			case "chunkCount":
				info.ChunkCount, _ = multi[i+1].IntegerValue()
			case "chunkSize":
				info.ChunkSize, _ = multi[i+1].IntegerValue()
			case "duplicatePolicy":
				info.DuplicatePolicy, _ = multi[i+1].StringValue()
			case "sourceKey":
				info.SourceKey, _ = multi[i+1].StringValue()
			case "labels":
				labelsMulti, _ := multi[i+1].MultiValue()
				for j := 0; j < len(labelsMulti); j += 2 {
					if j+1 < len(labelsMulti) {
						labelKey, _ := labelsMulti[j].StringValue()
						labelValue, _ := labelsMulti[j+1].StringValue()
						info.Labels[labelKey] = labelValue
					}
				}
			case "rules":
				rulesMulti, _ := multi[i+1].MultiValue()
				for _, ruleReply := range rulesMulti {
					ruleMulti, _ := ruleReply.MultiValue()
					if len(ruleMulti) >= 3 {
						rule := TSRule{}
						rule.DestKey, _ = ruleMulti[0].StringValue()
						rule.TimeBucket, _ = ruleMulti[1].IntegerValue()
						rule.Aggregation, _ = ruleMulti[2].StringValue()
						info.Rules = append(info.Rules, rule)
					}
				}
			}
		}
	}
	
	return info, nil
}