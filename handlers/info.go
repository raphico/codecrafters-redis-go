package handlers

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleInfo(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) > 1 {
		return protocol.NewErrorResponse("syntax error")
	}

	if len(r.Args) == 0 || strings.ToLower(r.Args[0]) == "replication" {
		info := s.Repl.View.Snapshot()
		return protocol.NewBulkStringResponse(fmt.Sprintf(
			"# Replication\r\nrole:%s\r\nmaster_replid:%s\r\nmaster_repl_offset:%d",
			info.Role,
			info.MasterReplID,
			info.MasterOffset,
		))
	}

	return protocol.NewErrorResponse(fmt.Sprintf("unknown INFO section '%s'", r.Args[0]))
}
