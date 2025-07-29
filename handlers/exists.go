package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleExists(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) < 1 {
		return protocol.NewErrorResponse("wrong number of arguments for 'exists' command")
	}

	count := 0
	for _, key := range r.Args {
		if _, err := s.Store.Get(key); err == nil {
			count++
		}
	}

	return protocol.NewIntegerResponse(count)
}
