package domain

const (
	IpOpertionDeny IpOperation = iota
	IpOperationAllow
)

type IpOperation int
