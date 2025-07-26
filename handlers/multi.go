package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleMulti(s *session.Session, r *protocol.Request) {
	s.SendSimpleString("OK")
}