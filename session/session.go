package session

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/store"
)

type Session struct {
	conn       net.Conn
	Store      *store.Store
	TxnContext *TxnContext
	Config     ConfigAccessor
}

func NewSession(conn net.Conn, store *store.Store, config ConfigAccessor) *Session {
	return &Session{
		conn:       conn,
		Store:      store,
		TxnContext: NewTxnContext(),
		Config:     config,
	}
}

func (s *Session) SendResponse(resp protocol.Response) {
	fmt.Fprint(s.conn, resp.Serialize())
}
