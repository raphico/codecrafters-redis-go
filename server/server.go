package server

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"

	"github.com/codecrafters-io/redis-starter-go/config"
	"github.com/codecrafters-io/redis-starter-go/protocol"
	"github.com/codecrafters-io/redis-starter-go/registry"
	"github.com/codecrafters-io/redis-starter-go/session"
	"github.com/codecrafters-io/redis-starter-go/store"
)

type Server struct {
	port     string
	logger   *slog.Logger
	registry *registry.Registry
	store    *store.Store
	config   *config.Config
}

func New(
	port string,
	logger *slog.Logger,
	registry *registry.Registry,
	store *store.Store,
	config *config.Config,
) *Server {
	return &Server{
		port,
		logger,
		registry,
		store,
		config,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%s", s.port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind to %s", s.port)
	}

	defer l.Close()

	s.logger.Info(fmt.Sprintf("Redis server running on port %s", s.port))

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
	session := session.NewSession(conn, s.store, s.config)

	for {
		request, err := protocol.ParseRequest(reader)

		if err != nil {
			break
		}

		s.logger.Debug("received command", "addr", addr, "command", request.Command)

		s.registry.Dispatch(session, request)
	}
}
