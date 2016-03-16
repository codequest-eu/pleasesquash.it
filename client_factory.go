package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	ghendpoint "golang.org/x/oauth2/github"
)

type ClientFactory interface {
	Client(token *oauth2.Token) *github.Client
	Exchange(code string) (*oauth2.Token, error)
	GetURL(state, repoName string) string
}

type clientFactoryImpl struct {
	config *oauth2.Config
}

func newOauthFactory(clientID, clientSecret string) ClientFactory {
	return &clientFactoryImpl{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"repo:status"},
			Endpoint:     ghendpoint.Endpoint,
		},
	}
}

func (f *clientFactoryImpl) Client(token *oauth2.Token) *github.Client {
	return github.NewClient(
		oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(token),
		),
	)
}

func (f *clientFactoryImpl) Exchange(code string) (*oauth2.Token, error) {
	return f.config.Exchange(context.Background(), code)
}

func (f *clientFactoryImpl) GetURL(state, repoName string) string {
	return f.config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)
}
