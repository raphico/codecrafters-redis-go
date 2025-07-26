package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandlePing(s *session.Session, r *protocol.Request) {
	if len(r.Args) > 1 {
		s.SendError("wrong number of arguments for 'ping' command")
		return
	}

	if len(r.Args) == 1 {
		msg := r.Args[0]
		s.SendSimpleString(msg)
	}

	s.SendSimpleString("PONG")
}
