package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleEcho(s *session.Session, r *protocol.Request) {
	if len(r.Args) != 1 {
		s.SendError("wrong number of arguments for 'echo' command")
		return
	}

	msg := r.Args[0]
	s.SendSimpleString(msg)
}
