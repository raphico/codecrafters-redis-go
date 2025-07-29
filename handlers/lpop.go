package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleLPOP(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 1 {
		return protocol.NewErrorResponse("wrong number of arguments for 'lpop' command")
	}

	key := r.Args[0]

	e, err := s.Store.Get(key)
	if err != nil {
		return protocol.NewNullBulkStringResponse()
	}

	if e.Kind != store.ListType {
		return protocol.NewErrorResponse("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	copy := append([]string{}, e.Value.([]string)...)
	popped := copy[0]
	curr := copy[1:]

	if err := s.Store.Update(key, curr); err != nil {
		return protocol.NewNullBulkStringResponse()
	}

	return protocol.NewBulkStringResponse(popped)
}
