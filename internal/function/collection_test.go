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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionCollection_GetJs(t *testing.T) {
	// Create a sample FunctionCollection
	collection := FunctionCollection{
		{FileName: "file1.js", FileContent: "//content1"},
		{FileName: "file2.js", FileContent: "//content2"},
	}

	// Test retrieving existing file content
	content := collection.GetJs("file1.js")
	assert.Equal(t, "//content1", content)

	// Test retrieving non-existing file content
	content = collection.GetJs("file3.js")
	assert.Equal(t, "", content)
}

func TestFunctionCollection_IsEmpty(t *testing.T) {
	// Test with an empty collection
	emptyCollection := FunctionCollection{}
	assert.True(t, emptyCollection.IsEmpty())

	// Test with a non-empty collection
	nonEmptyCollection := FunctionCollection{
		{FileName: "file1.js", FileContent: "content1"},
	}
	assert.False(t, nonEmptyCollection.IsEmpty())
}

func TestFunctionCollection_GetSetupFunctions(t *testing.T) {
	// Create a sample FunctionCollection
	collection := FunctionCollection{
		{FileName: "file1.js", FileContent: "//content1", HasFuncSetup: true},
		{FileName: "file2.js", FileContent: "//content2", HasFuncSetup: false},
		{FileName: "file3.js", FileContent: "//content3", HasFuncSetup: true},
	}

	// Test retrieving setup functions
	setupFunctions := collection.GetSetupFunctions()
	assert.Len(t, setupFunctions, 2)
	assert.Equal(t, "file1.js", setupFunctions[0].FileName)
	assert.Equal(t, "file3.js", setupFunctions[1].FileName)
}

func TestFunctionCollection_GetFunction(t *testing.T) {
	// Create a sample FunctionCollection
	collection := FunctionCollection{
		{FileName: "file1.js", FileContent: "//content1"},
		{FileName: "file2.js", FileContent: "//content2"},
	}

	// Test retrieving an existing function
	function, err := collection.GetFunction("file1.js")
	assert.NoError(t, err)
	assert.Equal(t, "file1.js", function.FileName)
	assert.Equal(t, "//content1", function.FileContent)

	// Test retrieving a non-existing function
	function, err = collection.GetFunction("file3.js")
	assert.Error(t, err)
	assert.Equal(t, Function{}, function)
	assert.EqualError(t, err, "no fonction files named file3.js")
}
