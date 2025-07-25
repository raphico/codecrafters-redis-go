package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type entry struct {
	value string
	expiryTime *time.Time
}

var store = map[string]entry{}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		request, err := parseRESPRequest(reader)
		if err != nil {
			break
		}

		command, args := request[0], request[1:]

		if strings.EqualFold(command, "ECHO") {
			if len(args) != 1 {
				fmt.Fprint(conn, "-ERR wrong number of arguments for 'echo' command\r\n")
				continue;
			}

			fmt.Fprintf(conn, "+%s\r\n", args[0])
		} else if strings.EqualFold(command, "SET") {
			if !isSetArgsValid(args) {
				fmt.Println(3)
				fmt.Fprint(conn, "-ERR syntax error\r\n")
				continue;
			}

			key, value := args[0], args[1]
			if len(args) == 2 {
				fmt.Println(0)
				store[key] = entry{value: value, expiryTime: nil}
				fmt.Fprint(conn, "+OK\r\n")
				continue
			}

			// this is already checked in isSetArgsValid
			ms, _ := strconv.Atoi(args[3])
			ttl := time.Duration(ms) * time.Millisecond
			expiry := time.Now().Add(ttl)
			store[key] = entry{value: value, expiryTime: &expiry}
			fmt.Fprint(conn, "+OK\r\n")
		} else if strings.EqualFold(command, "GET") {
			if len(args) != 1 {
				fmt.Fprint(conn, "ERR wrong number of arguments for 'get' command\r\n")
				continue;
			}
			
			key := args[0]
			entry, ok := store[key]
			if !ok {
				fmt.Fprint(conn, "-1\r\n")
				continue
			} else if entry.isExpired() {
				delete(store, key)
				fmt.Fprint(conn, "-1\r\n")
				continue
			}

			fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(entry.value), entry.value)
		} else {
			conn.Write([]byte("+PONG\r\n"))
		}

	}
}

func (e *entry) isExpired() bool {
	if e.expiryTime == nil {
		return false
	}

	return time.Now().After(*e.expiryTime)
}

func isSetArgsValid(args []string) bool {
	if len(args) < 2 {
		return false
	}

	if len(args) == 2 {
		return true
	}

	seenFlags := args[:2]

	i := 2
	for i < len(args) {
		arg := strings.ToUpper(args[i])

		if arg == "PX" {
			if i + 1 > len(args) {
				return false
			}

			if slices.Contains(seenFlags, arg) {
				return false
			}

			seenFlags = append(seenFlags, arg)

			if _, err := strconv.Atoi(args[i+1]); err!= nil {
				return false
			}

			i += 2
		} else {
			return false
		}
	}

	return true
}

func parseRESPRequest(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(line, "*") {
		return nil, fmt.Errorf("expected an array")
	}

	count, err := strconv.Atoi(strings.TrimSpace(line)[1:])
	if  err != nil {
		return nil, err
	}

	parts := make([]string, 0, count)

	for range count {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(line, "$") {
			return nil, fmt.Errorf("expected bulk string")
		}

		length, err := strconv.Atoi(strings.TrimSpace(line)[1:])
		if err != nil {
			return nil, err
		}

		// +2 for the CRLF
		buf := make([]byte, length+2)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, err
		}

		parts = append(parts, strings.TrimSpace(string(buf)))
	}

	return parts, nil
}