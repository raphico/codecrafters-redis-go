package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/codecrafters-io/redis-starter-go/handlers"
	"github.com/codecrafters-io/redis-starter-go/registry"
	"github.com/codecrafters-io/redis-starter-go/server"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func main() {
	port := flag.Int("port", 6379, "a custom port for running the redis server")

	flag.Parse()

	registry := registry.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := store.NewStore()
	s := server.New(*port, logger, registry, store)

	registry.Add("SET", handlers.HandleSet)
	registry.Add("GET", handlers.HandleGet)
	registry.Add("ECHO", handlers.HandleEcho)
	registry.Add("PING", handlers.HandlePing)
	registry.Add("INCR", handlers.HandleIncr)
	registry.Add("MULTI", handlers.HandleMulti)
	registry.Add("EXEC", handlers.MakeExecHandler(registry))
	registry.Add("DISCARD", handlers.HandleDiscard)

	err := s.Start()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
