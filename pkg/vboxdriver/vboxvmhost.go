package vboxdriver

import (
	"fmt"
	"strings"
)

const (
	propSSHRule = "/kutti/VMInfo/SSHForwardingRule"
)

// VBoxVMHost implements the VMHost interface for VirtualBox
type VBoxVMHost struct {
	driver *VBoxVMDriver

	name    string
	netname string
	status  string

	hostport int
}

// Name is the name of the host
func (vh *VBoxVMHost) Name() string {
	return vh.name
}

// Status can be "Created", "Fetched" or "Ready"
func (vh *VBoxVMHost) Status() string {
	return vh.status
}

// Start starts the VM
func (vh *VBoxVMHost) Start() error {
	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"startvm",
		vh.name,
		"--type",
		"headless",
	)

	if err != nil {
		return fmt.Errorf("Could not start the host '%s': %v", vh.name, err)
	}

	vh.status = "Running"
	return nil
}

// Stop stops the VM
func (vh *VBoxVMHost) Stop() error {
	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"controlvm",
		vh.name,
		"acpipowerbutton",
	)

	if err != nil {
		return fmt.Errorf("Could not stop the host '%s': %v", vh.name, err)
	}

	vh.status = "Stopped"
	return nil
}

func (vh *VBoxVMHost) ipAddress() string {
	// This guestproperty is only available if the VM is
	// running, and has the Virtual Machine additions enabled
	result, _ := vh.getproperty("/VirtualBox/GuestInfo/Net/1/V4/IP")
	return result
}

func (vh *VBoxVMHost) sshRule() string {
	// This guestproperty is only available if the VM is
	// running, and has the Virtual Machine additions enabled
	result, _ := vh.getproperty("/kutti/VMInfo/SSHForwardingRule")
	return result

}
func (vh *VBoxVMHost) forwardedPort() string {
	result := vh.sshRule()

	ruleparts := strings.Split(result, ":")
	if len(ruleparts) < 4 {
		return ""
	}

	return ruleparts[3]
}

func (vh *VBoxVMHost) forwardSSHPort(hostport int) error {
	// Modify NAT Network port forwarding rules to allow SSH to new host
	// SSH rule format is:
	// <rule name>:<protocol>:[<host ip>]:<host port>:[<guest ip>]:<guest port>
	// The brackets [] are to be taken literally.
	sshrule := fmt.Sprintf(
		"SSH Node %s:tcp:[]:%d:[%s]:22",
		vh.name,
		hostport,
		vh.ipAddress(),
	)

	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"natnetwork",
		"modify",
		"--netname",
		vh.netname,
		"--port-forward-4",
		sshrule,
	)

	if err != nil {
		return fmt.Errorf(
			"Could not create SSH port forwarding rule %s for node %s on network %s: %v",
			sshrule,
			vh.name,
			vh.netname,
			err,
		)
	}

	_, err = runwithresults(
		vh.driver.vboxmanagepath,
		"guestproperty",
		"set",
		vh.name,
		propSSHRule,
		sshrule,
	)
	if err != nil {
		return fmt.Errorf(
			"Could not save SSH port forwarding rule %s for node %s : %v",
			sshrule,
			vh.name,
			err,
		)
	}

	return nil
}

func (vh *VBoxVMHost) removeSSHPort() error {
	sshrule, _ := vh.getproperty(propSSHRule)
	if sshrule == "" {
		return nil
	}

	ruleparts := strings.Split(sshrule, ":")
	if len(ruleparts) < 1 {
		return fmt.Errorf(
			"Error while removing ssh rule %s for VM %s on network %s: the saved rule is invalid",
			sshrule,
			vh.name,
			vh.netname,
		)
	}

	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"natnetwork",
		"modify",
		"netname",
		vh.netname,
		"--port-forward-4",
		"delete",
		ruleparts[0],
	)
	if err != nil {
		return fmt.Errorf(
			"Error while removing ssh rule %s for VM %s on network %s: %v",
			sshrule,
			vh.name,
			vh.netname,
			err,
		)
	}

	return nil
}

func (vh *VBoxVMHost) getproperty(propname string) (string, bool) {
	output, err := runwithresults(
		vh.driver.vboxmanagepath,
		"guestproperty",
		"get",
		vh.name,
		propname,
	)

	// VBoxManage guestproperty gets the hardcoded value "No value set!"
	// if the property value cannot be retrieved
	if err != nil || output == "No value set!" {
		return "", false
	}

	return output, true
}
