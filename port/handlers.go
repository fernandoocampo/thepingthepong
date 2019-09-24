package port

import "net/http"

// RestHandler Defines behavior for any kind of handler in a REST mode.
type RestHandler interface {
	// GetAll get all records or those that matches a given criteria
	GetAll(w http.ResponseWriter, r *http.Request)
	// GetByID get record by id
	GetByID(w http.ResponseWriter, r *http.Request)
	// Create creates a new record
	Create(w http.ResponseWriter, r *http.Request)
	// Update updates the data of existing record.
	Update(w http.ResponseWriter, r *http.Request)
	// Delete deletes a record.
	Delete(w http.ResponseWriter, r *http.Request)
	// Health returns the health of this service
	Health(w http.ResponseWriter, r *http.Request)
}

// AuthHandler Defines behavior for authentication and authorization in REST mode.
type AuthHandler interface {
	// SignIn authenticates an user
	SignIn(w http.ResponseWriter, r *http.Request)
}
