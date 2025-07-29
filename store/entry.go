package store

import "time"

type Entry struct {
	Value      any
	ExpiryTime *time.Time
}

func (e *Entry) IsExpired() bool {
	if e.ExpiryTime == nil {
		return false
	}

	return time.Now().After(*e.ExpiryTime)
}
