package sync

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/google/go-github/github"
	"github.com/knqyf263/pet/config"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)
const (
	githubTokenEnvVariable = "PET_GITHUB_ACCESS_TOKEN"
)

// GistClient manages communication with Gist
type GistClient struct {
	Client *github.Client
	ID     string
}

// NewGistClient returns GistClient
func NewGistClient() (Client, error) {
	accessToken, err := getGithubAccessToken()
	if err != nil {
		return nil, fmt.Errorf(`access_token is empty.
Go https://github.com/settings/tokens/new and create access_token (only need "gist" scope).
Write access_token in config file (pet configure) or export $%v.
		`, githubTokenEnvVariable)
	}

	client := GistClient{
		Client: githubClient(accessToken),
		ID:     config.Conf.Gist.GistID,
	}
	return client, nil
}

func getGithubAccessToken() (string, error) {
	if config.Conf.Gist.AccessToken != "" {
		return config.Conf.Gist.AccessToken, nil
	} else if os.Getenv(githubTokenEnvVariable) != "" {
		return os.Getenv(githubTokenEnvVariable), nil
	}
	return "", errors.New("Github AccessToken not found in any source")
}

// GetSnippet returns the remote snippet
func (g GistClient) GetSnippet() (*Snippet, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Getting Gist..."
	defer s.Stop()

	if g.ID == "" {
		return &Snippet{}, nil
	}

	gist, res, err := g.Client.Gists.Get(context.Background(), g.ID)
	if err != nil {
		if res.StatusCode == 404 {
			return nil, errors.Wrapf(err, "No gist ID (%s)", g.ID)
		}
		return nil, errors.Wrapf(err, "Failed to get gist")
	}

	content := ""
	filename := config.Conf.Gist.FileName
	for _, file := range gist.Files {
		if *file.Filename == filename {
			content = *file.Content
		}
	}
	if content == "" {
		return nil, fmt.Errorf("%s is empty", filename)
	}

	return &Snippet{
		Content:   content,
		UpdatedAt: *gist.UpdatedAt,
	}, nil
}

// UploadSnippet uploads local snippets to Gist
func (g GistClient) UploadSnippet(content string) error {
	gist := &github.Gist{
		Description: github.String("description"),
		Public:      github.Bool(config.Conf.Gist.Public),
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(config.Conf.Gist.FileName): github.GistFile{
				Content: github.String(content),
			},
		},
	}

	if g.ID == "" {
		gistID, err := g.createGist(context.Background(), gist)
		if err != nil {
			return err
		}
		fmt.Printf("Gist ID: %s\n", *gistID)
	} else {
		if err := g.updateGist(context.Background(), gist); err != nil {
			return errors.Wrap(err, "Failed to update gist")
		}
	}
	return nil
}

func (g GistClient) createGist(ctx context.Context, gist *github.Gist) (gistID *string, err error) {
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

func (g GistClient) updateGist(ctx context.Context, gist *github.Gist) (err error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Updating Gist..."
	defer s.Stop()

	if _, _, err = g.Client.Gists.Edit(ctx, g.ID, gist); err != nil {
		return errors.Wrap(err, "Failed to edit gist")
	}
	return nil
}

func githubClient(accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return client
}
