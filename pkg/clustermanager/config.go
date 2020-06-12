package clustermanager

import (
	"encoding/json"

	"github.com/rajch/kutti/internal/pkg/configfilemanager"
)

type clusterManagerConfig struct {
	Clusters map[string]*Cluster
}

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
		cc.Clusters = loadedconfig.Clusters
	}

	return err
}

func (cc *clusterManagerConfig) Setdefaults() {
	cc = &clusterManagerConfig{
		Clusters: make(map[string]*Cluster),
	}
}

func init() {
	config = &clusterManagerConfig{
		Clusters: make(map[string]*Cluster),
	}
	clusterconfigmanager = configfilemanager.New(configFileName, config)
}
