package core

var (
	drivers map[string]func() (VMDriver, error)
	// DefaultCluster returns the default cluster
	DefaultCluster Cluster
	// Clusters returns currently defined clusters
	Clusters map[string]Cluster
)
