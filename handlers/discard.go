package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleDiscard(s *session.Session, r *protocol.Request) protocol.Response {
	if !s.TxnContext.InTransaction() {
		return protocol.NewErrorResponse("DISCARD without MULTI")
	}

	s.TxnContext.EndTransaction()
	return protocol.NewSimpleStringResponse("OK")
}
