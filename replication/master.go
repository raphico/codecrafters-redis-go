package replication

type Master struct {
	replID string
	offset int
}

func NewMaster() *Master {
	return &Master{
		replID: "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		offset: 0,
	}
}

func (m *Master) BroadCastCommands() {
	
}