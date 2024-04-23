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
	gheTokenEnvVariable = "PET_GHE_ACCESS_TOKEN"
)

// GHEGistClient manages communication with Gist
type GHEGistClient struct {
	Client *github.Client
	ID     string
}

// NewGHEGistClient returns GistClient
func NewGHEGistClient() (Client, error) {
	accessToken, err := getGHEAccessToken()
	if err != nil {
		return nil, fmt.Errorf(`access_token is empty.
Go https://github.com/settings/tokens/new and create access_token (only need "gist" scope).
Write access_token in config file (pet configure) or export $%v.
		`, gheTokenEnvVariable)
	}

	var baseUrl, uploadUrl string

	if config.Conf.GHEGist.BaseUrl != "" {
		fmt.Println(config.Conf.GHEGist.BaseUrl)
		baseUrl = config.Conf.GHEGist.BaseUrl
	}

	if config.Conf.GHEGist.UploadUrl != "" {
		fmt.Println(config.Conf.GHEGist.UploadUrl)
		uploadUrl = config.Conf.GHEGist.UploadUrl
	}

	client := GHEGistClient{
		Client: githubEnterpriseClient(accessToken, baseUrl, uploadUrl),
		ID:     config.Conf.GHEGist.GistID,
	}
	return client, nil
}

func getGHEAccessToken() (string, error) {
	if config.Conf.GHEGist.AccessToken != "" {
		return config.Conf.GHEGist.AccessToken, nil
	} else if os.Getenv(gheTokenEnvVariable) != "" {
		return os.Getenv(gheTokenEnvVariable), nil
	}
	return "", errors.New("GHE AccessToken not found in any source")
}

// GetSnippet returns the remote snippet
func (g GHEGistClient) GetSnippet() (*Snippet, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Getting Gist..."
	defer s.Stop()

	if g.ID == "" {
		return &Snippet{}, nil
	}

	gist, res, err := g.Client.Gists.Get(context.Background(), g.ID)
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			return nil, errors.Wrapf(err, "No gist ID (%s)", g.ID)
		}
		return nil, errors.Wrapf(err, "Failed to get gist")
	}

	content := ""
	filename := config.Conf.GHEGist.FileName
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
func (g GHEGistClient) UploadSnippet(content string) error {
	gist := &github.Gist{
		Description: github.String("description"),
		Public:      github.Bool(config.Conf.GHEGist.Public),
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(config.Conf.GHEGist.FileName): github.GistFile{
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

func (g GHEGistClient) createGist(ctx context.Context, gist *github.Gist) (gistID *string, err error) {
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

func (g GHEGistClient) updateGist(ctx context.Context, gist *github.Gist) (err error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Updating Gist..."
	defer s.Stop()

	if _, _, err = g.Client.Gists.Edit(ctx, g.ID, gist); err != nil {
		return errors.Wrap(err, "Failed to edit gist")
	}
	return nil
}

func githubEnterpriseClient(accessToken, baseURL, uploadURL string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client, _ := github.NewEnterpriseClient(baseURL, uploadURL, tc)
	return client
}
