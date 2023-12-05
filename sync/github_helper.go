package sync

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/knqyf263/pet/config"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"os"
)

const (
	githubTokenEnvVariable = "PET_GITHUB_ACCESS_TOKEN"
)

func getGithubAccessToken() (string, error) {
	if config.Conf.GitHub.AccessToken != "" {
		return config.Conf.Gist.AccessToken, nil
	} else if os.Getenv(githubTokenEnvVariable) != "" {
		return os.Getenv(githubTokenEnvVariable), nil
	}
	return "", errors.New("Github AccessToken not found in any source")
}

func getGithubGistAccessToken() (string, error) {
	if config.Conf.Gist.AccessToken != "" {
		return config.Conf.Gist.AccessToken, nil
	} else if os.Getenv(githubTokenEnvVariable) != "" {
		return os.Getenv(githubTokenEnvVariable), nil
	}
	return "", errors.New("Github AccessToken not found in any source")
}

func githubClient(accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	return client
}
