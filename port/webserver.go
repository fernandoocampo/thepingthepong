package port

import (
	"net/http"

	"log"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// WebServer defines behavior to start a http server
type WebServer interface {
	// Starts the http server in the given port number
	StartWebServer(port string)
}

type restServer struct {
	playerRestHandler RestHandler
}

// NewWebServer instance of a person handler
func NewWebServer(restHandler RestHandler) WebServer {
	return &restServer{
		playerRestHandler: restHandler,
	}
}

// StartWebServer starts the http server for this service
// on the given http port.
func (w *restServer) StartWebServer(port string) {
	router := newRouter(w.playerRestHandler)

	log.Println("Starting HTTP service at ", port)
	originsOk := handlers.AllowedOrigins([]string{"*"})
	err := http.ListenAndServe(":"+port, handlers.CORS(originsOk)(router))

	if err != nil {
		log.Panicf("An error occured starting HTTP listener at port %s, error %s", port, err)
	}
}

// NewRouter returns a pointer to a mux.Router we can use as a handler.
func newRouter(restHandler RestHandler) *mux.Router {
	// Create an instance of the Gorilla router
	// Gorilla router matches incoming requests against a list of
	// registered routes and calls a handler for the route that matches
	// the URL or other conditions
	router := mux.NewRouter().StrictSlash(true)

	// Get cash register by id
	router.Methods("GET").
		Path("/players").
		Name("getAllPlayers").
		HandlerFunc(restHandler.GetAll)

	// Get cash register by id
	router.Methods("GET").
		Path("/players/{playerid}").
		Name("getPlayerById").
		HandlerFunc(restHandler.GetByID)

	// Post to create a cash register
	router.Methods("POST").
		Path("/players").
		Name("createPlayer").
		HandlerFunc(restHandler.Create)

	// get health status of this service.
	router.Methods("GET").
		Path("/health").
		Name("health").
		HandlerFunc(restHandler.Health) // what's the health

	return router
}
