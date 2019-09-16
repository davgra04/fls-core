package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	"gitlab.devgru.cc/devgru/fls-core/common"
	"gitlab.devgru.cc/devgru/fls-core/fls"
)

// REST API handles
////////////////////////////////////////////////////////////////////////////////

// RouteRoot displays the name of the service
func RouteRoot(w http.ResponseWriter, r *http.Request) {
	common.Info.Printf("%v %v\n", r.Method, r.URL)
	fmt.Fprintf(w, "fls-core")
}

// RouteGetShows returns a json object containing upcoming shows
func RouteGetShows(w http.ResponseWriter, r *http.Request) {

	defaultRegion := "TX"

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	common.Warning.Printf("r.URL.Query(): %v", r.URL.Query())

	regions, ok := r.URL.Query()["region"]

	if !ok || len(regions) < 1 {
		// log.Println("Missing ")
		// return
		common.Info.Printf("Missing region parameter, defaulting to %v", defaultRegion)

		regions = []string{defaultRegion}

	}
	common.Info.Printf("%v %v?region=%v\n", r.Method, r.URL, regions)

	// Query()["key"] will return an array of items,
	// we only want the single item.
	region := regions[0]

	common.Info.Printf("region: %v", region)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	showDataJSON := fls.GetCachedShowsResponse(region)

	// fmt.Fprintf(w, ReadTestJSON())
	fmt.Fprintf(w, showDataJSON)
}

// RouteGetArtists returns a json list of followed artists
func RouteGetArtists(w http.ResponseWriter, r *http.Request) {
	common.Info.Printf("%v %v\n", r.Method, r.URL)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	artistJSON, err := json.Marshal(common.Cfg.Artists)
	if err != nil {
		fmt.Fprintf(w, `{"error": "could not marshal json"}`)
	}

	fmt.Fprintf(w, string(artistJSON))
}

// main
////////////////////////////////////////////////////////////////////////////////

func main() {

	// initialize logging
	common.InitLogging()

	// handle flags
	configPath := flag.String("c", "", "path to the fls-core config file")
	localOnly := flag.Bool("l", false, "listen on localhost only")
	flagNoColor := flag.Bool("no-color", false, "disable color output")
	flag.Parse()

	if *flagNoColor {
		color.NoColor = true // disables colorized output
	}

	if *configPath == "" {
		common.Error.Printf("Must provide path to config!")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// check for api key environment variable
	if os.Getenv("BANDSINTOWN_API_KEY") == "" {
		common.Error.Printf("Must set environment variable BANDSINTOWN_API_KEY!")
		os.Exit(1)
	}

	// read config
	common.ReadConfig(configPath)

	// launch goroutines
	// go fls.PollBandsInTown()

	// set routes
	http.HandleFunc("/", RouteRoot)
	http.HandleFunc("/v1/shows", RouteGetShows)
	http.HandleFunc("/v1/artists", RouteGetArtists)

	// serve
	port := ":8001"
	if *localOnly {
		port = "localhost" + port
	}
	common.Info.Printf("fls-core listening on %v\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		common.Error.Panicf("Failed to ListenAndServe: %v", err)
	}

}
