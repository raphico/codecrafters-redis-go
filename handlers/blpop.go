package handlers

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleBLPOP(s *session.Session, r *protocol.Request) protocol.Response {
	argsLen := len(r.Args)

	if argsLen != 2 {
		return protocol.NewErrorResponse("wrong number of arguments for 'blpop' command")
	}

	key := r.Args[0]

	if _, err := s.Store.Get(key); err == nil {
		resp := HandleLPOP(s, &protocol.Request{Command: "LPOP", Args: []string{key}})
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

	timeout := time.Duration(timeoutFloat) * time.Second

	popSignalChan := make(chan bool)

	s.Store.RegisterListWaiter(key, popSignalChan)

	if timeout == 0 {
		<-popSignalChan
		return protocol.NewArrayResponse([]protocol.Response{
			protocol.NewBulkStringResponse(key),
			HandleLPOP(s, &protocol.Request{Command: "LPOP", Args: []string{key}}),
		})
	} else {
		select {
		case <-popSignalChan:
			return protocol.NewArrayResponse([]protocol.Response{
				protocol.NewBulkStringResponse(key),
				HandleLPOP(s, &protocol.Request{Command: "LPOP", Args: []string{key}}),
			})
		case <-time.After(timeout):
			return protocol.NewNullBulkStringResponse()
		}
	}
}
