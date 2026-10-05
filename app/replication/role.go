package replication

type Role string

const (
	MasterRole  Role = "master"
	ReplicaRole Role = "slave"
)
