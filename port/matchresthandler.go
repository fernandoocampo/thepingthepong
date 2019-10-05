package port

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/fernandoocampo/thepingthepong/application/matchapp"
	"github.com/fernandoocampo/thepingthepong/domain"
)

// newMatch contains data to start a match
type newMatch struct {
	Player1ID string `json:"player1ID"`
	Player2ID string `json:"player2ID"`
}

// MatchRestHandler implements rest handler to expose matches logic
type matchRestHandler struct {
	service matchapp.MatchService
}

// NewMatchRestHandler creates a basic match rest handler
func NewMatchRestHandler(matchService matchapp.MatchService) RestHandler {
	log.Infof("creating match rest handler")
	return &matchRestHandler{
		service: matchService,
	}
}

// GetAll get all records or those that matches a given criteria
func (m *matchRestHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// GetByID get record by id
func (m *matchRestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// Create creates a new record
func (m *matchRestHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Info("starting create handler for match rest handler")
	// We can obtain the session token from the requests cookies, which come with every request
	status, ok := validateToken(r)
	if !ok {
		// If the cookie is not set, return an unauthorized status
		w.WriteHeader(status.StatusCode)
		return
	}
	// context constraint
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	defer r.Body.Close()

	var match newMatch
	log.Infof("request to play a match is: %v", r.Body)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&match); err != nil {
		log.Warnf("payload to create match is bad: %s", err.Error())
		RespondRestWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	log.Infof("consuming create from service to play a match: %v", match)
	savedMatch, err := m.service.Play(ctx, domain.Key(match.Player1ID),
		domain.Key(match.Player2ID))

	if err != nil {
		log.Errorf("something goes wront at service to play a match: %v, got: %s", match, err.Error())
		RespondRestWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondRestWithJSON(w, http.StatusOK, savedMatch)
}

// Update updates the data of existing record.
func (m *matchRestHandler) Update(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// Delete deletes a record.
func (m *matchRestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

// Health returns the health of this service
func (m *matchRestHandler) Health(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
