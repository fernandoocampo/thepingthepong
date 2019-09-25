package authapp_test

import (
	"context"
	"testing"

	"github.com/fernandoocampo/thepingthepong/application/authapp"
)

func TestAuthentication(t *testing.T) {
	ctx := context.Background()
	service := authapp.NewBasicAuthenticator()
	// GIven all these escenarios
	tests := []struct {
		usrname  string
		password string
		expected bool
		iserror  bool
	}{
		{"user1", "password1", true, false},
		{"user2", "password2", true, false},
		{"user3", "password3", false, false},
		{"user1", "password4", false, false},
	}
	// When an authentication is made
	for _, escenario := range tests {
		result, err := service.Authenticate(ctx, escenario.usrname, escenario.password)
		if result != escenario.expected {
			t.Errorf("expected that authenticate was %t but got %t", escenario.expected, result)
			if err != nil && escenario.iserror == false {
				t.Errorf("expected not error on authentication of user: %s with password: %s but got: %s",
					escenario.usrname, escenario.password, err.Error())
			}
		}
	}
}
