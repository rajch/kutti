package clustermanager

import (
	"fmt"

	"github.com/rajch/kutti/pkg/core"
)

// Cluster is a cluster
type Cluster struct {
	Name string

	DriverName string
	driver     core.VMDriver

	K8sVersion string

	NetworkName string
	network     core.VMNetwork

	Nodes map[string]*Node

	Type   string
	status string
}

func (c *Cluster) ensureDriver() error {
	var err error = nil
	var ok bool
	if c.driver == nil {
		c.driver, ok = core.GetDriver(c.DriverName)
		if !ok {
			c.status = "DriverNotPresent"
			return err
		}

		c.status = "Driver" + c.driver.Status()
	}

	return nil
}

func (c *Cluster) ensureNetwork() error {
	var err error
	if c.network == nil {
		c.network, err = c.driver.GetNetwork(c.NetworkName)
		if err != nil {
			c.status = "NetworkError"
			return err
		}

		c.status = "NetworkReady"
	}

	return nil
}

func (c *Cluster) createNetwork() error {
	c.NetworkName = c.Name + "net"
	nw, err := c.driver.CreateNetwork(c.NetworkName)
	if err != nil {
		c.status = "NetworkError"
		return err
	}
	c.network = nw
	c.status = "NetworkReady"
	return nil
}

// func (c *Cluster) ensurenodes() {
// 	for _, node := range c.Nodes {
// 		node.cluster = c
// 	}
// }

// func (c *Cluster) ensureHosts() error {
// 	if len(c.hosts) == 0 {
// 		for _, node := range c.Nodes {
// 			host, err := c.driver.GetHost(node.Name, c.NetworkName)
// 			if err != nil {
// 				node.status = "ERROR:" + err.Error()
// 			} else {
// 				node.status = host.Status()
// 			}
// 		}
// 	}

// 	c.status = "Ready"

// 	return nil
// }

func (c *Cluster) addnode(nodename string, nodetype string) (*Node, error) {
	err := c.ensureDriver()
	if err != nil {
		return nil, err
	}

	newnode := &Node{
		cluster:     c,
		ClusterName: c.Name,
		Name:        nodename,
		Type:        nodetype,
	}

	err = newnode.createHost()
	if err == nil {
		c.Nodes[nodename] = newnode
	}

	manager.Save()

	return newnode, err
}

func (c *Cluster) deletenodeentry(nodename string) error {
	delete(c.Nodes, nodename)
	return manager.Save()
}

func (c *Cluster) deletenode(nodename string) error {
	err := c.ensureDriver()
	if err != nil {
		return err
	}

	err = c.driver.DeleteHost(nodename, c.NetworkName, c.Name)
	if err == nil {
		err = c.deletenodeentry(nodename)
	}

	return err
}

// AddUninitializedNode adds a node, but does not start it or join it to the cluster
func (c *Cluster) AddUninitializedNode(nodename string) (*Node, error) {
	return c.addnode(nodename, "Unmanaged")
}

// DeleteNode deletes a node completely
func (c *Cluster) DeleteNode(nodename string, force bool) error {
	n, ok := c.Nodes[nodename]
	if !ok {
		return fmt.Errorf("node '%v' not found", nodename)
	}

	if n.Status() == "Unknown" {
		return c.deletenodeentry(nodename)
	}

	if n.Status() == "Running" {
		if !force {
			return fmt.Errorf("node '%v' is running. Cannot delete", nodename)
		}

		n.ForceStop()
	}

	return c.deletenode(nodename)
}

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

	// TODO: Ensure readiness of k8sversion

	newCluster.Type = "Unmanaged"
	newCluster.status = "Ready"
	return newCluster, nil

}
