package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleIncr(s *session.Session, r *protocol.Request) {
	key := r.Args[0]

	e, err := s.Store.Get(key)
	if err != nil {
		s.Store.Set(key, "1", nil)
		s.SendInteger(1)
		return
	}

	curr, err := strconv.Atoi(e.Value)
	if err != nil {
		s.SendError("value is not an integer or out of range")
		return
	}

	s.Store.Update(key, strconv.Itoa(curr+1))
	s.SendInteger(curr + 1)
}
