package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Conf is global config variable
var Conf Config

// Config is a struct of config
type Config struct {
	General GeneralConfig
	Gist    GistConfig
	GitLab  GitLabConfig
}

// GeneralConfig is a struct of general config
type GeneralConfig struct {
	SnippetFile string `toml:"snippetfile"`
	Editor      string `toml:"editor"`
	Column      int    `toml:"column"`
	SelectCmd   string `toml:"selectcmd"`
	Backend     string `toml:"backend"`
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
	ID          string `toml:"id"`
	Visibility  string `toml:"visibility"`
	AutoSync    bool   `toml:"auto_sync"`
}

// Flag is global flag variable
var Flag FlagConfig

// FlagConfig is a struct of flag
type FlagConfig struct {
	Debug     bool
	Query     string
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
	_, err = os.Create(cfg.General.SnippetFile)
	if err != nil {
		return errors.Wrap(err, "Failed to create a config file")
	}

	cfg.General.Editor = os.Getenv("EDITOR")
	if cfg.General.Editor == "" && runtime.GOOS != "windows" {
		cfg.General.Editor = "vim"
	}
	cfg.General.Column = 40
	cfg.General.SelectCmd = "fzf"
	cfg.General.Backend = "gist"

	cfg.Gist.FileName = "pet-snippet.toml"

	cfg.GitLab.FileName = "pet-snippet.toml"
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
