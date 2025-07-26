package session

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/store"
)

type Session struct {
	id string
	conn net.Conn
	Store *store.Store
}

func generateSessionId() string {
    return fmt.Sprintf("client-%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

func NewSession(conn net.Conn, store *store.Store) *Session {
	return &Session{
		id: generateSessionId(),
		conn: conn,
		Store: store,
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
