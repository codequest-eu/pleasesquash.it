package main

import "github.com/google/go-github/github"

func buildStatus(success bool) *github.RepoStatus {
	if success {
		return buildSuccess()
	}
	return buildError()
}

func buildSuccess() *github.RepoStatus {
	return buildStatusImpl("success", "Commits squashed")
}

func buildError() *github.RepoStatus {
	return buildStatusImpl("error", "Y U NO squash UR commits?")
}

func buildStatusImpl(state, description string) *github.RepoStatus {
	return &github.RepoStatus{
		State:       github.String(state),
		Description: github.String(description),
		Context:     github.String("pleasesquash.me"),
	}
}
