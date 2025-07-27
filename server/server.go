package server

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"

	"github.com/codecrafters-io/redis-starter-go/config"
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/registry"
	"github.com/codecrafters-io/redis-starter-go/replication"
	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

type Server struct {
	port      int
	logger    *slog.Logger
	registry  *registry.Registry
	store     *store.Store
	replicaof *config.ReplicaConfig
	master    *replication.Master
	replica   *replication.Replica
}

func New(
	port int,
	logger *slog.Logger,
	registry *registry.Registry,
	store *store.Store,
	replicaof *config.ReplicaConfig,
	master *replication.Master,
	replica *replication.Replica,
) *Server {
	return &Server{
		port,
		logger,
		registry,
		store,
		replicaof,
		master,
		replica,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind to %d", s.port)
	}

	defer l.Close()

	s.logger.Info(fmt.Sprintf("Redis server running on port %d", s.port))

	if s.replica != nil {
		// Handshake is initiated after the server startup and handle asynchronously
		go s.replica.Handshake(s.logger, s.port)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	s.logger.Info("new client connected", "addr", addr)

	defer func() {
		s.logger.Info("client disconnected", "addr", conn.RemoteAddr().String())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	var sess *session.Session
	if s.master != nil {
		sess = session.NewMasterSession(conn, s.store, s.master)
	} else {
		sess = session.NewReplicaSession(conn, s.store, s.replica)
	}

	for {
		request, err := protocol.ParseRequest(reader)

		if err != nil {
			break
		}

		s.logger.Debug("received command", "addr", addr, "command", request.Command)

		s.registry.Dispatch(sess, request)
	}
}
