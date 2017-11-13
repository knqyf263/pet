package sync

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/google/go-github/github"
	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// AutoSync syncs snippets automatically
func AutoSync(file string) error {
	gistID := config.Conf.Gist.GistID
	if config.Conf.Gist.GistID == "" {
		upload()
		return nil
	}

	gist, err := getGist(gistID)
	if err != nil {
		return err
	}

	fi, err := os.Stat(file)
	if os.IsNotExist(err) {
		return download(gist)
	} else if err != nil {
		return errors.Wrap(err, "Failed to get a FileInfo")
	}

	local := fi.ModTime().UTC()
	remote := gist.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return upload()
	case remote.After(local):
		return download(gist)
	default:
		return nil
	}
}

func getGist(gistID string) (*github.Gist, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Getting Gist..."
	defer s.Stop()

	client := githubClient()
	gist, res, err := client.Gists.Get(context.Background(), gistID)
	if err != nil {
		if res.StatusCode == 404 {
			return nil, errors.Wrapf(err, "No gist ID (%s)", gistID)
		}
		return nil, errors.Wrapf(err, "Failed to get gist")
	}
	return gist, nil
}

func upload() (err error) {
	ctx := context.Background()

	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return err
	}

	body, err := snippets.ToString()
	if err != nil {
		return err
	}

	client := githubClient()
	gist := &github.Gist{
		Description: github.String("description"),
		Public:      github.Bool(config.Conf.Gist.Public),
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(config.Conf.Gist.FileName): github.GistFile{
				Content: github.String(body),
			},
		},
	}

	gistID := config.Conf.Gist.GistID
	if gistID == "" {
		retGist, err := createGist(ctx, client, gist)
		if err != nil {
			return err
		}
		fmt.Printf("Gist ID: %s\n", retGist.GetID())
	} else {
		if err = updateGist(ctx, gistID, client, gist); err != nil {
			return err
		}
	}
	fmt.Println("Upload success")
	return nil
}

func download(gist *github.Gist) error {
	var (
		content     = ""
		snippetFile = config.Conf.General.SnippetFile
		filename    = config.Conf.Gist.FileName
	)
	for _, file := range gist.Files {
		if *file.Filename == filename {
			content = *file.Content
		}
	}
	if content == "" {
		return fmt.Errorf("%s is empty", filename)
	}

	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return err
	}
	body, err := snippets.ToString()
	if err != nil {
		return err
	}
	if content == body {
		// no need to download
		fmt.Println("Already up-to-date")
		return nil
	}

	fmt.Println("Download success")
	return ioutil.WriteFile(snippetFile, []byte(content), os.ModePerm)
}

func githubClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Conf.Gist.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return client
}

func createGist(ctx context.Context, client *github.Client, gist *github.Gist) (*github.Gist, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Creating Gist..."
	defer s.Stop()

	retGist, _, err := client.Gists.Create(ctx, gist)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create gist")
	}
	return retGist, nil
}

func updateGist(ctx context.Context, gistID string, client *github.Client, gist *github.Gist) (err error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	s.Suffix = " Updating Gist..."
	defer s.Stop()

	if _, _, err = client.Gists.Edit(ctx, gistID, gist); err != nil {
		return errors.Wrap(err, "Failed to edit gist")
	}
	return nil
}
