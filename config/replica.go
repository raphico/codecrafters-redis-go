package config

import (
	"fmt"
	"strconv"
	"strings"
)

type ReplicaConfig struct {
	Host string
	Port int
}

func ValidateReplicaof(config string) (*ReplicaConfig, error) {
	if strings.EqualFold(config, "no one") {
		return nil, nil
	}

	parts := strings.Fields(config)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid argument for --replicaof: you must provide host and port")
	}

	host, portStr := parts[0], parts[1]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("bad port number or port is not a number: %s", portStr)
	}

	return &ReplicaConfig{
		Host: host,
		Port: port,
	}, nil
}
