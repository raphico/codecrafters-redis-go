package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleMulti(s *session.Session, r *protocol.Request) protocol.Response {
	if s.TxnContext.InTransaction() {
		return protocol.NewErrorResponse("MULTI calls can not be nested")
	}

	s.TxnContext.BeginTransaction()
	return protocol.NewSimpleStringResponse("OK")
}
