package main

import (
	"image/color"
	"math"
	"math/cmplx"
)

func getColorForComplexNr(z0 complex128) color.RGBA {
	return getColorFromMandelbrot(runMandelbrot(z0))
}

func getColorFromMandelbrot(isUnlimited bool, magnitude float64, iterations int) color.RGBA {
	if isUnlimited {
		// adapted http://linas.org/art-gallery/escape/escape.html
		nu := math.Log(math.Log(magnitude)) / math.Log(2)
		hue := (float64(iterations)+1-nu)/float64(iterations) + imgConf.HueOffset
		return hslToRGB(hue, 1, 0.5)
	} else if imgConf.InsideBlack {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	} else {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
}

func runMandelbrot(z0 complex128) (bool, float64, int) {
	var z complex128

	for i := 0; i < imgConf.MaxIter; i++ {
		magnitude := cmplx.Abs(z)
		if magnitude > 2 {
			return true, magnitude, i
		}
		z = z*z + z0
	}
	return false, 0, 0
}
