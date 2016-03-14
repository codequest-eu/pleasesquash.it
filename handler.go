package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang/glog"
	"github.com/google/go-github/github"
	"github.com/satori/go.uuid"
)

type handler struct {
	state StateStore
	creds CredentialsStore
	oauth OAuthFactory
}

var errStateMismatch = errors.New("OAuth2 state comparison mismatch")

func (h *handler) submit(w http.ResponseWriter, r *http.Request) error {
	repoName := r.FormValue("repo_name")
	if _, _, err := getOwnerAndRepo(repoName); err != nil {
		return err
	}
	state := uuid.NewV4().String()
	if err := h.state.Set(w, state); err != nil {
		return err
	}
	url, err := h.oauth.GetURL(state, repoName)
	if err != nil {
		return err
	}
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func (h *handler) callback(w http.ResponseWriter, r *http.Request) (err error) {
	state, err := h.state.Get(r)
	if err != nil {
		return
	}
	if state != r.FormValue("state") {
		err = errStateMismatch
		return
	}
	token, err := h.oauth.Exchange(r.FormValue("code"))
	if err != nil {
		return
	}
	repoName := r.FormValue("repo_name")
	glog.Infof("Adding hook for %q", repoName)
	owner, repo, err := getOwnerAndRepo(repoName)
	if err != nil {
		return
	}
	client, err := h.oauth.Client(token)
	if err != nil {
		return
	}
	if _, _, err = client.Repositories.ListStatuses(owner, repo, "master", nil); err != nil {
		err = renderError(w, owner, repo)
		return
	}
	if err = h.creds.Set(repoName, token); err != nil {
		return
	}
	err = renderSuccess(w, owner, repo)
	return
}

func (h *handler) webhook(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() {
		if issue := recover(); issue != nil {
			err = recoveryError(issue)
		}
	}()
	evt := new(github.PullRequestEvent)
	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(evt); err != nil {
		return
	}
	token, err := h.creds.Get(*evt.Repo.FullName)
	if err != nil {
		return
	}
	client, err := h.oauth.Client(token)
	if err != nil {
		return
	}
	owner, repo, err := getOwnerAndRepo(*evt.Repo.FullName)
	if err != nil {
		return
	}
	_, _, err = client.Repositories.CreateStatus(
		owner,
		repo,
		*evt.PullRequest.Head.Ref,
		buildStatus(*evt.PullRequest.Commits == 1), // money shot!
	)
	return
}

func recoveryError(issue interface{}) error {
	switch x := issue.(type) {
	case string:
		return errors.New(x)
	case error:
		return x
	default:
		return errors.New("unknown panic")
	}
}
