package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleIncr(s *session.Session, r *protocol.Request) protocol.Response {
	key := r.Args[0]

	e,  err := s.Store.Get(key)
	if err != nil {
		s.Store.Set(key, "1", nil)
		return protocol.NewIntegerResponse(1)
	}

	value, ok := e.Value.(string)
	if !ok {
		return protocol.NewErrorResponse("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	curr, err := strconv.Atoi(value)
	if err != nil {
		return protocol.NewErrorResponse("value is not an integer or out of range")
	}

	s.Store.Update(key, strconv.Itoa(curr+1))
	return protocol.NewIntegerResponse(curr + 1)
}
