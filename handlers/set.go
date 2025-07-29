package handlers

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleSet(s *session.Session, r *protocol.Request) protocol.Response {
	if !isSetArgsValid(r.Args) {
		return protocol.NewErrorResponse("syntax error")
	}

	key, value := r.Args[0], r.Args[1]
	if len(r.Args) == 2 {
		s.Store.Set(key, store.StringType, value, nil)
		return protocol.NewSimpleStringResponse("OK")
	}

	// this is already checked in isSetArgsValid
	ms, _ := strconv.Atoi(r.Args[3])
	ttl := time.Duration(ms) * time.Millisecond

	s.Store.Set(key, store.StringType, value, &ttl)

	return protocol.NewSimpleStringResponse("OK")
}

func isSetArgsValid(args []string) bool {
	if len(args) < 2 {
		return false
	}

	if len(args) == 2 {
		return true
	}

	seenFlags := args[:2]

	i := 2
	for i < len(args) {
		arg := strings.ToUpper(args[i])

		if arg == "PX" {
			if i+1 > len(args) {
				return false
			}

			if slices.Contains(seenFlags, arg) {
				return false
			}

			seenFlags = append(seenFlags, arg)

			if _, err := strconv.Atoi(args[i+1]); err != nil {
				return false
			}

			i += 2
		} else {
			return false
		}
	}

	return true
}
