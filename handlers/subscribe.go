package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleSubscribe(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 1 {
		return protocol.NewErrorResponse("wrong number of arguments for 'subscribe' command")
	}

	channel := r.Args[0]
	s.Pubsub.Subscribe(channel, s)
	count := s.Pubsub.GetSubscribedCount(s)

	return protocol.NewArrayResponse([]protocol.Response{
		protocol.NewBulkStringResponse("subscribe"),
		protocol.NewBulkStringResponse(channel),
		protocol.NewIntegerResponse(count),
	})
}
