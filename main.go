package main

import (
	"encoding/json"
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

// Configuration
const (
	// Quality
	imgWidth       = 1024
	imgHeight      = 1024
	maxIter        = 1500
	samples        = 50
	hueOffset      = 0.0 // hsl color model; float in range [0,1)
	linearMixing   = true
	insideSetBlack = true

	scrapeLocations = false
	showProgress    = true
)

const (
	ratio = float64(imgWidth) / float64(imgHeight)
)

func main() {
	if scrapeLocations {
		scrapeLocationsToJSON()
	}

	log.Println("Reading location data...")
	file, err := os.ReadFile("locations.json")
	if err != nil {
		panic(err)
	}

	locs := LocationsFile{}
	_ = json.Unmarshal([]byte(file), &locs)

	zoom1Fractal := Location{
		XCenter: -0.75,
		YCenter: 0,
		Zoom:    1,
	}
	locs.Locations = append(locs.Locations, zoom1Fractal)

	log.Printf("Found %v locations.", len(locs.Locations))

	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("results", 0755)
	}

	for index, loc := range locs.Locations {

		log.Println("Allocating and rendering image ", index+1)
		img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
		render(img, loc)

		log.Println("Encoding image ", index+1)
		filename := "results/zoom" + strconv.FormatFloat(loc.Zoom, 'e', -1, 64) + "-iter" + strconv.Itoa(maxIter) + "-index" + strconv.Itoa(index+1)
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

func render(img *image.RGBA, loc Location) {
	jobs := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for y := range jobs {
				for x := 0; x < imgWidth; x++ {
					var r, g, b int
					for i := 0; i < samples; i++ {
						nx := 3*(1/loc.Zoom)*ratio*((float64(x)+RandFloat64())/float64(imgWidth)-0.5) + loc.XCenter
						ny := 3*(1/loc.Zoom)*((float64(y)+RandFloat64())/float64(imgHeight)-0.5) - loc.YCenter
						c := paint(mandelbrotIterComplex(nx, ny, maxIter))
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
		if showProgress {
			fmt.Printf("\r%d/%d (%d%%)", y, imgHeight, int(100*(float64(y)/float64(imgHeight))))
		}
	}
	if showProgress {
		fmt.Printf("\r%d/%[1]d (100%%)\n", imgHeight)
	}
}

func paint(magnitude float64, n int) color.RGBA {
	if magnitude > 2 {
		// adapted http://linas.org/art-gallery/escape/escape.html
		nu := math.Log(math.Log(magnitude)) / math.Log(2)
		hue := (float64(n)+1-nu)/float64(maxIter) + hueOffset
		return hslToRGB(hue, 1, 0.5)
	} else if insideSetBlack {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	} else {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
}

func mandelbrotIterComplex(px, py float64, maxIter int) (float64, int) {
	var current complex128
	pxpy := complex(px, py)

	for i := 0; i < maxIter; i++ {
		magnitude := cmplx.Abs(current)
		if magnitude > 2 {
			return magnitude, i
		}
		current = current*current + pxpy
	}

	magnitude := cmplx.Abs(current)
	return magnitude, maxIter
}
