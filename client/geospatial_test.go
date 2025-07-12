package client

import (
	"math"
	"testing"
)

// Test geospatial constants
func TestGeoConstants(t *testing.T) {
	if GeoUnitMeters != "M" {
		t.Error("GeoUnitMeters constant incorrect")
	}
	if GeoUnitKilometers != "KM" {
		t.Error("GeoUnitKilometers constant incorrect")
	}
	if GeoUnitFeet != "FT" {
		t.Error("GeoUnitFeet constant incorrect")
	}
	if GeoUnitMiles != "MI" {
		t.Error("GeoUnitMiles constant incorrect")
	}

	if GeoOrderAsc != "ASC" {
		t.Error("GeoOrderAsc constant incorrect")
	}
	if GeoOrderDesc != "DESC" {
		t.Error("GeoOrderDesc constant incorrect")
	}
}

// Test data structures
func TestGeoStructures(t *testing.T) {
	member := GeoMember{
		Longitude: -122.4194,
		Latitude:  37.7749,
		Member:    "San Francisco",
	}

	if member.Longitude != -122.4194 || member.Latitude != 37.7749 || member.Member != "San Francisco" {
		t.Error("GeoMember struct not working correctly")
	}

	coord := GeoCoordinate{
		Longitude: -74.0060,
		Latitude:  40.7128,
	}

	if coord.Longitude != -74.0060 || coord.Latitude != 40.7128 {
		t.Error("GeoCoordinate struct not working correctly")
	}
}

// Basic Geospatial Operations Tests

func TestGeoAdd(t *testing.T) {
	r.Del("cities")

	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
		{Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
		{Longitude: 2.3522, Latitude: 48.8566, Member: "Paris"},
	}

	added, err := r.GeoAdd("cities", members)
	if err != nil {
		t.Error(err)
	}
	if added != 3 {
		t.Errorf("Expected 3 added members, got %d", added)
	}

	// Adding same members again should return 0
	added, err = r.GeoAdd("cities", members)
	if err != nil {
		t.Error(err)
	}
	if added != 0 {
		t.Errorf("Expected 0 added members on duplicate, got %d", added)
	}
}

func TestGeoAddWithOptions(t *testing.T) {
	r.Del("cities")

	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
	}

	// Add with NX option
	opts := GeoAddOptions{NX: true}
	added, err := r.GeoAddWithOptions("cities", members, opts)
	if err != nil {
		t.Error(err)
	}
	if added != 1 {
		t.Errorf("Expected 1 added member with NX, got %d", added)
	}

	// Try to add same member again with NX - should fail
	added, err = r.GeoAddWithOptions("cities", members, opts)
	if err != nil {
		t.Error(err)
	}
	if added != 0 {
		t.Errorf("Expected 0 added members with NX on duplicate, got %d", added)
	}

	// Update with XX option
	updatedMembers := []GeoMember{
		{Longitude: -122.4195, Latitude: 37.7750, Member: "San Francisco"}, // Slightly different coords
	}
	optsXX := GeoAddOptions{XX: true, CH: true}
	changed, err := r.GeoAddWithOptions("cities", updatedMembers, optsXX)
	if err != nil {
		t.Error(err)
	}
	if changed != 1 {
		t.Errorf("Expected 1 changed member with XX+CH, got %d", changed)
	}
}

func TestGeoDist(t *testing.T) {
	r.Del("cities")

	// Add test cities
	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
		{Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
	}
	r.GeoAdd("cities", members)

	// Test distance in meters
	dist, err := r.GeoDist("cities", "San Francisco", "New York")
	if err != nil {
		t.Error(err)
	}
	// Distance should be approximately 4,135 km = 4,135,000 meters
	if dist < 4000000 || dist > 4300000 {
		t.Errorf("Expected distance ~4,135,000 meters, got %f", dist)
	}

	// Test distance in kilometers
	distKM, err := r.GeoDistWithUnit("cities", "San Francisco", "New York", GeoUnitKilometers)
	if err != nil {
		t.Error(err)
	}
	if distKM < 4000 || distKM > 4300 {
		t.Errorf("Expected distance ~4,135 km, got %f", distKM)
	}

	// Test non-existent member
	nonExistDist, err := r.GeoDist("cities", "San Francisco", "NonExistent")
	if err != nil {
		t.Error(err)
	}
	if nonExistDist != 0 {
		t.Errorf("Expected 0 distance for non-existent member, got %f", nonExistDist)
	}
}

func TestGeoHash(t *testing.T) {
	r.Del("cities")

	// Add test cities
	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
		{Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
	}
	r.GeoAdd("cities", members)

	// Get geohashes
	hashes, err := r.GeoHash("cities", "San Francisco", "New York", "NonExistent")
	if err != nil {
		t.Error(err)
	}
	if len(hashes) != 3 {
		t.Errorf("Expected 3 hash results, got %d", len(hashes))
	}

	// Valid geohashes should not be empty
	if hashes[0] == "" || hashes[1] == "" {
		t.Error("Expected non-empty geohashes for existing members")
	}

	// Non-existent member should have empty hash
	if hashes[2] != "" {
		t.Error("Expected empty geohash for non-existent member")
	}
}

func TestGeoPos(t *testing.T) {
	r.Del("cities")

	// Add test city
	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
	}
	r.GeoAdd("cities", members)

	// Get positions
	positions, err := r.GeoPos("cities", "San Francisco", "NonExistent")
	if err != nil {
		t.Error(err)
	}
	if len(positions) != 2 {
		t.Errorf("Expected 2 position results, got %d", len(positions))
	}

	// Valid position should not be nil
	if positions[0] == nil {
		t.Error("Expected non-nil position for existing member")
	} else {
		// Check coordinates are approximately correct
		if math.Abs(positions[0].Longitude-(-122.4194)) > 0.001 {
			t.Errorf("Longitude mismatch: expected ~-122.4194, got %f", positions[0].Longitude)
		}
		if math.Abs(positions[0].Latitude-37.7749) > 0.001 {
			t.Errorf("Latitude mismatch: expected ~37.7749, got %f", positions[0].Latitude)
		}
	}

	// Non-existent member should have nil position
	if positions[1] != nil {
		t.Error("Expected nil position for non-existent member")
	}
}

// Modern Search Commands Tests

func TestGeoSearch(t *testing.T) {
	r.Del("cities")

	// Add test cities
	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
		{Longitude: -122.2711, Latitude: 37.8044, Member: "Oakland"},
		{Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
	}
	r.GeoAdd("cities", members)

	// Search within radius from San Francisco
	sfCoord := &GeoCoordinate{Longitude: -122.4194, Latitude: 37.7749}
	radius := &GeoRadius{Radius: 50, Unit: GeoUnitKilometers}

	opts := GeoSearchOptions{
		FromLonLat: sfCoord,
		ByRadius:   radius,
		Order:      GeoOrderAsc,
	}

	locations, err := r.GeoSearch("cities", opts)
	if err != nil {
		t.Error(err)
	}

	// Should find San Francisco and Oakland, but not New York
	if len(locations) != 2 {
		t.Errorf("Expected 2 locations within 50km of SF, got %d", len(locations))
	}

	// Check that we got the expected cities
	found := make(map[string]bool)
	for _, loc := range locations {
		found[loc.Member] = true
	}
	if !found["San Francisco"] || !found["Oakland"] {
		t.Error("Expected to find San Francisco and Oakland")
	}
	if found["New York"] {
		t.Error("Did not expect to find New York within 50km of SF")
	}
}

func TestGeoSearchWithMember(t *testing.T) {
	r.Del("cities")

	// Add test cities
	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
		{Longitude: -122.2711, Latitude: 37.8044, Member: "Oakland"},
		{Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
	}
	r.GeoAdd("cities", members)

	// Search within radius from existing member
	memberName := "San Francisco"
	radius := &GeoRadius{Radius: 50, Unit: GeoUnitKilometers}

	opts := GeoSearchOptions{
		FromMember: &memberName,
		ByRadius:   radius,
		WithCoord:  true,
		WithDist:   true,
	}

	locations, err := r.GeoSearch("cities", opts)
	if err != nil {
		t.Error(err)
	}

	if len(locations) != 2 {
		t.Errorf("Expected 2 locations, got %d", len(locations))
	}

	// Verify that coordinates and distances are included
	for _, loc := range locations {
		if opts.WithCoord && loc.Coordinates == nil {
			t.Error("Expected coordinates to be included")
		}
		if opts.WithDist && loc.Distance == nil {
			t.Error("Expected distance to be included")
		}
	}
}

func TestGeoSearchStore(t *testing.T) {
	r.Del("cities", "nearby")

	// Add test cities
	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
		{Longitude: -122.2711, Latitude: 37.8044, Member: "Oakland"},
		{Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
	}
	r.GeoAdd("cities", members)

	// Search and store results
	sfCoord := &GeoCoordinate{Longitude: -122.4194, Latitude: 37.7749}
	radius := &GeoRadius{Radius: 50, Unit: GeoUnitKilometers}

	opts := GeoSearchStoreOptions{
		GeoSearchOptions: GeoSearchOptions{
			FromLonLat: sfCoord,
			ByRadius:   radius,
		},
	}

	stored, err := r.GeoSearchStore("nearby", "cities", opts)
	if err != nil {
		t.Error(err)
	}
	if stored != 2 {
		t.Errorf("Expected 2 stored locations, got %d", stored)
	}

	// Verify results were stored
	positions, err := r.GeoPos("nearby", "San Francisco", "Oakland")
	if err != nil {
		t.Error(err)
	}
	if positions[0] == nil || positions[1] == nil {
		t.Error("Expected stored locations to be retrievable")
	}
}

// Legacy Search Commands Tests

func TestGeoRadius(t *testing.T) {
	r.Del("cities")

	// Add test cities
	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
		{Longitude: -122.2711, Latitude: 37.8044, Member: "Oakland"},
		{Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
	}
	r.GeoAdd("cities", members)

	// Search within radius (legacy command)
	results, err := r.GeoRadius("cities", -122.4194, 37.7749, 50, GeoUnitKilometers)
	if err != nil {
		t.Error(err)
	}

	// Should find San Francisco and Oakland
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Check results
	found := make(map[string]bool)
	for _, result := range results {
		found[result] = true
	}
	if !found["San Francisco"] || !found["Oakland"] {
		t.Error("Expected to find San Francisco and Oakland")
	}
}

func TestGeoRadiusWithOptions(t *testing.T) {
	r.Del("cities")

	// Add test cities
	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
		{Longitude: -122.2711, Latitude: 37.8044, Member: "Oakland"},
	}
	r.GeoAdd("cities", members)

	// Search with additional options
	opts := GeoRadiusOptions{
		WithCoord: true,
		WithDist:  true,
		WithHash:  true,
		Order:     GeoOrderAsc,
		Count:     10,
	}

	locations, err := r.GeoRadiusWithOptions("cities", -122.4194, 37.7749, 50, GeoUnitKilometers, opts)
	if err != nil {
		t.Error(err)
	}

	if len(locations) != 2 {
		t.Errorf("Expected 2 locations, got %d", len(locations))
	}

	// Verify additional information is included
	for _, loc := range locations {
		if opts.WithCoord && loc.Coordinates == nil {
			t.Error("Expected coordinates to be included")
		}
		if opts.WithDist && loc.Distance == nil {
			t.Error("Expected distance to be included")
		}
		if opts.WithHash && loc.Hash == nil {
			t.Error("Expected hash to be included")
		}
	}
}

func TestGeoRadiusByMember(t *testing.T) {
	r.Del("cities")

	// Add test cities
	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
		{Longitude: -122.2711, Latitude: 37.8044, Member: "Oakland"},
		{Longitude: -74.0060, Latitude: 40.7128, Member: "New York"},
	}
	r.GeoAdd("cities", members)

	// Search around existing member
	results, err := r.GeoRadiusByMember("cities", "San Francisco", 50, GeoUnitKilometers)
	if err != nil {
		t.Error(err)
	}

	// Should find San Francisco and Oakland
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Check results
	found := make(map[string]bool)
	for _, result := range results {
		found[result] = true
	}
	if !found["San Francisco"] || !found["Oakland"] {
		t.Error("Expected to find San Francisco and Oakland")
	}
}

func TestGeoRadiusByMemberWithOptions(t *testing.T) {
	r.Del("cities")

	// Add test cities
	members := []GeoMember{
		{Longitude: -122.4194, Latitude: 37.7749, Member: "San Francisco"},
		{Longitude: -122.2711, Latitude: 37.8044, Member: "Oakland"},
	}
	r.GeoAdd("cities", members)

	// Search with options
	opts := GeoRadiusOptions{
		WithCoord: true,
		WithDist:  true,
		Count:     1,
		Order:     GeoOrderAsc,
	}

	locations, err := r.GeoRadiusByMemberWithOptions("cities", "San Francisco", 50, GeoUnitKilometers, opts)
	if err != nil {
		t.Error(err)
	}

	// With COUNT 1, should get only the closest result (San Francisco itself)
	if len(locations) != 1 {
		t.Errorf("Expected 1 location with COUNT 1, got %d", len(locations))
	}

	// Should be San Francisco (distance 0)
	if locations[0].Member != "San Francisco" {
		t.Errorf("Expected San Francisco as closest, got %s", locations[0].Member)
	}

	// Distance should be 0 (or very close to 0)
	if locations[0].Distance != nil && *locations[0].Distance > 1 {
		t.Errorf("Expected distance ~0, got %f", *locations[0].Distance)
	}
}

// Edge Cases and Error Handling

func TestGeoInvalidCoordinates(t *testing.T) {
	r.Del("invalid")

	// Test invalid longitude (should be -180 to 180)
	invalidMembers := []GeoMember{
		{Longitude: 200, Latitude: 45, Member: "invalid_lon"},
	}

	_, err := r.GeoAdd("invalid", invalidMembers)
	if err == nil {
		t.Error("Expected error for invalid longitude")
	}

	// Test invalid latitude (should be -85.05112878 to 85.05112878)
	invalidMembers2 := []GeoMember{
		{Longitude: 45, Latitude: 100, Member: "invalid_lat"},
	}

	_, err = r.GeoAdd("invalid", invalidMembers2)
	if err == nil {
		t.Error("Expected error for invalid latitude")
	}
}

func TestGeoEmptyKey(t *testing.T) {
	r.Del("empty")

	// Test operations on empty key
	dist, err := r.GeoDist("empty", "member1", "member2")
	if err != nil {
		t.Error(err)
	}
	if dist != 0 {
		t.Errorf("Expected 0 distance for empty key, got %f", dist)
	}

	hashes, err := r.GeoHash("empty", "member1")
	if err != nil {
		t.Error(err)
	}
	if len(hashes) != 1 || hashes[0] != "" {
		t.Error("Expected empty hash for non-existent member")
	}

	positions, err := r.GeoPos("empty", "member1")
	if err != nil {
		t.Error(err)
	}
	if len(positions) != 1 || positions[0] != nil {
		t.Error("Expected nil position for non-existent member")
	}
}

// Performance and Scale Tests

func TestGeoLargeDataset(t *testing.T) {
	r.Del("large_cities")

	// Add many cities
	members := make([]GeoMember, 0, 100)
	for i := 0; i < 100; i++ {
		// Generate coordinates around the world
		lon := float64(i)*3.6 - 180 // -180 to 180
		lat := float64(i)*1.7 - 85  // -85 to 85
		member := GeoMember{
			Longitude: lon,
			Latitude:  lat,
			Member:    "city_" + string(rune(i)),
		}
		members = append(members, member)
	}

	// Add all members
	added, err := r.GeoAdd("large_cities", members)
	if err != nil {
		t.Error(err)
	}
	if added != 100 {
		t.Errorf("Expected 100 added members, got %d", added)
	}

	// Search within a large radius from center
	opts := GeoSearchOptions{
		FromLonLat: &GeoCoordinate{Longitude: 0, Latitude: 0},
		ByRadius:   &GeoRadius{Radius: 20000, Unit: GeoUnitKilometers}, // Very large radius
		Count:      50,
	}

	locations, err := r.GeoSearch("large_cities", opts)
	if err != nil {
		t.Error(err)
	}

	// Should find up to 50 cities (COUNT limit)
	if len(locations) > 50 {
		t.Errorf("Expected at most 50 locations with COUNT 50, got %d", len(locations))
	}
}
