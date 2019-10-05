package playerapp

import (
	"context"
	"fmt"

	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/pkg/errors"
)

// PlayerStatistics groups all the statistics for winner and loser
type PlayerStatistics struct {
	WinnerID, LoserID domain.Key
	Wins, Losses      int
}

// PlayerService defines standard behavior for player capabilities.
type PlayerService interface {
	// Create creates a player with the given data and return id or and error
	Create(ctx context.Context, names string, wins, losses int) (domain.Key, error)
	// FindByID finds a player by id
	FindByID(ctx context.Context, key domain.Key) (domain.Player, error)
	// FindAll get all the players
	FindAll(ctx context.Context, sorted bool) ([]domain.Player, error)
	// UpdateStatistics updates the winner and loser counter for winner and loser players
	UpdateStatistics(ctx context.Context, statistics PlayerStatistics) error
}

// NewPlayerStatistics builds a stats data.
func NewPlayerStatistics(winnerID, loserID domain.Key, wins, losses int) *PlayerStatistics {
	return &PlayerStatistics{
		WinnerID: winnerID,
		LoserID:  loserID,
		Wins:     wins,
		Losses:   losses,
	}
}

// basicPlayerService implements the player service.
type basicPlayerService struct {
	repository domain.PlayerRepository
}

// NewBasicPlayerService build a basic implementation for playerservice.
func NewBasicPlayerService(repository *domain.PlayerRepository) PlayerService {
	log.Info("creating basic player service")
	return &basicPlayerService{
		repository: *repository,
	}
}

// Create creates a player
func (b basicPlayerService) Create(ctx context.Context, names string, wins, losses int) (domain.Key, error) {
	log.Infof("creating player with names: '%s', wins: %d, losses: %d", names, wins, losses)
	// check that the given parameter is valid
	player := domain.NewPlayerWithStatistics(names, wins, losses)
	ok, errvalidation := domain.ValidatePlayer(*player)
	if !ok {
		log.Infof("Player %v is not valid, returning from service.", player)
		return "", errvalidation
	}
	log.Infof("getting ready to save player %v on repository", player)
	errsave := b.repository.Save(ctx, player)
	if errsave != nil {
		log.Errorf("player %v cannot be stored because: %s", player, errsave.Error())
		return "", errors.Wrap(errsave, "Player cannot be stored")
	}
	log.Infof("player stored with ID: %s", player.ID)
	return player.ID, nil
}

// FindByID finds a player by id
func (b basicPlayerService) FindByID(ctx context.Context, id domain.Key) (domain.Player, error) {
	log.Infof("finding player with id: %s", id)
	if id == "" { // nothig to search
		log.Infof("provided id was empty... returning an empty player")
		return domain.Player{}, nil
	}

	log.Infof("getting ready to find the player with id: %s on repository", id)
	result, err := b.repository.FindByID(ctx, id)

	if err != nil {
		log.Errorf("something was going wrong searching player with id: %s, because: %s", id, err.Error())
		return domain.Player{}, errors.Wrap(err, fmt.Sprintf("player with id %s could not be searched", id))
	}

	log.Infof("returning found player %v with id %s", result, id)

	return result, nil
}

// FindAll get all the players
func (b basicPlayerService) FindAll(ctx context.Context, sorted bool) ([]domain.Player, error) {
	log.Infof("getting ready to find all players with sorted: %t", sorted)
	result, err := b.repository.FindAll(ctx, sorted)
	log.Debugf("after finding all players, got %+v", result)
	if err != nil {
		log.Errorf("something goes wrong trying to find all players: %s", err.Error())
		return nil, errors.Wrap(err, "all players could not be searched")
	}

	return result, nil
}

// UpdateStatistics updates the winner and loser counter for winner and loser players
func (b basicPlayerService) UpdateStatistics(ctx context.Context, stats PlayerStatistics) error {
	log.Infof("getting ready to update statistics for players: %v", stats)
	err := b.repository.UpdateWins(ctx, stats.WinnerID, stats.Wins)
	if err != nil {
		log.Errorf("player %s cannot be update wins because: %s", stats.WinnerID, err.Error())
		return errors.Wrap(err, "winner player could not be updated")
	}
	err = b.repository.UpdateDefeats(ctx, stats.LoserID, stats.Losses)
	if err != nil {
		log.Errorf("player %s cannot be update wins because: %s", stats.LoserID, err.Error())
		return errors.Wrap(err, "loser player could not be updated")
	}
	return nil
}
