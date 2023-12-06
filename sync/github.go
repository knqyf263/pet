package sync

import (
	"context"
	"fmt"
	"github.com/knqyf263/pet/config"
	"time"

	"github.com/briandowns/spinner"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// GithubClient manages communication with Gist
type GithubClient struct {
	Client *github.Client
}

// NewGithubClient returns GithubClient
func NewGithubClient() (Client, error) {
	accessToken, err := getGithubAccessToken()
	if err != nil {
		return nil, fmt.Errorf(`access_token is empty.
Go https://github.com/settings/tokens/new and create access_token (only need "repo" scope).
Write access_token in config file (pet configure) or export $%v.
		`, githubTokenEnvVariable)
	}

	client := GithubClient{
		Client: githubClient(accessToken),
	}
	return client, nil
}

// GetSnippet returns the remote snippet
func (g GithubClient) GetSnippet() (*Snippet, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Getting configuration from Github..."
	defer s.Stop()

	ghConfig := config.Conf.GitHub
	fileContent, _, resp, err := g.Client.Repositories.GetContents(context.Background(), ghConfig.RepoOwner, ghConfig.RepoName, ghConfig.FileName, nil)
	if err != nil {
		fmt.Printf("Error from Github: %s", resp.Status)
		return nil, errors.Wrapf(err, "Failed to get repo")
	}

	return &Snippet{
		Content: *(fileContent.Content),
	}, nil
}

// UploadSnippet uploads local snippets to Github
func (g GithubClient) UploadSnippet(content string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Uploading configuration to Github..."
	defer s.Stop()

	// input data
	owner, repoName, fileName := config.Conf.GitHub.RepoOwner, config.Conf.GitHub.RepoName, config.Conf.GitHub.FileName
	// we need the SHA to be able to update the file
	sha := getShaForFile(fileName, owner, repoName, g)
	// update content
	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Updating snippet configuration"),
		Content: []byte(content),
		Branch:  github.String("main"),
		SHA:     &sha,
	}
	_, _, err := g.Client.Repositories.CreateFile(context.Background(), owner, repoName, fileName, opts)
	if err != nil {
		return errors.Wrap(err, "Failed to upload changes to github")
	}
	return nil
}

func getShaForFile(fileName string, owner string, repoName string, g GithubClient) string {
	sha := ""
	opt := &github.CommitsListOptions{
		Path: *github.String(fileName),
	}
	commits, _, _ := g.Client.Repositories.ListCommits(context.Background(), owner, repoName, opt)
	t, _, _ := g.Client.Git.GetTree(context.Background(), owner, repoName, commits[0].GetSHA(), true)
	for _, entry := range t.Entries {
		if *entry.Path == *github.String(fileName) {
			sha = *entry.SHA
		}
	}
	return sha
}
