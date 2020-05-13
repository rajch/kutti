package vboxdriver

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	propSSHAddress    = "/kutti/VMInfo/SSHAddress"
	propIPAddress     = "/VirtualBox/GuestInfo/Net/0/V4/IP"
	propLoggedInUsers = "/VirtualBox/GuestInfo/OS/LoggedInUsers"

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

// Name is the name of the host
func (vh *VBoxVMHost) Name() string {
	return vh.name
}

// Status can be ~"Created", "Fetched" or "Ready"~
// Stopped, Running
func (vh *VBoxVMHost) Status() string {
	return vh.status
}

// SSHAddress returns the address and port number to SSH into this host
func (vh *VBoxVMHost) SSHAddress() string {
	result := strings.Trim(vh.sshAddress(), " \n")
	return result
	// ruleparts := strings.Split(result, ":")
	// if len(ruleparts) < 4 {
	// 	return ""
	// }

	// return fmt.Sprintf("localhost:%s", ruleparts[3])
}

// Start starts the VM
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

// Stop stops the VM
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

	vh.status = "Stopped"
	return nil
}

// ForceStop stops the VM forcibly
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

	vh.status = "Stopped"
	return nil
}

// WaitForStateChange waits the specified number of seconds
func (vh *VBoxVMHost) WaitForStateChange(timeoutinseconds int) {
	vh.vboxwaitforstatechange(timeoutinseconds)
	// time.Sleep(time.Duration(timeoutinseconds) * time.Second)
}

func (vh *VBoxVMHost) forwardingrulename(vmport int) string {
	return fmt.Sprintf("Node %s Port %d", vh.qname(), vmport)
}

// ForwardPort creates a rule to forward the specified vm port to the
// specified host port
func (vh *VBoxVMHost) ForwardPort(hostport int, vmport int) error {
	// Modify NAT Network port forwarding rules to forward vm port to host port
	// Port forwarding rule format is:
	// <rule name>:<protocol>:[<host ip>]:<host port>:[<guest ip>]:<guest port>
	// The brackets [] are to be taken literally.
	forwardingrule := fmt.Sprintf(
		"%s:tcp:[]:%d:[%s]:%d",
		vh.forwardingrulename(vmport),
		hostport,
		vh.ipAddress(),
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

// UnforwardPort removes the rule which forwarded the specified vm port
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

// ForwardSSHPort forwards the SSH port of this host to the specified host port
func (vh *VBoxVMHost) ForwardSSHPort(hostport int) error {
	// Modify NAT Network port forwarding rules to allow SSH to new host
	// SSH rule format is:
	// <rule name>:<protocol>:[<host ip>]:<host port>:[<guest ip>]:<guest port>
	// The brackets [] are to be taken literally.
	// sshrule := fmt.Sprintf(
	// 	"SSH Node %s:tcp:[]:%d:[%s]:22",
	// 	vh.qname(),
	// 	hostport,
	// 	vh.ipAddress(),
	// )

	// _, err := runwithresults(
	// 	vh.driver.vboxmanagepath,
	// 	"natnetwork",
	// 	"modify",
	// 	"--netname",
	// 	vh.netname,
	// 	"--port-forward-4",
	// 	sshrule,
	// )
	// sshrule := fmt.Sprintf(
	// 	"%s:tcp:[]:%d:[%s]:22",
	// 	vh.forwardingrulename(22),
	// 	hostport,
	// 	vh.ipAddress(),
	// )
	sshaddress := fmt.Sprintf("localhost:%d", hostport)

	err := vh.ForwardPort(hostport, 22)
	if err != nil {
		return fmt.Errorf(
			"Could not create SSH port forwarding rule for node %s on network %s: %v",
			vh.name,
			vh.netname,
			err,
		)
	}

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

func (vh *VBoxVMHost) ipAddress() string {
	// This guestproperty is only available if the VM is
	// running, and has the Virtual Machine additions enabled
	result, _ := vh.getproperty(propIPAddress)
	return result
}

func (vh *VBoxVMHost) sshAddress() string {
	result, _ := vh.getproperty(propSSHAddress)
	return result
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

	// fmt.Println("Not found")
	// fmt.Println("---------")
	// fmt.Println(r1.FindAllStringSubmatch(strErr, 1)[0][1])
	// fmt.Println("NOT Not found")
	// fmt.Println(r1.FindAllStringSubmatch(strStopped, 1))
	// fmt.Println("Stopped")
	// fmt.Println("---------")
	// results := r2.FindAllStringSubmatch(strStopped, -1)
	// for i := 0; i < len(results); i++ {
	// 	fmt.Printf("Name: %v, Value:%v\n", results[i][1], results[i][2])
	// }
}

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
	execname := fmt.Sprintf("/home/%s/rw-installscripts/set-hostname.sh", vboxUsername)

	_, err := vh.runwithresults(
		"/usr/bin/sudo",
		execname,
		newname,
	)

	return err
}

func (vh *VBoxVMHost) vboxwaitforstatechange(timeoutinseconds int) error {
	_, err := runwithresults(
		vh.driver.vboxmanagepath,
		"guestproperty",
		"wait",
		vh.qname(),
		propLoggedInUsers,
		"--timeout",
		fmt.Sprintf("%v", timeoutinseconds*1000),
		"--fail-on-timeout",
	)

	return err
}
