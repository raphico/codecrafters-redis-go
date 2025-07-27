package handlers

import (
	"fmt"
	"slices"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/replication"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleInfo(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) > 1 {
		return protocol.NewErrorResponse("syntax error")
	}

	availableSections := []string{"replication"}
	if len(r.Args) != 0 && !slices.Contains(availableSections, r.Args[0]) {
		return protocol.NewErrorResponse(fmt.Sprintf("unknown INFO section '%s'", r.Args[0]))
	}

	section := handleReplicationSection(s.Repl.View.Snapshot())
	return protocol.NewBulkStringResponse(section)
}

func handleReplicationSection(info replication.Info) string {
	if info.Role == replication.RoleMaster {
		return fmt.Sprintf(
			"# Replication\r\nrole:%s\r\nmaster_replid:%s\r\nmaster_repl_offset:%d",
			info.Role,
			info.MasterReplID,
			info.MasterOffset,
		)
	}

	return fmt.Sprintf(
		"# Replication\r\nrole:%s\r\nmaster_host:%s\r\nmaster_post:%d",
		info.Role,
		info.MasterHost,
		info.MasterPort,
	)
}
