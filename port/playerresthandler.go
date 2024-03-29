package port

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/fernandoocampo/thepingthepong/application/playerapp"
	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/gorilla/mux"
)

// newPlayer contains data to save a new player
type newPlayer struct {
	Names  string `json:"names"`
	Wins   int    `json:"wins,omitempty"`
	Losses int    `json:"losses,omitempty"`
}

type playerRestHandler struct {
	service playerapp.PlayerService
}

// NewPlayerRestHandler instance of a basic implementation of player rest handler
func NewPlayerRestHandler(playerService playerapp.PlayerService) RestHandler {
	log.Infof("creating player rest handler")
	return playerRestHandler{
		service: playerService,
	}
}

const timeout = time.Second * 5

// GetAll get all records or those that matches a given criteria
func (p playerRestHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Info("initializing player rest handler to get all")
	log.Info("checking if token is valid")

	// context constraint
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Read parameters in the query url
	filters := r.URL.Query()
	sortedparam := filters.Get("sorted")
	log.Infof("finding all players with parameter sorted: %s", sortedparam)

	var sorted bool
	var err error
	var players []domain.Player

	// not filters get all
	if strings.EqualFold("true", sortedparam) {
		sorted = true
	}

	log.Infof("getting ready to find all players")
	players, err = p.service.FindAll(ctx, sorted)

	if err != nil {
		log.Errorf("something goes wrong on service to get all players: %s", err.Error())
		RespondRestWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondRestWithJSON(w, http.StatusOK, players)
}

// GetByID get record by id
func (p playerRestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Info("starting get by id handler")
	// context constraint
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Read the 'playerid' path parameter from the mux map
	var playerid = mux.Vars(r)["playerid"]
	log.Infof("player id to get one player is: %s", playerid)
	log.Infof("getting ready to find player with id: %s on service", playerid)
	player, err := p.service.FindByID(ctx, domain.Key(playerid))
	if err != nil {
		RespondRestWithError(w, http.StatusInternalServerError, err.Error())
	}
	RespondRestWithJSON(w, http.StatusOK, player)
}

// Create creates a new record
func (p playerRestHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Info("starting create handler")
	// We can obtain the session token from the requests cookies, which come with every request
	status, ok := validateToken(r)
	if !ok {
		// If the cookie is not set, return an unauthorized status
		w.WriteHeader(status.StatusCode)
		return
	}
	// context constraint
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var player newPlayer
	// close the body buffer at the end of the function
	defer r.Body.Close()
	// Create the decoder for bank regarding to the body request
	decoder := json.NewDecoder(r.Body)
	// Get all the data of the request and map to player struct
	// if error we response with error message
	if err := decoder.Decode(&player); err != nil {
		log.Warnf("payload to create player is bad: %s", err.Error())
		RespondRestWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Infof("consuming create from service to create player: %v", player)
	_, err := p.service.Create(ctx, player.Names, player.Wins, player.Losses)
	if err != nil {
		log.Errorf("something goes wront at service to create player: %v, got: %s", player, err.Error())
		RespondRestWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondRestWithJSON(w, http.StatusOK, "created!")

}

// Update updates the data of existing record.
func (p playerRestHandler) Update(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// Delete deletes a record.
func (p playerRestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// Health returns the health of this service
func (p playerRestHandler) Health(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
