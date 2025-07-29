package store

import "time"

type EntryType int

const (
	StringType EntryType = iota
	ListType
)

type Entry struct {
	Value      any
	Kind       EntryType
	ExpiryTime *time.Time
}

func (e *Entry) IsExpired() bool {
	if e.ExpiryTime == nil {
		return false
	}

	return time.Now().After(*e.ExpiryTime)
}
