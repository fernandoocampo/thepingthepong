package domain_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/fernandoocampo/thepingthepong/domain"
)

type data struct {
	param  domain.Player
	result bool
	err    error
}

var givenData = []data{
	{
		param: domain.Player{
			ID:      "sfsfs-sf-sfs-sfsf",
			Names:   "Rafael Nadal",
			Wins:    0,
			Losses:  0,
			Created: time.Now(),
			Updated: time.Now(),
		},
		result: true,
		err:    nil,
	},
	{
		param: domain.Player{
			ID:      "sfsfs-sf-sfs-sfsf",
			Names:   "Rafael Nadal",
			Wins:    -1,
			Losses:  0,
			Created: time.Now(),
			Updated: time.Now(),
		},
		result: false,
		err:    errors.New("Player wins cannot be less than zero"),
	},
	{
		param: domain.Player{
			ID:      "sfsfs-sf-sfs-sfsf",
			Names:   "Rafael Nadal",
			Wins:    0,
			Losses:  -2,
			Created: time.Now(),
			Updated: time.Now(),
		},
		result: false,
		err:    errors.New("Player losses cannot be less than zero"),
	},
	{
		param: domain.Player{
			ID:      "sfssf-2342-sdfs-sdsds-sfs",
			Names:   "     ",
			Wins:    0,
			Losses:  0,
			Created: time.Now(),
			Updated: time.Now(),
		},
		result: false,
		err:    errors.New("Player names cannot contain only spaces"),
	},
	{
		param: domain.Player{
			ID:      "sfssf-2342-sdfs-sdsds-sfs",
			Names:   "     ",
			Wins:    0,
			Losses:  -3,
			Created: time.Now(),
			Updated: time.Now(),
		},
		result: false,
		err:    errors.New("Player names cannot contain only spaces\nPlayer losses cannot be less than zero"),
	},
	{
		param: domain.Player{
			ID:      "sfssf-2342-sdfs-sdsds-sfs",
			Names:   "     ",
			Wins:    -1,
			Losses:  0,
			Created: time.Now(),
			Updated: time.Now(),
		},
		result: false,
		err:    errors.New("Player names cannot contain only spaces\nPlayer wins cannot be less than zero"),
	},
	{
		param: domain.Player{
			ID:      "sfssf-2342-sdfs-sdsds-sfs",
			Names:   "     ",
			Wins:    -1,
			Losses:  -1,
			Created: time.Now(),
			Updated: time.Now(),
		},
		result: false,
		err:    errors.New("Player names cannot contain only spaces\nPlayer wins cannot be less than zero\nPlayer losses cannot be less than zero"),
	},
	{
		param: domain.Player{
			ID:      "sfssf-2342-sdfs-fssdsd-sfssds",
			Names:   "",
			Wins:    0,
			Losses:  0,
			Created: time.Now(),
			Updated: time.Now(),
		},
		result: false,
		err:    errors.New("Player names cannot be empty"),
	},
}

func TestValidatePlayer(t *testing.T) {
	for _, v := range givenData {
		ok, err := domain.ValidatePlayer(v.param)
		if ok != v.result {
			t.Errorf("Given parameter [%v], it expects [%t], but it got [%t]", v.param, v.result, ok)
		}
		if ok == false && ok == v.result {
			if err.Error() != v.err.Error() {
				t.Errorf("Given parameter [%v], it expects [%s], but it got [%s]", v.param, v.err.Error(), err.Error())
			}
		}
	}
}

func TestSimulateMatchh(t *testing.T) {
	player1 := domain.NewPlayer("Wang Hao")
	player2 := domain.NewPlayer("Zhang Jike")

	got := domain.SimulateMatch(*player1, *player2)

	if len(got.Narrative) == 0 {
		t.Errorf("a fulled narrative was expected, but got: %v", got.Narrative)
	}
	for index, sentence := range got.Narrative {
		if sentence == "" || strings.TrimSpace(sentence) == "" {
			t.Errorf("each sentence in the match must contains some text, but sentence: %d was empty", index+1)
		}
	}
	if got.ID == "" {
		t.Errorf("the match must contain an ID but it was empty")
	}
	if got.Winner == nil {
		t.Errorf("a winner between player: %q and player: %q was expected, but none won",
			player1.Names, player2.Names)
	}
	if got.Loser == nil {
		t.Errorf("a loser between player: %q and player: %q was expected, but none lost",
			player1.Names, player2.Names)
	}
}
