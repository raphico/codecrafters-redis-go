package registry

import (
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
