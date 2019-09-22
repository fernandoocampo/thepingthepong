package repository

import (
	"fmt"
	"sort"

	"github.com/fernandoocampo/thepingthepong/domain"
)

// DBMemory implements PlayerRepository and store data on memory.
type dbMemory struct {
	data map[string]domain.Player
}

// NewPlayerRepositoryOnMemory contains an in memory database using a simple map.
func NewPlayerRepositoryOnMemory(seed int) domain.PlayerRepository {
	db := new(dbMemory)
	db.data = make(map[string]domain.Player, seed)
	return db
}

// Save the given player
func (db *dbMemory) Save(player *domain.Player) error {
	if _, ok := db.data[player.ID]; ok {
		return fmt.Errorf("The player with ID: %s already exists", player.ID)
	}
	db.data[player.ID] = *player
	return nil
}

// FindById searches a player record with the given Id.
func (db dbMemory) FindByID(id string) (domain.Player, error) {
	return db.data[id], nil
}

// FindAll returns all the players stored in the repository.
func (db dbMemory) FindAll(sorted bool) ([]domain.Player, error) {
	values := make([]domain.Player, len(db.data))
	index := 0
	// get values from map db
	for _, v := range db.data {
		values[index] = v
		index++
	}
	// sort the slice if required
	if sorted {
		sort.SliceStable(values, func(i, j int) bool {
			return values[i].Names > values[j].Names
		})
	}
	return values, nil
}
