# 🧠 Build Your Own Redis in Go (from Scratch)

A fully functional Redis clone built entirely from scratch using the Go standard library. Built as part of the [CodeCrafters "Build Your Own Redis" Challenge](https://codecrafters.io/challenges/redis), this project explores core Redis features and systems programming concepts including RESP parsing, custom TCP servers, in-memory data storage, transactions, replication, and persistence.

> 🚀 Zero dependencies. Just Go, sockets, and deep protocol understanding.

## 🎯 Goals

- Understand and implement the Redis protocol (RESP)
- Build a robust, concurrent TCP server from scratch
- Reproduce core Redis functionality (PING, ECHO, SET, GET, INCR, etc.)
- Handle complex features like transactions, replication, and RDB persistence
- Strengthen skills in concurrency, systems design, and Go internals

## Implemented Features

### Core Commands

| Command           | Description                                                |
| ----------------- | ---------------------------------------------------------- |
| `PING`            | Returns `PONG` or a custom message                         |
| `ECHO <msg>`      | Echoes back the provided message                           |
| `SET <key> <val>` | Stores a value under a key                                 |
| `GET <key>`       | Retrieves the value for a key                              |
| `INCR <key>`      | Increments integer value of key                            |
| `MULTI`           | Starts a transaction block                                 |
| `EXEC`            | Executes queued transaction commands                       |
| `DISCARD`         | Cancels a queued transaction                               |
| `INFO <section>`  | Returns server information; optionally filtered by section |

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

## Not Yet Implemented

These are on the roadmap or part of the extended challenge, but **not yet implemented**:

- ❌ **Replication** (leader-follower sync, `PSYNC`, `INFO`)
- ❌ **Persistence** (RDB snapshot format, disk I/O)
- ❌ **WAIT, ACK, INFO replica behavior**
- ❌ **Advanced data types**: Lists, Sets, Streams

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
├── app/
│   └── main.go                # Entry point: command registration, server start
│
├── config/
│   └── replica.go             # Parse and validate --replicaof flag
│
├── handlers/                  # RESP command handlers (e.g., INFO, SET, GET, etc.)
│   └── info.go                # INFO command handler (more can be added)
│
├── protocol/
│   ├── request.go             # Parse RESP requests
│   └── response.go            # Format RESP responses
│
├── registry/
│   └── registry.go            # Command dispatch registry
│
├── replication/
│   ├── master.go              # State and logic for master role
│   ├── replica.go             # State and logic for replica role
│   └── view.go                # Unified read-only interface (used by handlers like INFO)
│
├── server/
│   └── server.go              # TCP server, manages client connections
│
├── session/
│   ├── session.go             # Per-client state (e.g., auth, transaction, replication info)
│   └── transaction.go         # Handles transaction logic per session
│
├── store/
│   ├── store.go               # In-memory key-value store with thread-safety
│   └── entry.go               # Defines entry: key-value + optional expiry

```
