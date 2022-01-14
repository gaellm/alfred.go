package collection

import (
	"errors"
	"path/filepath"
)

func FindFiles(targetDir string, pattern ...string) ([]string, error) {

	var mockFiles []string

	for _, v := range pattern {
		matches, err := filepath.Glob(targetDir + v)
		if err != nil {
			return matches, errors.New("incorrect file path to mocks dir")
		}

		if len(matches) != 0 {
			//fmt.Println("Found : ", matches)
			mockFiles = append(mockFiles, matches...)
		}
	}

	if len(mockFiles) < 1 {
		return mockFiles, errors.New("no mock to use")
	}

	return mockFiles, nil
}
