package replication

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/config"
	"github.com/codecrafters-io/redis-starter-go/protocol"
)

type Replica struct {
	masterHost string
	masterPort int
}

func NewReplica(config config.ReplicaConfig) *Replica {
	return &Replica{
		masterHost: config.Host,
		masterPort: config.Port,
	}
}

func (r *Replica) Handshake() error {
	addr := net.JoinHostPort(r.masterHost, fmt.Sprint(r.masterPort))

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	msg := []protocol.Response{protocol.NewBulkStringResponse("PING")}
	fmt.Fprint(conn, protocol.NewArrayResponse(msg).Serialize())

	return nil
}