package clustermanager

import (
	"github.com/rajch/kutti/pkg/core"
)

// Driver is a kutti driver
type Driver struct {
	vmdriver core.VMDriver
}

// Name is the name of the driver
func (d *Driver) Name() string {
	return d.vmdriver.Name()
}

// Description is a one-line decription of the driver
func (d *Driver) Description() string {
	return d.vmdriver.Description()
}

// RequiresPortForwarding specifies if the driver's networks use NAT,
// and therefore require host ports to be forwarded
func (d *Driver) RequiresPortForwarding() bool {
	return d.vmdriver.RequiresPortForwarding()
}

// Status returns the driver status
func (d *Driver) Status() string {
	return d.vmdriver.Status()
}

// ForEachVersion iterates over available versions for this driver
func (d *Driver) ForEachVersion(f func(*Version) bool) error {
	driver := d.vmdriver

	images, err := driver.ListImages()
	if err != nil {
		return err
	}

	for _, value := range images {
		version := &Version{image: value}
		cancel := f(version)
		if cancel {
			break
		}
	}

	return nil
}
