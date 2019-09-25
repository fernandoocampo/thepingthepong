package port

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key") //TODO configurable

// Claims Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// ValidToken contains result data after token validation.
type ValidToken struct {
	StatusCode int
	Claims     *Claims
}

// Token contains the token value and its expiration time
type Token struct {
	token          string
	expirationTime time.Time
}

// generateToken receive credentials and generate a token
func generateToken(creds Credentials, validTime time.Duration) (*Token, error) {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(validTime)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Errorf("Token could not be generated, because: %s", err.Error())
		return nil, err
	}
	return &Token{token: tokenString, expirationTime: expirationTime}, nil
}

func headerToken(r *http.Request) (string, bool) {
	var result string
	// check if it is in the header
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) >= 2 {
		result = splitToken[1]
		return result, true
	}
	return "", false
}

func cookieToken(r *http.Request) (string, bool) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", false
	}

	// Get the JWT string from the cookie
	return c.Value, true
}

func token(r *http.Request) (string, bool) {
	token, ok := cookieToken(r)
	if !ok {
		return headerToken(r)
	}
	return token, ok
}

func validateToken(r *http.Request) (ValidToken, bool) {
	// We can obtain the session token from the requests cookies, which come with every request
	tknStr, isok := token(r)
	if !isok {
		// If the cookie is not set, return an unauthorized status
		return ValidToken{StatusCode: http.StatusUnauthorized, Claims: nil}, false
	}

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return ValidToken{StatusCode: http.StatusUnauthorized, Claims: nil}, false
		}
		return ValidToken{StatusCode: http.StatusBadRequest, Claims: nil}, false
	}
	if !tkn.Valid {
		return ValidToken{StatusCode: http.StatusUnauthorized, Claims: nil}, false
	}
	return ValidToken{StatusCode: http.StatusOK, Claims: claims}, true
}
