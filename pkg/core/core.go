package core

// VMDriver describes common VM operations
type VMDriver interface {
	ListNetworks() ([]VMNetwork, error)
	CreateNetwork(netname string) (VMNetwork, error)
	GetNetwork(netname string) (VMNetwork, error)
	DeleteNetwork(netname string) error

	/*
		FetchMasterNodeImage() error
		FetchWorkerNodeImage() error
	*/
	ListHosts() ([]VMHost, error)
	CreateHost(hostname string, networkname string, position int) (VMHost, error)
	GetHost(hostname string, networkname string) (VMHost, error)
	DeleteHost(hostname string, networkname string) error
}

type driverregisterfunc func() (VMDriver, error)

// VMNetwork describes a virtual network
type VMNetwork interface {
	Name() string
	NetCIDR() string
}

// VMHost describes a node
type VMHost interface {
	Name() string
	Status() string
	SSHAddress() string

	Start() error
	Stop() error
	WaitForStateChange(int)
	ForwardSSHPort(int) error
}

// Cluster defines a kutti Kubernetes cluster
type Cluster interface {
	Name() string
	Status() string
}

// SSHClient defines a simple SSH client
type SSHClient interface {
	RunWithResults(string, string) (string, error)
}

func init() {
	drivers = make(map[string]func() (VMDriver, error))
	Clusters = make(map[string]Cluster)
	DefaultCluster = nil
	DefaultClient = newSSHClient("user1", "Pass@word1")
}
