package clustermanager

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/rajch/kutti/internal/pkg/configfilemanager"
	"github.com/rajch/kutti/pkg/core"
)

var (
	errInvalidName         = errors.New("invalid name. Valid names are up to 10 characters long, must start with a lowercase letter, and may contain lowercase letters and digits only")
	errClusterExists       = errors.New("cluster already exists")
	errClusterDoesNotExist = errors.New("cluster does not exist")
	errClusterNotEmpty     = errors.New("cluster is not empty")
	errDriverDoesNotExist  = errors.New("driver does not exist")
	errImageNotAvailable   = errors.New("image not available")
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

	// DOING: Validate driver
	driver, ok := core.GetDriver(drivername)
	if !ok {
		return errDriverDoesNotExist
	}

	// DOING: Validate k8sversion
	driverimage, err := driver.GetImage(k8sversion)
	if err != nil {
		return err
	}

	if driverimage.Status() != "Available" {
		return errImageNotAvailable
	}

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

	// If this was the deault cluster, clear that
	if manager.DefaultClusterName == clustername {
		ClearDefaultCluster()
	}

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
	if manager.DefaultClusterName == "" {
		return nil
	}

	result, ok := manager.Clusters[manager.DefaultClusterName]
	if !ok {
		ClearDefaultCluster()
	}
	return result
}

// SetDefaultCluster sets the default cluster name.
// It returns an error if the cluster does not exist.
func SetDefaultCluster(clustername string) error {
	_, ok := manager.Clusters[clustername]
	if !ok {
		return errClusterDoesNotExist
	}

	manager.DefaultClusterName = clustername
	return Save()
}

// ClearDefaultCluster clears the default cluster name.
func ClearDefaultCluster() {
	manager.DefaultClusterName = ""
	Save()
}

// GetDriver gets the specified driver OR an error
func GetDriver(drivername string) (core.VMDriver, bool) {
	return core.GetDriver(drivername)
}

// ForEachDriver iterates over drivers
func ForEachDriver(f func(core.VMDriver) bool) {
	core.ForEachDriver(f)
}

// ForEachImage iterates over images for the specified driver
func ForEachImage(drivername string, f func(core.VMImage) bool) error {
	driver, ok := core.GetDriver(drivername)
	if !ok {
		return errDriverDoesNotExist
	}

	images, err := driver.ListImages()
	if err != nil {
		return err
	}

	for _, value := range images {
		cancel := f(value)
		if cancel {
			break
		}
	}

	return nil
}

// Save saves the current state to the configuration file.
func Save() error {
	data, err := json.Marshal(manager)
	if err != nil {
		return err
	}

	return configfilemanager.Save(configFileName, data)
}

// Load loads the cluster configuration from the configuration file
func Load() error {
	data, notexist, err := configfilemanager.Load(configFileName)
	if notexist {
		setdefaultmanagervalue()
		return Save()
	}

	var cm clusterManager
	err = json.Unmarshal(data, &cm)
	if err != nil {
		setdefaultmanagervalue()
		Save()
		return err
	}

	manager = cm
	return nil
}

func setdefaultmanagervalue() {
	manager = clusterManager{
		Clusters:           make(map[string]*Cluster),
		DefaultClusterName: "",
	}
}

func init() {
	Load()
}
