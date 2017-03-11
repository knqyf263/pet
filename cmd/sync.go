package cmd

import (
	"context"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync snippets",
	Long:  `Sync snippets with gist`,
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) (err error) {
	if config.Conf.Gist.AccessToken == "" {
		return fmt.Errorf(`access_token is empty.
Go https://github.com/settings/tokens/new and create access_token (only need "gist" scope).
Write access_token in config file (pet configure).
		`)
	}

	if config.Flag.Upload {
		return upload()
	}
	return download()
}

func githubClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Conf.Gist.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return client
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
	gist := github.Gist{
		Description: github.String("description"),
		Public:      github.Bool(true),
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(config.Conf.Gist.FileName): github.GistFile{
				Content: github.String(body),
			},
		},
	}

	gistID := config.Conf.Gist.GistID
	if gistID == "" {
		var retGist *github.Gist
		retGist, _, err = client.Gists.Create(ctx, &gist)
		if err != nil {
			return err
		}
		fmt.Printf("Gist ID: %s\n", retGist.GetID())
	} else {
		_, _, err = client.Gists.Edit(ctx, gistID, &gist)
		if err != nil {
			return err
		}
	}
	fmt.Println("Upload success")
	return nil
}

func download() error {
	if config.Conf.Gist.GistID == "" {
		return fmt.Errorf("Gist ID is empty")
	}
	ctx := context.Background()
	client := githubClient()
	resGist, _, err := client.Gists.Get(ctx, config.Conf.Gist.GistID)
	if err != nil {
		return fmt.Errorf("Failed to download gist: %v", err)
	}
	content := resGist.Files[github.GistFilename(config.Conf.Gist.FileName)].Content

	var snippets snippet.Snippets
	toml.Decode(*content, &snippets)
	if err := snippets.Save(); err != nil {
		return err
	}
	fmt.Println("Download success")
	return nil
}

func init() {
	RootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&config.Flag.Upload, "upload", "u", false,
		`Upload snippets to gist`)
}
