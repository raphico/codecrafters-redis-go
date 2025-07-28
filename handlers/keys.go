package handlers

import (
	"path"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleKeys(s *session.Session, r *protocol.Request) protocol.Response {
	if len(r.Args) != 1 {
		return protocol.NewErrorResponse("wrong number of arguments for 'keys' command")
	}

	pattern := r.Args[0]
	keys := s.Store.Keys()

	var matches []protocol.Response
	for _, k := range keys {
		if matched, err := path.Match(pattern, k); err == nil && matched {
			matches = append(matches, protocol.NewBulkStringResponse(k))
		}
	}

	return protocol.NewArrayResponse(matches)
}
