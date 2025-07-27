package registry

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

type Handler func(s *session.Session, r *protocol.Request) protocol.Response

type Registry struct {
	handlers map[string]Handler
}

func New() *Registry {
	return &Registry{
		handlers: make(map[string]Handler),
	}
}

func canonical(command string) string {
	return strings.ToLower(command)
}

func (reg *Registry) Add(command string, handler Handler) {
	reg.handlers[canonical(command)] = handler
}

func (reg *Registry) Dispatch(s *session.Session, r *protocol.Request) {
	cmd := canonical(r.Command)
	handler, ok := reg.handlers[cmd]

	if s.TxnContext.InTransaction() && cmd != "exec" && cmd != "discard" {
		if !ok {
			s.TxnContext.MarkDirty()
		}

		s.TxnContext.QueueCommand(r.Command, r.Args)
		s.SendResponse(protocol.NewSimpleStringResponse("QUEUED"))
		return
	}

	if !ok {
		s.SendResponse(protocol.NewErrorResponse(fmt.Sprintf("unknown command '%s'", r.Command)))
		return
	}

	resp := handler(s, r)
	s.SendResponse(resp)
}

func (reg *Registry) RunAndReturnResponse(s *session.Session, r *protocol.Request) protocol.Response {
	handler, ok := reg.handlers[canonical(r.Command)]
	if !ok {
		panic("All queued commands should be valid")
	}

	return handler(s, r)
}
