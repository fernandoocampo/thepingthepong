package main

import (
	"fmt"
	"os"

	"github.com/fernandoocampo/thepingthepong/application/playerapp"
	"github.com/fernandoocampo/thepingthepong/common/logging"
	"github.com/fernandoocampo/thepingthepong/infra/repository"
	"github.com/fernandoocampo/thepingthepong/port"
	"github.com/sirupsen/logrus"
)

var log *logging.Handle
var webserver port.WebServer

func init() {
	// initialize logger
	initLogger()
	// initialize inversion of control
	initIoC()
}

// initLogger Initialize logger
func initLogger() {
	fmt.Println("... starting thepingthepong service logger")
	fmt.Println()
	loglevelval := "debug" // TODO add configuration
	if loglevelval == "" {
		fmt.Println("You have not defined the log level SERVICE_LOG_LEVEL environment variable, so we are going to use [Info] level")
		fmt.Println()
		loglevelval = "Info"
	}

	logformatval := "json" // TODO add configuration
	if logformatval == "" {
		fmt.Println("You have not defined the log level SERVICE_LOG_FORMAT environment variable, so we are going to use [text] value")
		logformatval = "json"
	}

	options := logging.Options{
		LogLevel:  loglevelval,
		LogFormat: logformatval,
		LogFields: logrus.Fields{"pkg": "main", "srv": "thepingthepong"},
	}
	fmt.Printf("\n%+v\n", options)
	fmt.Println()
	var err error
	// load log for dao package
	log, err = logging.NewLogger(options)
	if err != nil {
		fmt.Printf("cant load logger: %v", err)
		os.Exit(1)
	}
}

// initIoC initializes the dao and service used service and controller.
func initIoC() {
	// initialize repository layer
	repo := repository.NewPlayerRepositoryOnMemory(5)
	// initialize application layer
	service := playerapp.NewBasicPlayerService(repo)
	// initialize port layer
	// initialize rest handler
	resthandler := port.NewPlayerRestHandler(service)
	// initialize web server
	webserver = port.NewWebServer(resthandler)
}

// initHTTPServer start webserver on the configuration parameter host.
func initHTTPServer() {
	log.Println("Starting thepingpong service")
	portvalue := "8287"  // add setting parameter
	if portvalue == "" { // check service port
		fmt.Println("You have not defined the application port SERVICE_APP_PORT environment variable, so we are going to use [8287] value")
		fmt.Println()
		portvalue = "8287"
	}
	log.Println("Starting application on ", portvalue)

	webserver.StartWebServer(portvalue)
}

func main() {
	initHTTPServer()
}
