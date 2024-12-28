package storage

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdevsin/keyforge/internal/logger"
)

// Setup function for tests
func setupTestDB(tb testing.TB) *PebbleDB {
	if tb != nil {
		tb.Helper()
	}
	dbPath := path.Join(".", "testDb")
	err := os.RemoveAll(dbPath) // Ensure the database directory is clean
	if tb != nil {
		assert.NoError(tb, err, "Failed to clean up database directory")
	} else if err != nil {
		log.Fatalf("Failed to clean up database directory: %v", err)
	}

	db := GetDatabaseInstance(logger.GetLogger(), dbPath)
	if tb != nil {
		assert.NoError(tb, err, "Failed to open database")
	} else if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	return db
}

// Teardown function for tests and benchmarks
func teardownTestDB(tb testing.TB, db *PebbleDB) {
	if tb != nil {
		tb.Helper()
	}

	err := db.Close()
	if tb != nil {
		assert.NoError(tb, err, "Failed to close database")
	} else if err != nil {
		log.Printf("Failed to close database: %v", err)
	}

	err = os.RemoveAll(path.Join(".", "testDb")) // Cleanup database directory
	if tb != nil {
		assert.NoError(tb, err, "Failed to remove database directory")
	} else if err != nil {
		log.Printf("Failed to remove database directory: %v", err)
	}
}

func TestPebbleDB(t *testing.T) {

	// Test WriteKey
	t.Run("WriteKey", func(t *testing.T) {
		pebbleDB := setupTestDB(t)
		defer teardownTestDB(t, pebbleDB)

		key := []byte("test-key")
		value := []byte("test-value")
		err := pebbleDB.WriteKey(key, value)
		assert.NoError(t, err, "Failed to write key")
	})

	// Test ReadKey
	t.Run("ReadKey", func(t *testing.T) {
		pebbleDB := setupTestDB(t)
		defer teardownTestDB(t, pebbleDB)

		// Setup: Write a key
		key := []byte("test-key")
		value := []byte("test-value")
		err := pebbleDB.WriteKey(key, value)
		assert.NoError(t, err, "Failed to write key")

		// Test: Read the key
		readValue, err := pebbleDB.ReadKey(key)
		assert.NoError(t, err, "Failed to read key")
		assert.Equal(t, string(value), string(readValue), "Value mismatch")
	})

	// Test DeleteKey
	t.Run("DeleteKey", func(t *testing.T) {
		pebbleDB := setupTestDB(t)
		defer teardownTestDB(t, pebbleDB)

		// Setup: Write a key
		key := []byte("test-key")
		value := []byte("test-value")
		err := pebbleDB.WriteKey(key, value)
		assert.NoError(t, err, "Failed to write key")

		// Test: Delete the key
		err = pebbleDB.DeleteKey(key)
		assert.NoError(t, err, "Failed to delete key")

		// Verify: Key should not exist
		_, err = pebbleDB.ReadKey(key)
		assert.Error(t, err, "Key should not exist after deletion")
	})
}
