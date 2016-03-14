package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	ghendpoint "golang.org/x/oauth2/github"
)

type oauthFactoryImpl struct {
	config *oauth2.Config
}

func newOauthFactory(clientID, clientSecret string) OAuthFactory {
	return &oauthFactoryImpl{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"repo:status"},
			Endpoint:     ghendpoint.Endpoint,
		},
	}
}

func (f *oauthFactoryImpl) Client(token *oauth2.Token) *github.Client {
	return github.NewClient(
		oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(token),
		),
	)
}

func (f *oauthFactoryImpl) Exchange(code string) (*oauth2.Token, error) {
	return f.config.Exchange(context.Background(), code)
}

func (f *oauthFactoryImpl) GetURL(state, repoName string) string {
	return f.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}
