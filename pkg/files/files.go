/*
 * Copyright The Alfred.go Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package files

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

// FindFiles finds files using a pattern and a target directory.
func FindFiles(targetDir string, pattern ...string) ([]string, error) {
	// This is the files slice
	var files []string

	// Normalize the target directory path
	targetDir = filepath.Clean(targetDir)

	for _, v := range pattern {
		// Use filepath.Join to construct the full path
		searchPattern := filepath.Join(targetDir, v)

		// Perform the glob search
		matches, err := filepath.Glob(searchPattern)
		if err != nil {
			return nil, errors.New("incorrect file path")
		}

		if len(matches) != 0 {
			files = append(files, matches...)
		}
	}

	if len(files) < 1 {
		return nil, errors.New("no files or directory " + targetDir)
	}

	return files, nil
}

// FindAllFiles recursively finds files in a directory and its subdirectories that match the given patterns.
func FindAllFiles(targetDir string, patterns ...string) ([]string, error) {
	// This is the files slice
	var files []string

	// Normalize the target directory path
	targetDir = filepath.Clean(targetDir)

	// Walk through the directory and its subdirectories
	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Return the error if there's an issue accessing a file or directory
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file matches any of the patterns
		for _, pattern := range patterns {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				return err // Return the error if the pattern is invalid
			}
			if matched {
				files = append(files, path)
				break // Stop checking other patterns if a match is found
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Return an error if no files were found
	if len(files) == 0 {
		return nil, errors.New("no files or directory " + targetDir)
	}

	return files, nil
}

func GetFileName(filePath string) string {

	return path.Base(filePath)

}

func GetFileContent(filePath string) ([]byte, error) {

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fileContent, err
	}

	return fileContent, nil

}
