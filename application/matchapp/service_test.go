package matchapp_test

import (
	"context"
	"strings"
	"testing"

	"github.com/fernandoocampo/thepingthepong/application/matchapp"
	"github.com/fernandoocampo/thepingthepong/application/playerapp"
	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/fernandoocampo/thepingthepong/infra/repository"
)

func TestBasicPlay(t *testing.T) {
	player1Names := "Ma Long"
	player1InitialWins := 10
	player1InitialLoses := 13
	player2Names := "Xu Xin"
	player2InitialWins := 20
	player2InitialLoses := 5
	repo := repository.NewPlayerRepositoryOnMemory(10)
	playerService := playerapp.NewBasicPlayerService(&repo)
	ctx := context.TODO()
	player1ID, err := playerService.Create(ctx, player1Names, player1InitialWins, player1InitialLoses)
	assertNoError(t, err)
	player2ID, err := playerService.Create(ctx, player2Names, player2InitialWins, player2InitialLoses)
	assertNoError(t, err)

	basicMatchService := matchapp.NewBasicMatchService(playerService)

	got, err := basicMatchService.Play(ctx, player1ID, player2ID)
	assertNoError(t, err)

	if got.Winner == nil {
		t.Errorf("a winner between player: %q and player: %q was expected, but none won",
			player1Names, player2Names)
	}
	if got.Loser == nil {
		t.Errorf("a loser between player: %q and player: %q was expected, but none lost",
			player1Names, player2Names)
	}
	if len(got.Narrative) == 0 {
		t.Errorf("a fulled narrative was expected, but got: %v", got.Narrative)
	}
	for index, sentence := range got.Narrative {
		if sentence == "" || strings.TrimSpace(sentence) == "" {
			t.Errorf("each sentence in the match must contains some text, but sentence: %d was empty", index+1)
		}
	}
	winner, err := repo.FindByID(ctx, got.Winner.ID)
	assertNoError(t, err)
	loser, err := repo.FindByID(ctx, got.Loser.ID)
	assertNoError(t, err)

	assertAfterMatchPlayerState(t, player1InitialWins, player1InitialLoses, player1ID, winner, loser)
	assertAfterMatchPlayerState(t, player2InitialWins, player2InitialLoses, player2ID, winner, loser)

}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("error was not expected, but: %s", err)
	}
}

func assertAfterMatchPlayerState(t *testing.T, initialWins, initiaLosse int, playerID domain.Key, winner, loser domain.Player) {
	t.Helper()
	if loser.ID == playerID {
		if loser.Wins != initialWins {
			t.Errorf("loser %q had %d wins and must have %d wins", loser.Names, loser.Wins, initialWins)
		}
		if loser.Losses <= initiaLosse {
			t.Errorf("loser %q had %d losses and must have %d losses", loser.Names, loser.Losses, loser.Losses+1)
		}
	}
	if winner.ID == playerID {
		if winner.Losses != initiaLosse {
			t.Errorf("winner %q had %d losses and must have %d losses", winner.Names, winner.Losses, initiaLosse)
		}
		if winner.Wins <= initialWins {
			t.Errorf("winner %q had %d wins and must have %d wins", winner.Names, winner.Wins, initialWins+1)
		}
	}
}
