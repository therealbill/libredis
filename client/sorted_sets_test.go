package client

import (
	"testing"
)

func TestZAdd(t *testing.T) {
	r.Del("key")
	n, err := r.ZAdd("key", 1.0, "foo")
	if err != nil {
		t.Error(err)
	} else if n != 1 {
		t.Fail()
	}
}

func BenchmarkZAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r.ZAdd("key", 1.0, "foo")
	}
}

func TestZAddVariadic(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   1.0,
		"three": 3.0,
	}
	if n, err := r.ZAddVariadic("key", pairs); err != nil {
		t.Error(err)
	} else if n != 3 {
		t.Fail()
	}
	if n, _ := r.ZAddVariadic("key", map[string]float64{"two": 2.0}); n != 0 {
		t.Fail()
	}
}
func BenchmarkZAddVariadic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pairs := map[string]float64{
			"one":   1.0,
			"two":   1.0,
			"three": 3.0,
		}
		r.ZAddVariadic("key", pairs)
	}
}

func TestZCard(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   1.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if n, err := r.ZCard("key"); err != nil {
		t.Error(err)
	} else if n != 3 {
		t.Fail()
	}
}

func TestZCount(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if n, err := r.ZCount("key", "-inf", "+inf"); err != nil {
		t.Error(err)
	} else if n != 3 {
		t.Fail()
	}
	if n, _ := r.ZCount("key", "(1", "3"); n != 2 {
		t.Fail()
	}
}

func TestZIncrBy(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   1.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if n, err := r.ZIncrBy("key", 1.0, "two"); err != nil {
		t.Error(err)
	} else if n != 2.0 {
		t.Fail()
	}
}

func TestZRange(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if result, err := r.ZRange("key", 0, -1, false); err != nil {
		t.Error(err)
	} else if len(result) != 3 {
		t.Fail()
	} else if result[0] != "one" {
		t.Fail()
	}
	if result, err := r.ZRange("key", -2, -1, true); err != nil {
		t.Error(err)
	} else if len(result) != 4 {
		t.Fail()
	} else if result[0] != "two" {
		t.Fail()
	} else if result[1] != "2" {
		t.Fail()
	}
}

func TestZRank(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if n, err := r.ZRank("key", "three"); err != nil {
		t.Error(err)
	} else if n != 2 {
		t.Fail()
	}
	if n, err := r.ZRank("key", "four"); err != nil {
		t.Error(err)
	} else if n >= 0 {
		t.Fail()
	}
}

func TestZRem(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if n, err := r.ZRem("key", "three", "four"); err != nil {
		t.Error(err)
	} else if n != 1 {
		t.Fail()
	}
}

func TestZRemRangeByRank(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if n, err := r.ZRemRangeByRank("key", 0, 1); err != nil {
		t.Error(err)
	} else if n != 2 {
		t.Fail()
	}
}

func TestZRemRangeByScore(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if n, err := r.ZRemRangeByScore("key", "-inf", "(2"); err != nil {
		t.Error(err)
	} else if n != 1 {
		t.Fail()
	}
}

func TestZRevRange(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if result, err := r.ZRevRange("key", 0, -1, false); err != nil {
		t.Error(err)
	} else if len(result) != 3 {
		t.Fail()
	} else if result[0] != "three" {
		t.Fail()
	}
}

func TestZRevRank(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if n, err := r.ZRevRank("key", "three"); err != nil {
		t.Error(err)
	} else if n != 0 {
		t.Fail()
	}
	if n, err := r.ZRevRank("key", "four"); err != nil {
		t.Error(err)
	} else if n >= 0 {
		t.Fail()
	}
}

func TestZScore(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if result, err := r.ZScore("key", "member"); err != nil {
		t.Error(err)
	} else if result != nil {
		t.Fail()
	}
	if result, err := r.ZScore("key", "two"); err != nil {
		t.Error(err)
	} else if string(result) != "2" {
		t.Fail()
	}
}

func TestZScan(t *testing.T) {
	r.Del("key")
	pairs := map[string]float64{
		"one":   1.0,
		"two":   2.0,
		"three": 3.0,
	}
	r.ZAddVariadic("key", pairs)
	if _, list, err := r.ZScan("key", 0, "", 0); err != nil {
		t.Error(err)
	} else if len(list) == 0 {
		t.Fail()
	}
}

func TestZInterStore(t *testing.T) {
	r.Del("zset1", "zset2")
	r.ZAddVariadic("zset1", map[string]float64{
		"one": 1,
		"two": 2,
	})
	r.ZAddVariadic("zset2", map[string]float64{
		"one":   1,
		"two":   2,
		"three": 3,
	})
	if n, err := r.ZInterStore("out", []string{"zset1", "zset2"}, []int{2, 3}, ""); err != nil {
		t.Error(err)
	} else if n != 2 {
		t.Fail()
	}
}

func TestZUnionStore(t *testing.T) {
	r.Del("zset1", "zset2")
	r.ZAddVariadic("zset1", map[string]float64{
		"one": 1,
		"two": 2,
	})
	r.ZAddVariadic("zset2", map[string]float64{
		"one":   1,
		"two":   2,
		"three": 3,
	})
	if n, err := r.ZUnionStore("out", []string{"zset1", "zset2"}, []int{2, 3}, ""); err != nil {
		t.Error(err)
	} else if n != 3 {
		t.Fail()
	}
}

func TestZRangeByScore(t *testing.T) {
	r.Del("key")
	r.ZAddVariadic("key", map[string]float64{
		"one":   1,
		"two":   2,
		"three": 3,
	})
	if result, err := r.ZRangeByScore("key", "-inf", "+inf", false, false, 0, 0); err != nil {
		t.Error(err)
	} else if len(result) != 3 {
		t.Fail()
	}
}

func TestZRevRangeByScore(t *testing.T) {
	r.Del("key")
	r.ZAddVariadic("key", map[string]float64{
		"one":   1,
		"two":   2,
		"three": 3,
	})
	if result, err := r.ZRevRangeByScore("key", "(2", "(1", false, false, 0, 0); err != nil {
		t.Error(err)
	} else if len(result) != 0 {
		t.Fail()
	}
}

// Tests for new Phase 1 sorted set commands
func TestZPopMax(t *testing.T) {
	r.Del("zset")
	r.ZAddVariadic("zset", map[string]float64{"one": 1, "two": 2, "three": 3})
	
	member, err := r.ZPopMax("zset")
	if err != nil {
		t.Error(err)
	}
	if member.Member != "three" || member.Score != 3 {
		t.Error("Expected member 'three' with score 3")
	}
}

func TestZPopMaxCount(t *testing.T) {
	r.Del("zset")
	r.ZAddVariadic("zset", map[string]float64{"one": 1, "two": 2, "three": 3})
	
	members, err := r.ZPopMaxCount("zset", 2)
	if err != nil {
		t.Error(err)
	}
	if len(members) != 2 {
		t.Error("Expected 2 members, got", len(members))
	}
}

func TestZPopMin(t *testing.T) {
	r.Del("zset")
	r.ZAddVariadic("zset", map[string]float64{"one": 1, "two": 2, "three": 3})
	
	member, err := r.ZPopMin("zset")
	if err != nil {
		t.Error(err)
	}
	if member.Member != "one" || member.Score != 1 {
		t.Error("Expected member 'one' with score 1")
	}
}

func TestZPopMinCount(t *testing.T) {
	r.Del("zset")
	r.ZAddVariadic("zset", map[string]float64{"one": 1, "two": 2, "three": 3})
	
	members, err := r.ZPopMinCount("zset", 2)
	if err != nil {
		t.Error(err)
	}
	if len(members) != 2 {
		t.Error("Expected 2 members, got", len(members))
	}
}

func TestBZPopMax(t *testing.T) {
	r.Del("zset1", "zset2")
	r.ZAddVariadic("zset1", map[string]float64{"one": 1, "two": 2})
	
	result, err := r.BZPopMax([]string{"zset1", "zset2"}, 1)
	if err != nil {
		t.Error(err)
	}
	if result.Key != "zset1" {
		t.Error("Expected key 'zset1'")
	}
	if len(result.Members) == 0 {
		t.Error("Expected at least one member")
	}
}

func TestBZPopMin(t *testing.T) {
	r.Del("zset1", "zset2")
	r.ZAddVariadic("zset1", map[string]float64{"one": 1, "two": 2})
	
	result, err := r.BZPopMin([]string{"zset1", "zset2"}, 1)
	if err != nil {
		t.Error(err)
	}
	if result.Key != "zset1" {
		t.Error("Expected key 'zset1'")
	}
	if len(result.Members) == 0 {
		t.Error("Expected at least one member")
	}
}

func TestZRandMember(t *testing.T) {
	r.Del("zset")
	r.ZAddVariadic("zset", map[string]float64{"one": 1, "two": 2, "three": 3})
	
	member, err := r.ZRandMember("zset")
	if err != nil {
		t.Error(err)
	}
	if member == "" {
		t.Error("Expected non-empty member")
	}
}

func TestZMScore(t *testing.T) {
	r.Del("zset")
	r.ZAddVariadic("zset", map[string]float64{"one": 1, "two": 2, "three": 3})
	
	scores, err := r.ZMScore("zset", "one", "two", "nonexistent")
	if err != nil {
		t.Error(err)
	}
	if len(scores) != 3 {
		t.Error("Expected 3 scores, got", len(scores))
	}
}
