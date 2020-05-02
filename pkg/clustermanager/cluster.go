package clustermanager

import (
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
	if c.driver == nil {
		driver, ok := core.GetDriver(c.DriverName)
		if !ok {
			c.status = "DriverNotPresent"
			return errDriverDoesNotExist
		}

		c.driver = driver
		c.status = "Driver" + c.driver.Status()
	}

	return nil
}

func (c *Cluster) ensureNetwork() error {
	if c.network == nil {
		network, err := c.driver.GetNetwork(c.NetworkName)
		if err != nil {
			c.status = "NetworkError"
			return err
		}
		c.network = network
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

func (c *Cluster) deleteNetwork() error {
	c.ensureDriver()
	err := c.driver.DeleteNetwork(c.NetworkName)
	if err != nil {
		c.status = "NetworkDeleteError"
		return err
	}
	c.network = nil
	c.status = "NetworkDeleted"
	return nil
}

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

	err = newnode.createhost()
	if err == nil {
		c.Nodes[nodename] = newnode
	}

	Save()

	return newnode, err
}

func (c *Cluster) deletenodeentry(nodename string) error {
	delete(c.Nodes, nodename)
	return Save()
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
	if !IsValidName(nodename) {
		return nil, errInvalidName
	}

	return c.addnode(nodename, "Unmanaged")
}

// DeleteNode deletes a node completely
func (c *Cluster) DeleteNode(nodename string, force bool) error {
	n, ok := c.Nodes[nodename]
	if !ok {
		return errNodeNotFound
	}

	if n.Status() == "Unknown" {
		return c.deletenodeentry(nodename)
	}

	if n.Status() == "Running" {
		if !force {
			return errNodeIsRunning
		}

		n.ForceStop()
	}

	return c.deletenode(nodename)
}
