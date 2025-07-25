package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleExec(s *session.Session, r *protocol.Request) {
	if !s.TxnContext.InTransaction() {
		s.SendError("EXEC without MULTI")
		return
	}

	s.TxnContext.EndTransaction()
	s.SendArray()
}
