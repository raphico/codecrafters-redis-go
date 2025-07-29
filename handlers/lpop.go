package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleLpop(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 1 && len(r.Args) != 2 {
		return protocol.NewErrorResponse("wrong number of arguments for 'lpop' command")
	}

	key := r.Args[0]
	count := 1
	if len(r.Args) == 2 {
		v, err := strconv.Atoi(r.Args[1])
		if err != nil {
			return protocol.NewErrorResponse("value is not an integer or out of range")
		}

		if v <= 0 {
			return protocol.NewErrorResponse("value is out of range, must be positive")
		}

		count = v
	}

	e, err := s.Store.Get(key)
	if err != nil {
		if count == 1 {
			return protocol.NewNullBulkStringResponse()
		}

		return protocol.NewArrayResponse([]protocol.Response{})
	}

	if e.Kind != store.ListType {
		return protocol.NewErrorResponse("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	list, ok := e.Value.([]string)
	if !ok {
		panic("unexpected type: value is not a list")
	}

	if count > len(list) {
		count = len(list)
	}

	popped := list[0:count]
	remaining := list[count:]

	// Redis automatically deletes a key when the list becomes empty
	if len(remaining) == 0 {
		s.Store.Delete(key)
	} else {
		if err := s.Store.Update(key, remaining); err != nil {
			if count == 1 {
				return protocol.NewNullBulkStringResponse()
			}

			return protocol.NewArrayResponse([]protocol.Response{})
		}
	}

	if count == 1 {
		return protocol.NewBulkStringResponse(popped[0])
	}

	var response []protocol.Response
	for _, v := range popped {
		response = append(response, protocol.NewBulkStringResponse(v))
	}

	return protocol.NewArrayResponse(response)
}
