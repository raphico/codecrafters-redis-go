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

var allowedInSubscribeMode = map[string]bool{
	"SUBSCRIBE":    true,
	"UNSUBSCRIBE":  true,
	"PSUBSCRIBE":   true,
	"PUNSUBSCRIBE": true,
	"PING":         true,
	"QUIT":         true,
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
	if s.InSubscribeMode() && !allowedInSubscribeMode[strings.ToUpper(r.Command)] {
		s.SendResponse(protocol.NewErrorResponse("Can't execute 'echo': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context"))
		return
	}

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
		panic("unexpected error: all queued commands should be valid")
	}

	return handler(s, r)
}
