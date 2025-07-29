package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleGet(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 1 {
		return protocol.NewErrorResponse("wrong number of arguments for 'get' command")
	}

	key := r.Args[0]

	entry, err := s.Store.Get(key)
	if err != nil {
		return protocol.NewNullBulkStringResponse()
	}

	value, ok := entry.Value.(string)
	if !ok {
		return protocol.NewErrorResponse("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return protocol.NewBulkStringResponse(value)
}
