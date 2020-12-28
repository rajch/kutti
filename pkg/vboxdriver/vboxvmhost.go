package vboxdriver

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	propSSHAddress     = "/kutti/VMInfo/SSHAddress"
	propIPAddress      = "/VirtualBox/GuestInfo/Net/0/V4/IP"
	propIPAddress2     = "/VirtualBox/GuestInfo/Net/1/V4/IP"
	propIPAddress3     = "/VirtualBox/GuestInfo/Net/2/V4/IP"
	propLoggedInUsers  = "/VirtualBox/GuestInfo/OS/LoggedInUsers"
	propSavedIPAddress = "/kutti/VMInfo/SavedIPAddress"

	vboxUsername = "kuttiadmin"
	vboxPassword = "Pass@word1"
)

// VBoxVMHost implements the VMHost interface for VirtualBox
type VBoxVMHost struct {
	driver *VBoxVMDriver

	name        string
	netname     string
	clustername string
	status      string
}

func (vh *VBoxVMHost) qname() string {
	return vh.clustername + "-" + vh.name
}

// Name is the name of the VM host.
func (vh *VBoxVMHost) Name() string {
	return vh.name
}

// Status can be Stopped, Running or Error.
func (vh *VBoxVMHost) Status() string {
	return vh.status
}

// SSHAddress returns the address and port number to SSH into this VM host.
func (vh *VBoxVMHost) SSHAddress() string {
	result := vh.sshAddress()
	return result
}

// Start starts a VM host.
// It does this by running the command:
//   VBoxManage startvm <hostname> --type headless
// Note that a VM host may not be ready for further operations at the end of this.
// See WaitForStateChange.
func (vh *VBoxVMHost) Start() error {
	output, err := runwithresults(
		vh.driver.vboxmanagepath,
		"startvm",
		vh.qname(),
		"--type",
		"headless",
	)

	if err != nil {
		return fmt.Errorf("Could not start the host '%s': %v. Output was %s", vh.name, err, output)
	}

	vh.status = "Running"
	return nil
}

// Stop stops a VM host.
// It does this by running the command:
//   VBoxManage controlvm <hostname> acpipowerbutton
// Note that a VM host may not be ready for further operations at the end of this.
// See WaitForStateChange.
func (vh *VBoxVMHost) Stop() error {
	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"controlvm",
		vh.qname(),
		"acpipowerbutton",
	)

	if err != nil {
		return fmt.Errorf("Could not stop the host '%s': %v", vh.name, err)
	}

	// Big risk. Deleteing the LoggedInUser property that is used
	// to check running status. Should be ok, because starting a
	// VM host is supposed to recreate that property.
	vh.unsetproperty(propLoggedInUsers)

	vh.status = "Stopped"
	return nil
}

// ForceStop stops a VM host forcibly.
// It does this by running the command:
//   VBoxManage controlvm <hostname> poweroff
func (vh *VBoxVMHost) ForceStop() error {
	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"controlvm",
		vh.qname(),
		"poweroff",
	)

	if err != nil {
		return fmt.Errorf("Could not force stop the host '%s': %v", vh.name, err)
	}

	// Big risk. Deleteing the LoggedInUser property that is used
	// to check running status. Should be ok, because starting a
	// VM host is supposed to recreate that property.
	vh.unsetproperty(propLoggedInUsers)

	vh.status = "Stopped"
	return nil
}

// WaitForStateChange waits the specified number of seconds,
// or until the VM host status changes from Stopped to Running or vice versa.
// It does this by running the command:
//   VBoxManage guestproperty wait <hostname> /VirtualBox/GuestInfo/OS/LoggedInUsers --timeout <milliseconds> --fail-on-timeout
// WaitForStateChange should be called after a call to Start, before
// any other operation. From observation, it should not be called before Stop.
func (vh *VBoxVMHost) WaitForStateChange(timeoutinseconds int) {
	_, _ = runwithresults(
		vh.driver.vboxmanagepath,
		"guestproperty",
		"wait",
		vh.qname(),
		propLoggedInUsers,
		"--timeout",
		fmt.Sprintf("%v", timeoutinseconds*1000),
		"--fail-on-timeout",
	)
}

func (vh *VBoxVMHost) forwardingrulename(vmport int) string {
	return fmt.Sprintf("Node %s Port %d", vh.qname(), vmport)
}

// ForwardPort creates a rule to forward the specified VM host port to the
// specified physical host port. It does this by running the command:
//   VBoxManage natnetwork modify --netname <networkname> --port-forward-4 <rule>
// Port forwarding rule format is:
//   <rule name>:<protocol>:[<host ip>]:<host port>:[<guest ip>]:<guest port>
// The brackets [] are to be taken literally.
// This driver writes the rule name as "Node <nodename> Port <vmport>".
func (vh *VBoxVMHost) ForwardPort(hostport int, vmport int) error {
	// Modify NAT Network port forwarding rules to forward vm port to host port
	// Port forwarding rule format is:
	// <rule name>:<protocol>:[<host ip>]:<host port>:[<guest ip>]:<guest port>
	// The brackets [] are to be taken literally.
	forwardingrule := fmt.Sprintf(
		"%s:tcp:[]:%d:[%s]:%d",
		vh.forwardingrulename(vmport),
		hostport,
		vh.savedipAddress(),
		vmport,
	)

	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"natnetwork",
		"modify",
		"--netname",
		vh.netname,
		"--port-forward-4",
		forwardingrule,
	)

	if err != nil {
		return fmt.Errorf(
			"Could not create port forwarding rule %s for node %s on network %s: %v",
			forwardingrule,
			vh.name,
			vh.netname,
			err,
		)
	}

	return nil
}

// UnforwardPort removes the rule which forwarded the specified VM host port.
// It does this by running the command:
//   VBoxManage natnetwork modify --netname <networkname> --port-forward-4 delete <rulename>
// This driver writes the rule name as "Node <nodename> Port <vmport>".
func (vh *VBoxVMHost) UnforwardPort(vmport int) error {
	rulename := vh.forwardingrulename(vmport)
	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"natnetwork",
		"modify",
		"--netname",
		vh.netname,
		"--port-forward-4",
		"delete",
		rulename,
	)
	if err != nil {
		return fmt.Errorf(
			"Error while removing port forwarding rule %s for VM %s on network %s: %v",
			rulename,
			vh.name,
			vh.netname,
			err,
		)
	}

	return nil
}

// ForwardSSHPort forwards the SSH port of a VM host to the specified
// physical host port. See ForwardPort for details.
func (vh *VBoxVMHost) ForwardSSHPort(hostport int) error {

	err := vh.ForwardPort(hostport, 22)
	if err != nil {
		return fmt.Errorf(
			"Could not create SSH port forwarding rule for node %s on network %s: %v",
			vh.name,
			vh.netname,
			err,
		)
	}

	sshaddress := fmt.Sprintf("localhost:%d", hostport)

	err = vh.setproperty(
		propSSHAddress,
		sshaddress,
	)
	if err != nil {
		return fmt.Errorf(
			"Could not save SSH address for node %s : %v",
			vh.name,
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
		vh.qname(),
		propname,
	)

	// VBoxManage guestproperty gets the hardcoded value "No value set!"
	// if the property value cannot be retrieved
	if err != nil || output == "No value set!" || output == "No value set!\n" {
		return "", false
	}

	// Output is in the format
	// Value: <value>
	// So, 7th rune onwards
	return output[7:], true
}

func (vh *VBoxVMHost) setproperty(propname string, value string) error {
	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"guestproperty",
		"set",
		vh.qname(),
		propname,
		value,
	)

	if err != nil {
		return fmt.Errorf(
			"Could not set property %s for host %s: %v",
			propname,
			vh.name,
			err,
		)
	}

	return nil
}

func (vh *VBoxVMHost) unsetproperty(propname string) error {
	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"guestproperty",
		"unset",
		vh.qname(),
		propname,
	)

	if err != nil {
		return fmt.Errorf(
			"Could not unset property %s for host %s: %v",
			propname,
			vh.name,
			err,
		)
	}

	return nil
}

func trimpropend(s string) string {
	return strings.TrimSpace(s)
}

func (vh *VBoxVMHost) ipAddress() string {
	// This guestproperty is only available if the VM is
	// running, and has the Virtual Machine additions enabled
	result, _ := vh.getproperty(propIPAddress)
	return trimpropend(result)
}

func (vh *VBoxVMHost) savedipAddress() string {
	// This guestproperty is set when the VM is created
	result, _ := vh.getproperty(propSavedIPAddress)
	return trimpropend(result)
}

func (vh *VBoxVMHost) sshAddress() string {
	result, _ := vh.getproperty(propSSHAddress)
	return trimpropend(result)
}

// VirtualBox properties, and correspoding actions
var propMap = map[string]func(*VBoxVMHost, string){
	propLoggedInUsers: func(vh *VBoxVMHost, value string) {
		vh.status = "Running"
	},
}

func (vh *VBoxVMHost) parseProps(propstr string) {
	// There are two possibilities. Either:
	// VBoxManage: error: Could not find a registered machine named 'xxx'
	// ...
	// Or:
	// Name: /VirtualBox/GuestInfo/Net/0/V4/IP, value: 10.0.2.15, timestamp: 1568552111298588000, flags:
	// ...
	r1, _ := regexp.Compile("error: (.*)\n")
	r2, _ := regexp.Compile("Name: (.*), value: (.*), timestamp: (.*), flags:(.*)\n")

	// This should not have made it this far. Still,
	// belt and suspenders...
	errorsfound := r1.FindAllStringSubmatch(propstr, 1)
	if len(errorsfound) != 0 {
		// deal with the error with:
		// errorsfound[0][1]
		vh.status = "Error:" + errorsfound[0][1]
		return
	}

	results := r2.FindAllStringSubmatch(propstr, -1)
	for _, record := range results {
		// Parse each line with
		// record[1] - Name and record[2] - Value
		f, ok := propMap[record[1]]
		if ok {
			f(vh, record[2])
		}
	}
}

// runwithresults allows running commands inside a VM Host.
// It does this by running the command:
// - VBoxManage guestcontrol <hostname> --username <username> --password <password> run -- <command line>
// This requires Virtual Machine Additions to be running in the guest operating system.
// The guest OS should be fully booted up.
func (vh *VBoxVMHost) runwithresults(execpath string, paramarray ...string) (string, error) {
	params := []string{

		"guestcontrol",
		vh.qname(),
		"--username",
		vboxUsername,
		"--password",
		vboxPassword,
		"run",
		"--",
		execpath,
	}
	params = append(params, paramarray...)

	output, err := runwithresults(
		vh.driver.vboxmanagepath,
		params...,
	)

	return output, err
}

func (vh *VBoxVMHost) renamehost(newname string) error {
	execname := fmt.Sprintf("/home/%s/kutti-installscripts/set-hostname.sh", vboxUsername)

	_, err := vh.runwithresults(
		"/usr/bin/sudo",
		execname,
		newname,
	)

	return err
}
