package fls

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// config
////////////////////////////////////////////////////////////////////////////////

// Cfg is the global FLSConfig for whole program
var Cfg *FLSConfig

// FLSConfig represents the input configuration for fls-core
type FLSConfig struct {
	Artists               []string `json:"artists"`                             // List of followed artists
	RefreshPeriodSeconds  int      `json:"bandsintown_query_period_seconds"`    // number of seconds between FLSData refreshes
	RateLimitMillis       int      `json:"bandsintown_rate_limit_ms"`           // limit on the number of BandsInTown API requests per second
	MaxConcurrentRequests int      `json:"bandsintown_max_concurrent_requests"` // maximum number of concurrent requests to the BandsInTown API
}

// ReadConfig asdf
func ReadConfig(configPath string) {
	configJSON, err := ioutil.ReadFile(configPath)
	if err != nil {
		Error.Fatalf("Could not read %v: %v", configPath, err)
	}

	err = json.Unmarshal(configJSON, &Cfg)
	if err != nil {
		Error.Fatalf("Could not parse JSON in %v: %v", configPath, err)
	}
}

// logging
////////////////////////////////////////////////////////////////////////////////

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// InitLogging sets up loggers for each logging level (ERROR, WARN, INFO, TRACE)
func InitLogging() {

	logPath := "log.txt"
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not open logfile %v: %v", logPath, err))
	}

	traceHandle := ioutil.Discard
	infoHandle := io.MultiWriter(os.Stdout, logFile)
	warningHandle := io.MultiWriter(os.Stdout, logFile)
	errorHandle := io.MultiWriter(os.Stdout, logFile)

	Trace = log.New(traceHandle, "[TRACE] ", log.Ldate|log.Ltime)
	Info = log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime)
	Warning = log.New(warningHandle, "[WARNING] ", log.Ldate|log.Ltime)
	Error = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime)

}
