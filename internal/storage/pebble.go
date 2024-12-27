package storage

import (
	"log"
	"sync"

	"github.com/cockroachdb/pebble"
)

type PebbleDB struct {
	db *pebble.DB
}

var (
	once     sync.Once
	instance *PebbleDB
	initErr  error
)

// InitializeDatabase initializes and opens the Pebble database only once.
func InitializeDatabase(path string) error {
	once.Do(func() {
		db, err := pebble.Open(path, &pebble.Options{})
		if err != nil {
			initErr = err
			return
		}
		instance = &PebbleDB{db: db}
	})
	return initErr
}

// GetDatabaseInstance provides the initialized PebbleDB instance.
// Ensure `InitializeDatabase` is called before using this function.
func GetDatabaseInstance() *PebbleDB {
	if instance == nil {
		log.Fatal("Database has not been initialized. Call InitializeDatabase first.")
	}
	return instance
}

// Close closes the Pebble database.
func (p *PebbleDB) Close() error {
	return p.db.Close()
}

// WriteKey writes a key-value pair to the Pebble database.
func (p *PebbleDB) WriteKey(key, value []byte) error {
	if err := p.db.Set(key, value, pebble.Sync); err != nil {
		log.Printf("Failed to write key %s: %v", key, err)
		return err
	}
	return nil
}

// ReadKey reads the value of a given key from the Pebble database.
func (p *PebbleDB) ReadKey(key []byte) ([]byte, error) {
	value, closer, err := p.db.Get(key)
	if err != nil {
		log.Printf("Failed to read key %s: %v", key, err)
		return nil, err
	}
	defer closer.Close()
	return value, nil
}

// DeleteKey deletes a key-value pair from the Pebble database.
func (p *PebbleDB) DeleteKey(key []byte) error {
	if err := p.db.Delete(key, pebble.Sync); err != nil {
		log.Printf("Failed to delete key %s: %v", key, err)
		return err
	}
	return nil
}
