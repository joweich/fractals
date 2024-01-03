package main

import (
	"fmt"
	"image"
	"testing"
	"time"
)

type BenchmarkConfig struct {
	MaxIterValues []int
	ImgSizeValues []int
	loc           Location
}

func getBenchmarkConfig() BenchmarkConfig {
	imgConf = ImageConfig{
		Width:       -1,
		Height:      -1,
		Samples:     5,
		MaxIter:     -1,
		Offset:      0.0,
		Mixing:      true,
		InsideBlack: true,
		Grayscale:   false,
		RndGlobal:   uint64(time.Now().UnixNano()),
	}

	return BenchmarkConfig{
		MaxIterValues: []int{10, 100, 1000, 10000},
		ImgSizeValues: []int{100, 1000},
		loc: Location{
			XCenter: -0.5,
			YCenter: 0,
			Zoom:    1,
		},
	}

}

func BenchmarkRenderImage(b *testing.B) {
	benchmarkConfig := getBenchmarkConfig()

	for _, size := range benchmarkConfig.ImgSizeValues {
		imgConf.Height = size
		imgConf.Width = size
		for _, maxIter := range benchmarkConfig.MaxIterValues {
			imgConf.MaxIter = maxIter
			testId := fmt.Sprintf("Size_%d_MaxIter_%d", size, maxIter)
			b.Run(testId, func(subB *testing.B) {
				img := image.NewRGBA(image.Rect(0, 0, imgConf.Width, imgConf.Height))

				subB.ResetTimer()
				for i := 0; i < subB.N; i++ {
					renderImage(img, benchmarkConfig.loc)
				}
			})
		}
	}
}

func BenchmarkRenderRow(b *testing.B) {
	benchmarkConfig := getBenchmarkConfig()

	for _, size := range benchmarkConfig.ImgSizeValues {
		imgConf.Height = size
		imgConf.Width = size
		for _, maxIter := range benchmarkConfig.MaxIterValues {
			imgConf.MaxIter = maxIter
			testId := fmt.Sprintf("Size_%d_MaxIter_%d", size, maxIter)
			b.Run(testId, func(subB *testing.B) {
				img := image.NewRGBA(image.Rect(0, 0, imgConf.Width, imgConf.Height))
				subB.ResetTimer()

				rndLocal := RandUint64(&imgConf.RndGlobal)
				for i := 0; i < subB.N; i++ {
					renderRow(benchmarkConfig.loc, imgConf.Height/2, img, &rndLocal)
				}
			})
		}
	}
}

func BenchmarkGetColorForPixel(b *testing.B) {
	benchmarkConfig := getBenchmarkConfig()
	imgConf.Height = 1000
	imgConf.Width = 1000

	for _, maxIter := range benchmarkConfig.MaxIterValues {
		imgConf.MaxIter = maxIter
		testId := fmt.Sprintf("MaxIter_%d", maxIter)
		b.Run(testId, func(subB *testing.B) {
			subB.ResetTimer()

			for i := 0; i < subB.N; i++ {
				getColorForPixel(benchmarkConfig.loc, imgConf.Height/2, imgConf.Width/2, &imgConf.RndGlobal)
			}
		})
	}
}

func BenchmarkGetColorFromMandelbrotUnlimited(b *testing.B) {
	magnitude := 2.0
	iterations := 100
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		iterResult := MandelbrotIterResult{
			IsUnlimited: true,
			Magnitude:   magnitude,
			Iterations:  iterations,
		}
		getColorFromMandelbrot(iterResult)
	}
}

func BenchmarkRunMandelbrot(b *testing.B) {
	benchmarkConfig := getBenchmarkConfig()

	c := complex(-0.5, 0) // (-0.5, 0) is part of the mandelbrot set, i.e. zn is bounded for all zn
	for _, maxIter := range benchmarkConfig.MaxIterValues {
		testId := fmt.Sprintf("MaxIter_%d", maxIter)
		b.Run(testId, func(subB *testing.B) {
			imgConf.MaxIter = maxIter
			subB.ResetTimer()

			for i := 0; i < subB.N; i++ {
				runMandelbrot(c)
			}
		})
	}
}
