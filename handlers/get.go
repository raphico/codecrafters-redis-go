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

	entry, expired, err := s.Store.Get(key)
	if err != nil {
		return protocol.NewNullBulkStringResponse()
	}

	if expired {
		s.Store.Delete(key)
		return protocol.NewNullBulkStringResponse()
	}

	return protocol.NewBulkStringResponse(entry.Value)
}
