package main

import (
	"net/http"

	"github.com/google/go-github/github"

	"golang.org/x/oauth2"
)

// StateStore is a mechanism allowing temporary storage of OAuth2 state to
// prevent from CSRF attacks.
type StateStore interface {
	// Get allows getting state from an HTTP request.
	Get(r *http.Request) (string, error)

	// Set allows setting state on an HTTP response. Must be called
	// before the first byte of the response body is written.
	Set(w http.ResponseWriter, state string) error
} // type StateStore interface

// CredentialsStore is a pluggable credentials backend for the application.
type CredentialsStore interface {
	// Get allows getting OAuth2 token for a owner/repo combination.
	Get(repoName string) (*oauth2.Token, error)

	// Set allows setting OAuth2 token for a owner/repo combination.
	Set(repoName string, token *oauth2.Token) error
} // type StateStore interface

type OAuthFactory interface {
	Client(token *oauth2.Token) (*github.Client, error)
	Exchange(code string) (*oauth2.Token, error)
	GetURL(state, repoName string) (string, error)
}
