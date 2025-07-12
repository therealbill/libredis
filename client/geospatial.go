package client

import (
	"strconv"
)

// Geospatial constants
const (
	GeoUnitMeters     = "M"
	GeoUnitKilometers = "KM"
	GeoUnitFeet       = "FT"
	GeoUnitMiles      = "MI"

	GeoOrderAsc  = "ASC"
	GeoOrderDesc = "DESC"
)

// GeoMember represents a geospatial member with coordinates
type GeoMember struct {
	Longitude float64
	Latitude  float64
	Member    string
}

// GeoAddOptions represents options for GEOADD command
type GeoAddOptions struct {
	NX bool // NX option - only add new elements
	XX bool // XX option - only update existing elements
	CH bool // CH option - return count of changed elements
}

// GeoCoordinate represents longitude/latitude coordinates
type GeoCoordinate struct {
	Longitude float64
	Latitude  float64
}

// GeoSearchOptions represents options for GEOSEARCH command
type GeoSearchOptions struct {
	// Search center (exactly one must be specified)
	FromMember *string        // FROMMEMBER option
	FromLonLat *GeoCoordinate // FROMLONLAT option

	// Search area (exactly one must be specified)
	ByRadius *GeoRadius // BYRADIUS option
	ByBox    *GeoBox    // BYBOX option

	// Result options
	Order     string // ASC or DESC
	Count     int64  // COUNT option
	Any       bool   // ANY option
	WithCoord bool   // WITHCOORD option
	WithDist  bool   // WITHDIST option
	WithHash  bool   // WITHHASH option
}

// GeoSearchStoreOptions represents options for GEOSEARCHSTORE command
type GeoSearchStoreOptions struct {
	GeoSearchOptions
	StoreDist bool // STOREDIST option
}

// GeoRadius represents a radius search parameter
type GeoRadius struct {
	Radius float64
	Unit   string
}

// GeoBox represents a box search parameter
type GeoBox struct {
	Width  float64
	Height float64
	Unit   string
}

// GeoRadiusOptions represents options for legacy GEORADIUS commands
type GeoRadiusOptions struct {
	WithCoord bool   // WITHCOORD option
	WithDist  bool   // WITHDIST option
	WithHash  bool   // WITHHASH option
	Count     int64  // COUNT option
	Any       bool   // ANY option
	Order     string // ASC or DESC
	Store     string // STORE option
	StoreDist string // STOREDIST option
}

// GeoLocation represents a geospatial search result
type GeoLocation struct {
	Member      string
	Coordinates *GeoCoordinate
	Distance    *float64
	Hash        *int64
}

// Basic Operations

// GEOADD key [NX|XX] [CH] longitude latitude member [longitude latitude member ...]
// GeoAdd adds geospatial items to a geospatial index.
func (r *Redis) GeoAdd(key string, members []GeoMember) (int64, error) {
	args := []interface{}{"GEOADD", key}
	for _, member := range members {
		args = append(args, member.Longitude, member.Latitude, member.Member)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// GeoAddWithOptions adds geospatial items with additional options.
func (r *Redis) GeoAddWithOptions(key string, members []GeoMember, opts GeoAddOptions) (int64, error) {
	args := []interface{}{"GEOADD", key}

	if opts.NX {
		args = append(args, "NX")
	}
	if opts.XX {
		args = append(args, "XX")
	}
	if opts.CH {
		args = append(args, "CH")
	}

	for _, member := range members {
		args = append(args, member.Longitude, member.Latitude, member.Member)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// GEODIST key member1 member2 [M|KM|FT|MI]
// GeoDist returns the distance between two geospatial members.
func (r *Redis) GeoDist(key, member1, member2 string) (float64, error) {
	return r.GeoDistWithUnit(key, member1, member2, GeoUnitMeters)
}

// GeoDistWithUnit returns the distance with a specific unit.
func (r *Redis) GeoDistWithUnit(key, member1, member2, unit string) (float64, error) {
	args := packArgs("GEODIST", key, member1, member2, unit)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}

	if rp.Type == BulkReply && rp.Bulk == nil {
		return 0, nil // One or both members don't exist
	}

	distStr, err := rp.StringValue()
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(distStr, 64)
}

// GEOHASH key member [member ...]
// GeoHash returns geohash strings for the specified members.
func (r *Redis) GeoHash(key string, members ...string) ([]string, error) {
	args := packArgs("GEOHASH", key, members)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	if rp.Type == MultiReply {
		result := make([]string, len(rp.Multi))
		for i, item := range rp.Multi {
			if item.Type == BulkReply && item.Bulk != nil {
				result[i], _ = item.StringValue()
			} else {
				result[i] = "" // Member doesn't exist
			}
		}
		return result, nil
	}

	return nil, nil
}

// GEOPOS key member [member ...]
// GeoPos returns coordinates for the specified members.
func (r *Redis) GeoPos(key string, members ...string) ([]*GeoCoordinate, error) {
	args := packArgs("GEOPOS", key, members)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	if rp.Type == MultiReply {
		result := make([]*GeoCoordinate, len(rp.Multi))
		for i, item := range rp.Multi {
			if item.Type == MultiReply && item.Multi != nil && len(item.Multi) >= 2 {
				lonStr, _ := item.Multi[0].StringValue()
				latStr, _ := item.Multi[1].StringValue()
				lon, _ := strconv.ParseFloat(lonStr, 64)
				lat, _ := strconv.ParseFloat(latStr, 64)
				result[i] = &GeoCoordinate{
					Longitude: lon,
					Latitude:  lat,
				}
			} else {
				result[i] = nil // Member doesn't exist
			}
		}
		return result, nil
	}

	return nil, nil
}

// Modern Search Commands

// GEOSEARCH key [FROMMEMBER member] [FROMLONLAT longitude latitude] [BYRADIUS radius M|KM|FT|MI] [BYBOX width height M|KM|FT|MI] [ASC|DESC] [COUNT count [ANY]] [WITHCOORD] [WITHDIST] [WITHHASH]
// GeoSearch queries a geospatial index for members within a specified area.
func (r *Redis) GeoSearch(key string, opts GeoSearchOptions) ([]GeoLocation, error) {
	args := []interface{}{"GEOSEARCH", key}

	// Add search center
	if opts.FromMember != nil {
		args = append(args, "FROMMEMBER", *opts.FromMember)
	} else if opts.FromLonLat != nil {
		args = append(args, "FROMLONLAT", opts.FromLonLat.Longitude, opts.FromLonLat.Latitude)
	}

	// Add search area
	if opts.ByRadius != nil {
		args = append(args, "BYRADIUS", opts.ByRadius.Radius, opts.ByRadius.Unit)
	} else if opts.ByBox != nil {
		args = append(args, "BYBOX", opts.ByBox.Width, opts.ByBox.Height, opts.ByBox.Unit)
	}

	// Add result options
	if opts.Order != "" {
		args = append(args, opts.Order)
	}

	if opts.Count > 0 {
		args = append(args, "COUNT", opts.Count)
		if opts.Any {
			args = append(args, "ANY")
		}
	}

	if opts.WithCoord {
		args = append(args, "WITHCOORD")
	}
	if opts.WithDist {
		args = append(args, "WITHDIST")
	}
	if opts.WithHash {
		args = append(args, "WITHHASH")
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	return parseGeoLocations(rp.Multi, opts.WithCoord, opts.WithDist, opts.WithHash)
}

// GEOSEARCHSTORE destination source [FROMMEMBER member] [FROMLONLAT longitude latitude] [BYRADIUS radius M|KM|FT|MI] [BYBOX width height M|KM|FT|MI] [ASC|DESC] [COUNT count [ANY]] [STOREDIST]
// GeoSearchStore executes a geospatial search and stores results in another key.
func (r *Redis) GeoSearchStore(destination, source string, opts GeoSearchStoreOptions) (int64, error) {
	args := []interface{}{"GEOSEARCHSTORE", destination, source}

	// Add search center
	if opts.FromMember != nil {
		args = append(args, "FROMMEMBER", *opts.FromMember)
	} else if opts.FromLonLat != nil {
		args = append(args, "FROMLONLAT", opts.FromLonLat.Longitude, opts.FromLonLat.Latitude)
	}

	// Add search area
	if opts.ByRadius != nil {
		args = append(args, "BYRADIUS", opts.ByRadius.Radius, opts.ByRadius.Unit)
	} else if opts.ByBox != nil {
		args = append(args, "BYBOX", opts.ByBox.Width, opts.ByBox.Height, opts.ByBox.Unit)
	}

	// Add result options
	if opts.Order != "" {
		args = append(args, opts.Order)
	}

	if opts.Count > 0 {
		args = append(args, "COUNT", opts.Count)
		if opts.Any {
			args = append(args, "ANY")
		}
	}

	if opts.StoreDist {
		args = append(args, "STOREDIST")
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return 0, err
	}
	return rp.IntegerValue()
}

// Legacy Search Commands (Deprecated but still supported)

// GEORADIUS key longitude latitude radius M|KM|FT|MI [WITHCOORD] [WITHDIST] [WITHHASH] [COUNT count [ANY]] [ASC|DESC] [STORE key] [STOREDIST key]
// GeoRadius returns members within a radius from coordinates (deprecated - use GeoSearch).
func (r *Redis) GeoRadius(key string, longitude, latitude, radius float64, unit string) ([]string, error) {
	args := packArgs("GEORADIUS", key, longitude, latitude, radius, unit)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	return rp.ListValue()
}

// GeoRadiusWithOptions returns members with additional result information (deprecated).
func (r *Redis) GeoRadiusWithOptions(key string, longitude, latitude, radius float64, unit string, opts GeoRadiusOptions) ([]GeoLocation, error) {
	args := []interface{}{"GEORADIUS", key, longitude, latitude, radius, unit}

	if opts.WithCoord {
		args = append(args, "WITHCOORD")
	}
	if opts.WithDist {
		args = append(args, "WITHDIST")
	}
	if opts.WithHash {
		args = append(args, "WITHHASH")
	}

	if opts.Count > 0 {
		args = append(args, "COUNT", opts.Count)
		if opts.Any {
			args = append(args, "ANY")
		}
	}

	if opts.Order != "" {
		args = append(args, opts.Order)
	}

	if opts.Store != "" {
		args = append(args, "STORE", opts.Store)
	}
	if opts.StoreDist != "" {
		args = append(args, "STOREDIST", opts.StoreDist)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	// If STORE or STOREDIST is used, return count as single GeoLocation
	if opts.Store != "" || opts.StoreDist != "" {
		count, _ := rp.IntegerValue()
		return []GeoLocation{{Member: strconv.FormatInt(count, 10)}}, nil
	}

	return parseGeoLocations(rp.Multi, opts.WithCoord, opts.WithDist, opts.WithHash)
}

// GEORADIUSBYMEMBER key member radius M|KM|FT|MI [WITHCOORD] [WITHDIST] [WITHHASH] [COUNT count [ANY]] [ASC|DESC] [STORE key] [STOREDIST key]
// GeoRadiusByMember returns members within a radius from another member (deprecated).
func (r *Redis) GeoRadiusByMember(key, member string, radius float64, unit string) ([]string, error) {
	args := packArgs("GEORADIUSBYMEMBER", key, member, radius, unit)
	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	return rp.ListValue()
}

// GeoRadiusByMemberWithOptions returns members with additional information (deprecated).
func (r *Redis) GeoRadiusByMemberWithOptions(key, member string, radius float64, unit string, opts GeoRadiusOptions) ([]GeoLocation, error) {
	args := []interface{}{"GEORADIUSBYMEMBER", key, member, radius, unit}

	if opts.WithCoord {
		args = append(args, "WITHCOORD")
	}
	if opts.WithDist {
		args = append(args, "WITHDIST")
	}
	if opts.WithHash {
		args = append(args, "WITHHASH")
	}

	if opts.Count > 0 {
		args = append(args, "COUNT", opts.Count)
		if opts.Any {
			args = append(args, "ANY")
		}
	}

	if opts.Order != "" {
		args = append(args, opts.Order)
	}

	if opts.Store != "" {
		args = append(args, "STORE", opts.Store)
	}
	if opts.StoreDist != "" {
		args = append(args, "STOREDIST", opts.StoreDist)
	}

	rp, err := r.ExecuteCommand(args...)
	if err != nil {
		return nil, err
	}

	// If STORE or STOREDIST is used, return count as single GeoLocation
	if opts.Store != "" || opts.StoreDist != "" {
		count, _ := rp.IntegerValue()
		return []GeoLocation{{Member: strconv.FormatInt(count, 10)}}, nil
	}

	return parseGeoLocations(rp.Multi, opts.WithCoord, opts.WithDist, opts.WithHash)
}

// Helper functions

func parseGeoLocations(replies []*Reply, withCoord, withDist, withHash bool) ([]GeoLocation, error) {
	if replies == nil {
		return nil, nil
	}

	locations := make([]GeoLocation, len(replies))
	for i, reply := range replies {
		location := GeoLocation{}

		if !withCoord && !withDist && !withHash {
			// Simple member name only
			location.Member, _ = reply.StringValue()
		} else {
			// Complex result with additional information
			if reply.Type == MultiReply && reply.Multi != nil {
				idx := 0

				// First element is always the member name
				if idx < len(reply.Multi) {
					location.Member, _ = reply.Multi[idx].StringValue()
					idx++
				}

				// Parse additional fields based on options
				if withDist && idx < len(reply.Multi) {
					distStr, _ := reply.Multi[idx].StringValue()
					if dist, err := strconv.ParseFloat(distStr, 64); err == nil {
						location.Distance = &dist
					}
					idx++
				}

				if withHash && idx < len(reply.Multi) {
					if hash, err := reply.Multi[idx].IntegerValue(); err == nil {
						location.Hash = &hash
					}
					idx++
				}

				if withCoord && idx < len(reply.Multi) {
					coordReply := reply.Multi[idx]
					if coordReply.Type == MultiReply && len(coordReply.Multi) >= 2 {
						lonStr, _ := coordReply.Multi[0].StringValue()
						latStr, _ := coordReply.Multi[1].StringValue()
						lon, _ := strconv.ParseFloat(lonStr, 64)
						lat, _ := strconv.ParseFloat(latStr, 64)
						location.Coordinates = &GeoCoordinate{
							Longitude: lon,
							Latitude:  lat,
						}
					}
				}
			}
		}

		locations[i] = location
	}

	return locations, nil
}
