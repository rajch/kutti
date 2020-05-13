package clustermanager

import (
	"encoding/json"

	"github.com/rajch/kutti/internal/pkg/configfilemanager"
)

type clusterManagerConfig struct {
	Clusters map[string]*Cluster
	// DefaultClusterName string
}

// var (
// 	config clusterManagerConfig
// )

const (
	configFileName = "clusters.json"
)

var (
	clusterconfigmanager configfilemanager.ConfigManager
	config               *clusterManagerConfig
)

func (cc *clusterManagerConfig) Serialize() ([]byte, error) {
	return json.Marshal(cc)
}

func (cc *clusterManagerConfig) Deserialize(data []byte) error {
	var loadedconfig *clusterManagerConfig
	err := json.Unmarshal(data, &loadedconfig)
	if err == nil {
		cc = loadedconfig
	}

	return err
}

func (cc *clusterManagerConfig) Setdefaults() {
	cc = &clusterManagerConfig{
		Clusters: make(map[string]*Cluster),
	}
}

// Save saves the current state to the configuration file.
// func Save() error {
// 	data, err := json.Marshal(config)
// 	if err != nil {
// 		return err
// 	}

// 	return configfilemanager.SaveConfigfile(configFileName, data)
// }

// // Load loads the cluster configuration from the configuration file
// func Load() error {
// 	data, notexist, err := configfilemanager.LoadConfigfile(configFileName)
// 	if notexist {
// 		setdefaultmanagervalue()
// 		return Save()
// 	}

// 	var cm clusterManagerConfig
// 	err = json.Unmarshal(data, &cm)
// 	if err != nil {
// 		setdefaultmanagervalue()
// 		Save()
// 		return err
// 	}

// 	config = cm
// 	return nil
// }

// func setdefaultmanagervalue() {
// 	config = clusterManagerConfig{
// 		Clusters: make(map[string]*Cluster),
// 		//DefaultClusterName: "",
// 	}
// }

func init() {
	//Load()
	config = &clusterManagerConfig{
		Clusters: make(map[string]*Cluster),
	}
	clusterconfigmanager = configfilemanager.New(configFileName, config)
}
