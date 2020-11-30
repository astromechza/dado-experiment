package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDado(t *testing.T) {
	d := Dado{b: 4, seed: make(map[float64]int)}
	d.Add(1, 1)
	d.Add(2, 2)
	d.Add(3, 8)
	d.Add(4, 3)
	d.Add(5, 3)
	d.Add(6, 4)
	d.Add(7, 3)
	assert.Equal(t, map[float64]int{1:1, 2:2, 3:8, 4:3, 5:3, 6:4, 7:3}, d.seed)
	assert.Equal(t, ([]Bucket)(nil), d.buckets)
	d.Add(8, 3)
	assert.Equal(t, map[float64]int{1:1, 2:2, 3:8, 4:3, 5:3, 6:4, 7:3, 8:3}, d.seed)
	assert.Equal(t, []Bucket{
		Bucket{left:SubBucket{bound:1, count:1}, right:SubBucket{bound:2, count:2}},
		Bucket{left:SubBucket{bound:3, count:8}, right:SubBucket{bound:4, count:3}},
		Bucket{left:SubBucket{bound:5, count:3}, right:SubBucket{bound:6, count:4}},
		Bucket{left:SubBucket{bound:7, count:3}, right:SubBucket{bound:8, count:3}},
	}, d.buckets)
	d.Add(9, 1)
	assert.Equal(t, []Bucket{
		Bucket{left:SubBucket{bound:1, count:1}, right:SubBucket{bound:2, count:2}},
		Bucket{left:SubBucket{bound:3, count:8}, right:SubBucket{bound:4, count:3}},
		Bucket{left:SubBucket{bound:6, count:7}, right:SubBucket{bound:8, count:6}},
		Bucket{left:SubBucket{bound:8.5, count:0}, right:SubBucket{bound:9, count:1}},
	}, d.buckets)
	d.Add(0, 1)
	assert.Equal(t, []Bucket{
		Bucket{left:SubBucket{bound:1, count:2}, right:SubBucket{bound:2, count:2}},
		Bucket{left:SubBucket{bound:3, count:8}, right:SubBucket{bound:4, count:3}},
		Bucket{left:SubBucket{bound:6, count:7}, right:SubBucket{bound:8, count:6}},
		Bucket{left:SubBucket{bound:8.5, count:0}, right:SubBucket{bound:9, count:1}},
	}, d.buckets)
	d.Add(0, 1)
	assert.Equal(t, []Bucket{
		Bucket{left:SubBucket{bound:1, count:3}, right:SubBucket{bound:2, count:2}},
		Bucket{left:SubBucket{bound:3, count:8}, right:SubBucket{bound:4, count:3}},
		Bucket{left:SubBucket{bound:6, count:7}, right:SubBucket{bound:8, count:6}},
		Bucket{left:SubBucket{bound:8.5, count:0}, right:SubBucket{bound:9, count:1}},
	}, d.buckets)
}

func TestDado_Add_beyond(t *testing.T) {
	d := Dado{
		b: 2,
		seed: map[float64]int{},
		buckets: []Bucket{
			Bucket{left:SubBucket{bound:1, count:1}, right:SubBucket{bound:2, count:1}},
			Bucket{left:SubBucket{bound:3, count:1}, right:SubBucket{bound:4, count:1}},
		},
	}
	d.Add(5, 1)
	assert.Equal(t, []Bucket{
		Bucket{left:SubBucket{bound:2, count:2}, right:SubBucket{bound:4, count:2}},
		Bucket{left:SubBucket{bound:4.5, count:0}, right:SubBucket{bound:5, count:1}},
	}, d.buckets)
}

func TestDado_Add_split_before_merge(t *testing.T) {
	d := Dado{
		b: 3,
		seed: map[float64]int{},
		buckets: []Bucket{
			Bucket{left:SubBucket{bound:1, count:1}, right:SubBucket{bound:2, count:1}},
			Bucket{left:SubBucket{bound:3, count:1}, right:SubBucket{bound:4, count:1}},
			Bucket{left:SubBucket{bound:5, count:1}, right:SubBucket{bound:6, count:1}},
		},
	}
	d.Add(1, 10)
	assert.Equal(t, []Bucket{
		Bucket{left:SubBucket{bound:0.5, count:5}, right:SubBucket{bound:1, count:6}},
		Bucket{left:SubBucket{bound:1.5, count:0}, right:SubBucket{bound:2, count:1}},
		Bucket{left:SubBucket{bound:4, count:2}, right:SubBucket{bound:6, count:2}},
	}, d.buckets)
	assert.Equal(t, 0, d.lastSi)
	assert.Equal(t, 1, d.lastMi)
}

func BenchmarkDado_Add(b *testing.B) {
	d := Dado{b: 10, seed: make(map[float64]int)}
	for i := 0; i < b.N; i++ {
		d.Add(rand.Float64(), 1)
	}
	b.Logf("--")
	for i, bucket := range d.buckets {
		b.Logf("%d e = %f", i, bucket.e())
		b.Logf("%d.%d = %f (x%d)", i, 0, bucket.left.bound, bucket.left.count)
		b.Logf("%d.%d = %f (x%d)", i, 1, bucket.right.bound, bucket.right.count)
	}
	b.Logf("--")
}

func TestDado_Add_fuzz(t *testing.T) {
	d := Dado{b: 5, seed: make(map[float64]int)}
	for i := 0; i < 10; i++ {
		d.Add(rand.ExpFloat64(), 1)
	}
	for i := 0; i < 10000; i++ {
		before := fmt.Sprintf("%v", d.buckets)
		v := rand.ExpFloat64()
		d.Add(v, 1)
		after := fmt.Sprintf("%v", d.buckets)
		lastRightBound := math.Inf(-1)
		count := 0
		for i, b := range d.buckets {
			if lastRightBound >= b.left.bound || b.left.bound >= b.right.bound {
				t.Errorf("adding %v\nbefore=%s\nafter=%s\nfailure @ %d\nlast mi=%d si=%d", v, before, after, i, d.lastMi, d.lastSi)
				t.FailNow()
			}
			lastRightBound = b.right.bound
			count += b.left.count + b.right.count
		}
		require.Equal(t, i+10 + 1, count)
	}
	t.Logf("%v", d.buckets)
}
