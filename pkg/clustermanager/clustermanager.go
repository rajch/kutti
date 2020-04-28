package clustermanager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/rajch/kutti/pkg/core"
)

var (
	errInvalidName         = errors.New("invalid name. Valid names are up to 10 characters long, must start with a lowercase letter, and may contain lowercase letters and digits only")
	errClusterExists       = errors.New("cluster already exists")
	errClusterDoesNotExist = errors.New("cluster does not exist")
	errClusterNotEmpty     = errors.New("cluster is not empty")
	manager                clusterManager
)

const (
	configFileName = "clusters.json"
)

type clusterManager struct {
	Clusters           map[string]*Cluster
	DefaultClusterName string
}

// IsValidName checks for the validity of a name.
// Valid names are up to 10 characters long, must start with a lowercase letter, and may contain lowercase letters and digits only.
func IsValidName(name string) bool {
	matched, _ := regexp.MatchString("^[a-z]([a-z0-9]{1,9})$", name)
	return matched
}

// NewEmptyCluster creates a new, empty cluster
func NewEmptyCluster(name string, k8sversion string, drivername string) error {
	// Validate name
	if !IsValidName(name) {
		return errInvalidName
	}

	// Check if name exists
	_, ok := manager.Clusters[name]
	if ok {
		return errClusterExists
	}

	// TODO: Validate driver

	// TODO: Validate k8sversion

	newCluster, err := newEmptyCluster(name, k8sversion, drivername)
	if err != nil {
		return err
	}

	manager.Clusters[name] = newCluster
	return Save()
}

// DeleteCluster deletes a cluster.
// Currently, the cluster must be empty.
func DeleteCluster(clustername string) error {
	cluster, ok := GetCluster(clustername)
	if !ok {
		return errClusterDoesNotExist
	}

	// TODO: Temporary condition. Will fix later.
	if len(cluster.Nodes) > 0 {
		return errClusterNotEmpty
	}

	err := cluster.deleteNetwork()
	if err != nil {
		return err
	}

	delete(manager.Clusters, clustername)

	return Save()
}

// GetCluster gets a named cluster, or nil if not present
func GetCluster(name string) (*Cluster, bool) {
	cluster, ok := manager.Clusters[name]
	if !ok {
		return nil, ok
	}
	return cluster, true
}

// ForEachCluster iterates over clusters
func ForEachCluster(f func(*Cluster) bool) {
	for _, cluster := range manager.Clusters {
		if cancel := f(cluster); cancel {
			break
		}
	}
}

// DefaultCluster returns the default cluster, or nil if none has been set
func DefaultCluster() *Cluster {
	result, _ := manager.Clusters[manager.DefaultClusterName]
	return result
}

// ForEachDriver iterates over drivers
func ForEachDriver(f func(core.VMDriver) bool) {
	core.ForEachDriver(f)
}

func getconfigfilepath() (string, error) {
	configPath, err := core.ConfigDir()
	if err != nil {
		return "", err
	}

	datafilepath := path.Join(configPath, configFileName)
	return datafilepath, nil
}

// Save saves the current state to the configuration file.
func Save() error {
	data, err := json.Marshal(manager)
	if err != nil {
		return err
	}

	datafilepath, err := getconfigfilepath()
	if err != nil {
		return err
	}

	file, err := os.Create(datafilepath)
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Load loads the cluster configuration from the configuration file
func Load() error {
	datafilepath, err := getconfigfilepath()
	if err != nil {
		return err
	}
	_, err = os.Stat(datafilepath)
	if !(err == nil || os.IsNotExist(err)) {
		return err
	}

	if err == nil {
		data, err := ioutil.ReadFile(datafilepath)

		if err != nil {
			return err
		}

		var cm clusterManager
		err = json.Unmarshal(data, &cm)
		if err != nil {
			return err
		}

		manager = cm
		return nil
	}

	manager = clusterManager{
		Clusters:           make(map[string]*Cluster),
		DefaultClusterName: "",
	}

	return Save()

}

func init() {
	Load()
}
