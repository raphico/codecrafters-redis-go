package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleMulti(s *session.Session, r *protocol.Request) {
	if s.TxnContext.InTransaction() {
		s.SendError("MULTI calls can not be nested")
		return
	}

	s.TxnContext.BeginTransaction()
	s.SendSimpleString("OK")
}
