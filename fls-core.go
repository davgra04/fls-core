package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// logging
////////////////////////////////////////////////////////////////////////////////

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// InitLogging TODO TODO TODO
func InitLogging(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// bandsintown objects
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

// fls objects
////////////////////////////////////////////////////////////////////////////////

// FLSConfig represents the input configuration for fls-core
type FLSConfig struct {
	Artists               []string `json:"artists"`                             // List of followed artists
	RefreshPeriodSeconds  int      `json:"bandsintown_query_period_seconds"`    // number of seconds between FLSData refreshes
	RateLimitMillis       int      `json:"bandsintown_rate_limit_ms"`           // limit on the number of BandsInTown API requests per second
	MaxConcurrentRequests int      `json:"bandsintown_max_concurrent_requests"` // maximum number of concurrent requests to the BandsInTown API
}

// FLSData represents all of the non-cache data in fls-core
type FLSData struct {
	Config          *FLSConfig                 `json:"config"`           // Stores fls-core configuration, not sure what to put here yet
	BandsInTownData map[string]BandsInTownData `json:"bandsintown_data"` // maps artist_name -> bandsintown artist info
}

// GetShowsResponse represents the data returned by the RouteGetShows endpoint
type GetShowsResponse struct {
}

// GetArtistsResponse represents the data returned by the RouteGetArtists endpoint
type GetArtistsResponse struct {
}

// REST API handles
////////////////////////////////////////////////////////////////////////////////

// RouteRoot displays the name of the service
func RouteRoot(w http.ResponseWriter, r *http.Request) {
	Info.Printf("Hit root handler. %v %v\n", r.Method, r.URL)
	fmt.Fprintf(w, "fls-core")
}

// RouteGetShows returns a json object containing upcoming shows
func RouteGetShows(w http.ResponseWriter, r *http.Request) {
	Info.Printf("Hit shows handler. %v %v\n", r.Method, r.URL)

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
	Info.Printf("Hit artists handler. %v %v\n", r.Method, r.URL)

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `["Electric Light Orchestra","Gorillaz","King Gizzard and the Lizard Wizard","LCD Sound System","Pond","System of a Down","Tame Impala","The Beatles","Unknown Mortal Orchestra","Weezer"]`)
}

// bandsintown query goroutine
////////////////////////////////////////////////////////////////////////////////

// PollBandsInTown periodically polls BandsInTown for show data, saves to FLSData, and initiates cache rebuild goroutine
func PollBandsInTown(flscfg *FLSConfig) {
	fmt.Println("called QueryBandsInTown")
	apiKey := os.Getenv("BANDSINTOWN_API_KEY")

	// channel for limiting requests
	limiterDuration := time.Duration(int64(flscfg.RateLimitMillis) * time.Millisecond.Nanoseconds())
	limiter := time.Tick(limiterDuration)
	Info.Printf("limiter duration: %v [%T]", limiterDuration, limiterDuration)
	// Info.Printf("time.Millisecond.Nanoseconds(): %v [%T]", time.Millisecond.Nanoseconds(), time.Millisecond.Nanoseconds())
	// Info.Printf("1000000: 1000000")
	// Info.Printf("limiter: %v [%T]", limiter, limiter)

	bandsInTownClient := http.Client{
		Timeout: time.Second * 10,
	}

	// TODO: load from file
	flsdata := FLSData{}

	for {

		// set times for this polling period
		startTime := time.Now()
		nextPollTime := startTime.Add(time.Duration(flscfg.RefreshPeriodSeconds) * time.Second)

		// get data from api
		//////////////////////

		// load requests into requests channel
		requests := make(chan string, 1000)
		for _, artist := range flscfg.Artists {
			url := fmt.Sprintf("https://rest.bandsintown.com/artists/%s/events?app_id=%s", url.PathEscape(artist), apiKey)
			requests <- url
			Info.Printf("    preparing request for %-40v [%v]", artist, url)
		}
		close(requests)

		// make requests using limiter
		for url := range requests {
			<-limiter // wait for limiter

			// build and do request
			Info.Printf("    requesting %v", url)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				Error.Printf("Could not create request object: %v", err)
				continue
			}

			res, err := bandsInTownClient.Do(req)
			if err != nil {
				Error.Printf("Could not make request: %v", err)
				continue
			}

			// read response
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				Error.Printf("Could not read response body: %v", err)
				continue
			}

			Info.Printf("    response: %v", string(body))
		}

		// TODO: update the event id index to track a number of things about shows:
		//     * which shows are new
		//     * which shows have disappeared from bandsintown
		//     * which shows have changed information

		// save data to FLSData
		/////////////////////////

		// trigger cache rebuild goroutine
		////////////////////////////////////

		// sleep until next query
		///////////////////////////
		time.Sleep(time.Until(nextPollTime))

	}

}

// cache rebuild goroutine
////////////////////////////////////////////////////////////////////////////////

// RebuildCache updates the cache based on the current state of FLSData
func RebuildCache() {
	fmt.Println("called RebuildCache")
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

	// handle flags
	configPath := flag.String("c", "", "path to the fls-core config file")
	flag.Parse()

	if *configPath == "" {
		Error.Printf("Must provide path to config!")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// check for api key environment variable
	if os.Getenv("BANDSINTOWN_API_KEY") == "" {
		Error.Printf("Must set environment variable BANDSINTOWN_API_KEY!")
		os.Exit(1)
	}

	// read config
	configJSON, err := ioutil.ReadFile(*configPath)
	if err != nil {
		Error.Printf("Could not read %v: %v", *configPath, err)
	}
	fmt.Printf("%v\n\n", string(configJSON))

	var cfg FLSConfig
	err = json.Unmarshal(configJSON, &cfg)
	if err != nil {
		Error.Printf("Could not parse JSON in %v: %v", *configPath, err)
	}
	Info.Printf("cfg: %v", cfg)

	// launch goroutines
	PollBandsInTown(&cfg)
	os.Exit(0)

	// set routes
	http.HandleFunc("/", RouteRoot)
	http.HandleFunc("/v1/shows", RouteGetShows)
	http.HandleFunc("/v1/artists", RouteGetArtists)

	// start serving
	port := ":8001"
	Info.Printf("fls-core serving on port %v\n", port)
	http.ListenAndServe(port, nil)

}
