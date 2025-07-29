package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleRpush(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) < 2 {
		return protocol.NewErrorResponse("wrong number of arguments for 'rpush' command")
	}

	key, values := r.Args[0], r.Args[1:]

	e, err := s.Store.Get(key)
	if err != nil {
		s.Store.Set(key, store.ListType, values, nil)
		return protocol.NewIntegerResponse(len(values))
	}

	if e.Kind != store.ListType {
		return protocol.NewErrorResponse("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	// Clone to avoid shared slice mutation
	prevCopy := append([]string{}, e.Value.([]string)...)
	curr := append(prevCopy, values...)

	if err := s.Store.Update(key, curr); err != nil {
		s.Store.Set(key, store.ListType, values, nil)
		return protocol.NewIntegerResponse(len(values))
	}

	return protocol.NewIntegerResponse(len(curr))
}
