package domain_test

import (
	"errors"
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
