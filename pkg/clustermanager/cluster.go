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

// ValidateNodeName checks for the validity of a node name.
// It uses IsValidname to check name validity, and also checks if a node name
// already exists in the cluster.
func (c *Cluster) ValidateNodeName(name string) error {
	if !IsValidName(name) {
		return errInvalidName
	}

	// Check if name exists
	_, ok := c.Nodes[name]
	if ok {
		return errNodeExists
	}

	return nil
}

// NewUninitializedNode adds a node, but does not join it to a kubernetes cluster
func (c *Cluster) NewUninitializedNode(nodename string) (*Node, error) {
	err := c.ValidateNodeName(nodename)
	if err != nil {
		return nil, err
	}

	return c.addnode(nodename, "Unmanaged")
}

// DeleteNode deletes a node completely. By default, a node is not deleted
// if it is running. The force parameter causes the node to be stopped and
// deleted. In some rare cases for some drivers, manual cleanup may be
// needed after a forced delete.
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
	kuttilog.Println(2, "Unmapping ports...")
	for key := range n.Ports {
		err := n.host.UnforwardPort(key)
		if err != nil {
			kuttilog.Printf(2, "Error while unmapping ports for node '%s': %v.", nodename, err)
		}
	}
	kuttilog.Println(2, "Ports unmapped.")

	return c.deletenode(nodename, force)
}

// CheckHostport checks if a host port is occupied in the current cluster.
func (c *Cluster) CheckHostport(hostport int) error {
	for _, nodevalue := range c.Nodes {
		for _, hostportvalue := range nodevalue.Ports {
			if hostportvalue == hostport {
				return errPortAlreadyUsed
			}
		}
	}
	return nil
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

func (c *Cluster) deletenode(nodename string, force bool) error {
	err := c.ensuredriver()
	if err != nil {
		return err
	}

	err = c.driver.DeleteHost(nodename, c.NetworkName, c.Name)
	if err == nil || force {
		err = c.deletenodeentry(nodename)
	}

	return err
}
