package handlers

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/session"
)

func HandleSet(s *session.Session, r *protocol.Request) {
	if !isSetArgsValid(r.Args) {
		s.SendError("syntax error")
		return
	}

	key, value := r.Args[0], r.Args[1]
	if len(r.Args) == 2 {
		s.Store.Set(key, value, nil)
		s.SendSimpleString("OK")
		return
	}

	// this is already checked in isSetArgsValid
	ms, _ := strconv.Atoi(r.Args[3])
	ttl := time.Duration(ms) * time.Millisecond

	s.Store.Set(key, value, &ttl)

	s.SendSimpleString("OK")
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
