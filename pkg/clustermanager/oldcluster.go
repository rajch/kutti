package clustermanager

// import (
// 	"fmt"
// )

// // Cluster defines a kutti Kubernetes cluster
// type Cluster struct {
// 	name       string
// 	driver     VMDriver
// 	k8sversion string

// 	network VMNetwork
// 	hosts   []VMHost

// 	status string
// }

// func (c *Cluster) isinited() error {
// 	if c.name == "" {
// 		return fmt.Errorf("Cluster has no name")
// 	}

// 	if c.driver == nil {
// 		return fmt.Errorf("VM driver not configured")
// 	}

// 	return nil
// }

// func (c *Cluster) networkname() string {
// 	return c.name + "net"
// }

// func (c *Cluster) ensurenetwork() error {

// 	if c.network != nil {
// 		return nil
// 	}

// 	// Get network name
// 	netname := c.networkname()

// 	// Try to get network
// 	nw, err := c.driver.GetNetwork(netname)
// 	if err == nil {
// 		c.network = nw
// 		return nil
// 	}

// 	// Doesn't exist, create it
// 	nw, err = c.driver.CreateNetwork(netname)
// 	if err != nil {
// 		return err
// 	}

// 	c.network = nw
// 	return nil
// }

// func (c *Cluster) addnode(nodename string) (VMHost, error) {
// 	newnode, err := c.driver.CreateHost(nodename, c.networkname(), len(c.hosts), c.k8sversion)
// 	if err != nil {
// 		return nil, err
// 	}

// 	c.hosts = append(c.hosts, newnode)
// 	return newnode, nil
// }

// func (c *Cluster) Name() string {
// 	return c.name
// }

// func (c *Cluster) K8sVersion() string {
// 	return c.k8sversion
// }

// func (c *Cluster) Status() string {
// 	return c.status
// }

// // NewCluster creates a new cluster object
// func NewCluster(driver VMDriver, name string, k8sversion string, masternodename string, mastermappedport int) (ICluster, error) {
// 	result := &Cluster{name: name, driver: driver, status: "Created"}

// 	// Ensure network is created
// 	err := result.ensurenetwork()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create master node
// 	masternode, err := result.addnode(masternodename)
// 	if err != nil {
// 		return nil, fmt.Errorf("Could not create master node for cluster %s: %v", name, err)
// 	}

// 	// Start master node and wait
// 	err = masternode.Start()
// 	if err != nil {
// 		result.status = "Error"
// 		return result, fmt.Errorf("Cluster created, but could not start master node %s:%v", masternodename, err)
// 	}
// 	masternode.WaitForStateChange(20)

// 	// Forward SSH port
// 	err = masternode.ForwardSSHPort(mastermappedport)
// 	if err != nil {
// 		result.status = "Error"
// 		return result, fmt.Errorf("Cluster and master node created, but could not forward SSH port:%v", err)
// 	}

// 	// TODO: rename master node
// 	output, err := DefaultClient.RunWithResults(
// 		masternode.SSHAddress(),
// 		fmt.Sprintf(commandSetHostName, masternodename),
// 	)
// 	if err != nil {
// 		result.status = "Error"
// 		return result, fmt.Errorf("Cluster and master node created, but could not rename master node to %s at address %s:Error:%v:Output:%s", masternodename, masternode.SSHAddress(), err, output)
// 	}

// 	// TODO: run kubeadm init in master node

// 	// TODO: add network

// 	// TODO: add local provisioner

// 	result.status = "Ready"
// 	return result, nil
// }
