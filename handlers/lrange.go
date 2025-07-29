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

	list, ok := e.Value.([]string)
	if !ok {
		panic("list is corrupted")
	}

	length := len(list)

	// normalize negative indices
	if start < 0 {
		start += length
	}

	if stop < 0 {
		stop += length
	}

	// negative index is out of range (i.e. >= the length of the list)
	if start < 0 {
		start = 0
	}

	if stop < 0 {
		stop = 0
	}

	if stop >= length {
		stop = length - 1
	}

	if start > stop || start >= length {
		return protocol.NewArrayResponse([]protocol.Response{})
	}

	var resp []protocol.Response
	for i := start; i <= stop; i++ {
		resp = append(resp, protocol.NewBulkStringResponse(list[i]))
	}

	return protocol.NewArrayResponse(resp)

}
