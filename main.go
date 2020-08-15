package main

import (
	"flag"
	"net/http"
	"os"
	"time"

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

func initialize(configPath, dbPath string) {
	// read config
	fls.ReadConfig(configPath)

	// set bandsintown API key
	fls.APIKey = os.Getenv("BANDSINTOWN_API_KEY")

	// set up bandsintown HTTP client
	fls.BandsInTownClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	// set up DAO, init DB if necessary
	fls.DAO = fls.NewSqliteDAO(dbPath)
}

////////////////////////////////////////////////////////////////////////////////
// main
////////////////////////////////////////////////////////////////////////////////

func main() {
	// initialize logging
	fls.InitLogging()

	// handle flags and check env vars
	configPath, localOnly := handleFlagsAndEnv()
	dbPath := "data.sql"

	// initialize
	initialize(configPath, dbPath)
	defer fls.DAO.Close()

	// launch bandsintown polling goroutine
	// go fls.PollBandsInTown()

	// set routes
	http.HandleFunc("/", fls.RouteRoot)
	http.HandleFunc("/artists", fls.RouteArtists)
	// http.HandleFunc("/shows", fls.RouteShows)

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
