package common

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
)

// config
////////////////////////////////////////////////////////////////////////////////

// global FLSConfig for whole program
var Cfg *FLSConfig

// FLSConfig represents the input configuration for fls-core
type FLSConfig struct {
	Artists               []string `json:"artists"`                             // List of followed artists
	RefreshPeriodSeconds  int      `json:"bandsintown_query_period_seconds"`    // number of seconds between FLSData refreshes
	RateLimitMillis       int      `json:"bandsintown_rate_limit_ms"`           // limit on the number of BandsInTown API requests per second
	MaxConcurrentRequests int      `json:"bandsintown_max_concurrent_requests"` // maximum number of concurrent requests to the BandsInTown API
}

// logging
////////////////////////////////////////////////////////////////////////////////

// Logger todo
type Logger struct {
	logger *log.Logger
	color  func(format string, args ...interface{}) string
}

// NewLogger asdf
func NewLogger(_logger *log.Logger, _color func(format string, args ...interface{}) string) *Logger {
	return &Logger{logger: _logger, color: _color}
}

// Printf asdf
func (l *Logger) Printf(format string, args ...interface{}) {
	// fmt.Printf("called Printf(%v, %v)\n", format, args)
	colorizedFormat := l.color(format)
	// fmt.Printf("got colorizedFormat: %v\n", colorizedFormat)
	l.logger.Printf(colorizedFormat, args...)
	// fmt.Printf("got retval: %v\n", retval)
	// return retval
}

// Panicf asdf
func (l *Logger) Panicf(format string, args ...interface{}) int {
	return l.Panicf(l.color(format, args...))
}

var (
	Trace   *Logger
	Info    *Logger
	Warning *Logger
	Error   *Logger
)

// InitLogging TODO TODO TODO
func InitLogging() {

	logPath := "log.txt"
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Sprintf("Could not open logfile %v: %v", logPath, err))
	}

	// // Use your own io.Writer output
	// color.New(color.FgBlue).Fprintln(myWriter, "blue color!")

	// blue := color.New(color.FgBlue)
	// blue.Fprint(writer, "This will print text in blue.")

	traceHandle := ioutil.Discard
	infoHandle := io.MultiWriter(os.Stdout, logFile)
	warningHandle := io.MultiWriter(os.Stdout, logFile)
	errorHandle := io.MultiWriter(os.Stdout, logFile)

	Trace = NewLogger(
		log.New(traceHandle, color.GreenString("[TRACE] "), log.Ldate|log.Ltime),
		color.GreenString,
	)

	Info = NewLogger(
		log.New(infoHandle, color.CyanString("[INFO] "), log.Ldate|log.Ltime),
		color.CyanString,
	)

	Warning = NewLogger(
		log.New(warningHandle, color.YellowString("[WARNING] "), log.Ldate|log.Ltime),
		color.YellowString,
	)

	Error = NewLogger(
		log.New(errorHandle, color.RedString("[ERROR] "), log.Ldate|log.Ltime),
		color.RedString,
	)

}

// misc
////////////////////////////////////////////////////////////////////////////////

// ReadTestJSON reads and returns test.json
func ReadTestJSON() string {

	dat, err := ioutil.ReadFile("test.json")
	if err != nil {
		panic(err)
	}
	// common.Info.Println(string(dat))

	return string(dat)
}

// ReadConfig asdf
func ReadConfig(configPath *string) {
	configJSON, err := ioutil.ReadFile(*configPath)
	if err != nil {
		Error.Printf("Could not read %v: %v", *configPath, err)
	}
	// common.Info.Printf("%v\n\n", string(configJSON))

	err = json.Unmarshal(configJSON, &Cfg)
	if err != nil {
		Error.Printf("Could not parse JSON in %v: %v", *configPath, err)
	}
	// Info.Printf("common.Cfg: %v", Cfg)

}
