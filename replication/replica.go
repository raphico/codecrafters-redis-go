package replication

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/config"
	"github.com/codecrafters-io/redis-starter-go/protocol"
)

type ReplicaClient struct {
	masterHost string
	masterPort int
	masterConn net.Conn
	bufReader  *bufio.Reader
	bufWriter  *bufio.Writer
	logger     *slog.Logger
}

func NewReplicaClient(cfg config.ReplicaConfig, logger *slog.Logger, ) *ReplicaClient {
	return &ReplicaClient{
		masterHost: cfg.Host,
		masterPort: cfg.Port,
		logger:     logger,
	}
}

func (r *ReplicaClient) EstablishHandshake(listeningPort int) {
	addr := net.JoinHostPort(r.masterHost, fmt.Sprint(r.masterPort))

	r.logger.Info(fmt.Sprintf("Connecting to MASTER %s", addr))

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Connection with master failed %s", err.Error()))
		return
	}

	r.logger.Info("MASTER <-> REPLICA sync started")

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	r.bufReader = reader
	r.bufWriter = writer

	// 1. sends a PING to the master to check if its alive/responsive
	if err = r.sendCommand([]protocol.Response{
		protocol.NewBulkStringResponse("PING"),
	}); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to send PING: %s", err))
		return
	}
	if err = r.expectResponseContains("PING", "PONG"); err != nil {
		r.logger.Error(err.Error())
	}

	// 2. REPLCONF listening-port
	cmdArgs := []protocol.Response{
		protocol.NewBulkStringResponse("REPLCONF"),
		protocol.NewBulkStringResponse("listening-port"),
		protocol.NewBulkStringResponse(fmt.Sprint(listeningPort)),
	}
	if err := r.sendCommand(cmdArgs); err != nil {
		r.logger.Error(fmt.Sprintf("Failed to send REPLCONF listening-port: %s", err))
		return
	}
	if err = r.expectResponseContains("REPLCONF listening-port", "OK"); err != nil {
		r.logger.Error(err.Error())
	}

	// 3. REPLCONF capa
	cmdArgs = []protocol.Response{
		protocol.NewBulkStringResponse("REPLCONF"),
		protocol.NewBulkStringResponse("capa"),
		protocol.NewBulkStringResponse("psync2"),
	}
	if err := r.sendCommand(cmdArgs); err != nil {
		r.logger.Error("Failed to send REPLCONF capa psync2", "error", err)
		return
	}
	if err = r.expectResponseContains("REPLCONF capa", "OK"); err != nil {
		r.logger.Error(err.Error())
	}

	// 4. PSYNC
	cmdArgs = []protocol.Response{
		protocol.NewBulkStringResponse("PSYNC"),
		protocol.NewBulkStringResponse("?"),
		protocol.NewBulkStringResponse("-1"),
	}
	if err := r.sendCommand(cmdArgs); err != nil {
		r.logger.Error("Failed to send PSYNC", "error", err)
		return
	}
	if err = r.expectResponseContains("PSYNC", "FULLRESYNC"); err != nil {
		r.logger.Error(err.Error())
	}

	r.logger.Info("MASTER <-> REPLICA sync successful")

	r.masterConn = conn
	// go r.startSyncLoop()
}

func (r *ReplicaClient) sendCommand(cmdArgs []protocol.Response) error {
	serialized, ok := protocol.NewArrayResponse(cmdArgs).Serialize().(string)
	if !ok {
		panic("Serialization did not return a string")
	}

	if _, err := r.bufWriter.WriteString(serialized); err != nil {
		return err
	}

	return r.bufWriter.Flush()
}

func (r *ReplicaClient) expectResponseContains(cmd string, expected string) error {
	resp, err := r.readLine()
	if err != nil {
		return fmt.Errorf("failed to read %s response: %w", cmd, err)
	}

	if !strings.Contains(resp, expected) {
		return fmt.Errorf("unexpected response: got '%s', expected to contain '%s'", resp, expected)
	}

	return nil
}

func (r *ReplicaClient) readLine() (string, error) {
	line, err := r.bufReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// func (r *ReplicaClient) startSyncLoop() {
// 	for {
// 		request, err := protocol.ParseRequest(r.bufReader)
// 		if err != nil {
// 			r.logger.Error("PError reading from master: ", "protocol error", err)
// 			r.masterConn.Close()
// 			return
// 		}

// 		switch
// 	}
// }
