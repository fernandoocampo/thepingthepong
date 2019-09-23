package domain

// PlayerRepository defines standard behavior
type PlayerRepository interface {
	// Save the given player
	Save(player *Player) error
	// FindById searches a player record with the given Id.
	FindByID(id string) (Player, error)
	// FindAll returns all the players stored in the repository.
	FindAll(sorted bool) ([]Player, error)
}
