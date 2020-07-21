package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	// datafile = "city-list-small.json"
	datafile = "city.list.json"
)

// Coord represents the city's coordinates.
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// City represents a city structure.
type City struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Country string `json:"country"`
	Coord   Coord  `json:"coord"`
}

func main() {
	// ////////////////////////////////////////////
	// Handle flags
	var (
		country string
		state   string
		out     string
	)

	flag.StringVar(&country, "cnt", "", "two-character country abbreviation")
	flag.StringVar(&state, "st", "", "two-letter US state abbreviation")
	flag.StringVar(&out, "o", "", "output file; if not specified, will write to stdout")
	flag.Parse()

	cities, err := parseList(datafile, country)
	if err != nil {
		log.Fatal(err)
	}

	var shortList []City
	if state == "" {
		shortList = parseListByCountry(cities, country)
	} else {
		shortList = parseListByCountry(cities, country, state)
	}
	// for _, city := range shortList {
	// 	fmt.Println(city)
	// }
	var b []byte
	switch out == "" {
	case true:
		fmt.Println(shortList)
	default:
		b, err = json.MarshalIndent(&shortList, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(string(b))
		file, err := os.OpenFile(out, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.Write(b)
		if err != nil {
			log.Fatal(err)
		}
	}

}

// parseList parses the entire file.
func parseList(src string, country string) (cities []City, err error) {
	// Open our source file.
	jsonFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Defer the closing of the file so that we can parse it later on.
	defer jsonFile.Close()

	// Unmarshal file into slice of cities.
	var list []City
	val, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(val, &list)
	if err != nil {
		return nil, err
	}

	// Remove for production.
	// for _, city := range list {
	// 	fmt.Println(city)
	// }

	return list, nil
}

// parseListByCountry takes a list of cities and returns a parsed list based
// on the country specified.
func parseListByCountry(cities []City, country string, state ...string) (newCities []City) {
	var nList []City

	for _, city := range cities {
		switch len(state) > 0 {
		case true:
			if city.Country == country && city.State == state[0] {
				nList = append(nList, city)
			}
		default:
			if city.Country == country {
				nList = append(nList, city)
			}
		}

	}

	return nList
}
