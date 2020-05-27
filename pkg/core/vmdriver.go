package core

// VMDriver describes operations to manage VMNetworks, VMHosts and VMImages.
// The RequiresPortForwarding() method is particularly important. It is expected to
// return true if the VMNetworks of the driver use NAT. This means that ports of the
// VMHosts will need to be forwarded to physical ports for access.
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

	FetchImageList() error
	ListImages() ([]VMImage, error)
	GetImage(k8sversion string) (VMImage, error)
}
