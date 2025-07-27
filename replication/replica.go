package replication

import "github.com/codecrafters-io/redis-starter-go/config"

type Replica struct {
	masterHost string
	masterPort int
}

func NewReplica(config config.ReplicaConfig) *Replica {
	return &Replica{
		masterHost: config.Host,
		masterPort: config.Port,
	}
}
