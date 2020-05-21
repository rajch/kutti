package clustermanager

import (
	"github.com/rajch/kutti/internal/pkg/kuttilog"
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

// NewUninitializedNode adds a node, but does not join it to a kubernetes cluster
func (c *Cluster) NewUninitializedNode(nodename string) (*Node, error) {
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

		kuttilog.Printf(2, "Stopping node %s...", nodename)
		err := n.ForceStop()
		if err != nil {
			return err
		}
		kuttilog.Printf(2, "Node %s stopped.", nodename)
	}

	// Unmap ports
	for key := range n.Ports {
		err := n.host.UnforwardPort(key)
		if err != nil {
			kuttilog.Printf(3, "Error while unmapping ports for node '%s': %v.", nodename, err)
		}
	}

	return c.deletenode(nodename)
}

func (c *Cluster) ensuredriver() error {
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

func (c *Cluster) ensurenetwork() error {
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

func (c *Cluster) createnetwork() error {
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

func (c *Cluster) deletenetwork() error {
	c.ensuredriver()
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
	err := c.ensuredriver()
	if err != nil {
		return nil, err
	}

	newnode := &Node{
		cluster:     c,
		ClusterName: c.Name,
		Name:        nodename,
		Type:        nodetype,
		Ports:       make(map[int]int),
	}

	err = newnode.createhost()
	if err == nil {
		c.Nodes[nodename] = newnode
		clusterconfigmanager.Save()
	}

	return newnode, err
}

func (c *Cluster) deletenodeentry(nodename string) error {
	delete(c.Nodes, nodename)
	return clusterconfigmanager.Save()
}

func (c *Cluster) deletenode(nodename string) error {
	err := c.ensuredriver()
	if err != nil {
		return err
	}

	err = c.driver.DeleteHost(nodename, c.NetworkName, c.Name)
	if err == nil {
		err = c.deletenodeentry(nodename)
	}

	return err
}
