package state

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

const (
	repoCookieKey  = "oauth-repo"
	stateCookieKey = "oauth-state"
)

type secureCookieStore struct {
	cutter *securecookie.SecureCookie
}

func NewSecureCookieStore(hashKey, blockKey []byte) Store {
	return &secureCookieStore{securecookie.New(hashKey, blockKey)}
}

func (s *secureCookieStore) GetState(r *http.Request) (string, error) {
	return s.getCookie(r, stateCookieKey)
}

func (s *secureCookieStore) SetState(w http.ResponseWriter, state string) error {
	return s.setCookie(w, stateCookieKey, state)
}

func (s *secureCookieStore) GetRepo(r *http.Request) (string, error) {
	return s.getCookie(r, repoCookieKey)
}

func (s *secureCookieStore) SetRepo(w http.ResponseWriter, repo string) error {
	return s.setCookie(w, repoCookieKey, repo)
}

func (s *secureCookieStore) getCookie(r *http.Request, key string) (string, error) {
	var value string
	cookie, err := r.Cookie(key)
	if err != nil {
		return "", err
	}
	return value, s.cutter.Decode(key, cookie.Value, &value)
}

func (s *secureCookieStore) setCookie(w http.ResponseWriter, key, value string) error {
	encoded, err := s.cutter.Encode(key, value)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{Name: key, Value: encoded, Path: "/"})
	return nil
}
