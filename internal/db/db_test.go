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
package db

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDBManager_Singleton(t *testing.T) {
	// Ensure GetDBManager always returns the same instance
	manager1 := GetDBManager()
	manager2 := GetDBManager()

	assert.Equal(t, manager1, manager2)
}

func TestDBManager_Init(t *testing.T) {
	// Create a temporary directory for the database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "badger")

	// Initialize the database
	manager := GetDBManager()
	err := manager.Init(dbPath)
	assert.NoError(t, err)

	// Ensure the database path is set correctly
	assert.Equal(t, dbPath, manager.GetDBPath())

}

func TestDBManager_SetAndGet(t *testing.T) {
	// Create a temporary directory for the database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "badger")

	// Initialize the database
	manager := GetDBManager()
	err := manager.Init(dbPath)
	assert.NoError(t, err)

	// Set a key-value pair
	err = manager.Set("key1", "value1")
	assert.NoError(t, err)

	// Retrieve the value by key
	value, err := manager.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value)

}

func TestDBManager_Get_NonExistentKey(t *testing.T) {
	// Create a temporary directory for the database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "badger")

	// Initialize the database
	manager := GetDBManager()
	err := manager.Init(dbPath)
	assert.NoError(t, err)

	// Attempt to retrieve a non-existent key
	_, err = manager.Get("nonexistent")
	assert.Error(t, err)

}

func TestDBManager_Delete(t *testing.T) {
	// Create a temporary directory for the database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "badger")

	// Initialize the database
	manager := GetDBManager()
	err := manager.Init(dbPath)
	assert.NoError(t, err)

	// Set a key-value pair
	err = manager.Set("key1", "value1")
	assert.NoError(t, err)

	// Delete the key-value pair
	err = manager.Delete("key1")
	assert.NoError(t, err)

	// Attempt to retrieve the deleted key
	_, err = manager.Get("key1")
	assert.Error(t, err)

}
