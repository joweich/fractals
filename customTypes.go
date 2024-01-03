package main

type Location struct {
	XCenter float64
	YCenter float64
	Zoom    float64
}

type LocationsFile struct {
	Locations []Location
}

type ImageConfig struct {
	Width       int
	Height      int
	Samples     int
	MaxIter     int
	Offset      float64
	Mixing      bool
	InsideBlack bool
	Grayscale   bool
	RndGlobal   uint64
}

type MandelbrotIterResult struct {
	IsUnlimited bool
	Magnitude   float64
	Iterations  int
}
