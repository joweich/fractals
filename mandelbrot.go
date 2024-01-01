package main

import (
	"image/color"
	"math"
	"math/cmplx"
)

func getColorForComplexNr(c complex128) color.RGBA {
	return getColorFromMandelbrot(runMandelbrot(c))
}

func getColorFromMandelbrot(isUnlimited bool, magnitude float64, iterations int) color.RGBA {
	if isUnlimited {
		// adapted http://linas.org/art-gallery/escape/escape.html
		smooth := (float64(iterations) + 1 - math.Log(math.Log(magnitude))/math.Log(2)) / float64(imgConf.MaxIter)
		offset := smooth + imgConf.Offset
		mod := math.Mod(offset, 1)
		if imgConf.Grayscale {
			return hslToRGB(0, 0, mod)
		}
		return hslToRGB(mod, 1, 0.5)
	} else if imgConf.InsideBlack {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	} else {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
}

func runMandelbrot(c complex128) (bool, float64, int) {
	var z complex128

	for i := 0; i < imgConf.MaxIter; i++ {
		magnitude := cmplx.Abs(z)
		if magnitude > 2 {
			return true, magnitude, i
		}
		z = z*z + c
	}
	return false, 0, 0
}
