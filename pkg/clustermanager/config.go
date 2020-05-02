package clustermanager

import (
	"encoding/json"

	"github.com/rajch/kutti/internal/pkg/configfilemanager"
)

type clusterManagerConfig struct {
	Clusters           map[string]*Cluster
	DefaultClusterName string
}

var (
	config clusterManagerConfig
)

const (
	configFileName = "clusters.json"
)

// Save saves the current state to the configuration file.
func Save() error {
	data, err := json.Marshal(config)
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

	var cm clusterManagerConfig
	err = json.Unmarshal(data, &cm)
	if err != nil {
		setdefaultmanagervalue()
		Save()
		return err
	}

	config = cm
	return nil
}

func setdefaultmanagervalue() {
	config = clusterManagerConfig{
		Clusters:           make(map[string]*Cluster),
		DefaultClusterName: "",
	}
}

func init() {
	Load()
}
