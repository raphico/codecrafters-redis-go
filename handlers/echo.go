package handlers

import "github.com/codecrafters-io/redis-starter-go/protocol"

func HandleEcho(w protocol.Response, r *protocol.Request) {
	if len(r.Args) != 1 {
		w.SendError("wrong number of arguments for 'echo' command")
		return
	}

	msg := r.Args[0]
	w.SendSimpleString(msg)
}
