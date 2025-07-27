package session

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/config"
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/store"
)

type Session struct {
	conn       net.Conn
	Store      *store.Store
	TxnContext *TxnContext
	Info       info
}

type info struct {
	Role string
}

func NewSession(conn net.Conn, store *store.Store, replicaof *config.ReplicaConfig) *Session {
	info := info{Role: "master"}
	if replicaof != nil {
		info.Role = "slave"
	}

	return &Session{
		conn:       conn,
		Store:      store,
		TxnContext: NewTxnContext(),
		Info:       info,
	}
}

func (s *Session) SendResponse(resp protocol.Response) {
	fmt.Fprint(s.conn, resp.Serialize())
}
