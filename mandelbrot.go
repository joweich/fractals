package main

import (
	"image/color"
	"math"
)

func getColorForComplexNr(c complex128) color.RGBA {
	return getColorFromMandelbrot(runMandelbrot(c))
}

func getColorFromMandelbrot(iterResult MandelbrotIterResult) color.RGBA {
	if iterResult.IsUnlimited {
		// adapted http://linas.org/art-gallery/escape/escape.html
		smooth := (float64(iterResult.Iterations) + 1 - math.Log(math.Log(iterResult.Magnitude))/math.Log(2)) / float64(imgConf.MaxIter)
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

func runMandelbrot(c complex128) MandelbrotIterResult {
	var z complex128

	for i := 1; i < imgConf.MaxIter; i++ {
		z = z*z + c
		magnitudeSquared := real(z)*real(z) + imag(z)*imag(z)
		if magnitudeSquared > 4 {
			return MandelbrotIterResult{
				IsUnlimited: true,
				Magnitude:   math.Sqrt(magnitudeSquared),
				Iterations:  i,
			}
		}
	}
	return MandelbrotIterResult{
		IsUnlimited: false,
	}
}
