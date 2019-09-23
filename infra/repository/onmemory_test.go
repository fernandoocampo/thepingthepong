package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/fernandoocampo/thepingthepong/infra/repository"
)

func TestSavePlayerAndFindIt(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// given a new player
	inmemoryrepo := repository.NewPlayerRepositoryOnMemory(5)
	newplayer := domain.NewPlayer("Ma Long")

	// when we save the player in the inmemory db
	err := inmemoryrepo.Save(ctx, newplayer)

	// then we check that is not an error
	if err != nil {
		t.Errorf("An attempt was made to save the player [%+v], it was expected non errors but we got %s",
			newplayer, err.Error())
	}

	// query the new player and check if data was stored well
	savedplayer, errfind := inmemoryrepo.FindByID(ctx, newplayer.ID)

	if errfind != nil {
		t.Errorf("Querying the player with id [%s], it was expected a player, but we got %s",
			newplayer.ID, err.Error())
	}

	if savedplayer.Names != newplayer.Names {
		t.Errorf(`The stored player is not equal to the given as paremeter one, 
		it was expected a player with name %s", but we got %s`,
			savedplayer.Names, newplayer.ID)
	}
}

func TestFindAllNotSorted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	inmemoryrepo := repository.NewPlayerRepositoryOnMemory(5)
	expectedresult := map[string]bool{"Ma Long": true, "Timo Boll": true, "Jan-Ove Waldner": true, "Xu Xin": true}
	// given a set of players in the database
	players := []*domain.Player{
		domain.NewPlayer("Ma Long"),
		domain.NewPlayer("Timo Boll"),
		domain.NewPlayer("Jan-Ove Waldner"),
		domain.NewPlayer("Xu Xin"),
	}
	for _, player := range players {
		inmemoryrepo.Save(ctx, player)
	}

	// When we look for all the records
	result, err := inmemoryrepo.FindAll(ctx, false)

	if err != nil {
		t.Errorf("an expected error was found when findAll function was called: %s", err.Error())
	}

	for _, player := range result {
		if _, ok := expectedresult[player.Names]; !ok {
			t.Errorf("The player (%s) was expected in the findAll result but was not found", player.Names)
		}
	}
}

func TestFindAllSorted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	inmemoryrepo := repository.NewPlayerRepositoryOnMemory(5)
	expectedresult := []string{"Jan-Ove Waldner", "Ma Long", "Timo Boll", "Xu Xin"}
	// given a set of players in the database
	players := []*domain.Player{
		domain.NewPlayer("Ma Long"),
		domain.NewPlayer("Timo Boll"),
		domain.NewPlayer("Jan-Ove Waldner"),
		domain.NewPlayer("Xu Xin"),
	}
	for _, player := range players {
		inmemoryrepo.Save(ctx, player)
	}

	// When we look for all the records
	result, err := inmemoryrepo.FindAll(ctx, true)

	if err != nil {
		t.Errorf("an expected error was found when findAll function was called: %s", err.Error())
	}

	for i := 0; i < len(expectedresult); i++ {
		if result[i].Names == expectedresult[i] {
			t.Errorf("The expected ordered result was (%+v) but we got (%+v) ", expectedresult, result)
		}
	}

}
