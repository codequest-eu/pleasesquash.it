package main

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

const (
	repoCookieKey  = "oauth-repo"
	stateCookieKey = "oauth-state"
)

type cookieStateStore struct {
	cutter *securecookie.SecureCookie
}

func newSecureCookieStore(hashKey, blockKey []byte) StateStore {
	return &cookieStateStore{securecookie.New(hashKey, blockKey)}
}

func (s *cookieStateStore) GetState(r *http.Request) (string, error) {
	return s.getCookie(r, stateCookieKey)
}

func (s *cookieStateStore) SetState(w http.ResponseWriter, state string) error {
	return s.setCookie(w, stateCookieKey, state)
}

func (s *cookieStateStore) GetRepo(r *http.Request) (string, error) {
	return s.getCookie(r, repoCookieKey)
}

func (s *cookieStateStore) SetRepo(w http.ResponseWriter, repo string) error {
	return s.setCookie(w, repoCookieKey, repo)
}

func (s *cookieStateStore) getCookie(r *http.Request, key string) (string, error) {
	var value string
	cookie, err := r.Cookie(key)
	if err != nil {
		return "", err
	}
	return value, s.cutter.Decode(key, cookie.Value, &value)
}

func (s *cookieStateStore) setCookie(w http.ResponseWriter, key, value string) error {
	encoded, err := s.cutter.Encode(key, value)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{Name: key, Value: encoded, Path: "/"})
	return nil
}
