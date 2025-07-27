package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/codecrafters-io/redis-starter-go/config"
	"github.com/codecrafters-io/redis-starter-go/handlers"
	"github.com/codecrafters-io/redis-starter-go/registry"
	"github.com/codecrafters-io/redis-starter-go/replication"
	"github.com/codecrafters-io/redis-starter-go/server"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func main() {

	port := flag.Int("port", 6379, "a custom port for running the redis server")
	rawReplicaof := flag.String(
		"replicaof",
		"no one",
		"specifies that a redis server is a replica of another redis server",
	)

	flag.Parse()

	replicaof, err := config.ValidateReplicaof(*rawReplicaof)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	var master *replication.Master
	if replicaof == nil {
		master = replication.NewMaster()
	}

	registry := registry.New()
	store := store.NewStore()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := server.New(*port, logger, registry, store, replicaof, master)

	registry.Add("SET", handlers.HandleSet)
	registry.Add("GET", handlers.HandleGet)
	registry.Add("ECHO", handlers.HandleEcho)
	registry.Add("PING", handlers.HandlePing)
	registry.Add("INCR", handlers.HandleIncr)
	registry.Add("MULTI", handlers.HandleMulti)
	registry.Add("EXEC", handlers.MakeExecHandler(registry))
	registry.Add("DISCARD", handlers.HandleDiscard)
	registry.Add("INFO", handlers.HandleInfo)

	err = s.Start()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
