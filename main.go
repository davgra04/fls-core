package main

import (
	"flag"
	"net/http"
	"os"

	"gitlab.devgru.cc/devgru/fls-core/fls"
)

func handleFlagsAndEnv() (string, bool) {
	configPath := flag.String("c", "", "path to the fls-core config file")
	localOnly := flag.Bool("l", false, "listen on localhost only")
	flag.Parse()

	if *configPath == "" {
		flag.PrintDefaults()
		fls.Error.Fatal("Must provide path to config!")
	}

	// check for api key environment variable
	if os.Getenv("BANDSINTOWN_API_KEY") == "" {
		fls.Error.Fatalf("Must set environment variable BANDSINTOWN_API_KEY!")
	}

	return *configPath, *localOnly
}

// main
////////////////////////////////////////////////////////////////////////////////

func main() {
	// initialize logging
	fls.InitLogging()

	// handle flags and check env vars
	configPath, localOnly := handleFlagsAndEnv()

	// read config
	fls.ReadConfig(configPath)

	// init DB if necessary
	fls.InitializeDatabase()
	defer fls.DB.Close()

	// launch bandsintown polling goroutine
	// go fls.PollBandsInTown()

	// set routes
	http.HandleFunc("/", fls.RouteRoot)
	http.HandleFunc("/v1/artists", fls.RouteArtists)
	// http.HandleFunc("/v1/shows", fls.RouteShows)

	// serve
	addr := ":8001"
	if localOnly {
		addr = "localhost" + addr
	}
	fls.Info.Printf("fls-core listening on %v\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fls.Error.Fatalf("Failed to ListenAndServe: %v", err)
	}
}
