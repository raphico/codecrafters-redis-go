package handlers

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandlePSYNC(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 2 {
		return protocol.NewErrorResponse("wrong number of arguments for 'PSYNC' command")
	}

	info := s.Repl.View.Snapshot()
	resp := fmt.Sprintf("FULLRESYNC %s %d", info.MasterReplID, info.MasterOffset)
	return protocol.NewSimpleStringResponse(resp)
}
