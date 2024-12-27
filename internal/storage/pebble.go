package storage

import (
	"github.com/cockroachdb/pebble"
	"github.com/tdevsin/keyforge/internal/logger"
)

type PebbleDB struct {
	db *pebble.DB
}

func GetDatabaseInstance(logger *logger.Logger, path string) *PebbleDB {
	db, err := pebble.Open(path, &pebble.Options{
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}
	instance := &PebbleDB{db: db}
	return instance
}

// Close closes the Pebble database.
func (p *PebbleDB) Close() error {
	return p.db.Close()
}

// WriteKey writes a key-value pair to the Pebble database.
func (p *PebbleDB) WriteKey(key, value []byte) error {
	if err := p.db.Set(key, value, pebble.Sync); err != nil {
		return err
	}
	return nil
}

// ReadKey reads the value of a given key from the Pebble database.
func (p *PebbleDB) ReadKey(key []byte) ([]byte, error) {
	value, closer, err := p.db.Get(key)
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	return value, nil
}

// DeleteKey deletes a key-value pair from the Pebble database.
func (p *PebbleDB) DeleteKey(key []byte) error {
	if err := p.db.Delete(key, pebble.Sync); err != nil {
		return err
	}
	return nil
}
