package handlers

import "github.com/codecrafters-io/redis-starter-go/protocol"

func HandleMulti(w protocol.Response, r *protocol.Request) {
	w.SendSimpleString("OK")
}