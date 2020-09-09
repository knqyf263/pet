package sync

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/knqyf263/pet/config"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"time"
)

type GitClient struct{}

func NewGitClient() (Client, error) {
	return GitClient{}, nil
}

func (g GitClient) GetSnippet() (*Snippet, error) {
	directory, err := getDirectory()
	if err != nil {
		return nil, err
	}

	r, err := getRepository(directory)
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, fmt.Errorf("unable to open git work tree: %w", err)
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return nil, fmt.Errorf("unable to pull: %w", err)
	}

	filename := filepath.Join(directory, "snippet.toml")
	ref, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("unable to find HEAD ref: %w", err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("unable to create commit: %w", err)
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unabled to read file snippet.toml: %w", err)
	}

	return &Snippet{
		Content:   string(content),
		UpdatedAt: commit.Author.When,
	}, nil
}

func (g GitClient) UploadSnippet(body string) error {
	directory, err := getDirectory()
	if err != nil {
		return err
	}

	author, authorEmail, err := getGitSignature()
	if err != nil {
		return err
	}

	r, err := getRepository(directory)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("unable to open git work tree: %w", err)
	}

	_, err = w.Add("snippet.toml")
	if err != nil {
		return fmt.Errorf("unable to add snippet.toml to index: %w", err)
	}

	_, err = w.Commit("update snippets", &git.CommitOptions{
		Author: &object.Signature{
			Name:  author,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("unable to create commit: %w", err)
	}

	return r.Push(&git.PushOptions{})
}

func getDirectory() (string, error) {
	if config.Conf.Git.Directory != "" {
		return config.Conf.Git.Directory, nil
	}
	return "", errors.New("git directory not found")
}

func getRepository(directory string) (*git.Repository, error) {
	r, err := git.PlainOpen(directory)

	if err != nil {
		r, err = git.PlainClone(directory, false, &git.CloneOptions{
			URL:               config.Conf.Git.Repo,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})

		if err != nil {
			return nil, errors.Wrapf(err, "unable to create repository")
		}
	}

	return r, nil
}

func getGitSignature() (string, string, error) {
	if config.Conf.Git.Author != "" && config.Conf.Git.AuthorEmail != "" {
		return config.Conf.Git.Author, config.Conf.Git.AuthorEmail, nil
	}
	return "", "", errors.New("git signature (Author, AuthorEmail) not found")
}
