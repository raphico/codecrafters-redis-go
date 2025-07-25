package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var store = map[string]string{}

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
			if len(args) != 2 {
				fmt.Fprint(conn, "-ERR syntax error\r\n")
				continue;
			}

			store[args[0]] = args[1]
			fmt.Fprint(conn, "+OK\r\n")
		} else if strings.EqualFold(command, "GET") {
			if len(args) != 1 {
				fmt.Fprint(conn, "ERR wrong number of arguments for 'get' command\r\n")
				continue;
			}

			value, ok := store[args[0]]
			if !ok {
				fmt.Fprint(conn, "-1\r\n")
				continue
			}

			fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(value), value)
		} else {
			conn.Write([]byte("+PONG\r\n"))
		}

	}
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