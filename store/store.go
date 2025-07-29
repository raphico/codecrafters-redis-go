package store

import (
	"fmt"
	"sync"
	"time"
)

type ListWaitChans map[string][]chan bool

type Store struct {
	// to safely protect shared data from race conditions or crashes
	mu sync.RWMutex

	data          map[string]Entry
	listWaitChans ListWaitChans
}

func NewStore() *Store {
	return &Store{
		data:          make(map[string]Entry),
		listWaitChans: make(ListWaitChans),
	}
}

func (s *Store) Get(key string) (*Entry, error) {
	// allows multiple readers at the same time, but blocks if a writer is currently modifying
	s.mu.RLock()
	e, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("key '%s' does not exist", key)
	}

	if e.IsExpired() {
		s.mu.Lock()
		delete(s.data, key)
		s.mu.Unlock()
		return nil, fmt.Errorf("key '%s' does not exist", key)
	}

	return &e, nil
}

func (s *Store) Set(key string, kind EntryType, value any, ttl *time.Duration) {
	s.mu.Lock()

	var exp *time.Time
	if ttl != nil {
		t := time.Now().Add(*ttl)
		exp = &t
	}

	s.data[key] = Entry{Value: value, Kind: kind, ExpiryTime: exp}
	s.mu.Unlock()

	if ch := s.notifyListWaiter(kind, key); ch != nil {
		select {
		case ch <- true:
		default:
		}
	}
}

func (s *Store) Update(key string, value any) error {
	// blocks all other writers and readers
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.data[key]
	if !ok {
		return fmt.Errorf("key '%s' does not exist", key)
	}

	if e.IsExpired() {
		delete(s.data, key)
		return fmt.Errorf("key '%s' does not exist", key)
	}

	s.data[key] = Entry{Value: value, Kind: e.Kind, ExpiryTime: e.ExpiryTime}

	return nil
}

func (s *Store) Delete(key string) int {
	// blocks all other writers and readers
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[key]
	if exists {
		delete(s.data, key)
		return 1
	}

	return 0
}

func (s *Store) Keys() []string {
	// allows multiple readers at the same time, but blocks if a writer is currently modifying
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string
	for k, e := range s.data {
		if !e.IsExpired() {
			keys = append(keys, k)
		}
	}

	return keys
}

func (s *Store) RegisterListWaiter(key string, ch chan bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.listWaitChans[key] = append(s.listWaitChans[key], ch)
}

func (s *Store) notifyListWaiter(kind EntryType, key string) chan bool {
	if kind != ListType {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	waiters, ok := s.listWaitChans[key]
	if !ok || len(waiters) == 0 {
		return nil
	}

	waiterChan := waiters[0]
	if len(waiters) == 1 {
		delete(s.listWaitChans, key)
	} else {
		s.listWaitChans[key] = waiters[1:]
	}

	return waiterChan
}
