package snippet

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/knqyf263/pet/path"
)

var tomlRegEx = regexp.MustCompile(`^.+\.(toml)$`)

// getFiles returns a list of files in the specified directory.
func getFiles(dir string) (fileList []string) {
	absPath, err := path.NewAbsolutePath(dir)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(
		absPath.Get(),
		func(p string, f os.FileInfo, err error) error {
			if err == nil && tomlRegEx.MatchString(f.Name()) {
				fileList = append(fileList, p)
			}
			return nil
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return fileList
}
