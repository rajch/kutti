package clustermanager

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/rajch/kutti/pkg/core"
)

type clusterManager struct {
	Clusters           map[string]*Cluster
	DefaultClusterName string
}

func (cm *clusterManager) DefaultCluster() *Cluster {
	result, _ := cm.Clusters[cm.DefaultClusterName]
	return result
}

// NewEmptyCluster creates a new, empty cluster
func NewEmptyCluster(name string, k8sversion string, drivername string) error {
	newCluster, err := newEmptyCluster(name, k8sversion, drivername)
	if err != nil {
		return err
	}

	manager.Clusters[name] = newCluster
	return manager.Save()
}

func (cm *clusterManager) Save() error {
	data, err := json.Marshal(cm)
	if err != nil {
		return err
	}

	configPath, err := core.ConfigDir()
	if err != nil {
		return err
	}

	datafilepath := path.Join(configPath, "clusters.json")
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

// DeleteCluster deletes a cluster
// func DeleteCluster(clustername string) error {
// 	cluster, ok := GetCluster(clustername)
// 	if !ok {
// 		return fmt.Errorf("cluster '%s does not exist", clustername)
// 	}

// 	return nil
// }

// Load loads the cluster configuration from the configuration file
func Load() error {
	configPath, err := core.ConfigDir()
	if err != nil {
		return err
	}

	datafilepath := path.Join(configPath, "clusters.json")
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

	return manager.Save()

}

// Clusters returns clusters
func Clusters() map[string]*Cluster {
	return manager.Clusters
}

// GetCluster gets a named cluster, or nil if not present
func GetCluster(name string) (*Cluster, bool) {
	cluster, ok := manager.Clusters[name]
	if !ok {
		return nil, ok
	}
	return cluster, true
}

// DefaultCluster returns the default cluster, or nil if none has been set
func DefaultCluster() *Cluster {
	return manager.DefaultCluster()
}

func init() {
	Load()
}

var manager clusterManager