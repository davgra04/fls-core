package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// config
////////////////////////////////////////////////////////////////////////////////

// logging
////////////////////////////////////////////////////////////////////////////////

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func InitLogging(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// bandsintown
////////////////////////////////////////////////////////////////////////////////

// show data
////////////////////////////////////////////////////////////////////////////////

// REST api handles
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
	Info.Printf("fls-core serving on port %v\n", port)
	http.ListenAndServe(port, nil)

}
