package store

import (
	"fmt"
	"time"
)

type Store map[string]Entry

var MemStore = Store{}

func (s Store) Get(key string) (*Entry, error) {
	e, ok := s[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' does not exist", key)
	}

	return &e, nil
}

func (s Store) Set(key, value string, ttl *time.Duration) {
	var exp *time.Time
	if ttl != nil {
		t := time.Now().Add(*ttl)
		exp = &t
	}

	s[key] = Entry{Value: value, ExpiryTime: exp}
}

func (s Store) Update(key, value string) error {
	e, ok := s[key]
	if !ok {
		return fmt.Errorf("key '%s' does not exist", key)
	}

	s[key] = Entry{Value: value, ExpiryTime: e.ExpiryTime}

	return nil
}

func (s Store) Delete(key string) {
	delete(s, key)
}
