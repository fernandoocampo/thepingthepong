package authapp

import "context"

// Authenticator defines standard behavior for player capabilities.
type Authenticator interface {
	// Authenticate validate if the given user is who he says he is. Returns error if there is
	// something goes wrong and the bool value to say that it was authenticated or not.
	Authenticate(ctx context.Context, username, password string) (bool, error)
}

type basicAuthenticator struct {
	users map[string]string
}

// NewBasicAuthenticator creates a new basic authenticator
func NewBasicAuthenticator() Authenticator {
	return basicAuthenticator{
		users: map[string]string{
			"user1": "password1",
			"user2": "password2",
		},
	}
}

// Authenticate validate if the given user is who he says he is. Returns error if there is
// something goes wrong and the bool value to say that it was authenticated or not.
func (b basicAuthenticator) Authenticate(ctx context.Context, username, password string) (bool, error) {
	// Get the expected password from our in memory map
	expectedPassword, ok := b.users[username]
	if !ok || expectedPassword != password {
		return false, nil
	}
	return true, nil
}
