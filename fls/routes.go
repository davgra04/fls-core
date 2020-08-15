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
		artists, err := getArtistList()
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

		go func() {
			artistList := strings.Split(artists, ",")
			if err := addArtists(artistList); err != nil {
				Error.Printf("Failed to add artists: %v\n", err)
			}
		}()

		return
	}

	http.Error(w, "method not allowed", 405)
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
