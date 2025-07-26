package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleGet(s *session.Session, r *protocol.Request) {
	if len(r.Args) != 1 {
		s.SendError("wrong number of arguments for 'get' command")
		return
	}

	key := r.Args[0]

	entry, err := s.Store.Get(key)
	if err != nil {
		s.SendNullBulkString()
		return
	}

	s.SendBulkString(entry.Value)
}
