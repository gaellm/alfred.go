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

package mock

import (
	"alfred/internal/log"
	"alfred/pkg/files"
	"context"
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func getMocksFilesContent(filePaths []string) ([][]byte, error) {

	var filesContent [][]byte

	for _, filePath := range filePaths {

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return filesContent, err
		}

		filesContent = append(filesContent, fileContent)

	}

	return filesContent, nil
}

func CreateMockCollectionFromFolder(path string) MockCollection {

	//matches, err := files.FindFiles(path, "*.json")
	matches, err := files.FindAllFiles(path, "*.json")
	if err != nil {
		log.Error(context.Background(), "error during mocks load...", err)
		//panic("error during mocks load......" + err.Error())
	}

	filesContent, err := getMocksFilesContent(matches)
	if err != nil {
		log.Error(context.Background(), "Application Panic", errors.New("error during mocks content load..."+err.Error()))
		panic("error during mocks content load......" + err.Error())
	}

	var mockCollection MockCollection

	for _, fileContent := range filesContent {

		currentMock, err := BuildMockFromJson(fileContent)
		if err != nil {
			log.Error(context.Background(), "Error during mock build from json", err, zap.String("text-provided", string(fileContent)))
			panic(fmt.Errorf("fatal error, mock build from json: %w", err))
		}

		mockCollection.Mocks = append(mockCollection.Mocks, &currentMock)
	}

	return mockCollection
}
