package handlers

import "github.com/codecrafters-io/redis-starter-go/protocol"

func HandleExec(w protocol.Response, r *protocol.Request) {
	w.SendError("EXEC without MULTI")
}