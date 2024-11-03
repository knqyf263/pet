package snippet

import (
	"os"
	"path/filepath"
	"regexp"
)

var tomlRegEx = regexp.MustCompile("^.+\\.(toml)$")

func getFiles(path string) (fileList []string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err == nil && tomlRegEx.MatchString(f.Name()) {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return fileList
}
