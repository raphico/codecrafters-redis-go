package handlers

import (
	"encoding/base64"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

const emptyRdbFileInBase64 = "UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog=="

func HandlePSYNC(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 2 {
		return protocol.NewErrorResponse("wrong number of arguments for 'PSYNC' command")
	}

	info := s.Repl.View.Snapshot()
	resp := fmt.Sprintf("FULLRESYNC %s %d", info.MasterReplID, info.MasterOffset)

	rdbFile, err := base64.StdEncoding.DecodeString(emptyRdbFileInBase64)
	if err != nil {
		s.CloseConnection()
	}

	s.SendResponse(protocol.NewSimpleStringResponse(resp))
	return protocol.NewRawBytesResponse(rdbFile)
}
