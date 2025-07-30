package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleUnsubscribe(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 1 {
		return protocol.NewErrorResponse("wrong number of arguments for 'unsubscribe' command")
	}

	channel := r.Args[0]
	count := s.Pubsub.Unsubscribe(channel, s)

	return protocol.NewArrayResponse([]protocol.Response{
		protocol.NewBulkStringResponse("unsubscribe"),
		protocol.NewBulkStringResponse(channel),
		protocol.NewIntegerResponse(count),
	})
}
