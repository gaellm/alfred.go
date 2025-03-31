package db

import (
	"log"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

type DBManager struct {
	db   *badger.DB
	once sync.Once
}

var instance *DBManager

// GetDBManager returns a singleton instance of DBManager
func GetDBManager() *DBManager {
	if instance == nil {
		instance = &DBManager{}
	}
	return instance
}

// Init initializes the Badger database
func (m *DBManager) Init(dbPath string) error {
	var err error
	m.once.Do(func() {
		opts := badger.DefaultOptions(dbPath).WithLogger(nil) // Disable Badger's internal logging
		m.db, err = badger.Open(opts)
		if err != nil {
			log.Fatalf("Failed to open Badger DB: %v", err)
		}
	})
	return err
}

// Close closes the Badger database
func (m *DBManager) Close() {
	if m.db != nil {
		m.db.Close()
	}
}

// Set sets a key-value pair in the database
func (m *DBManager) Set(key, value string) error {
	return m.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
}

// Get retrieves a value by key from the database
func (m *DBManager) Get(key string) (string, error) {
	var value string
	err := m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err // This will return badger.ErrKeyNotFound if the key is missing
		}
		return item.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
	})
	return value, err
}

// Delete removes a key-value pair from the database
func (m *DBManager) Delete(key string) error {
	return m.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// GetDB returns the underlying Badger DB instance
func (m *DBManager) GetDB() *badger.DB {
	return m.db
}
