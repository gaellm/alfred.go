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
	"alfred/internal/db"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dop251/goja"
	"github.com/google/uuid"
)

var (
	initDbOnce sync.Once
)

func getDbWrapper() []wrapper {
	return []wrapper{
		{
			name:  "dbGet",
			value: dbGetGojaFunction,
		},
		{
			name:  "dbSet",
			value: dbSetGojaFunction,
		},
		{
			name:  "dbDelete",
			value: dbDeleteGojaFunction,
		},
		{
			name:  "dbLoadFile",
			value: dbLoadFileGojaFunction,
		},
	}
}

func getDBManager() (*db.DBManager, error) {

	// Initialize the database manager
	dbManager := db.GetDBManager()
	var err error

	initDbOnce.Do(func() {

		// Generate a random UUID-based path for the database
		randomUUID := uuid.New().String()
		dbPath := filepath.Join(os.TempDir(), fmt.Sprintf("alfred_badger_%s", randomUUID))

		err = dbManager.Init(dbPath)
	})

	return dbManager, err
}

func dbSetGojaFunction(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(vm.ToValue("dbSet requires a key and a value as arguments"))
		}
		key := call.Argument(0).String()
		value := call.Argument(1).String()

		dbManager, err := getDBManager()
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}

		err = dbManager.Set(key, value)
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}
		return goja.Undefined()
	}
}

func dbGetGojaFunction(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("dbGet requires a key as an argument"))
		}
		key := call.Argument(0).String()

		dbManager, err := getDBManager()
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}

		value, err := dbManager.Get(key)
		if err != nil {
			// Return undefined if the key is not found
			return goja.Undefined()
		}
		return vm.ToValue(value)
	}
}

func dbDeleteGojaFunction(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {

		if len(call.Arguments) < 1 {
			panic(vm.ToValue("dbDelete requires a key as an argument"))
		}
		key := call.Argument(0).String()

		dbManager, err := getDBManager()
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}

		err = dbManager.Delete(key)
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}
		return goja.Undefined()
	}
}

func dbLoadFileGojaFunction(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {

		if len(call.Arguments) < 1 {
			panic(vm.ToValue("dbLoadFile requires a file path as an argument"))
		}
		filePath := call.Argument(0).String()

		// Read the file contents using os.ReadFile
		content, err := os.ReadFile(filePath)
		if err != nil {
			panic(vm.ToValue("Failed to read file: " + err.Error()))
		}

		// Parse the JSON into a map
		var data map[string]string
		err = json.Unmarshal(content, &data)
		if err != nil {
			panic(vm.ToValue("Failed to parse JSON file: " + err.Error()))
		}

		dbManager, err := getDBManager()
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}

		// Use a batch write for faster imports
		writeBatch := dbManager.GetDB().NewWriteBatch() // Get the Badger DB instance
		defer writeBatch.Cancel()                       // Ensure the batch is canceled if something goes wrong

		for key, value := range data {
			err = writeBatch.Set([]byte(key), []byte(value))
			if err != nil {
				panic(vm.ToValue("Failed to add key-value pair to batch: " + err.Error()))
			}
		}

		// Commit the batch
		err = writeBatch.Flush()
		if err != nil {
			panic(vm.ToValue("Failed to commit batch: " + err.Error()))
		}

		return vm.ToValue("File loaded successfully using batch writes")
	}
}
