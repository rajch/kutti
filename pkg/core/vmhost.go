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
	ForwardPort(int, int) error
	UnforwardPort(int) error
	ForwardSSHPort(int) error
}
