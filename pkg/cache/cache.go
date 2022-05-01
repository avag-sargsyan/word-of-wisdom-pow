package cache

import (
	"fmt"
	"sync"
)

// General struct which should be stored in memory
type Store struct {
	table map[int]Data
	sync.RWMutex
}

// Actual data which hold a necessary information
type Data struct {
	Secret     string
}

// Returns new memory store
func NewStore() *Store {
	return &Store{
		table: make(map[int]Data, 0),
	}
}

// Get the value associated with the key or an error in case the
// key was not found, or any other error encountered
func (s *Store) Get(key int) (string, error) {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.table[key]
	if !ok {
		return "", fmt.Errorf("has no entry for the given key: %d", key)
	}

	return v.Secret, nil
}

// Store for the given key
func (s *Store) Put(key int, secret string) {
	s.Lock()
	defer s.Unlock()
	s.table[key] = Data{
		Secret: secret,
	}
}

// Remove the data from the store for the given key
func (s *Store) Delete(key int) {
	s.Lock()
	defer s.Unlock()
	delete(s.table, key)
}
