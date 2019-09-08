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
	"sort"
	"strings"
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
	// fields populated by BandsInTown
	ID             string                 `json:"id"`
	ArtistID       string                 `json:"artist_id"`
	URL            string                 `json:"url"`
	OnSaleDatetime string                 `json:"on_sale_datetime"` // 2017-03-01T18:00:00
	Datetime       string                 `json:"datetime"`
	Description    string                 `json:"description"`
	Venue          *BandsInTownVenueData  `json:"venue"`
	Offers         []BandsInTownOfferData `json:"offers"`
	Lineup         []string               `json:"lineup"`
	// fields populated by fls-data
	DateAdded   int `json:"date_added"`
	DateUpdated int `json:"date_updated"`
	DateRemoved int `json:"date_removed"`
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
	QueryDate int64                  `json:"query_date"` // UNIX timestamp for time of last bandsintown API call
	Artist    BandsInTownArtistData  `json:"artist"`
	Events    []BandsInTownEventData `json:"events"`
}

// fls objects
////////////////////////////////////////////////////////////////////////////////

// global FLSConfig for whole program
var flscfg *FLSConfig

// FLSConfig represents the input configuration for fls-core
type FLSConfig struct {
	Artists               []string `json:"artists"`                             // List of followed artists
	RefreshPeriodSeconds  int      `json:"bandsintown_query_period_seconds"`    // number of seconds between FLSData refreshes
	RateLimitMillis       int      `json:"bandsintown_rate_limit_ms"`           // limit on the number of BandsInTown API requests per second
	MaxConcurrentRequests int      `json:"bandsintown_max_concurrent_requests"` // maximum number of concurrent requests to the BandsInTown API
}

// FLSData represents all of the non-cache data in fls-core
type FLSData struct {
	BandsInTownData map[string]*BandsInTownData `json:"bandsintown_data"` // maps artist_name -> bandsintown artist info
}

// FLSEventData represents an artist's show data for a particular region
type FLSEventData struct {
	ShowID    string `json:"show_id"`
	Artist    string `json:"artist"`
	Date      string `json:"date"`
	DateAdded string `json:"date_added"`
	Venue     string `json:"venue"`
	Lineup    string `json:"lineup"`
	City      string `json:"city"`
	Region    string `json:"region"`
}

// GetShowsResponse represents an artist's show data for a particular region, used by RouteGetShows endpoint
type GetShowsResponse struct {
	QueryDate int64          `json:"query_date"`
	Region    string         `json:"region"`
	Shows     []FLSEventData `json:"shows"`
}

// GetCachedShowsResponse TODO TODO TODO
func GetCachedShowsResponse(region string) string {

	// try to read cached json value, generate if none is there
	// TODO

	// just generating every time for now...
	flsdata := ReadFLSData("data.json")
	Info.Printf("flsdata: %v", flsdata)

	showData := GetShowsResponse{QueryDate: time.Now().Unix(), Region: region}

	for _, artist := range flscfg.Artists {
		for _, event := range flsdata.BandsInTownData[artist].Events {
			// Info.Printf("venue: %v", event.Venue)
			// Info.Printf("region: %v", event.Venue.Region)

			if event.Venue.Region == region {
				showData.Shows = append(showData.Shows, FLSEventData{
					ShowID:    event.ID,
					Artist:    artist,
					Date:      event.Datetime,
					DateAdded: string(event.DateAdded),
					Venue:     event.Venue.Name,
					Lineup:    strings.Join(event.Lineup, ", "),
					City:      event.Venue.City,
					Region:    event.Venue.Region,
				})
				Info.Printf("Found a TX show: %v", showData.Shows[len(showData.Shows)-1])
			}
		}
	}

	// sort shows by date
	sort.SliceStable(showData.Shows, func(i, j int) bool {
		return showData.Shows[i].Date < showData.Shows[j].Date
	})

	// marshal to json
	showDataJSON, err := json.Marshal(showData)
	if err != nil {
		Error.Printf("Failed to marshal show data: %v", err)
	}

	return string(showDataJSON)

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

	showDataJSON := GetCachedShowsResponse("TX")

	// fmt.Fprintf(w, ReadTestJSON())
	fmt.Fprintf(w, showDataJSON)
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

// ArtistRequest TODO TODO TODO
type ArtistRequest struct {
	artist, url string
}

// PollBandsInTown periodically polls BandsInTown for show data, saves to FLSData, and initiates cache rebuild goroutine
func PollBandsInTown() {
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
	// flsdata := FLSData{}

	// initialize fresh flsdata

	flsdata := FLSData{BandsInTownData: make(map[string]*BandsInTownData)}

	// flsdata := make(map[string]*BandsInTownData)
	// for _, artist := range flscfg.Artists {
	// 	flsdata[artist] = &BandsInTownData{}
	// 	// flsdata[artist] = make(map[string]T)
	// }

	fmt.Print("HEYO")

	for {

		// set times for this polling period
		startTime := time.Now()
		nextPollTime := startTime.Add(time.Duration(flscfg.RefreshPeriodSeconds) * time.Second)

		// get data from api
		//////////////////////

		// load requests into requests channel
		requests := make(chan ArtistRequest, 1000)
		for _, artist := range flscfg.Artists {
			url := fmt.Sprintf("https://rest.bandsintown.com/artists/%s/events?app_id=%s", url.PathEscape(artist), apiKey)
			requests <- ArtistRequest{url: url, artist: artist}
			// Info.Printf("    preparing request for %-40v [%v]", artist, url)
		}
		close(requests)

		// make requests using limiter
		for ar := range requests {
			<-limiter // wait for limiter

			// build and do request
			Info.Printf("    requesting %-40v [%v]", ar.artist, ar.url)
			req, err := http.NewRequest(http.MethodGet, ar.url, nil)
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

			// Info.Printf("    response: %v", string(body))

			// parse response
			var events []BandsInTownEventData
			err = json.Unmarshal(body, &events)
			if err != nil {
				Error.Printf("Could not parse JSON in response: %v", err)
			}

			// write events to flsdata

			if _, ok := flsdata.BandsInTownData[ar.artist]; !ok {
				flsdata.BandsInTownData[ar.artist] = &BandsInTownData{QueryDate: time.Now().Unix()}
			}

			flsdata.BandsInTownData[ar.artist].Events = events

			// for _, e := range events {

			// 	flsdata[ar.artist].Events = append(flsdata[ar.artist].Events, e)
			// }

		}

		// save data to data.json
		///////////////////////////
		flsdataJSON, err := json.Marshal(flsdata)
		if err != nil {
			Error.Printf("UH OH")
			os.Exit(1)
		}
		err = ioutil.WriteFile("data.json", flsdataJSON, 0666)
		if err != nil {
			Error.Printf("UH OH")
			os.Exit(1)
		}

		// TODO: update the event id index to track a number of things about shows:
		//     * which shows are new
		//     * which shows have disappeared from bandsintown
		//     * which shows have changed information

		// flsdata.BandsInTownData.Events = ""

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

// ReadFLSData TODO TODO TODO
func ReadFLSData(flsdataPath string) *FLSData {

	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		Error.Panic(err)
	}

	flsdata := &FLSData{}
	err = json.Unmarshal(data, &flsdata)
	if err != nil {
		Error.Panicf("Failed to unmarshal %v: %v", flsdataPath, err)
	}

	return flsdata
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

	err = json.Unmarshal(configJSON, &flscfg)
	if err != nil {
		Error.Printf("Could not parse JSON in %v: %v", *configPath, err)
	}
	Info.Printf("flscfg: %v", flscfg)

	// launch goroutines
	// go PollBandsInTown()

	// set routes
	http.HandleFunc("/", RouteRoot)
	http.HandleFunc("/v1/shows", RouteGetShows)
	http.HandleFunc("/v1/artists", RouteGetArtists)

	// start serving
	port := "localhost:8001"
	Info.Printf("fls-core serving on port %v\n", port)
	http.ListenAndServe(port, nil)

}
