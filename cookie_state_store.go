package main

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

const cookieKey = "oauth-state"

type cookieStateStore struct {
	cutter *securecookie.SecureCookie
}

func newSecureCookieStore(hashKey, blockKey []byte) StateStore {
	return &cookieStateStore{securecookie.New(hashKey, blockKey)}
}

func (s *cookieStateStore) Get(r *http.Request) (string, error) {
	var ret string
	cookie, err := r.Cookie(cookieKey)
	if err != nil {
		return "", err
	}
	return ret, s.cutter.Decode(cookieKey, cookie.Value, &ret)
}

func (s *cookieStateStore) Set(w http.ResponseWriter, state string) error {
	encoded, err := s.cutter.Encode(cookieKey, state)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{Name: cookieKey, Value: encoded, Path: "/"})
	return nil
}
