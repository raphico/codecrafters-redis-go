package handlers

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleBlpop(s *session.Session, r *protocol.Request) protocol.Response {
	argsLen := len(r.Args)

	if argsLen != 2 {
		return protocol.NewErrorResponse("wrong number of arguments for 'blpop' command")
	}

	key := r.Args[0]

	if _, err := s.Store.Get(key); err == nil {
		resp := HandleLpop(s, &protocol.Request{Command: "LPOP", Args: []string{key}})
		if resp.Type != protocol.NullBulkStringType {
			return protocol.NewArrayResponse([]protocol.Response{
				protocol.NewBulkStringResponse(key),
				resp,
			})
		}
	}

	timeoutFloat, err := strconv.ParseFloat(r.Args[argsLen-1], 64)
	if err != nil || timeoutFloat < 0 {
		return protocol.NewErrorResponse("timeout is not a float or out of range")
	}

	timeout := time.Duration(timeoutFloat * float64(time.Second))

	// a buffered channel to prevent deadlock even if the client already timed out or disconnected
	popSignalChan := make(chan struct{}, 1)

	s.Store.RegisterListWaiter(key, popSignalChan)

	// recheck the list before waiting
	if entry, err := s.Store.Get(key); err == nil {
		if list, ok := entry.Value.([]string); ok && len(list) > 0 {
			// clear a signal that might have been sent by RPUSH
			select {
			case <-popSignalChan:
			default:
			}
			return protocol.NewArrayResponse([]protocol.Response{
				protocol.NewBulkStringResponse(key),
				HandleLpop(s, &protocol.Request{Command: "LPOP", Args: []string{key}}),
			})
		}
	}

	if timeout == 0 {
		<-popSignalChan
		return protocol.NewArrayResponse([]protocol.Response{
			protocol.NewBulkStringResponse(key),
			HandleLpop(s, &protocol.Request{Command: "LPOP", Args: []string{key}}),
		})
	} else {
		select {
		case <-popSignalChan:
			return protocol.NewArrayResponse([]protocol.Response{
				protocol.NewBulkStringResponse(key),
				HandleLpop(s, &protocol.Request{Command: "LPOP", Args: []string{key}}),
			})
		case <-time.After(timeout):
			return protocol.NewNullArrayResponse()
		}
	}
}
