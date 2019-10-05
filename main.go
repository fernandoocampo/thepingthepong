package main

import (
	"fmt"
	"os"

	"github.com/fernandoocampo/thepingthepong/application/authapp"
	"github.com/fernandoocampo/thepingthepong/application/matchapp"
	"github.com/fernandoocampo/thepingthepong/application/playerapp"
	"github.com/fernandoocampo/thepingthepong/common/logging"
	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/fernandoocampo/thepingthepong/infra/repository"
	"github.com/fernandoocampo/thepingthepong/port"
	"github.com/sirupsen/logrus"
)

var log *logging.Handle
var webserver port.WebServer

func init() {
	// load configuration
	loadConfiguration()
	// initialize logger
	initLogger()
	// initialize inversion of control
	initIoC()
}

func loadConfiguration() {
	domain.LoadConfiguration("conf/")
}

// initLogger Initialize logger
func initLogger() {
	fmt.Println("... starting thepingthepong service logger")
	fmt.Println()
	loglevelval := domain.Configuration.Log.Main.Level
	if loglevelval == "" {
		fmt.Println("You have not defined the log level for main log, so warn will be used instead")
		fmt.Println()
		loglevelval = "Warn"
	}

	logformatval := domain.Configuration.Log.Main.Format
	if logformatval == "" {
		fmt.Println("You have not defined the log format for main log, so json will be used instead")
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
	// initializes log for all modules
	port.InitLog(domain.Configuration.Log.Port)
	domain.InitLog(domain.Configuration.Log.Domain)
	repository.InitLog(domain.Configuration.Log.Repository)
	matchapp.InitLog(domain.Configuration.Log.Matchapp)
	playerapp.InitLog(domain.Configuration.Log.Playerapp)

}

// initIoC initializes the dao and service used service and controller.
func initIoC() {
	// initialize repository layer
	repo := repository.NewPlayerRepositoryOnMemory(5)
	// initialize application layer
	playerService := playerapp.NewBasicPlayerService(&repo)
	matchService := matchapp.NewBasicMatchService(playerService)
	authservice := authapp.NewBasicAuthenticator()
	// initialize port layer
	// initialize rest handler
	playerhandler := port.NewPlayerRestHandler(playerService)
	matchhandler := port.NewMatchRestHandler(matchService)
	authhandler := port.NewBasicAuthRestHandler(authservice)
	// initialize web server
	webserver = port.NewWebServer(playerhandler, matchhandler, authhandler)
}

// initHTTPServer start webserver on the configuration parameter host.
func initHTTPServer() {
	log.Println("Starting thepingpong service")
	portvalue := "8287"  // TODO add setting parameter
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
