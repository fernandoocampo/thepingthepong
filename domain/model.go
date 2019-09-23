package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Player models the ping pong player.
type Player struct {
	ID      string    `json:"id,omitempty"`    // internal id
	Names   string    `json:"names,omitempty"` // player names
	Wins    int       `json:"wins"`            // the number of wins of this player
	Losses  int       `json:"losses"`          // the number of losses of this player
	Created time.Time `json:"created"`         // The creation date
	Updated time.Time `json:"updated"`         // the update date
}

// NewPlayer creates a new player with a random uuid ID.
func NewPlayer(names string) *Player {
	return NewPlayerWithStatistics(names, 0, 0)
}

// NewPlayerWithStatistics creates a new player with a random uuid ID.
func NewPlayerWithStatistics(names string, wins, losses int) *Player {
	return &Player{
		ID:      uuid.New().String(),
		Names:   names,
		Wins:    wins,
		Losses:  losses,
		Created: time.Now(),
		Updated: time.Now(),
	}
}

// ValidatePlayer checks that the given player has not empty names and wins and losses
// are not negative.
func ValidatePlayer(player Player) (bool, error) {
	var result []string
	// check for empty parameter
	if &player == nil {
		return false, errors.New("There in not data in the player parameter")
	}
	// check for a valid names
	if player.Names == "" {
		result = append(result, "Player names cannot be empty")
	}
	// check for a valid names value without just spaces
	if player.Names != "" && strings.TrimSpace(player.Names) == "" {
		result = append(result, "Player names cannot contain only spaces")
	}
	// check that wins value cannot be negative
	if player.Wins < 0 {
		result = append(result, "Player wins cannot be less than zero")
	}
	// check that wins value cannot be negative
	if player.Losses < 0 {
		result = append(result, "Player losses cannot be less than zero")
	}

	if len(result) > 0 {
		return false, errors.New(strings.Join(result, "\n"))
	}
	return true, nil
}
