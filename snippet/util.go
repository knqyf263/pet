package snippet


import (
	"os"
	"log"
	"regexp"
	"path/filepath"
)

func getFiles(path string) ([]string, error) {
	tomlRegEx, err := regexp.Compile("^.+\\.(toml)$")
	if err != nil {
		log.Fatal(err)
	}

	fileList := make([]string, 0)
	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err == nil && tomlRegEx.MatchString(f.Name()) {
			fileList = append(fileList, path)
		}
		return nil
	})
	
	if err != nil {
		panic(err)
	}

	return fileList, nil
}
