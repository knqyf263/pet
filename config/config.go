package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Conf is global config variable
var Conf Config

// Config is a struct of config
type Config struct {
	General   GeneralConfig `toml:"General"`
	Gist      GistConfig    `toml:"Gist"`
	EnvGist   GistConfig    `toml:"EnvGist"`
	GitLab    GitLabConfig  `toml:"GitLab"`
	EnvGitLab GitLabConfig  `toml:"EnvGitLab"`
}

// GeneralConfig is a struct of general config
type GeneralConfig struct {
	SnippetFile string `toml:"snippetfile"`
	EnvFile     string `toml:"envfile"`
	Editor      string `toml:"editor"`
	Column      int    `toml:"column"`
	SelectCmd   string `toml:"selectcmd"`
	Backend     string `toml:"backend"`
	SortBy      string `toml:"sortby"`
}

// GistConfig is a struct of config for Gist
type GistConfig struct {
	FileName    string `toml:"file_name"`
	AccessToken string `toml:"access_token"`
	GistID      string `toml:"gist_id"`
	Public      bool   `toml:"public"`
	AutoSync    bool   `toml:"auto_sync"`
}

// GitLabConfig is a struct of config for GitLabSnippet
type GitLabConfig struct {
	FileName    string `toml:"file_name"`
	AccessToken string `toml:"access_token"`
	Url         string `toml:"url"`
	ID          string `toml:"id"`
	Visibility  string `toml:"visibility"`
	AutoSync    bool   `toml:"auto_sync"`
	Insecure    bool   `toml:"skip_ssl"`
}

// Flag is global flag variable
var Flag FlagConfig

// FlagConfig is a struct of flag
type FlagConfig struct {
	Debug     bool
	Query     string
	FilterTag string
	Command   bool
	Delimiter string
	OneLine   bool
	Color     bool
	Tag       bool
}

// Load loads a config toml
func (cfg *Config) Load(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		cfg.General.SnippetFile = expandPath(cfg.General.SnippetFile)
		cfg.General.EnvFile = expandPath(cfg.General.EnvFile)
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	dir, err := GetDefaultConfigDir()
	if err != nil {
		return errors.Wrap(err, "Failed to get the default config directory")
	}
	cfg.General.SnippetFile = filepath.Join(dir, "snippet.toml")
	cfg.General.EnvFile = filepath.Join(dir, "env.toml")

	_, err = os.Create(cfg.General.SnippetFile)
	if err != nil {
		return errors.Wrap(err, "Failed to create a config file")
	}

	_, err = os.Create(cfg.General.EnvFile)
	if err != nil {
		return errors.Wrap(err, "Failed to create an env config file")
	}

	cfg.General.Editor = os.Getenv("EDITOR")
	if cfg.General.Editor == "" && runtime.GOOS != "windows" {
		if isCommandAvailable("sensible-editor") {
			cfg.General.Editor = "sensible-editor"
		} else {
			cfg.General.Editor = "vim"
		}
	}
	cfg.General.Column = 40
	cfg.General.SelectCmd = "fzf"
	cfg.General.Backend = "gist"

	cfg.Gist.FileName = "pet-snippet.toml"
	cfg.EnvGist.FileName = "pet-env.toml"

	cfg.GitLab.FileName = "pet-snippet.toml"
	cfg.EnvGitLab.FileName = "pet-env.toml"
	cfg.GitLab.Visibility = "private"

	return toml.NewEncoder(f).Encode(cfg)
}

// GetDefaultConfigDir returns the default config directory
func GetDefaultConfigDir() (dir string, err error) {
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "pet")
		}
		dir = filepath.Join(dir, "pet")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "pet")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("cannot create directory: %v", err)
	}
	return dir, nil
}

func expandPath(s string) string {
	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = filepath.Join(os.Getenv("USERPROFILE"), s[2:])
		} else {
			s = filepath.Join(os.Getenv("HOME"), s[2:])
		}
	}
	return os.Expand(s, os.Getenv)
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
