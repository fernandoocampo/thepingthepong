package domain

import "context"

// PlayerRepository defines standard behavior
type PlayerRepository interface {
	// Save the given player
	Save(ctx context.Context, player *Player) error
	// FindById searches a player record with the given Id.
	FindByID(ctx context.Context, id string) (Player, error)
	// FindAll returns all the players stored in the repository.
	FindAll(ctx context.Context, sorted bool) ([]Player, error)
}
