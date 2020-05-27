package core

// VMNetwork describes a virtual network.
// A cluster of VMHosts are connected via a VMNetwork. A VMNetwork
// manages VMHost IP addresses.
type VMNetwork interface {
	Name() string
	NetCIDR() string
}
