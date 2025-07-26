package registry

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/protocol"
)

type Handler func(w protocol.Response, r *protocol.Request)

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

func (reg *Registry) Dispatch(w protocol.Response, r *protocol.Request) {
	handler, ok := reg.handlers[canonical(r.Command)]
	if !ok {
		w.SendError("unknown command '" + r.Command + "'")
		return
	}

	handler(w, r)
}
