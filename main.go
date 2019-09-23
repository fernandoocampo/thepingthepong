package main

import (
	"fmt"
	"log"

	"github.com/fernandoocampo/thepingthepong/application/playerapp"
	"github.com/fernandoocampo/thepingthepong/infra/repository"
	"github.com/fernandoocampo/thepingthepong/port"
)

var webserver port.WebServer

func init() {
	// initialize inversion of control
	initIoC()
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
	log.Println("Starting bank service")
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
