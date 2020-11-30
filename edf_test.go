package main

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEDF_empty(t *testing.T) {
	assert.Equal(t, []float64{}, EqualWidthEDF([]float64{}, 1))
}

func TestEDF_one(t *testing.T) {
	assert.Equal(t, []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, EqualWidthEDF([]float64{42}, 10))
}

func TestEDF_two(t *testing.T) {
	assert.Equal(t, []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 1}, EqualWidthEDF([]float64{42, 17}, 10))
}

func TestEDF_three(t *testing.T) {
	assert.Equal(t, []float64{0.25, 0.25, 0.25, 0.25, 0.25, 0.5, 0.5, 0.5, 0.5, 1}, EqualWidthEDF([]float64{42, 17, 40, 30}, 10))
}

func benchmarkEdfUniformM(b *testing.B, numSamples int, numBins int) {
	x := make([]float64, numSamples)
	scale := rand.Float64() * 1000
	for i := 0; i < numSamples; i++ {
		x[i] = rand.Float64() * scale
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EqualWidthEDF(x, numBins)
	}
}

func BenchmarkEDF_1000x50(b *testing.B) {
	benchmarkEdfUniformM(b, 1000, 50)
}

func BenchmarkEDF_10000x50(b *testing.B) {
	benchmarkEdfUniformM(b, 10000, 50)
}
