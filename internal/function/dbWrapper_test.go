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
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

func TestGetDbWrapper(t *testing.T) {
	wrappers := getDbWrapper()

	assert.Len(t, wrappers, 4)
	assert.Equal(t, "dbGet", wrappers[0].name)
	assert.Equal(t, "dbSet", wrappers[1].name)
	assert.Equal(t, "dbDelete", wrappers[2].name)
	assert.Equal(t, "dbLoadFile", wrappers[3].name)
}

func TestGetDBManager(t *testing.T) {
	dbManager, err := getDBManager()

	assert.NoError(t, err)
	assert.NotNil(t, dbManager)

	// Ensure the database path is created
	assert.DirExists(t, dbManager.GetDBPath())
}

func TestDbSetGojaFunction(t *testing.T) {
	vm := goja.New()
	dbSet := dbSetGojaFunction(vm)

	// Call dbSet with a key and value
	dbSet(goja.FunctionCall{
		Arguments: []goja.Value{
			vm.ToValue("key1"),
			vm.ToValue("value1"),
		},
	})

	// Verify the value was set
	dbManager, _ := getDBManager()
	value, err := dbManager.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value)
}

func TestDbGetGojaFunction(t *testing.T) {
	vm := goja.New()
	dbSet := dbSetGojaFunction(vm)
	dbGet := dbGetGojaFunction(vm)

	// Set a key-value pair
	dbSet(goja.FunctionCall{
		Arguments: []goja.Value{
			vm.ToValue("key1"),
			vm.ToValue("value1"),
		},
	})

	// Get the value
	result := dbGet(goja.FunctionCall{
		Arguments: []goja.Value{
			vm.ToValue("key1"),
		},
	})

	assert.Equal(t, "value1", result.String())

	// Test getting a non-existent key
	result = dbGet(goja.FunctionCall{
		Arguments: []goja.Value{
			vm.ToValue("nonexistent"),
		},
	})

	assert.True(t, goja.IsUndefined(result))
}

func TestDbDeleteGojaFunction(t *testing.T) {
	vm := goja.New()
	dbSet := dbSetGojaFunction(vm)
	dbDelete := dbDeleteGojaFunction(vm)

	// Set a key-value pair
	dbSet(goja.FunctionCall{
		Arguments: []goja.Value{
			vm.ToValue("key1"),
			vm.ToValue("value1"),
		},
	})

	// Delete the key
	dbDelete(goja.FunctionCall{
		Arguments: []goja.Value{
			vm.ToValue("key1"),
		},
	})

	// Verify the key was deleted
	dbManager, _ := getDBManager()
	_, err := dbManager.Get("key1")
	assert.Error(t, err)
}

func TestDbLoadFileGojaFunction(t *testing.T) {
	vm := goja.New()
	dbLoadFile := dbLoadFileGojaFunction(vm)

	// Create a temporary JSON file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "data.json")
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	content, _ := json.Marshal(data)
	err := os.WriteFile(filePath, content, 0644)
	assert.NoError(t, err)

	// Load the file into the database
	dbLoadFile(goja.FunctionCall{
		Arguments: []goja.Value{
			vm.ToValue(filePath),
		},
	})

	// Verify the data was loaded
	dbManager, _ := getDBManager()
	value1, err := dbManager.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value1)

	value2, err := dbManager.Get("key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", value2)
}
