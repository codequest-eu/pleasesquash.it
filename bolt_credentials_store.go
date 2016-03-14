package main

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
	"golang.org/x/oauth2"
)

var (
	boltBucketName         = []byte("credentials")
	errCredentialsNotFound = errors.New("credentials not found")
)

type boltCredentialsStore struct {
	db *bolt.DB
}

func createBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists(boltBucketName)
	return err
}

func newBoltCredentialsStore(path string) (CredentialsStore, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	if err := db.Update(createBucket); err != nil {
		return nil, err
	}
	return &boltCredentialsStore{db}, nil
}

func (s *boltCredentialsStore) Get(repoName string) (*oauth2.Token, error) {
	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		if data = tx.Bucket(boltBucketName).Get([]byte(repoName)); data == nil {
			return errCredentialsNotFound
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	ret := new(oauth2.Token)
	return ret, json.Unmarshal(data, ret)
}

func (s *boltCredentialsStore) Set(repoName string, token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(boltBucketName).Put([]byte(repoName), data)
	})
}
