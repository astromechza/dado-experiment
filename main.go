package main

import (
	"log"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)


func main() {
	rand.Seed(time.Now().UnixNano())

	data := ExponentialDistribution(make([]float64, 10000), 0.1)
	xx := make([]float64, 0)
	for x := 0.0; x < 1; x += 0.01 {
		xx = append(xx, x)
	}

	dado := &Dado{
		b: 10,
		seed: make(map[float64]int),
	}
	log.Printf("len(data): %d", len(data))
	for _, d := range data {
		dado.Add(d, 1)
	}
	log.Printf("dado%v\n\n", dado.buckets)

	dadoEDF := dado.EDF()
	log.Printf("eft%+v\n\n", dadoEDF)

	p, _ := plot.New()
	_ = plotutil.AddLines(p,
		"5", EqualWidthEDF(data, 5).FXYs(xx),
		"10", EqualWidthEDF(data, 10).FXYs(xx),
		"20", EqualWidthEDF(data, 20).FXYs(xx),
		"dado", dadoEDF.XYs(),
		"*", AllPointsEDF(data).FXYs(xx),
	)
	if err := p.Save(12*vg.Inch, 8*vg.Inch, "bins.png"); err != nil {
		panic(err)
	}

	data2 := ExponentialDistribution(make([]float64, 100000), 0.5)
	data3 := NormalDistribution(make([]float64, 100000), 0.25)
	data4 := UniformDistribution(make([]float64, 100000))

	p, _ = plot.New()
	_ = plotutil.AddLines(p,
		"edf n1", EqualWidthEDF(data, 100).XYs(),
		"edf e1", EqualWidthEDF(data2, 100).XYs(),
		"edf n2", EqualWidthEDF(data3, 100).XYs(),
		"edf u1", EqualWidthEDF(data4, 100).XYs(),
	)
	if err := p.Save(12*vg.Inch, 8*vg.Inch, "edf.png"); err != nil {
		panic(err)
	}
	//
	//xs := []int{5, 10, 15, 20, 30, 40, 50, 75, 100, 150, 200, 300, 400, 600, 800}
	//
	//self := make(plotter.XYs, 0)
	//edf := EqualWidthEDF(data, make([]float64, len(data)))
	//for _, i := range xs {
	//	partialEdf := EqualWidthEDF(data, make([]float64, i))
	//	self = append(self, plotter.XY{
	//		X: float64(i),
	//		Y: KolmogorovSmirnovStat(edf, partialEdf),
	//	})
	//}
	//
	//normVsExpo := make(plotter.XYs, 0)
	//for _, i := range xs {
	//	partialEdf := EqualWidthEDF(data2, make([]float64, i))
	//	normVsExpo = append(normVsExpo, plotter.XY{
	//		X: float64(i),
	//		Y: KolmogorovSmirnovStat(edf, partialEdf),
	//	})
	//}
	//
	//normVsNorm := make(plotter.XYs, 0)
	//for _, i := range xs {
	//	partialEdf := EqualWidthEDF(data3, make([]float64, i))
	//	normVsNorm = append(normVsNorm, plotter.XY{
	//		X: float64(i),
	//		Y: KolmogorovSmirnovStat(edf, partialEdf),
	//	})
	//}
	//
	//normVsUniform := make(plotter.XYs, 0)
	//for _, i := range xs {
	//	partialEdf := EqualWidthEDF(data4, make([]float64, i))
	//	normVsUniform = append(normVsUniform, plotter.XY{
	//		X: float64(i),
	//		Y: KolmogorovSmirnovStat(edf, partialEdf),
	//	})
	//}
	//
	//p, _ = plot.New()
	//p.Title.Text = "KS stat of low res histogram"
	//p.X.Label.Text = "# of buckets"
	//p.Y.Label.Text = "KS stat"
	//_ = plotutil.AddLinePoints(p,
	//	"self", self,
	//	"norm vs expo", normVsExpo,
	//	"norm vs norm", normVsNorm,
	//	"norm vs uniform", normVsUniform,
	//)
	//
	//if err := p.Save(12*vg.Inch, 8*vg.Inch, "points.png"); err != nil {
	//	panic(err)
	//}
}
