package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/pubsub"
	"github.com/codecrafters-io/redis-starter-go/store"
)

type Session struct {
	ID         string
	conn       net.Conn
	Store      *store.Store
	TxnContext *TxnContext
	Config     ConfigAccessor
	Pubsub     *pubsub.PubsubManager
}

func NewSession(
	conn net.Conn,
	store *store.Store,
	config ConfigAccessor,
	ps *pubsub.PubsubManager,
) *Session {
	return &Session{
		ID:         newSessionID(),
		conn:       conn,
		Store:      store,
		TxnContext: NewTxnContext(),
		Config:     config,
		Pubsub:     ps,
	}
}

func newSessionID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("unexpected: error generating session id")
	}
	return "sess-" + hex.EncodeToString(b)
}

func (s *Session) GetSessionID() string {
	return s.ID
}

func (s *Session) InSubscribeMode() bool {
	count := s.Pubsub.GetSubscribedCount(s)
	return count > 0
}

func (s *Session) SendResponse(resp protocol.Response) {
	fmt.Fprint(s.conn, resp.Serialize())
}

func (s *Session) SendMessage(channel, message string) {
	response := protocol.NewArrayResponse([]protocol.Response{
		protocol.NewBulkStringResponse("message"),
		protocol.NewBulkStringResponse(channel),
		protocol.NewBulkStringResponse(message),
	})

	s.SendResponse(response)
}
