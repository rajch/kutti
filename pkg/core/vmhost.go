package core

// VMHost describes a virtual host. In kutti, a VMHost can be
// started, stoppped normally and stopped forcibly. If the VMDriver and
// VMNetwork use NAT, then a VMHost can also have its ports forwarded
// to physical host ports.
type VMHost interface {
	Name() string
	Status() string
	SSHAddress() string

	Start() error
	Stop() error
	ForceStop() error
	WaitForStateChange(timeoutinseconds int)
	ForwardPort(hostport int, vmport int) error
	UnforwardPort(vmport int) error
	ForwardSSHPort(hostport int) error
}
