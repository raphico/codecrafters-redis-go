package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/persistence"
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleSave(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 0 {
		return protocol.NewErrorResponse("wrong number of arguments for 'save' command")
	}

	err := persistence.SaveRDB(s.Config, s.Store)
	if err != nil {
		return protocol.NewErrorResponse("")
	}

	return protocol.NewSimpleStringResponse("OK")
}
