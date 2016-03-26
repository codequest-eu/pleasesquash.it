package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/google/go-github/github"
	"github.com/satori/go.uuid"

	"github.com/codequest-eu/pleasesquash.me/credentials"
	"github.com/codequest-eu/pleasesquash.me/state"
)

type handler struct {
	state state.Store
	creds credentials.Store
	oauth ClientFactory
}

var errStateMismatch = errors.New("OAuth2 state comparison mismatch")

func (h *handler) submit(w http.ResponseWriter, r *http.Request) error {
	repoName := r.FormValue("repo_name")
	if _, _, err := getOwnerAndRepo(repoName); err != nil {
		return err
	}
	state := uuid.NewV4().String()
	if err := h.state.SetState(w, state); err != nil {
		return err
	}
	if err := h.state.SetRepo(w, repoName); err != nil {
		return err
	}
	glog.Infof("Hook setup requested for %q", repoName)
	http.Redirect(w, r, h.oauth.GetURL(state, repoName), http.StatusFound)
	return nil
}

func (h *handler) callback(w http.ResponseWriter, r *http.Request) (err error) {
	state, err := h.state.GetState(r)
	if err != nil {
		return
	}
	if state != r.FormValue("state") {
		err = errStateMismatch
		return
	}
	repoName, err := h.state.GetRepo(r)
	if err != nil {
		return
	}
	token, err := h.oauth.Exchange(r.FormValue("code"))
	if err != nil {
		return
	}
	owner, repo, err := getOwnerAndRepo(repoName)
	if err != nil {
		return
	}
	if _, _, err = h.oauth.Client(token).Repositories.ListStatuses(owner, repo, "master", nil); err != nil {
		err = renderError(w, owner, repo)
		return
	}
	if err = h.creds.Set(repoName, token); err != nil {
		return
	}
	glog.Infof("Added hook for %q", repoName)
	err = renderSuccess(w, owner, repo)
	return
}

func (h *handler) webhook(w http.ResponseWriter, r *http.Request) (err error) {
	if r.Header.Get("X-GitHub-Event") != "pull_request" {
		_, err = fmt.Fprint(w, "Whatever 🐶")
		return
	}
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
	owner, repo, err := getOwnerAndRepo(*evt.Repo.FullName)
	if err != nil {
		return
	}
	_, _, err = h.oauth.Client(token).Repositories.CreateStatus(
		owner,
		repo,
		*evt.PullRequest.Head.SHA,
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
