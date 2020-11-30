package main

import (
	"math"
	"sort"
)

type Bucket struct {
	left SubBucket
	right SubBucket
}

func (b *Bucket) e() float64 {
	avg := float64(b.left.count + b.right.count) / 2
	e := math.Pow(float64(b.left.count) - avg, 2) + math.Pow(float64(b.right.count) - avg, 2)
	return e
}

func (b *Bucket) split(leftHint *Bucket) (Bucket, Bucket) {
	right := Bucket{
		left: SubBucket{bound: (b.left.bound + b.right.bound) / 2, count:  b.right.count / 2},
		right: SubBucket{bound: b.right.bound, count: b.right.count - b.right.count / 2},
	}
	left := Bucket{
		left: SubBucket{bound: b.left.bound * 2 - right.left.bound, count:  b.left.count / 2},
		right: SubBucket{bound: b.left.bound, count: b.left.count - b.left.count / 2},
	}
	if leftHint != nil {
		left.left.bound = (leftHint.right.bound + b.left.bound) / 2
	}
	return left, right
}

func merge(bucket Bucket, bucket2 Bucket) Bucket {
	return Bucket{
		left: SubBucket{
			bound: bucket.right.bound,
			count: bucket.left.count + bucket.right.count,
		},
		right: SubBucket{
			bound: bucket2.right.bound,
			count: bucket2.left.count + bucket2.right.count,
		},
	}
}

type SubBucket struct {
	bound float64
	count int
}

type Dado struct {
	b int
	seed map[float64]int
	buckets []Bucket
	lastMi int
	lastSi int
}

func (d *Dado) FindBestToSplit() (int, float64) {
	bestI := 0
	bestE := math.Inf(-1)
	for i := 0; i < d.b; i++ {
		e := d.buckets[i].e()
		if e > bestE {
			bestI = i
			bestE = e
		}
	}
	return bestI, bestE
}

func (d *Dado) FindBestToMerge(ignore int) (int, float64) {
	bestI := 0
	bestE := math.Inf(+1)
	for i := 0; i < d.b - 1; i++ {
		if i == ignore || i + 1 == ignore {
			continue
		}
		vb := merge(d.buckets[i], d.buckets[i+1])
		e := vb.e()
		if e < bestE {
			bestI = i
			bestE = e
		}
	}
	return bestI, bestE
}

func (d *Dado) Add(v float64, c int) {
	if d.buckets == nil {
		d.seed[v] = d.seed[v] + c
		if len(d.seed) == d.b * 2 {
			// convert seed into bounds/counts
			subBuckets := make([]SubBucket, 0, d.b * 2)
			for k, v := range d.seed {
				subBuckets = append(subBuckets, SubBucket{bound: k, count: v})
			}
			sort.Slice(subBuckets, func(i, j int) bool {
				return subBuckets[i].bound < subBuckets[j].bound
			})
			// group into buckets
			d.buckets = make([]Bucket, 0, d.b)
			for i := 0; i < d.b; i++ {
				d.buckets = append(d.buckets, Bucket{
					left: subBuckets[i*2],
					right: subBuckets[i*2+1],
				})
			}
		}
	} else {
		bi := sort.Search(d.b, func(i int) bool {
			return d.buckets[i].right.bound >= v
		})
		if bi == d.b {
			// merge
			mi, _ := d.FindBestToMerge(-2)
			d.lastMi = mi
			d.lastSi = -1
			vb := merge(d.buckets[mi], d.buckets[mi+1])
			d.buckets[mi] = vb
			// a, b, mi, c, d, e
			// a, b, mi, d, e
			copy(d.buckets[mi+1:], d.buckets[mi+2:])

			// beyond end
			lr := d.buckets[d.b-2].right.bound
			d.buckets[d.b-1] = Bucket{
				left: SubBucket{bound: (v + lr) / 2, count: 0},
				right: SubBucket{bound: v, count: c},
			}
		} else {
			// left or right
			if v <= d.buckets[bi].left.bound {
				d.buckets[bi].left.count += c
			} else {
				d.buckets[bi].right.count += c
			}
			si, se := d.FindBestToSplit()
			mi, me := d.FindBestToMerge(si)
			//log.Printf("si=%d se=%g mi=%d me=%g", si, se, mi, me)
			if se > me {
				d.lastMi = mi
				d.lastSi = si
				vb := merge(d.buckets[mi], d.buckets[mi+1])
				d.buckets[mi] = vb
				if mi < si {
					// a, mi, b, c, si, e, f
					// a, mi, c, si, e, e, f
					copy(d.buckets[mi+1:si], d.buckets[mi+2:si+1])
					si--
				} else {
					// a, si, b, c, mi, e
					// a, si, b, b, c, mi
					copy(d.buckets[si+1:mi+2], d.buckets[si:mi+1])
				}

				var hint *Bucket
				if si > 0 {
					hint = &d.buckets[si-1]
				}
				d.buckets[si], d.buckets[si+1] = d.buckets[si].split(hint)
			}
		}
	}

}

func (d *Dado) EDF() *EDF {
	e := &EDF{
		bounds: make([]float64, d.b*2),
		counts: make([]int, d.b*2),
		total: 0,
	}
	for i, b := range d.buckets {
		e.bounds[i*2] = b.left.bound
		e.counts[i*2] = b.left.count
		e.bounds[i*2+1] = b.right.bound
		e.counts[i*2+1] = b.right.count
		e.total += b.left.count + b.right.count
	}
	return e
}
