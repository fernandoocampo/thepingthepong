package port_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fernandoocampo/thepingthepong/application/matchapp"
	"github.com/fernandoocampo/thepingthepong/application/playerapp"
	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/fernandoocampo/thepingthepong/infra/repository"
	"github.com/fernandoocampo/thepingthepong/port"
	"github.com/gorilla/mux"
)

func TestCreateAMatch(t *testing.T) {
	repo := repository.NewPlayerRepositoryOnMemory(1)
	playerService := playerapp.NewBasicPlayerService(&repo)
	matchService := matchapp.NewBasicMatchService(playerService)
	matchhandler := port.NewMatchRestHandler(matchService)

	// Given a the following players to start a match.
	player1ID, err := playerService.Create(context.TODO(), "Jan-Ove Waldner", 0, 0)
	assertNoError(t, err)
	player2ID, err := playerService.Create(context.TODO(), "Timo Boll", 0, 0)
	assertNoError(t, err)

	strjson := fmt.Sprintf(`{"player1ID": "%s", "player2ID": "%s"}`, player1ID, player2ID)
	req, errreq := http.NewRequest("POST", "/matches", bytes.NewBuffer([]byte(strjson)))
	assertNoError(t, errreq)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/matches", matchhandler.Create).Methods("POST")

	// generates token
	tokencookie, tokenok := generateToken(t)
	if !tokenok {
		t.Fatalf("token cannot be generated, we got this token")
	}

	req.AddCookie(tokencookie)

	// When client consumes a rest api.
	r.ServeHTTP(rr, req)

	// Then we check the result of the player found.

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var got domain.MatchReport
	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&got)

	assertNoError(t, err)

	if got.Winner == nil {
		t.Errorf("a winner between player: %q and player: %q was expected, but none won",
			player1ID, player2ID)
	}
	if got.Loser == nil {
		t.Errorf("a loser between player: %q and player: %q was expected, but none loose",
			player1ID, player2ID)
	}
	if len(got.Narrative) == 0 {
		t.Errorf("a fulled narrative was expected, but got: %v", got.Narrative)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("error was not expected, but: %s", err)
	}
}
