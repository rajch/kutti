package vboxdriver

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rajch/kutti/internal/pkg/kuttilog"

	"github.com/rajch/kutti/pkg/core"
)

const (
	driverName        = "vbox"
	driverDescription = "Kutti driver for VirtualBox >=6.0"
)

// DefaultNetCIDR is the address range used by virtual networks.
var DefaultNetCIDR = "192.168.125.0/24"

var (
	dhcpaddress       = "192.168.125.3"
	dhcpnetmask       = "255.255.255.0"
	ipNetAddr         = "192.168.125"
	iphostbase        = 10
	forwardedPortBase = 10000
)

// VBoxVMDriver implements the VMDriver interface for VirtualBox.
type VBoxVMDriver struct {
	vboxmanagepath string
	status         string
}

// Name returns the driver identifier string.
func (vd *VBoxVMDriver) Name() string {
	return driverName
}

// Description returns the driver description.
func (vd *VBoxVMDriver) Description() string {
	return driverDescription
}

// RequiresPortForwarding specifies if the driver's networks use NAT,
// and therefore require host ports to be forwarded.
func (vd *VBoxVMDriver) RequiresPortForwarding() bool {
	return true
}

// Status returns the driver status.
func (vd *VBoxVMDriver) Status() string {
	return vd.status
}

/*ListNetworks parses the list of NAT networks returned by
    VBoxManage natnetwork list
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

// CreateNetwork creates a new VirtualBox NAT network.
// It uses the CIDR common to all Kutti networks, and is dhcp-enabled at start.
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

// GetNetwork returns a network, or an error.
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

// DeleteNetwork deletes a network.
// It does this by running the command:
//   VBoxManage natnetwork remove --netname <networkname>
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
    VBoxManage list vms
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

	result := []core.VMHost{}
	lines := strings.Split(output, "\n")
	if len(lines) < 1 {
		return result, nil
	}

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
// It also starts the VM, changes the hostname, saves the IP address, and stops
// it again.
// It runs the following two VBoxManage commands, in order:
//   VBoxManage import <nodeimageovafile> --vsys 0 --vmname "<hostname>"
//   VBoxManage modifyvm "<hostname>" --nic1 natnetwork --nat-network1 <networkname>
// The first imports from an .ova file (easiest way to get fully configured VM), while
// setting the VM name. The second connects the first network interface card to
// the NAT network.
func (vd *VBoxVMDriver) CreateHost(hostname string, networkname string, clustername string, k8sversion string) (core.VMHost, error) {
	kuttilog.Println(2, "Importing image...")

	ovafile, err := imagepathfromk8sversion(k8sversion)
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(ovafile); err != nil {
		return nil, fmt.Errorf("Could not retrieve image %s: %v", ovafile, err)
	}

	l, err := runwithresults(
		vd.vboxmanagepath,
		"import",
		ovafile,
		"--vsys",
		"0",
		"--vmname",
		clustername+"-"+hostname,
		"--vsys",
		"0",
		"--group",
		"/"+clustername,
	)

	if err != nil {
		return nil, fmt.Errorf("Could not import ovafile %s: %v(%v)", ovafile, err, l)
	}

	// Attach newly created VM to NAT Network
	kuttilog.Println(2, "Attaching host to network...")
	newhost := &VBoxVMHost{driver: vd, name: hostname, netname: networkname, clustername: clustername, status: "Created"}

	_, err = runwithresults(
		vd.vboxmanagepath,
		"modifyvm",
		newhost.qname(),
		"--nic1",
		"natnetwork",
		"--nat-network1",
		networkname,
	)

	if err != nil {
		newhost.status = "Error"
		return newhost, fmt.Errorf("Could not attach node %s to network %s: %v", hostname, networkname, err)
	}

	// Start the host
	kuttilog.Println(2, "Starting host...")
	err = newhost.Start()
	if err != nil {
		return newhost, err
	}
	newhost.WaitForStateChange(25)

	// Change the name
	for renameretries := 1; renameretries < 4; renameretries++ {
		kuttilog.Printf(2, "Renaming host (attempt %v/3)...", renameretries)
		err = newhost.renamehost(hostname)
		if err == nil {
			break
		}
	}

	if err != nil {
		return newhost, err
	}
	kuttilog.Println(2, "Host renamed.")
	// Save the IP Address
	ipaddress := newhost.ipAddress()
	kuttilog.Printf(2, "Obtained IP address '%v'", ipaddress)
	newhost.setproperty(propSavedIPAddress, ipaddress)

	kuttilog.Println(2, "Stopping host...")
	newhost.Stop()
	newhost.WaitForStateChange(25)

	newhost.status = "Stopped"

	return newhost, nil
}

// GetHost returns the named host, or an error.
// It does this by running the command:
//   VBoxManage guestprorty enumerate <hostname> --patterns "/VirtualBox/GuestInfo/Net/0/*|/kutti/*|/VirtualBox/GuestInfo/OS/LoggedInUsers"
func (vd *VBoxVMDriver) GetHost(hostname string, networkname string, clustername string) (core.VMHost, error) {
	output, err := runwithresults(
		vd.vboxmanagepath,
		"guestproperty",
		"enumerate",
		clustername+"-"+hostname,
		"--patterns",
		"/VirtualBox/GuestInfo/Net/0/*|/kutti/*|/VirtualBox/GuestInfo/OS/LoggedInUsers",
	)

	if err != nil {
		return nil, fmt.Errorf("Host %s not found", hostname)
	}

	foundhost := &VBoxVMHost{driver: vd, name: hostname, netname: networkname, clustername: clustername, status: "Stopped"}

	if output != "" {
		foundhost.parseProps(output)
	}

	return foundhost, nil
}

// DeleteHost deletes a VM.
// It does this by running the command:
//   VBoxManage unregistervm "<hostname>" --delete
func (vd *VBoxVMDriver) DeleteHost(hostname string, networkname string, clustername string) error {

	output, err := runwithresults(
		vd.vboxmanagepath,
		"unregistervm",
		clustername+"-"+hostname,
		"--delete",
	)

	if err != nil {
		return fmt.Errorf("Could not delete host %s: %v:%s", hostname, err, output)
	}

	return nil
}

// ListK8sVersions lists the known Kubernetes versions.
func (vd *VBoxVMDriver) ListK8sVersions() ([]string, error) {
	imageconfigmanager.Load()

	result := make([]string, len(imagedata.images))

	i := 0
	for key := range imagedata.images {
		result[i] = key
		i++
	}

	return result, nil
}

// FetchImageList fetches the latest list of VM images.
func (vd *VBoxVMDriver) FetchImageList() error {
	return fetchimagelist()
}

// ListImages lists the known VM images.
func (vd *VBoxVMDriver) ListImages() ([]core.VMImage, error) {
	imageconfigmanager.Load()

	result := make([]core.VMImage, len(imagedata.images))

	i := 0
	for _, image := range imagedata.images {
		result[i] = image
		i++
	}

	return result, nil
}

// GetImage returns an image corresponding to a Kubernetes version,
// or an error.
func (vd *VBoxVMDriver) GetImage(k8sversion string) (core.VMImage, error) {
	imageconfigmanager.Load()
	result, ok := imagedata.images[k8sversion]
	if !ok {
		return nil, fmt.Errorf("no image found for version %v", k8sversion)
	}

	return result, nil
}

// New returns a pointer to a new VBoxVMDriver or an error.
func New() (*VBoxVMDriver, error) {
	result := &VBoxVMDriver{}

	// find VBoxManage tool and set it
	vbmpath, err := findvboxmanage()
	if err != nil {
		result.status = "Error:" + err.Error()
		return result, err
	}
	result.vboxmanagepath = vbmpath

	// test VBoxManage version
	vbmversion, err := runwithresults(vbmpath, "--version")
	if err != nil {
		result.status = "Error:" + err.Error()
		return result, err
	}
	var majorversion int
	_, err = fmt.Sscanf(vbmversion, "%d", &majorversion)
	if err != nil || majorversion < 6 {
		err = fmt.Errorf("unsupported VBoxManage version %v. 6.0 and above are supported", vbmversion)
		result.status = "Error:" + err.Error()
		return result, err
	}

	result.status = "Ready"
	return result, nil
}

func runwithresults(execpath string, paramarray ...string) (result string, err error) {
	if kuttilog.V(4) {
		kuttilog.Println(4, "[DEBUG]------------------")
		kuttilog.Println(4, "[DEBUG]Executing command:")
		kuttilog.Println(4, execpath, strings.Join(paramarray, " "))
		kuttilog.Println(4, "[DEBUG]")
	}
	cmd := exec.Command(execpath, paramarray...)
	output, err := cmd.CombinedOutput()
	/*if err != nil {
		err = fmt.Errorf("%v(%v)", string(cmd.Stderr))
		return
	}*/
	if kuttilog.V(4) {
		kuttilog.Println(4, "[DEBUG]Execution results:")
		kuttilog.Println(4, string(output))
		if err != nil {
			kuttilog.Printf(4, "Error: %v\n", err)
		}
		kuttilog.Println(4, "[DEBUG]==================")
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

func init() {
	newdriver, _ := New()
	core.RegisterDriver(driverName, newdriver)
}
