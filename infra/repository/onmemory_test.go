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

func TestUpdatePlayer(t *testing.T) {
	t.Run("increase wins to player", func(t *testing.T) {
		// given a new player
		repo := repository.NewPlayerRepositoryOnMemory(5)
		newplayer := domain.NewPlayer("Ma Long")
		saveAPlayer(t, repo, newplayer)
		// increases
		err := repo.UpdateWins(context.TODO(), newplayer.ID, 1)
		assertNoError(t, err)
		// query the new player and check if data was updated well
		savedplayer, errfind := repo.FindByID(context.TODO(), newplayer.ID)
		assertNoError(t, errfind)
		assertStatistics(t, "winner", savedplayer.Wins, newplayer.Wins, newplayer.Wins+1)
	})

	t.Run("increase defeats to player", func(t *testing.T) {
		// given a new player
		repo := repository.NewPlayerRepositoryOnMemory(5)
		newplayer := domain.NewPlayer("Ma Long")
		saveAPlayer(t, repo, newplayer)
		// increases
		err := repo.UpdateDefeats(context.TODO(), newplayer.ID, 1)
		assertNoError(t, err)
		// query the new player and check if data was updated well
		savedplayer, errfind := repo.FindByID(context.TODO(), newplayer.ID)
		assertNoError(t, errfind)
		assertStatistics(t, "loser", savedplayer.Losses, newplayer.Losses, newplayer.Losses+1)
	})

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

func saveAPlayer(t *testing.T, repo domain.PlayerRepository, newplayer *domain.Player) {
	t.Helper()
	ctx := context.TODO()

	// when we save the player in the inmemory db
	err := repo.Save(ctx, newplayer)

	// then we check that is not an error
	if err != nil {
		t.Fatalf("An attempt was made to save the player [%+v], it was expected non errors but we got %s",
			newplayer, err.Error())
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("error was not expected, but got: %s", err)
	}
}

func assertStatistics(t *testing.T, statistic string, got, oldvalue, want int) {
	t.Helper()
	if got <= oldvalue {
		t.Errorf(`The stored player %s counter is: %d, should be: %d`,
			statistic, got, want)
	}
}
