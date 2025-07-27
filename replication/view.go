package replication

// Provides a unified way to access replication-related metadata, such as role
// allowing handlers to be agnostic to whether the server is a replica/master

type Role string

type Info struct {
	Role         Role
	MasterReplID string
	MasterOffset int
}

type View interface {
	Snapshot() Info
}

const (
	RoleMaster  Role = "master"
	RoleReplica Role = "slave"
)

type MasterView struct {
	m *Master
}

func NewMasterView(m *Master) *MasterView {
	return &MasterView{m}
}

func (v *MasterView) Snapshot() Info {
	return Info{
		Role:         RoleMaster,
		MasterReplID: v.m.masterReplID,
		MasterOffset: v.m.masterOffset,
	}
}
