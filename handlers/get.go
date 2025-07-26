package handlers

import (
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleGet(w protocol.Response, r *protocol.Request) {
	if len(r.Args) != 1 {
		w.SendError("wrong number of arguments for 'get' command")
		return
	}

	key := r.Args[0]

	entry, err := store.MemStore.Get(key)
	if err != nil {
		w.SendNullBulkString()
		return
	}

	if entry.IsExpired() {
		store.MemStore.Delete(key)
		w.SendNullBulkString()
		return
	}

	w.SendBulkString(entry.Value)
}
