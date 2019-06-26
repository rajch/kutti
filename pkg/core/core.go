package core

import (
	"os"
	"path"
)

// VMDriver describes common VM operations
type VMDriver interface {
	ListNetworks() ([]VMNetwork, error)
	CreateNetwork(netname string) (VMNetwork, error)
	DeleteNetwork(netname string) error

	FetchMasterNodeImage() error
	FetchWorkerNodeImage() error

	ListNodes() error
	CreateNode(nodename string) error
	DeleteNode(nodename string) error
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

// CacheDir returns the location where the kutti cache should reside
func CacheDir() (result string, err error) {
	result, err = os.UserCacheDir()
	if err != nil {
		return
	}

	result = path.Join(result, "kutti")
	_, err = os.Stat(result)
	if os.IsNotExist(err) {
		err = os.Mkdir(result, 0755)
		if err != nil {
			result = ""
		}
	}

	return
}
