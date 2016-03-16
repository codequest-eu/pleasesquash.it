package credentials

import (
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
)

type googleStore struct {
	client *datastore.Client
}

func NewGoogleStore(keyPath, projectID string) (Store, error) {
	keyContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(
		keyContent,
		datastore.ScopeDatastore,
		datastore.ScopeUserEmail,
	)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	client, err := datastore.NewClient(
		ctx,
		projectID,
		cloud.WithTokenSource(conf.TokenSource(ctx)),
	)
	if err != nil {
		return nil, err
	}
	return &googleStore{client}, nil
}

func (s *googleStore) Get(repoName string) (*oauth2.Token, error) {
	ret := new(oauth2.Token)
	return ret, s.client.Get(context.Background(), keyFor(repoName), ret)
}

func (s *googleStore) Set(repoName string, token *oauth2.Token) error {
	_, err := s.client.Put(context.Background(), keyFor(repoName), token)
	return err
}

func keyFor(repoName string) *datastore.Key {
	return datastore.NewKey(
		context.Background(), // context
		"Tokens",             // kind
		repoName,             // name
		0,                    // id,
		nil,                  //parent
	)
}
