package main

import (
	"encoding/json"
	"log"
	"os"
)

func getLocations() LocationsFile {
	log.Println("Reading location data...")
	file, err := os.ReadFile("locations.json")
	if err != nil {
		panic(err)
	}

	locs := LocationsFile{}
	_ = json.Unmarshal(
		file,
		&locs,
	)

	zoom1Fractal := Location{
		XCenter: -0.75,
		YCenter: 0,
		Zoom:    1,
	}
	locs.Locations = append(locs.Locations, zoom1Fractal)

	log.Printf("Found %v locations.", len(locs.Locations))
	return locs
}
