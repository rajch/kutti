package core

import "fmt"

// Cluster defines a kutti Kubernetes cluster
type Cluster interface {
	Name() string
	Status() string
}

type kuttiNode struct {
	name   string
	status string
}

func (kn *kuttiNode) Name() string {
	return kn.name
}

func (kn *kuttiNode) Status() string {
	return kn.status
}

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

func (c *kuttiCluster) addnode(nodename string) error {
	newnode, err := c.driver.CreateHost(nodename, c.networkname(), len(c.hosts))
	if err != nil {
		return err
	}

	c.hosts = append(c.hosts, newnode)
	return nil
}

func (c *kuttiCluster) getnode(nodename string, position int) error {
	newnode, err := c.driver.GetHost(nodename)
	if err != nil {
		return err
	}

	c.hosts[position] = newnode
	return nil
}

func (c *kuttiCluster) Name() string {
	return c.name
}

func (c *kuttiCluster) Status() string {
	return c.status
}

// NewCluster creates a new cluster object
func NewCluster(name string, driver VMDriver) (Cluster, error) {
	result := &kuttiCluster{name: name, driver: driver, status: "Created"}

	err := result.ensurenetwork()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// LoadCluster creates an existing cluster
func LoadCluster(name string, driver VMDriver, hostnames ...string) (Cluster, error) {
	result := &kuttiCluster{name: name, driver: driver, status: "Fetched"}

	err := result.ensurenetwork()
	if err != nil {
		return nil, err
	}

	// TODO: make this better
	result.hosts = make([]VMHost, len(hostnames))
	for i, value := range hostnames {
		if err := result.getnode(value, i); err != nil {
			return nil, err
		}
	}

	return result, nil
}
