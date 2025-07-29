package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleType(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 1 {
		return protocol.NewErrorResponse("wrong number of arguments for 'type' command")
	}

	key := r.Args[0]

	e, err := s.Store.Get(key)
	if err != nil {
		return protocol.NewSimpleStringResponse("none")
	}

	var t string
	switch e.Kind {
	case store.StringType:
		t = "string"
	case store.ListType:
		t = "list"
	default:
		t = "none"
	}

	return protocol.NewSimpleStringResponse(t)
}
