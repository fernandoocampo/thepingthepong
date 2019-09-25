package port

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fernandoocampo/thepingthepong/application/authapp"
)

// Credentials Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type basicAuthRestHandler struct {
	service authapp.Authenticator
}

// NewBasicAuthRestHandler instance of a basic implementation of auth rest handler
func NewBasicAuthRestHandler(authService authapp.Authenticator) AuthHandler {
	log.Infof("creating auth rest handler")
	return basicAuthRestHandler{
		service: authService,
	}
}

// // SignIn authenticates an user
func (b basicAuthRestHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// context constraint
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	ok, err := b.service.Authenticate(ctx, creds.Username, creds.Password)

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	timeToExpire := 5 * time.Minute
	token, err := generateToken(creds, timeToExpire)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token.token,
		Expires: token.expirationTime,
	})
}
