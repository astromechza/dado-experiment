package main

import (
	"math"
	"sort"
	"gonum.org/v1/plot/plotter"
)

type EDF struct {
	bounds []float64
	counts []int
	total int
}

func (e *EDF) Add(v float64) {
	b := sort.SearchFloat64s(e.bounds, v)
	if b < len(e.bounds) {
		e.counts[b]++
		e.total++
	}
}

func (e *EDF) F(x float64) float64 {
	c := 0
	for i, b := range e.bounds {
		if b <= x {
			c += e.counts[i]
		} else {
			break
		}
	}
	return float64(c) / float64(e.total)
}

func (e *EDF) XYs() plotter.XYs {
	return e.FXYs(e.bounds)
}

func (e *EDF) FXYs(xx []float64) plotter.XYs {
	xys := make(plotter.XYs, 0)
	lastX := 0.0
	for i, x := range xx {
		y := e.F(x)
		if i > 0 {
			xys = append(xys, plotter.XY{X: lastX, Y: y})
		}
		xys = append(xys, plotter.XY{X: x, Y: y})
		lastX = x
	}
	return xys
}


func AllPointsEDF(rawSamples []float64) *EDF {
	sortedPoints := make([]float64, len(rawSamples))
	copy(sortedPoints, rawSamples)
	sort.Float64s(sortedPoints)
	o := &EDF{
		bounds: make([]float64, 0),
		counts: make([]int, 0),
		total: len(rawSamples),
	}
	v := math.NaN()
	count := 0
	for _, s := range sortedPoints {
		if s == v {
			count++
		} else {
			if count > 0 {
				o.bounds = append(o.bounds, v)
				o.counts = append(o.counts, count)
			}
			v = s
			count = 1
		}
	}
	if count > 0 {
		o.bounds = append(o.bounds, v)
		o.counts = append(o.counts, count)
	}
	return o
}

// EqualWidthEDF is empirical distribution function - a distribution function associated with empirical measures into
// an array of bins equally spaced between min and max. Complexity is assumed to be O(n log n + n).
// You must provide the raw samples, as well as the array of bins to cumulate into so that the min is places in the
// first bin and the max is placed in the last bin. If bins is nil or empty, a new slice will be assigned with the same
// length as the input samples.
func EqualWidthEDF(rawSamples []float64, bins int) *EDF {
	// first traversal to get min and max
	min, max := rawSamples[0], rawSamples[0]
	for _, s := range rawSamples {
		if s < min {
			min = s
		} else if s > max {
			max = s
		}
	}

	binSize := (max - min) / float64(bins)
	edf := &EDF{
		bounds: make([]float64, bins),
		counts: make([]int, bins),
	}
	for i := 0; i < bins; i++ {
		edf.bounds[i] = binSize + float64(i) * binSize
	}
	edf.bounds[bins-1] = max

	// second traversal to put samples into buckets
	for _, s := range rawSamples {
		edf.Add(s)
	}
	return edf
}
