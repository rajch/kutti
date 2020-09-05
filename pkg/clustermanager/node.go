package clustermanager

import (
	"github.com/rajch/kutti/internal/pkg/kuttilog"
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
	Ports       map[int]int
}

// Cluster returns the cluster this node belongs to
func (n *Node) Cluster() *Cluster {
	//fmt.Printf("BEFORE: Node:%+v]\n", n)
	if n.cluster == nil {
		n.cluster = config.Clusters[n.ClusterName]
		n.cluster.ensuredriver()
	}
	//fmt.Printf("AFTER: Node:%+v]\n", n)
	return n.cluster
}

// Status returns the current node status
func (n *Node) Status() string {
	err := n.ensurehost()
	if err != nil {
		return "Unknown"
	}

	return n.host.Status()
}

// Start starts a node
func (n *Node) Start() error {
	err := n.ensurehost()
	if err != nil {
		return err
	}

	if n.Status() == "Stopped" {
		err = n.host.Start()
		if err != nil {
			return err
		}

		n.host.WaitForStateChange(25)
		return nil
	}

	return errNodeCannotStart
}

// Stop starts a node
func (n *Node) Stop() error {
	err := n.ensurehost()
	if err != nil {
		return err
	}

	if n.Status() == "Running" {
		err = n.host.Stop()
		if err != nil {
			return err
		}

		n.host.WaitForStateChange(25)
		return nil
	}

	return errNodeCannotStop

}

// ForceStop stops a node forcibly
func (n *Node) ForceStop() error {
	err := n.ensurehost()
	if err != nil {
		return err
	}

	if n.Status() == "Running" {
		err = n.host.ForceStop()
		if err != nil {
			return err
		}

		// TODO: Consider moving this wait, or standardize the duration
		kuttilog.Print(2, "Waiting for node to stop...")
		n.host.WaitForStateChange(25)
		kuttilog.Println(2, "Done.")
		return nil
	}

	return errNodeCannotStop
}

// ForwardSSHPort forwards the node's SSH port
func (n *Node) ForwardSSHPort(hostport int) error {
	return n.ForwardPort(hostport, 22)
}

// ForwardPort forwards a port of the node to the specified host port
func (n *Node) ForwardPort(hostport int, nodeport int) error {
	err := n.Cluster().ensuredriver()
	if err != nil {
		return err
	}

	if !n.Cluster().driver.RequiresPortForwarding() {
		return errPortForwardNotSupported
	}

	err = n.ensurehost()
	if err != nil {
		return err
	}

	err = n.CheckHostport(hostport)
	if err != nil {
		return err
	}

	_, ok := n.Ports[nodeport]
	if ok {
		return errPortAlreadyUsed
	}

	err = n.host.ForwardPort(hostport, nodeport)
	if err != nil {
		return err
	}

	n.Ports[nodeport] = hostport
	return clusterconfigmanager.Save()
}

// UnforwardPort removes any mapping of the specified node port
func (n *Node) UnforwardPort(nodeport int) error {
	cluster := n.Cluster()
	err := cluster.ensuredriver()
	if err != nil {
		return err
	}

	if !cluster.driver.RequiresPortForwarding() {
		return errPortForwardNotSupported
	}

	if nodeport == 22 {
		return errPortCannotUnmap
	}

	_, ok := n.Ports[nodeport]
	if !ok {
		return errPortNotForwarded
	}

	err = n.ensurehost()
	if err != nil {
		return err
	}

	err = n.host.UnforwardPort(nodeport)
	if err != nil {
		return err
	}

	delete(n.Ports, nodeport)
	return clusterconfigmanager.Save()
}

func (n *Node) createhost() error {
	c := n.Cluster()
	host, err := c.driver.CreateHost(n.Name, c.NetworkName, c.Name, c.K8sVersion)
	if err != nil {
		n.host = nil
		return err
	}
	n.host = host
	return nil
}

func (n *Node) ensurehost() error {
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

// CheckHostport checks if a host port is occupied in the current cluster.
func (n *Node) CheckHostport(hostport int) error {
	c := n.Cluster()
	return c.CheckHostport(hostport)
}
