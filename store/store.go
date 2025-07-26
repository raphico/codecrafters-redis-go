package store

import "time"

type Store map[string]Entry

var MemStore = Store{}

func (s Store) Get(key string) (Entry, bool) {
	e, ok := s[key]
	return e, ok
}

func (s Store) Set(key string, value string, ttl *time.Duration) {
	var exp *time.Time
	if ttl != nil {
		t := time.Now().Add(*ttl)
		exp = &t
	}

	s[key] = Entry{Value: value, ExpiryTime: exp}
}

func (s Store) Delete(key string) {
	delete(s, key)
}
