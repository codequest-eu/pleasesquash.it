package state

import "net/http"

// Store is a mechanism allowing temporary storage of:
// - OAuth2 state to prevent from CSRF attacks;
// - repo/user pair to store between requests;
type Store interface {
	// GetState allows getting state from an HTTP request.
	GetState(r *http.Request) (string, error)

	// SetState allows setting state on an HTTP response. Must be called
	// before the first byte of the response body is written.
	SetState(w http.ResponseWriter, state string) error

	// GetRepo allows getting repo name from an HTTP request.
	GetRepo(r *http.Request) (string, error)

	// SetRepo allows setting state on an HTTP response. Must be called
	// before the first byte of the response body is written.
	SetRepo(w http.ResponseWriter, repo string) error
}
