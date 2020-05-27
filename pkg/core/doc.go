// Package core contains interfaces which define core kutti functionality,
// and some utility functions.
// Package core also provides a central place for drivers to register themselves. All drivers
// should call the RegisterDriver function with a unique name on init.
//
// The interfaces are:
//
// VMDriver
//
// This defines the interface for kutti "drivers". Each driver should be able to
// manage:
//
// - VM Hosts, which represent Kubernetes nodes
//
// - VM Networks, which connect VM Hosts and manage DHCP etc
//
// - VM Images, which allow templated creation of VM Hosts
//
// VMNetwork
//
// This defines a private network to which all nodes in a cluster will be connected.
// The network should allow connectivity between nodes, and public internet connectivity.
// For now, only IPv4 capability is assumed.
//
// VMHost
//
// This defines a host that will act as a Kubernetes node. The host should allow start,
// stop, force stop, and wait operations, and provide a way to connect to it via SSH.
//
// VMImage
//
// This defines an "image" from which a VMHost can be created. An image should have a
// unique name, a Kubernetes version, and a checksum facility.
//
// Driver Management
//
package core
