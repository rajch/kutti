package core

// VMHost describes a node
type VMHost interface {
	Name() string
	Status() string
	SSHAddress() string

	Start() error
	Stop() error
	ForceStop() error
	WaitForStateChange(int)
	ForwardSSHPort(int) error
}
