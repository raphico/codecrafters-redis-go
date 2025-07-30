package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandlePublish(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 2 {
		return protocol.NewErrorResponse("wrong number of arguments for 'publish' command")
	}

	channel, message := r.Args[0], r.Args[1]

	subscribersCount := s.Pubsub.Publish(channel, message)
	return protocol.NewIntegerResponse(subscribersCount)
}
