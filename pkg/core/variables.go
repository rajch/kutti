package core

var (
	drivers map[string]func() (VMDriver, error)
	// DefaultClient returns the default SSH client
	DefaultClient SSHClient
	// DefaultCluster returns the default cluster
	DefaultCluster Cluster
	// Clusters returns currently defined clusters
	Clusters map[string]Cluster
)
