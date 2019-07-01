package core

import (
	"os"
	"path"
)

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

	GetSSHAddressForNode(nodepostion int) string
}

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
	ForwardSSHPort(int) error
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
