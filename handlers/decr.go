package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleDecr(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 1 {
		return protocol.NewErrorResponse("wrong number of arguments for 'decr' command")
	}

	key := r.Args[0]

	e, err := s.Store.Get(key)
	if err != nil {
		s.Store.Set(key, store.StringType, "-1", nil)
		return protocol.NewIntegerResponse(-1)
	}

	if e.Kind != store.StringType {
		return protocol.NewErrorResponse("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	strVal, ok := e.Value.(string)
	if !ok {
		panic("unexpected type: value is not a string")
	}

	curr, err := strconv.Atoi(strVal)
	if err != nil {
		return protocol.NewErrorResponse("value is not an integer or out of range")
	}

	s.Store.Update(key, strconv.Itoa(curr-1))
	return protocol.NewIntegerResponse(curr - 1)
}
