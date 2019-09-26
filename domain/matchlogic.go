package domain

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	// PlayerHitSentence sets narrative when a player hit a ball
	PlayerHitSentence = "%q hit the ball"
	// PlayerWonSentence sets narrative when a player wins a match
	PlayerWonSentence = "Player %q won"
	// PlayerFailSentence sets narrative when a player fail a ball
	PlayerFailSentence = "%q fail the ball"
	// FatalNumber is the number that if a referee randon number is equals to it, that
	// player loses the match
	FatalNumber = 13
)

// referee defines the winner of the match randomly
var referee *rand.Rand

// MatchReport models a report of a match played between two ping pong players
type MatchReport struct {
	ID        Key       `json:"id,omitempty"`     // internal id
	Narrative []string  `json:"narrative"`        // match narrative
	Winner    *Player   `json:"winner,omitempty"` // player who wins
	Loser     *Player   `json:"loser,omitempty"`  // player who loses
	Created   time.Time `json:"created"`          // The creation date
}

func init() {
	referee = createReferee()
}

// NewMatchReport creates a new match report with a ID and Created date
func NewMatchReport() *MatchReport {
	return &MatchReport{
		ID:      GenerateUUIDKey(),
		Created: time.Now(),
	}
}

// NewMatchReportWithData creates a new match report from their players
func NewMatchReportWithData(winner, loser Player, narrative []string) *MatchReport {
	log.Debugf("creating a match report with winner: %v and loser: %v", winner, loser)
	return &MatchReport{
		ID:      GenerateUUIDKey(),
		Winner:  &winner,
		Loser:   &loser,
		Created: time.Now(),
	}
}

// SimulateMatch simulates a ping pong match between player1 and player2
func SimulateMatch(player1, player2 Player) *MatchReport {
	match := NewMatchReport()
	table := make(chan int)
	narrative := make(chan string, 2)
	player1Won := make(chan bool)
	player2Won := make(chan bool)
	finishNarrative := make(chan bool)
	go player1.move(narrative, table, player1Won)
	go player2.move(narrative, table, player2Won)
	go match.addSentenceToNarrative(narrative, finishNarrative)
	table <- 1
	select {
	case <-player1Won:
		match.setWinnerAndLoser(&player1, &player2)
	case <-player2Won:
		match.setWinnerAndLoser(&player2, &player1)
	}
	<-finishNarrative
	if log.LevelLabel == "debug" {
		for i, val := range match.Narrative {
			fmt.Printf("%d - %s\n", i, val)
		}
	}
	return match
}

func (m *MatchReport) setWinnerAndLoser(winner, losser *Player) {
	m.Winner = winner
	m.Loser = losser
}

// addSentenceToNarrative adds sentences about the narrative of the match
func (m *MatchReport) addSentenceToNarrative(sentences chan string, finish chan<- bool) {
	for sentence := range sentences {
		m.Narrative = append(m.Narrative, sentence)
	}
	finish <- true
}

// move defines a player behavior regarding to a match, here the match is narrated
// and identifies if the player wins or loses the game.
func (p Player) move(narrative chan<- string, table chan int, winner chan bool) {
	for {
		ball, ok := <-table
		if !ok {
			// if the channel is closed, we win
			narrative <- fmt.Sprintf(PlayerWonSentence, p.Names)
			close(narrative)
			winner <- true
			return
		}
		if referee.Intn(100) == FatalNumber {
			narrative <- fmt.Sprintf(PlayerFailSentence, p.Names)
			close(table)
			return
		}
		narrative <- fmt.Sprintf(PlayerHitSentence, p.Names)
		ball++
		table <- ball
	}
}

// createReferee creates a referee that is an int random generator
// to compare with the FatalNumbere
func createReferee() *rand.Rand {
	sourceForRandom := rand.NewSource(time.Now().UnixNano())
	return rand.New(sourceForRandom)
}
