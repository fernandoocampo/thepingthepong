package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Key is the primary key for every entity in the domain.
type Key string

// Player models the ping pong player.
type Player struct {
	ID      Key       `json:"id,omitempty"`    // internal id
	Names   string    `json:"names,omitempty"` // player names
	Wins    int       `json:"wins"`            // the number of wins of this player
	Losses  int       `json:"losses"`          // the number of losses of this player
	Created time.Time `json:"created"`         // The creation date
	Updated time.Time `json:"updated"`         // the update date
}

// GenerateUUIDKey generates a uuid key
func GenerateUUIDKey() Key {
	return Key(uuid.New().String())
}

// NewPlayer creates a new player with a random uuid ID.
func NewPlayer(names string) *Player {
	return NewPlayerWithStatistics(names, 0, 0)
}

// NewPlayerWithStatistics creates a new player with a random uuid ID.
func NewPlayerWithStatistics(names string, wins, losses int) *Player {
	log.Debugf("creating player with names: '%s', wins: %d, losses: %d", names, wins, losses)
	return &Player{
		ID:      GenerateUUIDKey(),
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
	log.Debugf("validating player %v", player)
	// check for empty parameter
	if &player == nil {
		return false, errors.New("There in not data in the player parameter")
	}
	// check for a valid names
	if player.Names == "" {
		log.Debugf("player has not valid names because it is empty")
		result = append(result, "Player names cannot be empty")
	}
	// check for a valid names value without just spaces
	if player.Names != "" && strings.TrimSpace(player.Names) == "" {
		log.Debugf("player has not valid names because it is just spaces")
		result = append(result, "Player names cannot contain only spaces")
	}
	// check that wins value cannot be negative
	if player.Wins < 0 {
		log.Debugf("player %s has not valid wins because it is negative: %d", player.Names, player.Wins)
		result = append(result, "Player wins cannot be less than zero")
	}
	// check that wins value cannot be negative
	if player.Losses < 0 {
		log.Debugf("player %s has not valid losses because it is negative: %d", player.Names, player.Losses)
		result = append(result, "Player losses cannot be less than zero")
	}

	if len(result) > 0 {
		strresult := strings.Join(result, "\n")
		log.Debugf("player %v has not valid data, because: %s \n", player, strresult)
		return false, errors.New(strresult)
	}
	return true, nil
}
