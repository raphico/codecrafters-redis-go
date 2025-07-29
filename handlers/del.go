package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleDel(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) == 0 {
		return protocol.NewErrorResponse("wrong number of arguments for 'del' command")
	}

	totalDeleted := 0
	for i := range len(r.Args) {
		totalDeleted += s.Store.Delete(r.Args[i])
	}

	return protocol.NewIntegerResponse(totalDeleted)
}
