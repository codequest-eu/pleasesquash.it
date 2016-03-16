package credentials

import "golang.org/x/oauth2"

// Store is a pluggable credentials backend for the application.
type Store interface {
	// Get allows getting OAuth2 token for a owner/repo combination.
	Get(repoName string) (*oauth2.Token, error)

	// Set allows setting OAuth2 token for a owner/repo combination.
	Set(repoName string, token *oauth2.Token) error
}
