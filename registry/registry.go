package registry

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

type Handler func(s *session.Session, r *protocol.Request)

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
	handler, ok := reg.handlers[canonical(r.Command)]
	if !ok {
		s.SendError("unknown command '" + r.Command + "'")
		if s.TxnContext.InTransaction() {
			s.TxnContext.MarkDirty()
		}
		return
	}

	if s.TxnContext.InTransaction() {
		s.TxnContext.QueueCommand(r.Command, r.Args)
		s.SendSimpleString("QUEUED")
		return
	}

	handler(s, r)
}
