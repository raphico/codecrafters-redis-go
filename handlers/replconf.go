package handlers

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleREPLCONF(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 2 {
		return protocol.NewErrorResponse("wrong number of arguments for 'REPLCONF' command")
	}

	option, _ := r.Args[0], r.Args[1]

	switch strings.ToLower(option) {
	case "listening-port":
		return protocol.NewBulkStringResponse("OK")
	case "capa":
		return protocol.NewBulkStringResponse("OK")
	default:
		return protocol.NewErrorResponse(fmt.Sprintf("Unknown REPLCONF option: %s", option))
	}
}
