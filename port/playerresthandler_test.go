package port_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fernandoocampo/thepingthepong/application/playerapp"
	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/fernandoocampo/thepingthepong/infra/repository"
	"github.com/fernandoocampo/thepingthepong/port"
	"github.com/gorilla/mux"
)

func TestCreateAPlayer(t *testing.T) {
	repo := repository.NewPlayerRepositoryOnMemory(1)
	service := playerapp.NewBasicPlayerService(repo)
	playerhandler := port.NewPlayerRestHandler(service)

	// Given a get request to find by Id a player.
	strnames := "Hugo Calderano"
	wins := 5
	losses := 3
	strjson := fmt.Sprintf(`{"names" : "%s", "wins": %d, "losses": %d}`, strnames, wins, losses)
	req, errreq := http.NewRequest("POST", "/players", bytes.NewBuffer([]byte(strjson)))

	if errreq != nil {
		t.Fatal(errreq)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(playerhandler.GetByID)
	r := mux.NewRouter()
	r.HandleFunc("/players", playerhandler.Create).Methods("POST")

	// When client consumes a rest api.
	// handler.ServeHTTP(rr, req)
	r.ServeHTTP(rr, req)

	// Then we check the result of the player found.

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check at the repository if the player was saved
	// context constraint
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	players, errfindall := repo.FindAll(ctx, false)
	if errfindall != nil {
		t.Fatal(errfindall)
	}
	exists := false
	for _, player := range players {
		if player.Names == strnames {
			exists = true
			break
		}
	}
	if !exists {
		t.Errorf("The player %s was not stored in the database", strnames)
	}
}

func TestGetAnExistingPlayer(t *testing.T) {
	repo := repository.NewPlayerRepositoryOnMemory(1)
	service := playerapp.NewBasicPlayerService(repo)
	playerhandler := port.NewPlayerRestHandler(service)
	// save a player in the db.
	newplayer := domain.NewPlayer("Wang Liqin")
	// context constraint
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	errsave := repo.Save(ctx, newplayer)
	if errsave != nil {
		t.Fatalf("A player cannot be saved because of: %s", errsave.Error())
	}
	// Given a get request to find by Id a player.
	req, errreq := http.NewRequest("GET", "/players/"+newplayer.ID, nil)
	if errreq != nil {
		t.Fatal(errreq)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(playerhandler.GetByID)
	r := mux.NewRouter()
	r.HandleFunc("/players/{playerid}", playerhandler.GetByID).Methods("GET")

	// When client consumes a rest api.
	// handler.ServeHTTP(rr, req)
	r.ServeHTTP(rr, req)

	// Then we check the result of the player found.

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := fmt.Sprintf(`{"id":"%s","names":"%s","wins":%d,"losses":%d,"created":"%s","updated":"%s"}`, newplayer.ID,
		newplayer.Names, newplayer.Wins, newplayer.Losses, newplayer.Created.Format("2006-01-02T15:04:05.999999-07:00"),
		newplayer.Updated.Format("2006-01-02T15:04:05.999999-07:00"))
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetAllPlayers(t *testing.T) {
	repo := repository.NewPlayerRepositoryOnMemory(1)
	service := playerapp.NewBasicPlayerService(repo)
	playerhandler := port.NewPlayerRestHandler(service)
	// save a player in the db.
	newplayernames := []string{"Wang Liqin", "Liu Guoliang", "Ding Ning"}
	// context constraint
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, newplayername := range newplayernames {
		newplayer := domain.NewPlayer(newplayername)
		errsave := repo.Save(ctx, newplayer)
		if errsave != nil {
			t.Fatalf("A player cannot be saved because of: %s", errsave.Error())
		}
	}
	// Given a get request to find by Id a player.
	req, errreq := http.NewRequest("GET", "/players", nil)
	if errreq != nil {
		t.Fatal(errreq)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(playerhandler.GetByID)
	r := mux.NewRouter()
	r.HandleFunc("/players", playerhandler.GetAll).Methods("GET")

	// When client consumes a rest api.
	// handler.ServeHTTP(rr, req)
	r.ServeHTTP(rr, req)

	// Then we check the result of the player found.

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	for _, newplayername := range newplayernames {

		if !strings.Contains(rr.Body.String(), newplayername) {
			t.Errorf("player %s was not found in the result body: got %v",
				newplayername, rr.Body.String())
		}
	}
}
