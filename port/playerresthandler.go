package port

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fernandoocampo/thepingthepong/application/playerapp"
	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/gorilla/mux"
)

// newPlayer contains data to save a new player
type newPlayer struct {
	Names  string `json:"names"`
	Wins   int    `json:"wins,omitempty"`
	Losses int    `json:"wins,omitempty"`
}

type playerRestHandler struct {
	service playerapp.PlayerService
}

// NewPlayerRestHandler instance of a basic implementation of player rest handler
func NewPlayerRestHandler(playerService playerapp.PlayerService) RestHandler {
	return playerRestHandler{
		service: playerService,
	}
}

// GetAll get all records or those that matches a given criteria
func (p playerRestHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Read parameters in the query url
	filters := r.URL.Query()
	sortedparam := filters.Get("sorted")

	var sorted bool
	var err error
	var players []domain.Player

	// not filters get all
	if strings.EqualFold("true", sortedparam) {
		sorted = true
	}

	players, err = p.service.FindAll(sorted)

	if err != nil {
		RespondRestWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondRestWithJSON(w, http.StatusOK, players)

}

// GetByID get record by id
func (p playerRestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Read the 'playerid' path parameter from the mux map
	var playerid = mux.Vars(r)["playerid"]
	player, err := p.service.FindByID(playerid)
	if err != nil {
		RespondRestWithError(w, http.StatusInternalServerError, err.Error())
	}
	RespondRestWithJSON(w, http.StatusOK, player)
}

// Create creates a new record
func (p playerRestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var player newPlayer
	// close the body buffer at the end of the function
	defer r.Body.Close()
	// Create the decoder for bank regarding to the body request
	decoder := json.NewDecoder(r.Body)
	// Get all the data of the request and map to player struct
	// if error we response with error message
	if err := decoder.Decode(&player); err != nil {
		RespondRestWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	_, err := p.service.Create(player.Names, player.Wins, player.Losses)
	if err != nil {
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
