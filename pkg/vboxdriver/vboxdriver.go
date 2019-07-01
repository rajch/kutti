package vboxdriver

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/rajch/kutti/pkg/core"
)

var (
	// DefaultNetCIDR is the address range used by virtual networks
	DefaultNetCIDR = "10.0.3.0/24"

	dhcpaddress       = "10.0.3.3"
	dhcpnetmask       = "255.255.255.0"
	ipNetAddr         = "10.0.3"
	iphostbase        = 10
	forwardedPortBase = 10000
)

// VBoxVMDriver implements the VMDriver interface for VirtualBox
type VBoxVMDriver struct {
	vboxmanagepath string
}

/*ListNetworks parses the list of NAT networks returned by
`VBoxManage natnetwork list`.
As of VBoxManage 6.0.8r130520, the format is:

  NAT Networks:

  Name:        KubeNet
  Network:     10.0.2.0/24
  Gateway:     10.0.2.1
  IPv6:        No
  Enabled:     Yes


  Name:        NatNetwork
  Network:     10.0.2.0/24
  Gateway:     10.0.2.1
  IPv6:        No
  Enabled:     Yes


  Name:        NatNetwork1
  Network:     10.0.2.0/24
  Gateway:     10.0.2.1
  IPv6:        No
  Enabled:     Yes

  3 networks found

Note the blank lines: one before and after
each network. If there are zero networks, the output is:

  NAT Networks:

  0 networks found


*/
func (vd *VBoxVMDriver) ListNetworks() ([]core.VMNetwork, error) {
	// TODO: consider a default pattern for all our network names
	output, err := runwithresults(vd.vboxmanagepath, "natnetwork", "list")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	numlines := len(lines)
	if numlines < 4 {
		// Bare mininum output should be
		//   NAT Networks:
		//
		//   0 networks found
		//
		return nil, errors.New("Could not recognise VBoxManage output for natnetworks list while getting lines")
	}

	var numnetworks int

	_, err = fmt.Sscanf(lines[numlines-2], "%d", &numnetworks)
	if err != nil {
		return nil, errors.New("Could not recognise VBoxManage output for natnetworks list while getting count")
	}

	justlines := lines[2 : numlines-2]
	numlines = len(justlines)

	result := make([]core.VMNetwork, numnetworks, numnetworks)

	for i, j := 0, 0; i < numlines; i, j = i+7, j+1 {
		result[j] = &VBoxVMNetwork{name: justlines[i][13:], netCIDR: justlines[i+1][13:]}
	}

	return result, nil
}

/*CreateNetwork creates a new VirtualBox NAT network.
It uses the CIDR common to all Kutti networks, and is dhcp-enabled at start.
*/
func (vd *VBoxVMDriver) CreateNetwork(netname string) (core.VMNetwork, error) {
	// Multiple VirtualBox NAT Networks can have the same IP range
	// So, all Kutti networks will use the same network CIDR
	// We start with dhcp enabled.
	output, err := runwithresults(vd.vboxmanagepath,
		"natnetwork",
		"add",
		"--netname",
		netname,
		"--network",
		DefaultNetCIDR,
		"--enable",
		"--dhcp",
		"on",
	)
	if err != nil {
		return nil, fmt.Errorf("Could not create NAT network %s:%v:%s", netname, err, output)
	}

	// Manually create the associated DHCP server
	// Hard-coding a thirty-node limit for now
	output, err = runwithresults(
		vd.vboxmanagepath,
		"dhcpserver",
		"add",
		"--netname",
		netname,
		"--ip",
		dhcpaddress,
		"--netmask",
		dhcpnetmask,
		"--lowerip",
		fmt.Sprintf("%s.%d", ipNetAddr, iphostbase),
		"--upperip",
		fmt.Sprintf("%s.%d", ipNetAddr, iphostbase+29),
		"--enable",
	)
	if err != nil {
		return nil, fmt.Errorf("Could not create DHCP server for network %s:%v:%s", netname, err, output)
	}

	newnetwork := &VBoxVMNetwork{name: netname, netCIDR: DefaultNetCIDR}

	return newnetwork, err
}

// GetNetwork returns a network, or an error
func (vd *VBoxVMDriver) GetNetwork(netname string) (core.VMNetwork, error) {
	networks, err := vd.ListNetworks()
	if err != nil {
		return nil, err
	}

	for _, network := range networks {
		if network.Name() == netname {
			return network, nil
		}
	}

	return nil, fmt.Errorf("Network %s not found", netname)
}

// DeleteNetwork deletes a network
func (vd *VBoxVMDriver) DeleteNetwork(netname string) error {
	output, err := runwithresults(vd.vboxmanagepath,
		"natnetwork",
		"remove",
		"--netname",
		netname,
	)
	if err != nil {
		return fmt.Errorf("Could not delete NAT network %s:%v:%s", netname, err, output)
	}

	output, err = runwithresults(
		vd.vboxmanagepath,
		"dhcpserver",
		"remove",
		"--netname",
		netname,
	)
	if err != nil {
		return fmt.Errorf("Could not delete DHCP server %s:%v:%s", netname, err, output)
	}

	return nil
}

/*ListHosts parses the list of VMs returned by
`VBoxManage list vms`
As of VBoxManage 6.0.8r130520, the format is:

  "Matsya" {e3509073-d188-4cca-8eaf-cb9f3be7ac4a}
  "Krishna" {5d9b1b16-5059-42ae-a160-e93b470f940e}
  "one" {06748689-7f4e-4915-8fbf-6111596f85a2}
  "two" {eee169a7-09eb-473e-96be-5d37868c5d5e}
  "minikube" {5bf78b43-3240-4f50-911b-fbc111d4d085}
  "Node 1" {53b82a61-ae52-44c2-86d5-4c686502dd64}

*/
func (vd *VBoxVMDriver) ListHosts() ([]core.VMHost, error) {
	output, err := runwithresults(
		vd.vboxmanagepath,
		"list",
		"vms",
	)
	if err != nil {
		return nil, fmt.Errorf("Could not get list of VMs: %v", err)
	}

	lines := strings.Split(output, "\n")
	if len(lines) < 1 {
		return []core.VMHost{}, nil
	}

	result := []core.VMHost{}
	actualcount := 0
	for _, value := range lines {
		line := strings.Split(value, " ")
		if len(line) == 2 {
			result = append(result, &VBoxVMHost{
				driver:  vd,
				name:    trimQuotes(line[0]),
				netname: "",
				status:  "Fetched",
			})
			actualcount++
		}
	}

	return result[:actualcount], err

}

// CreateHost creates a VM, and connects it to a previously created NAT network.
// It also starts the VM, and creates a port fowarding rule on the network to
// forward the SSH port.
func (vd *VBoxVMDriver) CreateHost(hostname string, networkname string, position int) (core.VMHost, error) {
	/*
		We need to run the following two VBoxManage commands, in order:

		- VBoxManage import <nodeimageovafile> --vsys 0 --vmname "<hostname>"
		- VBoxManage modifyvm "<hostname>" --nic1 natnetwork --nat-network1 <networkname>

		The first imports from an .ova file (easiest way to get fully configured VM), while
		setting the VM name. The second connects the first network interface card to
		the NAT network.
	*/

	cachedir, err := core.CacheDir()
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve CacheDir: %v", err)
	}

	// TODO: ovafile hardcoded here. Correct.
	ovafile := path.Join(cachedir, "Krishna-1.0.ova")
	if _, err = os.Stat(ovafile); err != nil {
		return nil, fmt.Errorf("Could not retrieve ovafile %s: %v", ovafile, err)
	}

	_, err = runwithresults(
		vd.vboxmanagepath,
		"import",
		ovafile,
		"--vsys",
		"0",
		"--vmname",
		hostname,
	)

	if err != nil {
		return nil, fmt.Errorf("Could not import ovafile %s: %v", ovafile, err)
	}

	// Attach newly created VM to NAT Network
	newhost := &VBoxVMHost{driver: vd, name: hostname, netname: networkname, status: "Created"}

	_, err = runwithresults(
		vd.vboxmanagepath,
		"modifyvm",
		hostname,
		"--nic1",
		"natnetwork",
		"--nat-network1",
		networkname,
	)

	if err != nil {
		return newhost, fmt.Errorf("Could not attach node %s to network %s: %v", hostname, networkname, err)
	}
	newhost.status = "NetworkAttached"

	/*
		// Start the host
		err = newhost.Start()
		if err != nil {
			return newhost, err
		}
		newhost.status = "Started"

		// Forward the SSH port
		err = newhost.ForwardSSHPort(10000 + position)
		if err != nil {
			return newhost, err
		}
	*/
	newhost.status = "Ready"

	return newhost, nil
}

// GetHost returns the named host, or an error.
func (vd *VBoxVMDriver) GetHost(hostname string, networkname string) (core.VMHost, error) {
	output, err := runwithresults(
		vd.vboxmanagepath,
		"guestproperty",
		"enumerate",
		hostname,
		"--patterns",
		"/VirtualBox/GuestInfo/Net/0/*",
	)

	if err != nil {
		return nil, fmt.Errorf("Host %s not found", hostname)
	}

	foundhost := &VBoxVMHost{driver: vd, name: hostname, netname: networkname, status: "Fetched"}

	if output != "" {
		// Parse output
		foundhost.status = "Ready"
	}

	return foundhost, nil
}

// DeleteHost deletes a VM.
func (vd *VBoxVMDriver) DeleteHost(hostname string, networkname string) error {
	/*
		We need to run:

		- VBoxManage unregistervm "<nodename>" --delete
	*/
	output, err := runwithresults(
		vd.vboxmanagepath,
		"unregistervm",
		hostname,
		"--delete",
	)

	if err != nil {
		return fmt.Errorf("Could not delete host %s: %v:%s", hostname, err, output)
	}

	return nil
}

// New returns a pointer to a new VBoxVMDriver OR an error
func New() (*VBoxVMDriver, error) {
	result := &VBoxVMDriver{}

	// find VBoxManage tool and set it
	vbmpath, err := findvboxmanage()
	if err != nil {
		return nil, err
	}
	result.vboxmanagepath = vbmpath

	// test VBoxManage version
	vbmversion, err := runwithresults(vbmpath, "--version")
	if err != nil {
		return nil, err
	}
	var majorversion int
	_, err = fmt.Sscanf(vbmversion, "%d", &majorversion)
	if err != nil || majorversion < 6 {
		return nil, errors.New("Unsupported VBoxManage version " + vbmversion + ". 6.0 and above are supported.")
	}

	return result, nil
}

func findvboxmanage() (path string, err error) {

	// First, try looking it up on the path
	path, err = exec.LookPath("VBoxManage")
	if err == nil {
		return
	}

	// TODO: Then try looking in well-known places

	// Give up
	return
}

func runwithresults(execpath string, paramarray ...string) (result string, err error) {
	cmd := exec.Command(execpath, paramarray...)
	output, err := cmd.Output()
	if err != nil {
		return
	}

	result = string(output)
	return
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
