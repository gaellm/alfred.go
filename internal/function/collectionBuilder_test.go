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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateFunctionCollectionFromFolder tests the CreateFunctionCollectionFromFolder function
func TestCreateFunctionCollectionFromFolder(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create mock JavaScript files in the temporary directory
	file1Path := filepath.Join(tempDir, "file1.js")
	t.Log("file1Path: ", file1Path)
	file2Path := filepath.Join(tempDir, "file2.js")
	err := os.WriteFile(file1Path, []byte("//content1"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(file2Path, []byte("//content2"), 0644)
	assert.NoError(t, err)

	// Call the function under test
	functionCollection, err := CreateFunctionCollectionFromFolder(tempDir)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, functionCollection, 2)

	// Validate the first function
	assert.Equal(t, "file1.js", functionCollection[0].FileName)
	assert.Equal(t, "//content1", functionCollection[0].FileContent)
	assert.False(t, functionCollection[0].HasFuncUpdateHelpers)
	assert.False(t, functionCollection[0].HasFuncAlfred)
	assert.False(t, functionCollection[0].HasFuncSetup)

	// Validate the second function
	assert.Equal(t, "file2.js", functionCollection[1].FileName)
	assert.Equal(t, "//content2", functionCollection[1].FileContent)
	assert.False(t, functionCollection[1].HasFuncUpdateHelpers)
	assert.False(t, functionCollection[1].HasFuncAlfred)
	assert.False(t, functionCollection[1].HasFuncSetup)
}

func TestCreateFunctionCollectionFromFolder_FileError(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create a mock JavaScript file with no read permissions
	filePath := filepath.Join(tempDir, "file1.js")
	err := os.WriteFile(filePath, []byte("content1"), 0000) // No permissions
	assert.NoError(t, err)

	// Call the function under test
	functionCollection, err := CreateFunctionCollectionFromFolder(tempDir)

	// Assertions
	assert.Error(t, err)
	assert.Len(t, functionCollection, 0)
}

func TestCreateFunctionCollectionFromFolder_NoFiles(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Call the function under test
	functionCollection, err := CreateFunctionCollectionFromFolder(tempDir)

	// Assertions
	assert.Error(t, err)
	assert.Len(t, functionCollection, 0)
}
