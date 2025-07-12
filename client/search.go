package client

import (
	"strconv"
	"strings"
)

// FTCreateOptions represents options for FT.CREATE command
type FTCreateOptions struct {
	OnHash      bool                   // Create index on hash documents
	OnJSON      bool                   // Create index on JSON documents
	Prefix      []string               // Document key prefixes
	Filter      string                 // Boolean filter expression
	Language    string                 // Default language
	LanguageField string               // Language field name
	Score       float64                // Default score
	ScoreField  string                 // Score field name
	PayloadField string               // Payload field name
	MaxTextFields bool                 // Enable MAXTEXTFIELDS
	NoOffsets   bool                   // Disable term offset vectors
	Temporary   int                    // Temporary index (seconds)
	NoHL        bool                   // Disable highlighting
	NoFields    bool                   // Don't store field contents
	NoFreqs     bool                   // Disable term frequency vectors
	StopWords   []string               // Custom stop words
	SkipInitialScan bool               // Skip initial document scan
}

// FTFieldSchema represents a field schema for FT.CREATE
type FTFieldSchema struct {
	Name       string
	Type       string // TEXT, NUMERIC, GEO, TAG, VECTOR
	Sortable   bool
	NoStem     bool
	NoIndex    bool
	PhoneticMatcher string
	Weight     float64
	Separator  string // For TAG fields
	Geometry   string // For GEO fields
}

// FTSearchOptions represents options for FT.SEARCH command
type FTSearchOptions struct {
	NoContent     bool     // Don't return document contents
	Verbatim      bool     // Don't use stemming
	NoStopWords   bool     // Don't filter stop words
	WithScores    bool     // Return document scores
	WithPayloads  bool     // Return document payloads
	WithSortKeys  bool     // Return sort keys
	Filter        []FTNumericFilter // Numeric filters
	GeoFilter     *FTGeoFilter      // Geographic filter
	InKeys        []string          // Limit to specific keys
	InFields      []string          // Limit to specific fields
	Return        []string          // Fields to return
	Summarize     *FTSummarizeOptions // Summarization options
	Highlight     *FTHighlightOptions // Highlighting options
	Slop          int               // Query slop for phrase queries
	Timeout       int               // Query timeout (milliseconds)
	InOrder       bool              // Terms must appear in order
	Language      string            // Query language
	Expander      string            // Query expander
	Scorer        string            // Scoring function
	ExplainScore  bool              // Explain score calculation
	Payload       string            // Document payload
	SortBy        string            // Sort by field
	SortOrder     string            // ASC or DESC
	Limit         *FTLimit          // Result pagination
}

// FTNumericFilter represents a numeric filter
type FTNumericFilter struct {
	Field string
	Min   interface{} // number or "-inf"
	Max   interface{} // number or "+inf"
}

// FTGeoFilter represents a geographic filter
type FTGeoFilter struct {
	Field     string
	Longitude float64
	Latitude  float64
	Radius    float64
	Unit      string // m, km, mi, ft
}

// FTSummarizeOptions represents summarization options
type FTSummarizeOptions struct {
	Fields    []string
	Frags     int
	Len       int
	Separator string
}

// FTHighlightOptions represents highlighting options
type FTHighlightOptions struct {
	Fields []string
	Tags   *FTHighlightTags
}

// FTHighlightTags represents highlight tags
type FTHighlightTags struct {
	Open  string
	Close string
}

// FTLimit represents result pagination
type FTLimit struct {
	Offset int
	Num    int
}

// FTAggregateOptions represents options for FT.AGGREGATE command
type FTAggregateOptions struct {
	Verbatim    bool
	Load        []string
	Timeout     int
	GroupBy     *FTGroupBy
	SortBy      []FTSortBy
	Apply       []FTApply
	Limit       *FTLimit
	Filter      string
}

// FTGroupBy represents GROUP BY clause
type FTGroupBy struct {
	Fields  []string
	Reduce  []FTReduce
}

// FTReduce represents reduction function
type FTReduce struct {
	Function string
	Args     []string
	As       string
}

// FTSortBy represents sort criteria
type FTSortBy struct {
	Property string
	Order    string // ASC or DESC
}

// FTApply represents APPLY clause
type FTApply struct {
	Expression string
	As         string
}

// Index Management Commands

// FTCreate command:
// Create an index with the given specification
// FT.CREATE index [ON HASH | JSON] [PREFIX count prefix [prefix ...]] schema
func (r *Redis) FTCreate(index string, schema []FTFieldSchema, options ...*FTCreateOptions) (string, error) {
	args := []interface{}{"FT.CREATE", index}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.OnHash {
			args = append(args, "ON", "HASH")
		} else if opt.OnJSON {
			args = append(args, "ON", "JSON")
		}
		
		if len(opt.Prefix) > 0 {
			args = append(args, "PREFIX", len(opt.Prefix))
			for _, prefix := range opt.Prefix {
				args = append(args, prefix)
			}
		}
		
		if opt.Filter != "" {
			args = append(args, "FILTER", opt.Filter)
		}
		
		if opt.Language != "" {
			args = append(args, "LANGUAGE", opt.Language)
		}
		
		if opt.LanguageField != "" {
			args = append(args, "LANGUAGE_FIELD", opt.LanguageField)
		}
		
		if opt.Score > 0 {
			args = append(args, "SCORE", opt.Score)
		}
		
		if opt.ScoreField != "" {
			args = append(args, "SCORE_FIELD", opt.ScoreField)
		}
		
		if opt.PayloadField != "" {
			args = append(args, "PAYLOAD_FIELD", opt.PayloadField)
		}
		
		if opt.MaxTextFields {
			args = append(args, "MAXTEXTFIELDS")
		}
		
		if opt.NoOffsets {
			args = append(args, "NOOFFSETS")
		}
		
		if opt.Temporary > 0 {
			args = append(args, "TEMPORARY", opt.Temporary)
		}
		
		if opt.NoHL {
			args = append(args, "NOHL")
		}
		
		if opt.NoFields {
			args = append(args, "NOFIELDS")
		}
		
		if opt.NoFreqs {
			args = append(args, "NOFREQS")
		}
		
		if len(opt.StopWords) > 0 {
			args = append(args, "STOPWORDS", len(opt.StopWords))
			for _, word := range opt.StopWords {
				args = append(args, word)
			}
		}
		
		if opt.SkipInitialScan {
			args = append(args, "SKIPINITIALSCAN")
		}
	}
	
	// Add schema
	args = append(args, "SCHEMA")
	for _, field := range schema {
		args = append(args, field.Name, field.Type)
		
		if field.Sortable {
			args = append(args, "SORTABLE")
		}
		if field.NoStem {
			args = append(args, "NOSTEM")
		}
		if field.NoIndex {
			args = append(args, "NOINDEX")
		}
		if field.PhoneticMatcher != "" {
			args = append(args, "PHONETIC", field.PhoneticMatcher)
		}
		if field.Weight > 0 {
			args = append(args, "WEIGHT", field.Weight)
		}
		if field.Separator != "" {
			args = append(args, "SEPARATOR", field.Separator)
		}
		if field.Geometry != "" {
			args = append(args, "GEOMETRY", field.Geometry)
		}
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// FTDropIndex command:
// Delete an index
// FT.DROPINDEX index [DD]
func (r *Redis) FTDropIndex(index string, deleteDocuments ...bool) (string, error) {
	args := []interface{}{"FT.DROPINDEX", index}
	if len(deleteDocuments) > 0 && deleteDocuments[0] {
		args = append(args, "DD")
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// FTInfo command:
// Return information and statistics on the index
// FT.INFO index
func (r *Redis) FTInfo(index string) (map[string]interface{}, error) {
	rp, err := r.ExecuteCommand("FT.INFO", index)
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
			value, _ := multi[i+1].StringValue()
			result[key] = value
		}
	}
	
	return result, nil
}

// Search Operations

// FTSearch command:
// Search the index with a textual query
// FT.SEARCH index query [options...]
func (r *Redis) FTSearch(index, query string, options ...*FTSearchOptions) ([]interface{}, error) {
	args := []interface{}{"FT.SEARCH", index, query}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.NoContent {
			args = append(args, "NOCONTENT")
		}
		if opt.Verbatim {
			args = append(args, "VERBATIM")
		}
		if opt.NoStopWords {
			args = append(args, "NOSTOPWORDS")
		}
		if opt.WithScores {
			args = append(args, "WITHSCORES")
		}
		if opt.WithPayloads {
			args = append(args, "WITHPAYLOADS")
		}
		if opt.WithSortKeys {
			args = append(args, "WITHSORTKEYS")
		}
		
		// Add filters
		for _, filter := range opt.Filter {
			args = append(args, "FILTER", filter.Field, filter.Min, filter.Max)
		}
		
		if opt.GeoFilter != nil {
			args = append(args, "GEOFILTER", opt.GeoFilter.Field, 
				opt.GeoFilter.Longitude, opt.GeoFilter.Latitude,
				opt.GeoFilter.Radius, opt.GeoFilter.Unit)
		}
		
		if len(opt.InKeys) > 0 {
			args = append(args, "INKEYS", len(opt.InKeys))
			for _, key := range opt.InKeys {
				args = append(args, key)
			}
		}
		
		if len(opt.InFields) > 0 {
			args = append(args, "INFIELDS", len(opt.InFields))
			for _, field := range opt.InFields {
				args = append(args, field)
			}
		}
		
		if len(opt.Return) > 0 {
			args = append(args, "RETURN", len(opt.Return))
			for _, field := range opt.Return {
				args = append(args, field)
			}
		}
		
		if opt.Summarize != nil {
			args = append(args, "SUMMARIZE")
			if len(opt.Summarize.Fields) > 0 {
				args = append(args, "FIELDS", len(opt.Summarize.Fields))
				for _, field := range opt.Summarize.Fields {
					args = append(args, field)
				}
			}
			if opt.Summarize.Frags > 0 {
				args = append(args, "FRAGS", opt.Summarize.Frags)
			}
			if opt.Summarize.Len > 0 {
				args = append(args, "LEN", opt.Summarize.Len)
			}
			if opt.Summarize.Separator != "" {
				args = append(args, "SEPARATOR", opt.Summarize.Separator)
			}
		}
		
		if opt.Highlight != nil {
			args = append(args, "HIGHLIGHT")
			if len(opt.Highlight.Fields) > 0 {
				args = append(args, "FIELDS", len(opt.Highlight.Fields))
				for _, field := range opt.Highlight.Fields {
					args = append(args, field)
				}
			}
			if opt.Highlight.Tags != nil {
				args = append(args, "TAGS", opt.Highlight.Tags.Open, opt.Highlight.Tags.Close)
			}
		}
		
		if opt.Slop > 0 {
			args = append(args, "SLOP", opt.Slop)
		}
		if opt.Timeout > 0 {
			args = append(args, "TIMEOUT", opt.Timeout)
		}
		if opt.InOrder {
			args = append(args, "INORDER")
		}
		if opt.Language != "" {
			args = append(args, "LANGUAGE", opt.Language)
		}
		if opt.Expander != "" {
			args = append(args, "EXPANDER", opt.Expander)
		}
		if opt.Scorer != "" {
			args = append(args, "SCORER", opt.Scorer)
		}
		if opt.ExplainScore {
			args = append(args, "EXPLAINSCORE")
		}
		if opt.Payload != "" {
			args = append(args, "PAYLOAD", opt.Payload)
		}
		if opt.SortBy != "" {
			args = append(args, "SORTBY", opt.SortBy)
			if opt.SortOrder != "" {
				args = append(args, opt.SortOrder)
			}
		}
		if opt.Limit != nil {
			args = append(args, "LIMIT", opt.Limit.Offset, opt.Limit.Num)
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
	
	result := make([]interface{}, len(multi))
	for i, reply := range multi {
		if reply.Type == 1 { // BulkReply
			result[i], _ = reply.StringValue()
		} else if reply.Type == 2 { // IntegerReply  
			result[i] = reply.Integer
		} else if reply.Type == 4 { // MultiReply
			subResult := make([]interface{}, len(reply.Multi))
			for j, subReply := range reply.Multi {
				subResult[j], _ = subReply.StringValue()
			}
			result[i] = subResult
		} else {
			result[i], _ = reply.StringValue()
		}
	}
	
	return result, nil
}

// FTAggregate command:
// Run a search query and perform aggregate transformations on the results
// FT.AGGREGATE index query [options...]
func (r *Redis) FTAggregate(index, query string, options ...*FTAggregateOptions) ([]interface{}, error) {
	args := []interface{}{"FT.AGGREGATE", index, query}
	
	if len(options) > 0 && options[0] != nil {
		opt := options[0]
		
		if opt.Verbatim {
			args = append(args, "VERBATIM")
		}
		
		if len(opt.Load) > 0 {
			args = append(args, "LOAD", len(opt.Load))
			for _, field := range opt.Load {
				args = append(args, field)
			}
		}
		
		if opt.Timeout > 0 {
			args = append(args, "TIMEOUT", opt.Timeout)
		}
		
		if opt.GroupBy != nil {
			args = append(args, "GROUPBY", len(opt.GroupBy.Fields))
			for _, field := range opt.GroupBy.Fields {
				args = append(args, field)
			}
			
			for _, reduce := range opt.GroupBy.Reduce {
				args = append(args, "REDUCE", reduce.Function)
				args = append(args, len(reduce.Args))
				for _, arg := range reduce.Args {
					args = append(args, arg)
				}
				if reduce.As != "" {
					args = append(args, "AS", reduce.As)
				}
			}
		}
		
		if len(opt.SortBy) > 0 {
			args = append(args, "SORTBY", len(opt.SortBy)*2)
			for _, sort := range opt.SortBy {
				args = append(args, sort.Property)
				if sort.Order != "" {
					args = append(args, sort.Order)
				} else {
					args = append(args, "ASC")
				}
			}
		}
		
		for _, apply := range opt.Apply {
			args = append(args, "APPLY", apply.Expression)
			if apply.As != "" {
				args = append(args, "AS", apply.As)
			}
		}
		
		if opt.Limit != nil {
			args = append(args, "LIMIT", opt.Limit.Offset, opt.Limit.Num)
		}
		
		if opt.Filter != "" {
			args = append(args, "FILTER", opt.Filter)
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
	
	result := make([]interface{}, len(multi))
	for i, reply := range multi {
		if reply.Type == 4 { // MultiReply
			subResult := make([]interface{}, len(reply.Multi))
			for j, subReply := range reply.Multi {
				subResult[j], _ = subReply.StringValue()
			}
			result[i] = subResult
		} else {
			result[i], _ = reply.StringValue()
		}
	}
	
	return result, nil
}

// FTExplain command:
// Return the execution plan for a complex query
// FT.EXPLAIN index query [DIALECT dialect]
func (r *Redis) FTExplain(index, query string, dialect ...int) (string, error) {
	args := []interface{}{"FT.EXPLAIN", index, query}
	if len(dialect) > 0 {
		args = append(args, "DIALECT", dialect[0])
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// Document Management (deprecated in RediSearch 2.0+)

// FTAdd command:
// Add a document to the index (deprecated)
// FT.ADD index docId score [NOSAVE] [REPLACE] [PARTIAL] [LANGUAGE language] [PAYLOAD payload] [IF condition] FIELDS field value [field value ...]
func (r *Redis) FTAdd(index, docID string, score float64, fields map[string]interface{}, options ...string) (string, error) {
	args := []interface{}{"FT.ADD", index, docID, score}
	
	// Add options
	for _, option := range options {
		switch strings.ToUpper(option) {
		case "NOSAVE", "REPLACE", "PARTIAL":
			args = append(args, option)
		}
	}
	
	// Add fields
	args = append(args, "FIELDS")
	for field, value := range fields {
		args = append(args, field, value)
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return "", err
	}
	return rp.StringValue()
}

// FTDel command:
// Delete a document from the index (deprecated)
// FT.DEL index docId [DD]
func (r *Redis) FTDel(index, docID string, deleteDocument ...bool) (int64, error) {
	args := []interface{}{"FT.DEL", index, docID}
	if len(deleteDocument) > 0 && deleteDocument[0] {
		args = append(args, "DD")
	}
	
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}