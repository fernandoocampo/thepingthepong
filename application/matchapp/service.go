package matchapp

import (
	"context"

	"github.com/fernandoocampo/thepingthepong/application/playerapp"
	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/pkg/errors"
)

// MatchService defines contract to execute a ping pong match
type MatchService interface {
	// Play simulates a match between player1 and player2 and returns a narrative about the event.
	Play(ctx context.Context, player1ID, player2ID domain.Key) (*domain.MatchReport, error)
}

// basicMatchService implements the Match service.
type basicMatchService struct {
	playerService playerapp.PlayerService
}

// NewBasicMatchService build a basic implementation for matchservice.
func NewBasicMatchService(playerService playerapp.PlayerService) MatchService {
	log.Info("creating basic player service")
	return &basicMatchService{
		playerService: playerService,
	}
}

// Play simulates a match between player1 and player2 and returns a narrative about the event.
func (b *basicMatchService) Play(ctx context.Context, player1ID, player2ID domain.Key) (*domain.MatchReport, error) {
	log.Infof("the match between %q and %q has began", player1ID, player2ID)
	log.Infof("finding player with id: %q", player1ID)
	player1, err := b.playerService.FindByID(ctx, player1ID)
	if err != nil { // just the logs
		log.Errorf("player 1: %s cannot be found because: %s", player1ID, err.Error())
		return nil, errors.Wrap(err, "player 1 not found at the match")
	}
	log.Infof("finding player with id: %q", player2ID)
	player2, err := b.playerService.FindByID(ctx, player2ID)
	if err != nil { // just the logs
		log.Errorf("player 2: %s cannot be found because: %s", player2ID, err.Error())
		return nil, errors.Wrap(err, "player 2 not found at the match")
	}
	match := domain.SimulateMatch(player1, player2)
	stats := playerapp.NewPlayerStatistics(match.Winner.ID, match.Loser.ID, 1, 1)
	err = b.playerService.UpdateStatistics(ctx, *stats)
	if err != nil { // just the logs
		log.Errorf("player statistics: %v cannot be updatedbecause: %s", stats, err.Error())
	}
	return match, nil
}
