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
	"path/filepath"
)

// Find files using a pattern and a target directory.
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
		return files, errors.New("no files or directory " + targetDir)
	}

	return files, nil
}

func FindAllFiles(targetDir string, pattern ...string) ([]string, error) {

	//This is the files slice
	var files []string

	err := filepath.Walk(targetDir,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() {

				path += "/"

				for _, v := range pattern {
					matches, err := filepath.Glob(path + v)
					if err != nil {
						return err
					}

					if len(matches) != 0 {
						files = append(files, matches...)
					}
				}
			}

			return nil
		})
	if err != nil {
		return files, err
	}

	if len(files) < 1 {
		return files, errors.New("no files or directory " + targetDir)
	}

	return files, nil
}
