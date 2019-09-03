package main

import (
	"fmt"
	"net/http"
)

// config
////////////////////////////////////////////////////////////////////////////////

// bandsintown
////////////////////////////////////////////////////////////////////////////////

// show data
////////////////////////////////////////////////////////////////////////////////

// REST api handles
////////////////////////////////////////////////////////////////////////////////

// RouteRoot displays the name of the service
func RouteRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "fls-core")
}

// RouteGetShows returns a json object containing upcoming shows
func RouteGetShows(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "{\"shows\":[{\"show_id\":0,\"artist\":\"Tame Impala\",\"date\":\"Sat Aug31 7pm\",\"venue\":\"White Oak Music Hall\",\"lineup\":\"White Denim\",\"city\":\"Houston\",\"region\":\"TX\"},{\"show_id\":1,\"artist\":\"King Gizzard and the Lizard Wizard\",\"date\":\"Sun Sep01 7pm\",\"venue\":\"White Oak Music Hall\",\"lineup\":\"Mildlife, Orb\",\"city\":\"Houston\",\"region\":\"TX\"},{\"show_id\":2,\"artist\":\"Unknown Mortal Orchestra\",\"date\":\"Mon Sep02 7pm\",\"venue\":\"White Oak Music Hall\",\"lineup\":\"Shakey Graves\",\"city\":\"Houston\",\"region\":\"TX\"}]}")
}

// main
////////////////////////////////////////////////////////////////////////////////

func main() {
	http.HandleFunc("/", RouteRoot)
	http.HandleFunc("/v1/shows", RouteGetShows)

	port := ":8001"

	fmt.Printf("fls-core serving on port %v\n", port)
	http.ListenAndServe(port, nil)
}
