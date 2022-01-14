package files

import (
	"errors"
	"path/filepath"
)

//Find files using a pattern and a target directory.
func FindFiles(targetDir string, pattern ...string) ([]string, error) {

	//This is the files slice
	var files []string

	for _, v := range pattern {
		matches, err := filepath.Glob(targetDir + v)
		if err != nil {
			return matches, errors.New("incorrect file path")
		}

		if len(matches) != 0 {
			//fmt.Println("Found : ", matches)
			files = append(files, matches...)
		}
	}

	if len(files) < 1 {
		return files, errors.New("no files")
	}

	return files, nil
}
