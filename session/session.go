package session

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/replication"
	"github.com/codecrafters-io/redis-starter-go/store"
)

type ReplicationState struct {
	View replication.View
}

type Session struct {
	conn       net.Conn
	Store      *store.Store
	TxnContext *TxnContext
	Repl       *ReplicationState
}

func NewReplicaSession(conn net.Conn, store *store.Store, replica *replication.Replica) *Session {
	return &Session{
		conn:       conn,
		Store:      store,
		TxnContext: NewTxnContext(),
		Repl: &ReplicationState{
			View: replication.NewReplicaView(replica),
		},
	}
}

func NewMasterSession(conn net.Conn, store *store.Store, master *replication.Master) *Session {
	return &Session{
		conn:       conn,
		Store:      store,
		TxnContext: NewTxnContext(),
		Repl: &ReplicationState{
			View: replication.NewMasterView(master),
		},
	}

}

func (s *Session) SendResponse(resp protocol.Response) {
	fmt.Fprint(s.conn, resp.Serialize())
}
