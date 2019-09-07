package core

import (
	"fmt"
)

type kuttiCluster struct {
	name   string
	driver VMDriver

	network VMNetwork
	hosts   []VMHost

	status string
}

func (c *kuttiCluster) isinited() error {
	if c.name == "" {
		return fmt.Errorf("Cluster has no name")
	}

	if c.driver == nil {
		return fmt.Errorf("VM driver not configured")
	}

	return nil
}

func (c *kuttiCluster) networkname() string {
	return c.name + "net"
}

func (c *kuttiCluster) ensurenetwork() error {

	if c.network != nil {
		return nil
	}

	// Get network name
	netname := c.networkname()

	// Try to get network
	nw, err := c.driver.GetNetwork(netname)
	if err == nil {
		c.network = nw
		return nil
	}

	// Doesn't exist, create it
	nw, err = c.driver.CreateNetwork(netname)
	if err != nil {
		return err
	}

	c.network = nw
	return nil
}

func (c *kuttiCluster) addnode(nodename string) (VMHost, error) {
	newnode, err := c.driver.CreateHost(nodename, c.networkname(), len(c.hosts))
	if err != nil {
		return nil, err
	}

	c.hosts = append(c.hosts, newnode)
	return newnode, nil
}

func (c *kuttiCluster) Name() string {
	return c.name
}

func (c *kuttiCluster) Status() string {
	return c.status
}

// NewCluster creates a new cluster object
func NewCluster(driver VMDriver, name string, masternodename string, mastermappedport int) (Cluster, error) {
	result := &kuttiCluster{name: name, driver: driver, status: "Created"}

	// Ensure network is created
	err := result.ensurenetwork()
	if err != nil {
		return nil, err
	}

	// Create master node
	masternode, err := result.addnode(masternodename)
	if err != nil {
		return nil, fmt.Errorf("Could not create master node for cluster %s: %v", name, err)
	}

	// Start master node and wait
	err = masternode.Start()
	if err != nil {
		return result, fmt.Errorf("Cluster created, but could not start master node %s:%v", masternodename, err)
	}
	masternode.WaitForStateChange(20)

	// Forward SSH port
	err = masternode.ForwardSSHPort(mastermappedport)
	if err != nil {
		return result, fmt.Errorf("Could not forward SSH port:%v", err)
	}

	// TODO: rename master node
	output, err := DefaultClient.RunWithResults(
		masternode.SSHAddress(),
		fmt.Sprintf(commandSetHostName, masternodename),
	)
	if err != nil {
		return result, fmt.Errorf("Could not rename master node to %s at address %s:Error:%v:Output:%s", masternodename, masternode.SSHAddress(), err, output)
	}

	// TODO: run kubeadm init in master node

	// TODO: add network

	// TODO: add local provisioner

	return result, nil
}
