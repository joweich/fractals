package main

type Location struct {
	XCenter float64
	YCenter float64
	Zoom    float64
}

type LocationsFile struct {
	Locations []Location
}
