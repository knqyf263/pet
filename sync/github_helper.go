package sync

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	githubTokenEnvVariable = "PET_GITHUB_ACCESS_TOKEN"
)

func githubClient(accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	return client
}
