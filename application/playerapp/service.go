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
	return basicPlayerService{
		repository: repository,
	}
}

// Create creates a player
func (b basicPlayerService) Create(names string, wins, losses int) (string, error) {
	// check that the given parameter is valid
	player := domain.NewPlayerWithStatistics(names, wins, losses)
	ok, errvalidation := domain.ValidatePlayer(*player)
	if !ok {
		return "", errvalidation
	}
	errsave := b.repository.Save(player)
	if errsave != nil {
		return "", errors.Wrap(errsave, "Player cannot be stored")
	}
	return player.ID, nil
}

// FindByID finds a player by id
func (b basicPlayerService) FindByID(id string) (domain.Player, error) {
	if id == "" { // nothig to search
		return domain.Player{}, nil
	}
	result, err := b.repository.FindByID(id)

	if err != nil {
		return domain.Player{}, errors.Wrap(err, fmt.Sprintf("player with id %s could not be searched", id))
	}

	return result, nil
}

// FindAll get all the players
func (b basicPlayerService) FindAll(sorted bool) ([]domain.Player, error) {
	result, err := b.repository.FindAll(sorted)

	if err != nil {
		return nil, errors.Wrap(err, "all players could not be searched")
	}

	return result, nil
}
