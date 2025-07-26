package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	Command string
	Args    []string
}

func ParseRequest(reader *bufio.Reader) (*Request, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(line, "*") {
		return nil, fmt.Errorf("expected an array")
	}

	count, err := strconv.Atoi(strings.TrimSpace(line)[1:])
	if err != nil {
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

	return &Request{
		Command: parts[0],
		Args:    parts[1:],
	}, nil
}
