package main

import (
	"math"
)

// Links
// https://www.real-statistics.com/non-parametric-tests/goodness-of-fit-tests/two-sample-kolmogorov-smirnov-test/

// KolmogorovSmirnovStat calculates a 2-sample Kolmogorov Smirnov statistic between 2 cdf results. We don't have the
// true continuous CDF data here, and can only compare 2 empirical distributions.
func KolmogorovSmirnovStat(x, y []float64) float64 {
	lx, ly := len(x), len(y)

	supD := 0.0

	for i := 0; i < lx; i++ {
		j := int(float64(ly) * float64(i) / float64(lx))
		if j >= ly {
			j = ly-1
		}
		d := math.Abs(x[i] - y[j])
		if d > supD {
			supD = d
		}
	}

	for i := 0; i < ly; i++ {
		j := int(float64(lx) * float64(i) / float64(ly))
		if j >= lx {
			j = lx-1
		}
		d := math.Abs(y[i] - x[j])
		if d > supD {
			supD = d
		}
	}

	return supD
}

func KolmogorovSmirnovP(x, y []float64) float64 {
	d := KolmogorovSmirnovStat(x, y)
	t := d * math.Sqrt(float64(len(x) * len(y)) / float64(len(x) + len(y)))
	return 1.0 - ksum(t, 1e-20, 100000)
}

func ksum(t float64, tolerance float64, iterations int) float64 {
	x := t * t * -2.0
	sign := -1.0
	i := 1
	partial := 0.5
	delta := 1.0
	for delta > tolerance && i < iterations {
		delta = math.Exp(x * float64(i) * float64(i))
		partial += sign * delta
		sign *= -1
		i++
	}
	return partial * 2
}
