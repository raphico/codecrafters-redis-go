package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleIncr(s *session.Session, r *protocol.Request) protocol.Response {
	key := r.Args[0]

	e, err := s.Store.Get(key)
	if err != nil {
		s.Store.Set(key, "1", nil)
		return protocol.NewIntegerResponse(1)
	}

	curr, err := strconv.Atoi(e.Value)
	if err != nil {
		return protocol.NewErrorResponse("value is not an integer or out of range")
	}

	s.Store.Update(key, strconv.Itoa(curr+1))
	return protocol.NewIntegerResponse(curr + 1)
}
