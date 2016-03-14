package main

import (
	"errors"
	"strings"
)

var errNoRepoData = errors.New("username/repo not provided")

func getOwnerAndRepo(input string) (string, string, error) {
	tokens := strings.Split(input, "/")
	if len(tokens) != 2 {
		return "", "", errNoRepoData
	}
	return tokens[0], tokens[1], nil
}
