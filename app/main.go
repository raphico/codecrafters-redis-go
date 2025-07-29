package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/codecrafters-io/redis-starter-go/config"
	"github.com/codecrafters-io/redis-starter-go/handlers"
	"github.com/codecrafters-io/redis-starter-go/persistence"
	"github.com/codecrafters-io/redis-starter-go/registry"
	"github.com/codecrafters-io/redis-starter-go/server"
	"github.com/codecrafters-io/redis-starter-go/store"
)

const port = "6379"

func main() {
	dir := flag.String("dir", "/tmp/redis-files", "the path of the rdb file")
	dbfilename := flag.String("dbfilename", "dump.rdb", "the filename of the rdb file")

	flag.Parse()

	registry := registry.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := store.NewStore()

	config, err := config.NewConfig(*dbfilename, *dir)
	if err == nil {
		persistence.LoadRDB(config, store)
	}

	s := server.New(port, logger, registry, store, config)

	registry.Add("SET", handlers.HandleSet)
	registry.Add("GET", handlers.HandleGet)
	registry.Add("ECHO", handlers.HandleEcho)
	registry.Add("PING", handlers.HandlePing)
	registry.Add("INCR", handlers.HandleIncr)
	registry.Add("DECR", handlers.HandleDecr)
	registry.Add("EXISTS", handlers.HandleExists)
	registry.Add("TYPE", handlers.HandleType)
	registry.Add("TTL", handlers.HandleTtl)
	registry.Add("MULTI", handlers.HandleMulti)
	registry.Add("EXEC", handlers.MakeExecHandler(registry))
	registry.Add("DISCARD", handlers.HandleDiscard)
	registry.Add("CONFIG", handlers.HandleConfig)
	registry.Add("SAVE", handlers.HandleSave)
	registry.Add("KEYS", handlers.HandleKeys)
	registry.Add("RPUSH", handlers.HandleRpush)
	registry.Add("LPUSH", handlers.HandleLpush)
	registry.Add("LRANGE", handlers.HandleLrange)
	registry.Add("LLEN", handlers.HandleLlen)
	registry.Add("LPOP", handlers.HandleLpop)
	registry.Add("RPOP", handlers.HandleRpop)
	registry.Add("BLPOP", handlers.HandleBlpop)
	registry.Add("DEL", handlers.HandleDel)

	err = s.Start()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
