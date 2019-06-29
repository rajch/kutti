package vboxdriver

// VBoxVMNetwork implements the VMNetwork interface for VirtualBox
type VBoxVMNetwork struct {
	name    string
	netCIDR string
}

// Name is the name of the network
func (vn *VBoxVMNetwork) Name() string {
	return vn.name
}

// NetCIDR is the network's IPv4 address range
func (vn *VBoxVMNetwork) NetCIDR() string {
	return vn.netCIDR
}
