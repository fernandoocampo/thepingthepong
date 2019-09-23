package playerapp

import (
	"fmt"

	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/pkg/errors"
)

// PlayerService defines standard behavior for player capabilities.
type PlayerService interface {
	// Create creates a player with the given data and return id or and error
	Create(names string, wins, losses int) (string, error)
	// FindByID finds a player by id
	FindByID(id string) (domain.Player, error)
	// FindAll get all the players
	FindAll(sorted bool) ([]domain.Player, error)
}

// basicPlayerService implements the player service.
type basicPlayerService struct {
	repository domain.PlayerRepository
}

// NewBasicPlayerService build a basic implementation for playerservice.
func NewBasicPlayerService(repository domain.PlayerRepository) PlayerService {
	log.Info("creating basic player service")
	return basicPlayerService{
		repository: repository,
	}
}

// Create creates a player
func (b basicPlayerService) Create(names string, wins, losses int) (string, error) {
	log.Infof("creating player with names: '%s', wins: %d, losses: %d", names, wins, losses)
	// check that the given parameter is valid
	player := domain.NewPlayerWithStatistics(names, wins, losses)
	ok, errvalidation := domain.ValidatePlayer(*player)
	if !ok {
		log.Infof("Player %v is not valid, returning from service.", player)
		return "", errvalidation
	}
	log.Infof("getting ready to save player %v on repository", player)
	errsave := b.repository.Save(player)
	if errsave != nil {
		log.Errorf("player %v cannot be stored because: %s", player, errsave.Error())
		return "", errors.Wrap(errsave, "Player cannot be stored")
	}
	log.Infof("player stored with ID: %s", player.ID)
	return player.ID, nil
}

// FindByID finds a player by id
func (b basicPlayerService) FindByID(id string) (domain.Player, error) {
	log.Infof("finding player with id: %s", id)
	if id == "" { // nothig to search
		log.Infof("provided id was empty... returning an empty player")
		return domain.Player{}, nil
	}

	log.Infof("getting ready to find the player with id: %s on repository", id)
	result, err := b.repository.FindByID(id)

	if err != nil {
		log.Errorf("something was going wrong searching player with id: %s, because: %s", id, err.Error())
		return domain.Player{}, errors.Wrap(err, fmt.Sprintf("player with id %s could not be searched", id))
	}

	log.Infof("returning found player %v with id %s", result, id)

	return result, nil
}

// FindAll get all the players
func (b basicPlayerService) FindAll(sorted bool) ([]domain.Player, error) {
	log.Infof("getting ready to find all players with sorted: %t", sorted)
	result, err := b.repository.FindAll(sorted)
	log.Debugf("after finding all players, got %+v", result)
	if err != nil {
		log.Errorf("something goes wrong trying to find all players: %s", err.Error())
		return nil, errors.Wrap(err, "all players could not be searched")
	}

	return result, nil
}
