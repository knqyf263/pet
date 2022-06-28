package envvar

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/knqyf263/pet/config"
)

type EnvVar struct {
	EnvVars []EnvVarInfo `toml:"env"`
}

type EnvVarInfo struct {
	Description string   `toml:"description"`
	Variables   []string `toml:"variables"`
	Tag         []string `toml:"tag"`
}

func (envvar *EnvVarInfo) GetVariables() []string {
	var variables []string
	for _, variable := range envvar.Variables {
		variables = append(variables, strings.Split(variable, "=")[0])
	}
	return variables
}

// Load reads toml file.
func (envvars *EnvVar) Load() error {
	configFile := config.Conf.General.SnippetFile
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil
	}
	if _, err := toml.DecodeFile(configFile, envvars); err != nil {
		return fmt.Errorf("failed to load snippet file. %v", err)
	}
	envvars.Order()
	return nil
}

// Save saves the snippets to toml file.
func (envvar *EnvVar) Save() error {
	snippetFile := config.Conf.General.SnippetFile
	f, err := os.Create(snippetFile)
	if err != nil {
		return fmt.Errorf("failed to save snippet file. err: %s", err)
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(envvar)
}

// ToString returns the contents of toml file.
func (snippets *EnvVar) ToString() (string, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(snippets)
	if err != nil {
		return "", fmt.Errorf("failed to convert struct to TOML string: %v", err)
	}
	return buffer.String(), nil
}

// Order snippets regarding SortBy option defined in config toml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (envvars *EnvVar) Order() {
	sortBy := config.Conf.General.SortBy
	switch {
	case sortBy == "description" || sortBy == "+description":
		sort.Sort(ByDescription(envvars.EnvVars))
	case sortBy == "-description":
		sort.Sort(sort.Reverse(ByDescription(envvars.EnvVars)))

	case sortBy == "-recency":
		envvars.reverse()
	}
}

func (envvars *EnvVar) reverse() {
	for i, j := 0, len(envvars.EnvVars)-1; i < j; i, j = i+1, j-1 {
		envvars.EnvVars[i], envvars.EnvVars[j] = envvars.EnvVars[j], envvars.EnvVars[i]
	}
}

type ByDescription []EnvVarInfo

func (a ByDescription) Len() int           { return len(a) }
func (a ByDescription) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool { return a[i].Description > a[j].Description }
