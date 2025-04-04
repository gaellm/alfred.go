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

package function

import (
	"alfred/pkg/files"
	"os"
	"path"
)

func CreateFunctionCollectionFromFolder(path string) (FunctionCollection, error) {

	functionCollection := FunctionCollection{}

	matches, err := files.FindFiles(path, "*.js")
	if err != nil {
		return functionCollection, err
	}

	for _, path := range matches {

		fileName := getFileName(path)
		fileContent, err := getFileContent(path)
		if err != nil {
			return functionCollection, err
		}

		function, err := CreateFunction(fileName, fileContent)
		if err != nil {
			return functionCollection, err
		}

		functionCollection = append(functionCollection, function)
	}

	return functionCollection, nil

}

func getFileName(filePath string) string {

	return path.Base(filePath)

}

func getFileContent(filePath string) ([]byte, error) {

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fileContent, err
	}

	return fileContent, nil

}
