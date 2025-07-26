package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
)

func HandlePing(w protocol.Response, r *protocol.Request) {
	if len(r.Args) > 1 {
		w.SendError("wrong number of arguments for 'ping' command")
		return
	}

	if len(r.Args) == 1 {
		msg := r.Args[0]
		w.SendSimpleString(msg)
	}

	w.SendSimpleString("PONG")
}
