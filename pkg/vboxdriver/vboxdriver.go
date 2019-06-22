package vboxdriver

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rajch/kutti/pkg/core"
)

var (
	// DefaultNetCIDR is the address range used by virtual networks
	DefaultNetCIDR = "10.0.2.0/24"
)

// VBoxVMDriver implements the VMDriver interface for VirtualBox
type VBoxVMDriver struct {
	vboxmanagepath string
}

/*ListNetworks parses the list of NAT networks returned by VBoxManage
  As of VBoxManage 6.0.8r130520, the format is as between the lines:
  ------------------------------------------------------------------
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
  ------------------------------------------------------------------
  If there are zero networks, the output is:
  ------------------------------------------------------------------
  NAT Networks:

  0 networks found

  ------------------------------------------------------------------
  The dotted lines shown here are not part of the output- they delimit
  it for documentation. Note the blank lines: one before and after
  each network.
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
		result[j] = core.VMNetwork{Name: justlines[i][13:], NetCIDR: justlines[i+1][13:]}
	}

	return result, nil
}

// CreateNetwork creates a network
func (vd *VBoxVMDriver) CreateNetwork(netname string) (result core.VMNetwork, err error) {
	// Multiple VirtualBox NAT Networks can have the same IP range
	_, err = runwithresults(vd.vboxmanagepath,
		"natnetwork",
		"add",
		"--netname",
		netname,
		"--network",
		DefaultNetCIDR,
	)
	if err != nil {
		return
	}

	result.Name = netname
	result.NetCIDR = DefaultNetCIDR

	return
}

// DeleteNetwork deletes a network
func (vd *VBoxVMDriver) DeleteNetwork(netname string) error {
	_, err := runwithresults(vd.vboxmanagepath,
		"natnetwork",
		"remove",
		"--netname",
		netname,
	)

	return err
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

	// Then try looking in well-known places

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
