package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// logging
////////////////////////////////////////////////////////////////////////////////

var (
	trace   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
)

// InitLogging TODO TODO TODO
func InitLogging(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	trace = log.New(traceHandle, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	info = log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	warning = log.New(warningHandle, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	error = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// config
////////////////////////////////////////////////////////////////////////////////

// FLSConfig TODO TODO TODO
type FLSConfig struct {
	// TODO
}

// bandsintown
////////////////////////////////////////////////////////////////////////////////

// BandsInTownEventData TODO TODO TODO
type BandsInTownEventData struct {
	ID             int                    `json:"id"`
	ArtistID       int                    `json:"artist_id"`
	URL            string                 `json:"url"`
	OnSaleDatetime string                 `json:"on_sale_datetime"` // 2017-03-01T18:00:00
	Datetime       string                 `json:"datetime"`
	Description    string                 `json:"description"`
	Venue          *BandsInTownVenueData  `json:"venue"`
	Offers         []BandsInTownOfferData `json:"offers"`
	Lineup         []string               `json:"lineup"`
}

// BandsInTownVenueData TODO TODO TODO
type BandsInTownVenueData struct {
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	City      string `json:"city"`
	Region    string `json:"region"`
	Country   string `json:"country"`
}

// BandsInTownOfferData TODO TODO TODO
type BandsInTownOfferData struct {
	Type   string `json:"type"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

// BandsInTownArtistData TODO TODO TODO
type BandsInTownArtistData struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	ImageURL        string `json:"image_url"`
	ThumbURL        string `json:"thumb_url"`
	FacebookPageURL string `json:"facebook_page_url"`
	MBID            string `json:"mbid"`
	TrackerCount    int    `json:"tracker_count"`
}

// BandsInTownData represents the full, raw picture from BandsInTown of an artist and their events
type BandsInTownData struct {
	QueryDate int                    `json:"query_date"` // UNIX timestamp for time of last bandsintown API call
	Artist    BandsInTownArtistData  `json:"artist"`
	Events    []BandsInTownEventData `json:"events"`
}

// show data
////////////////////////////////////////////////////////////////////////////////

// FLSData represents all of the non-cache data in fls-core
type FLSData struct {
	Config          *FLSConfig                 `json:"config"`           // Stores fls-core configuration, not sure what to put here yet
	BandsInTownData map[string]BandsInTownData `json:"bandsintown_data"` // maps artist_name -> bandsintown artist info
	FollowedArtists []string                   `json:"followed_artists"` // List of followed artists
}

// REST API handles
////////////////////////////////////////////////////////////////////////////////

// RouteRoot displays the name of the service
func RouteRoot(w http.ResponseWriter, r *http.Request) {
	info.Printf("Hit root handler. %v %v\n", r.Method, r.URL)
	fmt.Fprintf(w, "fls-core")
}

// RouteGetShows returns a json object containing upcoming shows
func RouteGetShows(w http.ResponseWriter, r *http.Request) {
	info.Printf("Hit shows handler. %v %v\n", r.Method, r.URL)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, ReadTestJSON())
}

// RouteGetArtists returns a json list of followed artists
func RouteGetArtists(w http.ResponseWriter, r *http.Request) {
	info.Printf("Hit artists handler. %v %v\n", r.Method, r.URL)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `["Electric Light Orchestra","Gorillaz","King Gizzard and the Lizard Wizard","LCD Sound System","Pond","System of a Down","Tame Impala","The Beatles","Unknown Mortal Orchestra","Weezer"]`)
}

// misc
////////////////////////////////////////////////////////////////////////////////

// ReadTestJSON reads and returns test.json
func ReadTestJSON() string {

	dat, err := ioutil.ReadFile("test.json")
	if err != nil {
		panic(err)
	}
	// Info.Println(string(dat))

	return string(dat)
}

// main
////////////////////////////////////////////////////////////////////////////////

func main() {

	// initialize logging
	InitLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	// set routes
	http.HandleFunc("/", RouteRoot)
	http.HandleFunc("/v1/shows", RouteGetShows)
	http.HandleFunc("/v1/artists", RouteGetArtists)

	// start serving
	port := ":8001"
	info.Printf("fls-core serving on port %v\n", port)
	http.ListenAndServe(port, nil)

}
