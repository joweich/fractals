package main

// locations from http://www.cuug.ab.ca/dewara/mandelbrot/images.html

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

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
		_ = ioutil.WriteFile("locations.json", res, 0644)
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
