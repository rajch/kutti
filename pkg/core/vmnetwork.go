package core

// VMNetwork describes a virtual network
type VMNetwork interface {
	Name() string
	NetCIDR() string
}
