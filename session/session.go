package session

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/store"
)

type Session struct {
	conn       net.Conn
	Store      *store.Store
	TxnContext *TxnContext
}

func NewSession(conn net.Conn, store *store.Store) *Session {
	return &Session{
		conn:       conn,
		Store:      store,
		TxnContext: NewTxnContext(),
	}
}

func (s *Session) SendSimpleString(reply string) {
	fmt.Fprintf(s.conn, "+%s\r\n", reply)
}

func (s *Session) SendError(reply string) {
	fmt.Fprintf(s.conn, "-ERR %s\r\n", reply)
}

func (s *Session) SendNullBulkString() {
	fmt.Fprint(s.conn, "$-1\r\n")
}

func (s *Session) SendBulkString(reply string) {
	fmt.Fprintf(s.conn, "$%d\r\n%s\r\n", len(reply), reply)
}

func (s *Session) SendInteger(reply int) {
	fmt.Fprintf(s.conn, ":%d\r\n", reply)
}

func (s *Session) SendArray() {
	fmt.Fprint(s.conn, "*0\r\n")
}
