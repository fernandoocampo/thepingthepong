package playerapp_test

import (
	"context"
	"testing"

	"github.com/fernandoocampo/thepingthepong/application/playerapp"
	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/fernandoocampo/thepingthepong/infra/repository"
)

func TestSaveValidPlayer(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewPlayerRepositoryOnMemory(1)
	service := playerapp.NewBasicPlayerService(&repo)
	// Given a player to save
	newplayername := "Ma Lin"
	newplayerwins := 4
	newplayerlosses := 1

	// When we want to store the new player
	newid, err := service.Create(ctx, newplayername, newplayerwins, newplayerlosses)

	// then we check if there is not any error
	if err != nil {
		t.Errorf("the player %s with wins %d and losses %d could not be created because: %s",
			newplayername, newplayerwins, newplayerlosses, err)
	}

	storedplayer, errfindbyid := repo.FindByID(ctx, newid)

	if errfindbyid != nil {
		t.Errorf("the player with id %s could not be searched because: %s",
			newid, errfindbyid.Error())
	}

	if storedplayer.Names != newplayername {
		t.Errorf("The player with ID %s has the name %s and should be %s", newid, storedplayer.Names, newplayername)
	}

	if storedplayer.Wins != newplayerwins {
		t.Errorf("The player with ID %s has the wins value: %d and should be %d", newid, storedplayer.Wins, newplayerwins)
	}

	if storedplayer.Losses != newplayerlosses {
		t.Errorf("The player with ID %s has the losses value: %d and should be %d", newid, storedplayer.Losses, newplayerlosses)
	}
}

func TestFindByID(t *testing.T) {
	ctx := context.Background()
	playertosave := domain.NewPlayer("Wang Hao")
	repo := repository.NewPlayerRepositoryOnMemory(1)
	repo.Save(ctx, playertosave)
	service := playerapp.NewBasicPlayerService(&repo)
	// Given an id to find
	playerID := playertosave.ID

	// When we want to find the player using the given id
	storedplayer, err := service.FindByID(ctx, playerID)

	// Then we check that there is not an error
	if err != nil {
		t.Errorf("The player with ID %s could be searched because: %s", playerID, err.Error())
	}

	if storedplayer.Names != playertosave.Names ||
		storedplayer.Wins != playertosave.Wins ||
		storedplayer.Losses != playertosave.Losses {
		t.Errorf(`the stored player should have the name %s, wins %d and losses %d
			 but have name %s, wins %d and losses %d`, playertosave.Names,
			playertosave.Wins, playertosave.Losses, storedplayer.Names, storedplayer.Wins,
			storedplayer.Losses)
	}

}

func TestFindAll(t *testing.T) {
	ctx := context.Background()
	// Given this repo and service
	repo := repository.NewPlayerRepositoryOnMemory(5)
	expectedresult := map[string]bool{"Ma Long": true, "Timo Boll": true, "Jan-Ove Waldner": true, "Xu Xin": true}
	players := []*domain.Player{
		domain.NewPlayer("Ma Long"),
		domain.NewPlayer("Timo Boll"),
		domain.NewPlayer("Jan-Ove Waldner"),
		domain.NewPlayer("Xu Xin"),
	}

	for _, player := range players {
		repo.Save(ctx, player)
	}
	service := playerapp.NewBasicPlayerService(&repo)

	// When we want to find the player using the given id
	result, err := service.FindAll(ctx, false)

	// Then we check that there is not an error
	if err != nil {
		t.Errorf("The players could be searched because: %s", err.Error())
	}

	for _, player := range result {
		if _, ok := expectedresult[player.Names]; !ok {
			t.Errorf("The player (%s) was expected in the findAll result but was not found", player.Names)
		}
	}

}
