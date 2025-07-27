package replication

import (
	"fmt"
	"log/slog"
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

func (r *Replica) Handshake(logger *slog.Logger) {
	addr := net.JoinHostPort(r.masterHost, fmt.Sprint(r.masterPort))

	logger.Info(fmt.Sprintf("Connecting to MASTER %s", addr))

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		logger.Error(fmt.Sprintf("Connection with master failed %s", err.Error()))
		return
	}

	logger.Info("MASTER <-> REPLICA sync started")

	// 1. sends a PING to the master to check if its alive/responsive
	msg := []protocol.Response{protocol.NewBulkStringResponse("PING")}
	_, err = fmt.Fprint(conn, protocol.NewArrayResponse(msg).Serialize())
	if err != nil {
		logger.Error(fmt.Sprintf("Error sending PING to master %s", err.Error()))
		return
	}

	// 2. sends two REPLCONF commands to configure the replication
	cmd := protocol.NewBulkStringResponse("REPLCONF")

	// 2.1 REPLCONF 1: notify the master what port it's listening on
	arg := protocol.NewBulkStringResponse("listening-port")
	port := protocol.NewBulkStringResponse(fmt.Sprint(r.masterPort))

	msg = []protocol.Response{cmd, arg, port}
	_, err = fmt.Fprint(conn, protocol.NewArrayResponse(msg).Serialize())
	if err != nil {
		logger.Error(fmt.Sprintf("Error sending REPLCONF listening-port %s", err.Error()))
		return
	}

	// 2.2 REPLCONF 2: notify the master of its capabilities
	arg = protocol.NewBulkStringResponse("capa")
	capa := protocol.NewBulkStringResponse("psync2")

	msg = []protocol.Response{cmd, arg, capa}
	_, err = fmt.Fprint(conn, protocol.NewArrayResponse(msg).Serialize())
	if err != nil {
		logger.Error(fmt.Sprintf("Error sending REPLCONF capa psync2 %s", err.Error()))
		return
	}

	logger.Info("MASTER <-> REPLICA sync successful")
}
