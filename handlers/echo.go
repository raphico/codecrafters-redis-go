package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleEcho(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 1 {
		return protocol.NewErrorResponse("wrong number of arguments for 'echo' command")
	}

	msg := r.Args[0]
	return protocol.NewSimpleStringResponse(msg)
}
