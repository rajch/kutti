package vboxdriver

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// VBoxVMDriver implements the VMDriver interface for VirtualBox
type VBoxVMDriver struct {
	vboxmanagepath string
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

/*
	This parses the list of NAT networks returned by VBoxManage
	As of VBoxManage 6.0.8r130520, the format is as between the lines:
	-----------------------------------------------------------------
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
*/
func listnetworks(vboxpath string) error {
	output, err := runwithresults(vboxpath, "natnetwork", "list", "pu")
	if err != nil {
		return err
	}

	fmt.Println(output)
	lines := strings.Split(output, "\n")
	numlines := len(lines)
	if numlines < 4 {
		// Bare mininum output should be
		//   NAT Networks:
		//
		//   0 networks found
		//
		return errors.New("Could not recognise VBoxManage output for natnetworks list while getting lines")
	}

	var numnetworks int
	fmt.Println(lines, "\n", lines[numlines-2])
	_, err = fmt.Sscanf(lines[numlines-2], "%d", &numnetworks)
	if err != nil {
		return errors.New("Could not recognise VBoxManage output for natnetworks list while getting count")
	}

	justlines := lines[2 : numlines-2]

	for i, line := range justlines {
		log.Printf("Line %d: %v\n", i, line)
	}

	return nil
}
