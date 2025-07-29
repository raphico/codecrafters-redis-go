package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleLRANGE(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 3 {
		return protocol.NewErrorResponse("wrong number of arguments for 'lrange' command")
	}

	key := r.Args[0]

	start, err := strconv.Atoi(r.Args[1])
	if err != nil {
		return protocol.NewErrorResponse("value is not an integer or out of range")
	}

	stop, err := strconv.Atoi(r.Args[2])
	if err != nil {
		return protocol.NewErrorResponse("value is not an integer or out of range")
	}

	e, err := s.Store.Get(key)
	if err != nil {
		return protocol.NewArrayResponse([]protocol.Response{})
	}

	if e.Kind != store.ListType {
		return protocol.NewErrorResponse("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	var resp []protocol.Response
	for i := start; i <= stop; i++ {
		resp = append(resp, protocol.NewBulkStringResponse(e.Value.([]string)[i]))
	}

	return protocol.NewArrayResponse(resp)

}
