package core

// VMDriver describes common VM operations
type VMDriver interface {
	ListNetworks() ([]VMNetwork, error)
	CreateNetwork(netname string) (VMNetwork, error)
	DeleteNetwork(netname string) error
}

// VMNetwork describes a virtual network
type VMNetwork struct {
	Name    string
	NetCIDR string
}

// Cluster defines a kutti Kubernetes cluster
type Cluster struct {
	Network VMNetwork
}
