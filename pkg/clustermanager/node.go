package clustermanager

import (
	"github.com/rajch/kutti/pkg/core"
)

// Node is a node
type Node struct {
	cluster *Cluster
	Name    string
	host    core.VMHost
	status  string
}

func (n *Node) createHost() error {
	c := n.cluster
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
		c := n.cluster
		host, err := c.driver.GetHost(n.Name, c.NetworkName, c.Name)
		if err != nil {
			return err
		}

		n.host = host
	}
	return nil
}
