package snippet

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/knqyf263/pet/config"
)

func getFiles(path string) (fileList []string) {
	tomlRegEx, err := regexp.Compile("^.+\\.(toml)$")
	if err != nil {
		log.Fatal(err)
	}

    expandedPath := config.Expand(path)
	err = filepath.Walk(
        expandedPath, 
        func(p string, f os.FileInfo, err error) error {
            if err == nil && tomlRegEx.MatchString(f.Name()) {
                fileList = append(fileList, p)
            }
            return nil
        },
    )

	if err != nil {
		panic(err)
	}

	return fileList
}
