package sync

import (
	"bytes"
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
Go https://github.com/settings/tokens/new and create access_token (only need "gist" scope).
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
	// todo add timeout?
	content, err := g.Client.Repositories.DownloadContents(context.Background(), ghConfig.RepoOwner, ghConfig.RepoName, ghConfig.FileName, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get repo")
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(content)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed read github response")
	}

	return &Snippet{
		Content: buf.String(),
	}, nil
}

// UploadSnippet uploads local snippets to Github
func (g GithubClient) UploadSnippet(content string) error {
	ghConfig := config.Conf.GitHub
	// get SHA to be able to update the file
	sha := ""
	opt := &github.CommitsListOptions{
		Path: *github.String(ghConfig.FileName),
	}
	commits, _, _ := g.Client.Repositories.ListCommits(context.Background(), ghConfig.RepoOwner, ghConfig.RepoName, opt)
	t, _, _ := g.Client.Git.GetTree(context.Background(), ghConfig.RepoOwner, ghConfig.RepoName, commits[0].GetSHA(), true)
	for _, entry := range t.Entries {
		if *entry.Path == *github.String(ghConfig.FileName) {
			sha = *entry.SHA
		}
	}
	// update content
	fileContent := []byte(content)
	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Updating snippet configuration"),
		Content: fileContent,
		Branch:  github.String("main"),
		//Committer: &github.CommitAuthor{Name: github.String("FirstName LastName"), Email: github.String("user@example.com")},
		SHA: &sha,
	}
	_, _, err := g.Client.Repositories.CreateFile(context.Background(), ghConfig.RepoOwner, ghConfig.RepoName, ghConfig.FileName, opts)
	if err != nil {
		return errors.Wrap(err, "Failed to upload changes to github")
	}
	return nil
}

/**
func (g GithubClient) createGist(ctx context.Context, gist *github.Gist) (gistID *string, err error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Creating Gist..."
	defer s.Stop()

	retGist, _, err := g.Client.Gists.Create(ctx, gist)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create gist")
	}
	return retGist.ID, nil
}

func (g GithubClient) updateGist(ctx context.Context, gist *github.Gist) (err error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Updating Gist..."
	defer s.Stop()

	if _, _, err = g.Client.Gists.Edit(ctx, g.ID, gist); err != nil {
		return errors.Wrap(err, "Failed to edit gist")
	}
	return nil
}
*/
