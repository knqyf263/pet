package sync

import (
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
	r, err := getRepository(directory)
	w, err := r.Worktree()

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	filename := filepath.Join(directory, "snippet.toml")
	ref, err := r.Head()
	commit, err := r.CommitObject(ref.Hash())
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, errors.Wrapf(err, "snippet.toml")
	}

	return &Snippet{
		Content:   string(content),
		UpdatedAt: commit.Author.When,
	}, nil
}

func (g GitClient) UploadSnippet(body string) error {
	directory, err := getDirectory()
	author, authorEmail, err := getGitSignature()
	r, err := getRepository(directory)
	w, err := r.Worktree()

	if err != nil {
		return errors.Wrapf(err, "Failed to load repository")
	}

	_, err = w.Add("snippet.toml")
	_, err = w.Commit("update snippets", &git.CommitOptions{
		Author: &object.Signature{
			Name:  author,
			Email: authorEmail,
			When:  time.Now(),
		},
	})

	return r.Push(&git.PushOptions{})
}

func getDirectory() (string, error) {
	if config.Conf.Git.Directory != "" {
		return config.Conf.Git.Directory, nil
	}
	return "", errors.New("Git Directory not found")
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
	return "", "", errors.New("Git Signature (Author, AuthorEmail) not found")
}
