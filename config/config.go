package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

// Conf is global config variable
var Conf Config

// Config is a struct of config
type Config struct {
	General GeneralConfig
	Gist    GistConfig
	GitLab  GitLabConfig
	GHEGist GHEGistConfig
}

// GeneralConfig is a struct of general config
type GeneralConfig struct {
	SnippetFile string
	SnippetDirs []string
	Editor      string
	Column      int
	SelectCmd   string
	Backend     string
	SortBy      string
	Color       bool
	Format      string
	Cmd         []string
}

// GistConfig is a struct of config for Gist
type GistConfig struct {
	FileName    string `toml:"file_name"`
	AccessToken string `toml:"access_token"`
	GistID      string `toml:"gist_id"`
	Public      bool
	AutoSync    bool `toml:"auto_sync"`
}

// GitLabConfig is a struct of config for GitLabSnippet
type GitLabConfig struct {
	FileName    string `toml:"file_name"`
	AccessToken string `toml:"access_token"`
	Url         string
	ID          string
	Visibility  string
	AutoSync    bool `toml:"auto_sync"`
	SkipSsl     bool `toml:"skip_ssl"`
}

// GHEGistConfig is a struct of config for Gist of Github Enterprise
type GHEGistConfig struct {
	BaseUrl     string `toml:"base_url"`
	UploadUrl   string `toml:"upload_url"`
	FileName    string `toml:"file_name"`
	AccessToken string `toml:"access_token"`
	GistID      string `toml:"gist_id"`
	Public      bool
	AutoSync    bool `toml:"auto_sync"`
}

// Flag is global flag variable
var Flag FlagConfig

// FlagConfig is a struct of flag
type FlagConfig struct {
	Debug        bool
	Query        string
	FilterTag    string
	Command      bool
	Delimiter    string
	OneLine      bool
	Color        bool
	Tag          bool
	UseMultiLine bool
	UseEditor    bool
}

// Load loads a config toml
func (cfg *Config) Load(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		f, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		err = toml.Unmarshal(f, cfg)
		if err != nil {
			return err
		}
		var snippetdirs []string
		for _, dir := range cfg.General.SnippetDirs {
			if !strings.HasSuffix(dir, "/") {
				dir = dir + "/"
			}

			snippetdirs = append(snippetdirs, dir) // note the = instead of :=
		}
		cfg.General.SnippetDirs = snippetdirs
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
	file_path, err := ExpandPath(cfg.General.SnippetFile)
	if err != nil {
		return errors.Wrap(err, "SnippetFile path is invalid: %v")
	}

	_, err = os.Create(file_path)
	if err != nil {
		return errors.Wrap(err, "Failed to create a snippet file")
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
	cfg.General.SelectCmd = "fzf --ansi --layout=reverse --border --height=90% --pointer=* --cycle --prompt=Snippets:"
	cfg.General.Backend = "gist"
	cfg.General.Color = false
	cfg.General.Format = "[$description]: $command $tags"

	cfg.Gist.FileName = "pet-snippet.toml"

	cfg.GitLab.FileName = "pet-snippet.toml"
	cfg.GitLab.Visibility = "private"

	return toml.NewEncoder(f).Encode(cfg)
}

// GetDefaultConfigDir returns the default config directory
func GetDefaultConfigDir() (dir string, err error) {
	if env, ok := os.LookupEnv("PET_CONFIG_DIR"); ok {
		dir = env
	} else if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "pet")
		}
		dir = filepath.Join(dir, "pet")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "pet")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("cannot create directory: %v", err)
	}
	return dir, nil
}

// Given a path to either a file or directory, returns its absolute path format.
// ExpandPath resolves "~/" prefix in a given system path.
// Raise error if path is an empty string as it
func ExpandPath(path string) (string, error) {
	if path == "" {
		error := errors.New("path to file/directory is not set.")
		return path, error
	}

	if strings.HasPrefix(path, "~") && os.IsPathSeparator(path[1]) {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return path, err
		}

		relativePath := path[2:]
		return filepath.Join(homedir, relativePath), nil
	}

	return path, nil
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
