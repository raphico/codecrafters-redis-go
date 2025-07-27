package store

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	// to safely protect shared data from race conditions or crashes
	mu sync.RWMutex

	data map[string]Entry
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Entry),
	}
}

func (s *Store) Get(key string) (*Entry, error) {
	// allows multiple readers at the same time, but blocks if a writer is currently modifying
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' does not exist", key)
	}

	if e.IsExpired() {
		s.Delete(key)
		return nil, fmt.Errorf("key '%s' has expired", key)
	}

	return &e, nil
}

func (s *Store) Set(key, value string, ttl *time.Duration) {
	// blocks all other writers and readers
	s.mu.Lock()
	defer s.mu.Unlock()

	var exp *time.Time
	if ttl != nil {
		t := time.Now().Add(*ttl)
		exp = &t
	}

	s.data[key] = Entry{Value: value, ExpiryTime: exp}
}

func (s *Store) Update(key, value string) error {
	// blocks all other writers and readers
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.data[key]
	if !ok {
		return fmt.Errorf("key '%s' does not exist", key)
	}

	s.data[key] = Entry{Value: value, ExpiryTime: e.ExpiryTime}

	return nil
}

func (s *Store) Delete(key string) {
	// blocks all other writers and readers
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
