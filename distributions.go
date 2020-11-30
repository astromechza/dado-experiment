package main

import (
	"math/rand"
)

func UniformDistribution(output []float64) []float64 {
	lo := len(output)
	for i := 0; i < lo; i++ {
		output[i] = rand.Float64()
	}
	return output
}

func NormalDistribution(output []float64, stdDev float64) []float64 {
	if stdDev < 0 || stdDev >= 1 {
		panic("bad std dev")
	}
	lo := len(output)
	for i := 0; i < lo; i++ {
		for {
			n := rand.NormFloat64() * stdDev + 0.5
			if n >= 0 && n <= 1 {
				output[i] = n
				break
			}
		}
	}
	return output
}


func ExponentialDistribution(output []float64, rate float64) []float64 {
	if rate < 0 || rate >= 1 {
		panic("bad rate")
	}
	lo := len(output)
	for i := 0; i < lo; i++ {
		for {
			n := rand.ExpFloat64() * rate
			if n >= 0 && n <= 1 {
				output[i] = n
				break
			}
		}
	}
	return output
}
