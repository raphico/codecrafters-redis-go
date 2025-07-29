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

	flag := strings.ToUpper(r.Args[2])
	expiryMs, _ := strconv.Atoi(r.Args[3]) // already validated in isSetArgsValid

	var ttl time.Duration
	if flag == "PX" {
		ttl = time.Duration(expiryMs) * time.Millisecond
	} else { // EX
		ttl = time.Duration(expiryMs) * time.Second
	}

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

	if len(args) != 4 {
		return false
	}

	flag := strings.ToUpper(args[2])
	if !slices.Contains([]string{"PX", "EX"}, flag) {
		return false
	}

	if _, err := strconv.Atoi(args[3]); err != nil {
		return false
	}

	return true
}
