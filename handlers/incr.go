package handlers

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleIncr(w protocol.Response, r *protocol.Request) {
	key := r.Args[0]

	e, err := store.MemStore.Get(key)
	if err != nil {
		store.MemStore.Set(key, "1", nil)
		w.SendInteger(1)
		return
	}

	curr, err := strconv.Atoi(e.Value)
	if err != nil {
		w.SendError("value is not an integer or out of range")
	}

	store.MemStore.Update(key, strconv.Itoa(curr+1))
	w.SendInteger(curr + 1)
}
