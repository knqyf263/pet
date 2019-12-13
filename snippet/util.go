package snippet

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func getFiles(path string) (fileList []string) {
	tomlRegEx, err := regexp.Compile("^.+\\.(toml)$")
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
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
