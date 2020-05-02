package clustermanager

import "github.com/rajch/kutti/pkg/core"

func newEmptyCluster(name string, k8sversion string, drivername string) (*Cluster, error) {
	newCluster := &Cluster{
		Name:       name,
		K8sVersion: k8sversion,
		DriverName: drivername,
		//hosts:      make(map[string]core.VMHost),
		Nodes:  make(map[string]*Node),
		status: "UnInitialzed",
	}

	// Ensure presence of VMdriver
	err := newCluster.ensureDriver()
	if err != nil {
		return newCluster, err
	}

	// Create VM Network
	err = newCluster.createNetwork()
	if err != nil {
		return newCluster, err
	}

	newCluster.Type = "Unmanaged"
	newCluster.status = "Ready"
	return newCluster, nil

}

// ForEachCluster iterates over clusters
func ForEachCluster(f func(*Cluster) bool) {
	for _, cluster := range config.Clusters {
		if cancel := f(cluster); cancel {
			break
		}
	}
}

// NewEmptyCluster creates a new, empty cluster
func NewEmptyCluster(name string, k8sversion string, drivername string) error {
	// Validate name
	if !IsValidName(name) {
		return errInvalidName
	}

	// Check if name exists
	_, ok := config.Clusters[name]
	if ok {
		return errClusterExists
	}

	// Validate driver
	driver, ok := core.GetDriver(drivername)
	if !ok {
		return errDriverDoesNotExist
	}

	// Validate k8sversion
	driverimage, err := driver.GetImage(k8sversion)
	if err != nil {
		return err
	}

	if driverimage.Status() != "Available" {
		return errImageNotAvailable
	}

	// Create cluster
	newCluster, err := newEmptyCluster(name, k8sversion, drivername)
	if err != nil {
		return err
	}

	config.Clusters[name] = newCluster
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

	delete(config.Clusters, clustername)

	// If this was the deault cluster, clear that
	if config.DefaultClusterName == clustername {
		ClearDefaultCluster()
	}

	return Save()
}

// GetCluster gets a named cluster, or nil if not present
func GetCluster(name string) (*Cluster, bool) {
	cluster, ok := config.Clusters[name]
	if !ok {
		return nil, ok
	}
	return cluster, true
}
