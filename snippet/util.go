package snippet

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/knqyf263/pet/path"
)

func getFiles(dir string) (fileList []string) {
	tomlRegEx, err := regexp.Compile(`^.+\.(toml)$`)
	if err != nil {
		log.Fatal(err)
	}

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
