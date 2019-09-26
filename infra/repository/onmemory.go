package repository

import (
	"context"
	"fmt"
	"sort"

	"github.com/fernandoocampo/thepingthepong/domain"
	"github.com/pkg/errors"
)

// DBMemory implements PlayerRepository and store data on memory.
type dbMemory struct {
	data map[domain.Key]domain.Player
}

// NewPlayerRepositoryOnMemory contains an in memory database using a simple map.
func NewPlayerRepositoryOnMemory(seed int) domain.PlayerRepository {
	log.Infof("creating on memory map repository for players with seed: %d", seed)
	db := new(dbMemory)
	db.data = make(map[domain.Key]domain.Player, seed)
	return db
}

// Save the given player
func (db *dbMemory) Save(ctx context.Context, player *domain.Player) error {
	log.Infof("receiven player: %v to store", player)
	var iserror error
	chanresult := make(chan error)

	go func() {
		if _, ok := db.data[player.ID]; ok {
			log.Errorf("record with id: %s already exists on db", player.ID)
			chanresult <- fmt.Errorf("The player with ID: %s already exists", player.ID)
			return
		}
		db.data[player.ID] = *player
		log.Infof("saving player: %v on database", player)
		chanresult <- nil
	}()
	select {
	case <-ctx.Done():
		log.Errorf("Operation take a long to time to finish: %s", ctx.Err())
		return errors.Wrap(ctx.Err(), "Could not finish save operation at time")
	case iserror = <-chanresult:
		return iserror
	}
}

// FindById searches a player record with the given Id.
func (db dbMemory) FindByID(ctx context.Context, id domain.Key) (domain.Player, error) {
	log.Infof("looking for player with id: %s", id)
	var result domain.Player
	resultchan := make(chan domain.Player)
	go func() {
		resultchan <- db.data[id]
	}()
	select {
	case <-ctx.Done():
		log.Errorf("Operation take a long to time to finish: %s", ctx.Err())
		return domain.Player{}, errors.Wrap(ctx.Err(), "Could not finish the find by id at time")
	case result = <-resultchan:
		log.Infof("player was found on repository: %v", result)
	}
	return result, nil
}

// FindAll returns all the players stored in the repository.
func (db dbMemory) FindAll(ctx context.Context, sorted bool) ([]domain.Player, error) {
	log.Infof("finding all players with sorted: %t", sorted)
	var result []domain.Player
	resultchan := make(chan []domain.Player)

	go func() {
		values := make([]domain.Player, len(db.data))
		index := 0
		// get values from map db
		for _, v := range db.data {
			values[index] = v
			index++
		}
		log.Debugf("All players without sorted are: %+v", values)
		// sort the slice if required
		if sorted {
			sort.SliceStable(values, func(i, j int) bool {
				return values[i].Names > values[j].Names
			})
			log.Debugf("All players sorted are: %+v", values)
		}
		resultchan <- values
	}()
	select {
	case <-ctx.Done():
		log.Errorf("Operation take a long to time to finish: %s", ctx.Err())
		return nil, errors.Wrap(ctx.Err(), "Could not finish the findAll at time")
	case result = <-resultchan:
		log.Infof("player was found on repository: %v", result)
	}
	return result, nil
}

// UpdateWins increases the value on field wins
func (db dbMemory) UpdateWins(ctx context.Context, playerID domain.Key, wins int) error {
	log.Infof("looking for player with id: %s", playerID)
	resultchan := make(chan bool)
	go func() {
		player := db.data[playerID]
		player.Wins += wins
		db.data[playerID] = player
		resultchan <- true
	}()
	select {
	case <-ctx.Done():
		log.Errorf("Operation take a long to time to finish: %s", ctx.Err())
		return errors.Wrap(ctx.Err(), "Could not finish the find by id at time")
	case <-resultchan:
		log.Infof("player %q was updated on repository", playerID)
	}
	return nil
}

// UpdateDefeats increases the value on field loses
func (db dbMemory) UpdateDefeats(ctx context.Context, playerID domain.Key, defeats int) error {
	log.Infof("looking for player with id: %s", playerID)
	resultchan := make(chan bool)
	go func() {
		player := db.data[playerID]
		player.Losses += defeats
		db.data[playerID] = player
		resultchan <- true
	}()
	select {
	case <-ctx.Done():
		log.Errorf("Operation take a long to time to finish: %s", ctx.Err())
		return errors.Wrap(ctx.Err(), "Could not finish the find by id at time")
	case <-resultchan:
		log.Infof("player %q was updated on repository", playerID)
	}
	return nil
}
