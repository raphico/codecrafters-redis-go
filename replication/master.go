package replication

type Master struct {
	masterReplID string
	masterOffset int
}

func NewMaster() *Master {
	return &Master{
		masterReplID: "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		masterOffset: 0,
	}
}
