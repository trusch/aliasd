package storage

import (
	"errors"

	"github.com/trusch/storage"
)

// Storage is used to store aliases
type Storage struct {
	aliasMap map[string]string
	bucket   string
	db       storage.Storage
}

// NewStorage creates a new storage instance
func NewStorage(db storage.Storage, bucket string) (*Storage, error) {
	store := &Storage{make(map[string]string), bucket, db}
	return store, store.loadFromDB()
}

func (store *Storage) loadFromDB() error {
	if err := store.db.CreateBucket(store.bucket); err != nil {
		return err
	}
	ch, err := store.db.List(store.bucket, nil)
	if err != nil {
		return err
	}
	for info := range ch {
		store.aliasMap[info.Key] = string(info.Value)
	}
	return nil
}

// Get retrieves a entry
func (store *Storage) Get(key string) (string, error) {
	if alias, ok := store.aliasMap[key]; ok {
		return alias, nil
	}
	return "", errors.New("no such key")
}

// Set sets an entry and writes to db
func (store *Storage) Set(key, value string) error {
	store.aliasMap[key] = value
	return store.db.Put(store.bucket, key, []byte(value))
}

// GetAll returns the whole alias map
func (store *Storage) GetAll() map[string]string {
	return store.aliasMap
}

// Close closes the leveldb
func (store *Storage) Close() error {
	return store.db.Close()
}

func (store *Storage) Del(key string) error {
	delete(store.aliasMap, key)
	return store.db.Delete(store.bucket, key)
}

func (store *Storage) Has(key string) bool {
	_, ok := store.aliasMap[key]
	return ok
}
