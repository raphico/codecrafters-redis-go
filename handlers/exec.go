package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/registry"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func MakeExecHandler(reg *registry.Registry) registry.Handler {
	return func(s *session.Session, r *protocol.Request) protocol.Response {
		responses := []protocol.Response{}

		if !s.TxnContext.InTransaction() {
			return protocol.NewErrorResponse("EXEC without MULTI")
		}

		if s.TxnContext.IsDirty() {
			s.TxnContext.EndTransaction()
			return protocol.NewErrorResponse("EXECABORT Transaction discarded because of previous errors.")
		}

		cmds := s.TxnContext.GetQueuedCommands()

		for _, c := range cmds {
			resp := reg.RunAndReturnResponse(
				s,
				&protocol.Request{
					Command: c.Name,
					Args:    c.Args,
				})

			responses = append(responses, resp)
		}

		s.TxnContext.EndTransaction()
		return protocol.NewArrayResponse(responses)
	}
}
