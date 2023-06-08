package main

/*
Locations from http://www.cuug.ab.ca/dewara/mandelbrot/images.html.
NOTE: As of 2023, the webpage seems to be offline. The scrapper is thus deprecated.
Please use the locations.json in the repo.
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func getLocations() LocationsFile {
	log.Println("Reading location data...")
	file, err := os.ReadFile("locations.json")
	if err != nil {
		if os.IsNotExist(err) {
			scrapeLocationsToJSON()
			file, _ = os.ReadFile("locations.json")
		} else {
			panic(err)
		}
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
	return locs
}

func scrapeLocationsToJSON() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	log.Println("Getting response...")
	resp, err := client.Get("http://www.cuug.ab.ca/dewara/mandelbrot/images.html")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println("Parsing HTML...")
	locFile := parseHTML(resp.Body)

	log.Println("Writing location data to JSON...")
	res, err := json.MarshalIndent(locFile, "", " ")
	if err != nil {
		fmt.Println(err)
	} else {
		_ = os.WriteFile("locations.json", res, 0644)
	}
	log.Println("Scraping location data successfull.")
}

func parseHTML(body io.Reader) LocationsFile {
	htmlTokens := html.NewTokenizer(body)
	locFile := LocationsFile{}
	currLoc := Location{}
	reachedEOF := false

	for !reachedEOF {
		tt := htmlTokens.Next()
		switch tt {
		case html.ErrorToken:
			err := htmlTokens.Err()
			if err == io.EOF {
				reachedEOF = true
			}
		case html.TextToken:
			t := htmlTokens.Token()
			parseTextToken(t.Data, &currLoc)
		case html.EndTagToken:
			t := htmlTokens.Token()
			if t.Data == "p" && currLoc.Zoom > 0 {
				locFile.Locations = append(locFile.Locations, currLoc)
			}
		}
	}
	return locFile
}

func parseTextToken(text string, currLoc *Location) {
	splitted := strings.Split(text, " ")
	if len(splitted) != 3 {
		return
	}
	prop := strings.Replace(splitted[0], "\n", "", -1)
	switch prop {
	case "X":
		currLoc.XCenter, _ = strconv.ParseFloat(splitted[2], 64)
	case "Y":
		currLoc.YCenter, _ = strconv.ParseFloat(splitted[2], 64)
	case "R":
		rezZoom, _ := strconv.ParseFloat(splitted[2], 64)
		currLoc.Zoom = 1 / rezZoom
	}
}
