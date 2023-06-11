package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

var imgConf ImageConfig

func main() {
	parseImageConfigArgs()
	generateIagesFromLocations(getLocations())
}

func parseImageConfigArgs() {
	imgWidthPtr := flag.Int("width", 1920, "The width of the image in pixels.")
	imgHeightPtr := flag.Int("height", 1024, "The height of the image in pixels.")
	samplesPtr := flag.Int("samples", 50, "The number of samples")
	maxIterPtr := flag.Int("iter", 500, "The max. number of iterations.")
	hueOffsetPtr := flag.Float64("hue", 0.0, "The hsl hue offset in the range [0, 1)")
	mixingPtr := flag.Bool("mixing", true, "Use linear color mixing.")
	insideBlackPtr := flag.Bool("black", true, "Paint area inside in black.")

	flag.Parse()

	imgConf = ImageConfig{
		Width:       *imgWidthPtr,
		Height:      *imgHeightPtr,
		Samples:     *samplesPtr,
		MaxIter:     *maxIterPtr,
		HueOffset:   *hueOffsetPtr,
		Mixing:      *mixingPtr,
		InsideBlack: *insideBlackPtr,
		RndGlobal:	 uint64(time.Now().UnixNano()),
	}
}

func generateIagesFromLocations(locs LocationsFile) {
	if _, err := os.Stat("results/" + strconv.Itoa(imgConf.MaxIter)); os.IsNotExist(err) {
		os.Mkdir("results/"+strconv.Itoa(imgConf.MaxIter), 0755)
	}

	for index, loc := range locs.Locations {
		log.Println("Allocating and rendering image", index+1)
		img := image.NewRGBA(image.Rect(0, 0, imgConf.Width, imgConf.Height))
		renderImage(img, loc)

		log.Println("Encoding image", index+1)
		filename := "results/" + strconv.Itoa(imgConf.MaxIter) + "/" + strconv.Itoa(index+1)
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

func renderImage(img *image.RGBA, loc Location) {
	jobs := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		rndLocal := RandUint64(&imgConf.RndGlobal)
		go func() {
			for y := range jobs {
				renderRow(loc, y, img, &rndLocal)
			}
		}()
	}

	for y := 0; y < imgConf.Height; y++ {
		jobs <- y
		fmt.Printf("\r%d/%d (%d%%)", y, imgConf.Height, int(100*(float64(y)/float64(imgConf.Height))))
	}
	fmt.Printf("\r%d/%[1]d (100%%)\n", imgConf.Height)
}

func renderRow(loc Location, y int, img *image.RGBA, rndLocal *uint64) {
	for x := 0; x < imgConf.Width; x++ {
		cr, cg, cb := getColorForPixel(loc, x, y, rndLocal)
		img.SetRGBA(x, y, color.RGBA{R: cr, G: cg, B: cb, A: 255})
	}
}

func getColorForPixel(loc Location, x int, y int, rndLocal *uint64) (uint8, uint8, uint8) {
	var r, g, b int
	for i := 0; i < imgConf.Samples; i++ {
		c := getColorForComplexNr(convertPixelToComplexNr(loc, x, y, rndLocal))

		if imgConf.Mixing {
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
	if imgConf.Mixing {
		cr = LinearToRGB(uint16(float64(r) / float64(imgConf.Samples)))
		cg = LinearToRGB(uint16(float64(g) / float64(imgConf.Samples)))
		cb = LinearToRGB(uint16(float64(b) / float64(imgConf.Samples)))
	} else {
		cr = uint8(float64(r) / float64(imgConf.Samples))
		cg = uint8(float64(g) / float64(imgConf.Samples))
		cb = uint8(float64(b) / float64(imgConf.Samples))
	}
	return cr, cg, cb
}

func convertPixelToComplexNr(loc Location, x int, y int, rndLocal *uint64) complex128 {
	ratio := float64(imgConf.Width) / float64(imgConf.Height)

	// RandFload64() is added for anti-aliasing
	nx := (1/loc.Zoom)*ratio*((float64(x)+RandFloat64(rndLocal))/float64(imgConf.Width)-0.5) + loc.XCenter
	ny := (1/loc.Zoom)*((float64(y)+RandFloat64(rndLocal))/float64(imgConf.Height)-0.5) - loc.YCenter
	return complex(nx, ny)
}
