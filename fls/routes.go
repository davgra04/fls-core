package fls

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// RouteRoot displays the name of the service
func RouteRoot(w http.ResponseWriter, r *http.Request) {
	// Info.Printf("%v %v\n", r.Method, r.URL)
	fmt.Fprintf(w, "{'name': 'fls-core'}")
}

// RouteArtists returns a json list of followed artists
func RouteArtists(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		Info.Println("GET artists")

		// get artist list
		artists, err := getArtists()
		if err != nil {
			Error.Printf("Failed to read artists: %v\n", err)
			http.Error(w, "failed to get artists", 500)
			return
		}

		// encode to json
		json, err := json.Marshal(artists)
		if err != nil {
			Error.Printf("Failed to marshal json: %v\n", err)
			http.Error(w, "failed to get artists", 500)
			return
		}

		Info.Printf("json: %s\n", json)

		fmt.Fprint(w, string(json))
		return
	} else if r.Method == "POST" {
		Info.Println("POST artists")

		artists := r.PostFormValue("artists")
		Info.Printf("    artists: %v\n", artists)

		// TODO: validate input

		artistList := strings.Split(artists, ",")
		if err := addArtists(artistList); err != nil {
			Error.Printf("Failed to add artists: %v\n", err)
		}

		return
	}

	http.Error(w, "method not allowed", 405)

	// Info.Printf("%v %v\n", r.Method, r.URL)

	// if r.Method != "GET" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.WriteHeader(http.StatusOK)

	// artistJSON, err := json.Marshal(Cfg.Artists)
	// if err != nil {
	// 	fmt.Fprintf(w, `{"error": "could not marshal json"}`)
	// }

	// fmt.Fprintf(w, string(artistJSON))
}

// // RouteShows returns a json object containing upcoming shows
// func RouteShows(w http.ResponseWriter, r *http.Request) {

// 	defaultRegion := "TX"

// 	if r.Method != "GET" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	Warning.Printf("r.URL.Query(): %v", r.URL.Query())

// 	regions, ok := r.URL.Query()["region"]

// 	if !ok || len(regions) < 1 {
// 		// log.Println("Missing ")
// 		// return
// 		Info.Printf("Missing region parameter, defaulting to %v", defaultRegion)

// 		regions = []string{defaultRegion}

// 	}
// 	Info.Printf("%v %v?region=%v\n", r.Method, r.URL, regions)

// 	// Query()["key"] will return an array of items,
// 	// we only want the single item.
// 	region := regions[0]

// 	Info.Printf("region: %v", region)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.WriteHeader(http.StatusOK)

// 	// showDataJSON := GetCachedShowsResponse(region)
// 	showDataJSON := "{}"

// 	// fmt.Fprintf(w, ReadTestJSON())
// 	fmt.Fprintf(w, showDataJSON)
// }
