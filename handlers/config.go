package handlers

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleConfig(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) == 0 {
		return protocol.NewErrorResponse("wrong number of arguments for 'config' command")
	}

	subCmd := strings.ToLower(r.Args[0])
	switch subCmd {
	case "get":
		return handleConfigGet(s, r)
	default:
		return protocol.NewErrorResponse(fmt.Sprintf("unknown subcommand `%s`", subCmd))
	}
}

func handleConfigGet(s *session.Session, r *protocol.Request) protocol.Response  {
	if len(r.Args) != 2 {
		return protocol.NewErrorResponse("wrong number of arguments for 'config|get' command")
	}

	requestedKey := r.Args[1]
	configMap := s.Config.GetConfig()

	if value, ok := configMap[requestedKey]; ok {
		return protocol.NewArrayResponse([]protocol.Response{
			protocol.NewBulkStringResponse("dir"),
			protocol.NewBulkStringResponse(value),
		})
	}

	return protocol.NewArrayResponse([]protocol.Response{})
}