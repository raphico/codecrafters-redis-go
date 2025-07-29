package handlers

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleTtl(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 1 {
		return protocol.NewErrorResponse("wrong number of arguments for 'ttl' command")
	}

	key := r.Args[0]

	e, err := s.Store.Get(key)
	if err != nil {
		return protocol.NewIntegerResponse(-2)
	}

	if e.ExpiryTime == nil {
		return protocol.NewIntegerResponse(-1)
	}

	ttl := int(time.Until(*e.ExpiryTime).Seconds())
	if ttl < 0 {
		return protocol.NewIntegerResponse(-2)
	}

	return protocol.NewIntegerResponse(ttl)
}
