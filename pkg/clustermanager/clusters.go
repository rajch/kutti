package clustermanager

import (
	klog "github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/core"
)

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
	newCluster, err := newunmanagedcluster(name, k8sversion, drivername)
	if err != nil {
		return err
	}

	config.Clusters[name] = newCluster
	return clusterconfigmanager.Save()
}

// GetCluster gets a named cluster, or nil if not present
func GetCluster(name string) (*Cluster, bool) {
	cluster, ok := config.Clusters[name]
	if !ok {
		return nil, ok
	}
	return cluster, true
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

	klog.Println(2, "Deleting network...")
	err := cluster.deletenetwork()
	if err != nil {
		return err
	}
	klog.Println(2, "Network deleted.")

	delete(config.Clusters, clustername)

	return clusterconfigmanager.Save()
}

func newunmanagedcluster(name string, k8sversion string, drivername string) (*Cluster, error) {
	newCluster := &Cluster{
		Name:       name,
		K8sVersion: k8sversion,
		DriverName: drivername,
		//hosts:      make(map[string]core.VMHost),
		Nodes:  make(map[string]*Node),
		status: "UnInitialized",
	}

	// Ensure presence of VMdriver
	err := newCluster.ensuredriver()
	if err != nil {
		return newCluster, err
	}

	// Create VM Network
	klog.Println(2, "Creating network...")
	err = newCluster.createnetwork()
	if err != nil {
		return newCluster, err
	}

	newCluster.Type = "Unmanaged"
	newCluster.status = "Ready"
	klog.Println(2, "Network created.")

	return newCluster, nil
}
