package main

import (
	"log/slog"
	"os"

	"github.com/codecrafters-io/redis-starter-go/handlers"
	"github.com/codecrafters-io/redis-starter-go/registry"
	"github.com/codecrafters-io/redis-starter-go/server"
)

const port = "6379"

func main() {
	registry := registry.New()
	registry.Add("SET", handlers.HandleSet)
	registry.Add("GET", handlers.HandleGet)
	registry.Add("ECHO", handlers.HandleEcho)
	registry.Add("PING", handlers.HandlePing)
	registry.Add("INCR", handlers.HandleIncr)
	registry.Add("MULTI", handlers.HandleMulti)
	registry.Add("EXEC", handlers.HandleExec)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := server.New(port, logger, registry)

	err := s.Start()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
