package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleLpush(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) < 2 {
		return protocol.NewErrorResponse("wrong number of arguments for 'lpush' command")
	}

	key, values := r.Args[0], r.Args[1:]

	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	e, err := s.Store.Get(key)
	if err != nil {
		s.Store.Set(key, store.ListType, values, nil)
		return protocol.NewIntegerResponse(len(values))
	}

	if e.Kind != store.ListType {
		return protocol.NewErrorResponse("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	list, ok := e.Value.([]string)
	if !ok {
		panic("unexpected type: value is not a list")
	}

	// Clone to avoid shared slice mutation
	prevCopy := append([]string{}, list...)
	curr := append(values, prevCopy...)

	if err := s.Store.Update(key, curr); err != nil {
		s.Store.Set(key, store.ListType, values, nil)
		return protocol.NewIntegerResponse(len(values))
	}

	return protocol.NewIntegerResponse(len(curr))
}
