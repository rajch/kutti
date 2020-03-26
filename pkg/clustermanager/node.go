package clustermanager

import (
	"github.com/rajch/kutti/pkg/core"
)

// Node is a node
type Node struct {
	cluster     *Cluster
	ClusterName string
	Name        string
	Type        string
	host        core.VMHost
	status      string
}

// Cluster returns the cluster this node belongs to
func (n *Node) Cluster() *Cluster {
	//fmt.Printf("BEFORE: Node:%+v]\n", n)
	if n.cluster == nil {
		n.cluster = manager.Clusters[n.ClusterName]
		n.cluster.ensureDriver()
	}
	//fmt.Printf("AFTER: Node:%+v]\n", n)
	return n.cluster
}

func (n *Node) createHost() error {
	c := n.Cluster()
	host, err := c.driver.CreateHost(n.Name, c.NetworkName, c.Name, c.K8sVersion)
	if err != nil {
		n.host = nil
		return err
	}
	n.host = host
	return nil
}

func (n *Node) ensureHost() error {
	if n.host == nil {
		c := n.Cluster()
		host, err := c.driver.GetHost(n.Name, c.NetworkName, c.Name)
		if err != nil {
			return err
		}

		n.host = host
	}
	return nil
}

// Status returns the current node status
func (n *Node) Status() string {
	n.ensureHost()
	// TODO: Get complete status here
	return n.host.Status()
}
