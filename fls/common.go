package fls

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// config
////////////////////////////////////////////////////////////////////////////////

// Cfg is the global FLSConfig for whole program
var Cfg *Config

// APIKey is the bandsintown REST API key
var APIKey string

// Config represents the input configuration for fls-core
type Config struct {
	RefreshPeriodSeconds int           `json:"bandsintown_refresh_period_seconds"` // number of seconds between events refresh
	RequestsPerSecond    float32       `json:"bandsintown_requests_per_second"`    // limit on the number of BandsInTown API requests per second
	RPSPeriod            time.Duration `json:"-"`                                  // period for RequestsPerSecond
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

	Cfg.RPSPeriod = time.Duration(1000.0/Cfg.RequestsPerSecond) * time.Millisecond
	Info.Printf("Cfg.RPSPeriod: %v\n", Cfg.RPSPeriod)
}

////////////////////////////////////////////////////////////////////////////////
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

////////////////////////////////////////////////////////////////////////////////
// database access object
////////////////////////////////////////////////////////////////////////////////

// DAO is the global database access object
var DAO *DBAccessObject
