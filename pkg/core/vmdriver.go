package core

// VMDriver describes common VM operations
type VMDriver interface {
	Name() string
	Description() string
	RequiresPortForwarding() bool
	Status() string

	ListNetworks() ([]VMNetwork, error)
	CreateNetwork(netname string) (VMNetwork, error)
	GetNetwork(netname string) (VMNetwork, error)
	DeleteNetwork(netname string) error

	ListHosts() ([]VMHost, error)
	CreateHost(hostname string, networkname string, clustername string, k8sversion string) (VMHost, error)
	GetHost(hostname string, networkname string, clustername string) (VMHost, error)
	DeleteHost(hostname string, networkname string, clustername string) error

	ListImages() ([]VMImage, error)
	GetImage(k8sversion string) (VMImage, error)
}
