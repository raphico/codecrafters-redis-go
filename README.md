# Build Your Own Redis in Go (from Scratch)

A fully functional Redis clone built entirely from scratch using the Go standard library. Built as part of the [CodeCrafters "Build Your Own Redis" Challenge](https://codecrafters.io/challenges/redis), this project explores core Redis features and systems programming concepts including RESP parsing, custom TCP servers, in-memory data storage, transactions, and persistence.

> 🚀 Zero dependencies. Just Go, sockets, and deep protocol understanding.

## Goals

- Understand and implement the Redis protocol (RESP)
- Build a robust, concurrent TCP server from scratch
- Reproduce core Redis functionality (PING, ECHO, SET, GET, INCR, etc.)
- Implement advanced Redis data types like Lists and Streams
- Handle complex features like transactions and RDB persistence
- Strengthen skills in concurrency, systems design, and Go internals

## Implemented Features

### Commands

| Command                       | Description                                                                                                             |
| ----------------------------- | ----------------------------------------------------------------------------------------------------------------------- |
| `PING`                        | Returns `PONG`, or a custom message if provided (`PING hello` → `hello`)                                                |
| `ECHO <msg>`                  | Returns the exact message sent                                                                                          |
| `SET <key> <val>`             | Stores `val` under `key`, overwrites if it exists                                                                       |
| `GET <key>`                   | Retrieves value for `key`, or nil if it doesn’t exist                                                                   |
| `DEL <key> [key ...]`         | Deletes one or more the key-value pairs                                                                                 |
| `INCR <key>`                  | Increments an integer value (creates it if missing, starts at 0)                                                        |
| `DECR <key>`                  | Decrements an integer value (creates it if missing, starts at -1)                                                       |
| `MULTI`                       | Begins transaction mode, queues following commands                                                                      |
| `EXEC`                        | Executes queued transaction commands                                                                                    |
| `KEYS <pattern>`              | Returns all keys matching glob-style pattern (`*`, etc)                                                                 |
| `CONFIG GET`                  | Returns RDB config like filename and directory                                                                          |
| `RPUSH <key> <val> [val ...]` | Appends value(s) to list at `key`, creates list if it doesn’t exist                                                     |
| `LPUSH <key> <val> [val ...]` | Prepends value(s) to the start of the list at `key` (creates list if it doesn’t exist)                                  |
| `LRANGE <key> <start> <stop>` | Returns elements in the list from index start to stop (inclusive, supports negative indices)                            |
| `LLEN <key>`                  | Returns the length of the list stored at `key`                                                                          |
| `LPOP  <key>`                 | Removes and returns the first element of the list at `key`                                                              |
| `BLPOP key timeout`           | Removes and returns the first element of the list at the given key, block if empty until a timeout or new data arrives. |

### Concurrency & Networking

- Manual TCP server built with `net.Listen` and `Accept`
- Goroutine-per-connection concurrency model
- Graceful client connection and disconnection
- RESP protocol parser (array, bulk strings, simple strings, errors, integers)
- Structured logging (with `slog`)

### In-Memory Store

- Thread-safe key-value store using `sync.RWMutex`
- Optional TTL support via expiry timestamps
- Clean handling of expired keys

### Transaction System

- Command queueing with `MULTI`
- Conditional execution with `EXEC`
- Aborting transactions via `DISCARD`
- Dirty flag detection for invalid commands within a transaction

### Persistence

- Loads data from RDB snapshot on startup
- Configurable file location via dbfilename and dir
- Parses string keys and values from standard RDB format
- Supports optional TTL and skips expired keys
- Validates file structure and handles malformed input

## Not Yet Implemented

These are on the roadmap or part of the extended challenge, but **not yet implemented**:

- ❌ **Advanced data types**: Streams

## How to Run

1. Clone the repository

```bash
git clone git@github.com:raphico/codecrafters-redis-go.git
cd codecrafters-redis-go
```

2. Run the server

```bash
go run app/main.go
```

3. Test with redis-cli

```bash
redis-cli PING
```

## Folder structure

```bash
├── app/main.go                # Entry point, registers commands, starts server
├── config/config.go           # Defines and validates configuration, such as RDB path
├── handlers/                  # Individual command handlers
├── persistence/
│   ├── load.go                # Handles restoring dataset from RDB file
│   └── save.go                # Handles saving dataset to RDB file
├── protocol/
│   ├── request.go             # Parses incoming RESP requests
│   └── response.go            # Serializes and formats RESP responses
├── registry/registry.go       # Dispatches commands to handlers
├── server/server.go           # TCP server, handles concurrent clients
├── session/
│   ├── session.go             # Per-client connection state
│   └── transaction.go         # Handles transaction queueing and context
├── store/
│   ├── store.go               # Thread-safe key-value store
│   └── entry.go               # Defines key-value pair with optional expiry
```
