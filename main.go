package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"os"
	"runtime"
	"strconv"
)

func main() {
	imgWidthPtr := flag.Int("width", 1920, "The width of the image in pixels.")
	imgHeightPtr := flag.Int("height", 1024, "The height of the image in pixels.")
	samplesPtr := flag.Int("samples", 50, "The number of samples")
	maxIterPtr := flag.Int("iter", 500, "The max. number of iterations.")
	hueOffsetPtr := flag.Float64("hue", 0.0, "The hsl hue offset in the range [0, 1)")
	mixingPtr := flag.Bool("mixing", true, "Use linear color mixing.")
	insideBlackPtr := flag.Bool("black", true, "Paint area inside in black.")

	flag.Parse()

	locs := getLocations()

	if _, err := os.Stat("results/" + strconv.Itoa(*maxIterPtr)); os.IsNotExist(err) {
		os.Mkdir("results/"+strconv.Itoa(*maxIterPtr), 0755)
	}

	for index, loc := range locs.Locations {
		log.Println("Allocating and rendering image ", index+1)
		img := image.NewRGBA(image.Rect(0, 0, *imgWidthPtr, *imgHeightPtr))
		render(img, loc, *samplesPtr, *maxIterPtr, *hueOffsetPtr, *mixingPtr, *insideBlackPtr)

		log.Println("Encoding image ", index+1)
		filename := "results/" + strconv.Itoa(*maxIterPtr) + "/" + strconv.Itoa(index+1) + strconv.FormatFloat(*hueOffsetPtr, 'b', -1, 64)
		f, err := os.Create(filename + ".png")
		if err != nil {
			panic(err)
		}
		err = png.Encode(f, img)
		if err != nil {
			panic(err)
		}
	}
}

func render(img *image.RGBA, loc Location, samples int, maxIter int, hueOffset float64, linearMixing bool, insideBlack bool) {
	imgWidth := img.Rect.Max.X
	imgHeight := img.Rect.Max.Y
	ratio := float64(imgWidth) / float64(imgHeight)

	jobs := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for y := range jobs {
				for x := 0; x < imgWidth; x++ {
					var r, g, b int
					for i := 0; i < samples; i++ {
						nx := 3*(1/loc.Zoom)*ratio*((float64(x)+RandFloat64())/float64(imgWidth)-0.5) + loc.XCenter
						ny := 3*(1/loc.Zoom)*((float64(y)+RandFloat64())/float64(imgHeight)-0.5) - loc.YCenter
						c := paint(mandelbrotIterComplex(nx, ny, maxIter, hueOffset, insideBlack))
						if linearMixing {
							r += int(RGBToLinear(c.R))
							g += int(RGBToLinear(c.G))
							b += int(RGBToLinear(c.B))
						} else {
							r += int(c.R)
							g += int(c.G)
							b += int(c.B)
						}
					}
					var cr, cg, cb uint8
					if linearMixing {
						cr = LinearToRGB(uint16(float64(r) / float64(samples)))
						cg = LinearToRGB(uint16(float64(g) / float64(samples)))
						cb = LinearToRGB(uint16(float64(b) / float64(samples)))
					} else {
						cr = uint8(float64(r) / float64(samples))
						cg = uint8(float64(g) / float64(samples))
						cb = uint8(float64(b) / float64(samples))
					}
					img.SetRGBA(x, y, color.RGBA{R: cr, G: cg, B: cb, A: 255})
				}
			}
		}()
	}

	for y := 0; y < imgHeight; y++ {
		jobs <- y
		fmt.Printf("\r%d/%d (%d%%)", y, imgHeight, int(100*(float64(y)/float64(imgHeight))))
	}
	fmt.Printf("\r%d/%[1]d (100%%)\n", imgHeight)
}

func paint(magnitude float64, maxIter int, hueOffset float64, insideBlack bool) color.RGBA {
	if magnitude > 2 {
		// adapted http://linas.org/art-gallery/escape/escape.html
		nu := math.Log(math.Log(magnitude)) / math.Log(2)
		hue := (float64(maxIter)+1-nu)/float64(maxIter) + hueOffset
		return hslToRGB(hue, 1, 0.5)
	} else if insideBlack {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	} else {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
}

func mandelbrotIterComplex(px, py float64, maxIter int, hueOffset float64, insideBlack bool) (float64, int, float64, bool) {
	var current complex128
	pxpy := complex(px, py)

	for i := 0; i < maxIter; i++ {
		magnitude := cmplx.Abs(current)
		if magnitude > 2 {
			return magnitude, i, float64(hueOffset), insideBlack
		}
		current = current*current + pxpy
	}

	magnitude := cmplx.Abs(current)
	return magnitude, maxIter, hueOffset, insideBlack
}
