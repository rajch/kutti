package clustermanager

import "github.com/rajch/kutti/pkg/core"

// ForEachDriver iterates over drivers
func ForEachDriver(f func(core.VMDriver) bool) {
	core.ForEachDriver(f)
}

// GetDriver gets the specified driver OR an error
func GetDriver(drivername string) (core.VMDriver, bool) {
	return core.GetDriver(drivername)
}

// ForEachImage iterates over images for the specified driver
func ForEachImage(drivername string, f func(core.VMImage) bool) error {
	driver, ok := core.GetDriver(drivername)
	if !ok {
		return errDriverDoesNotExist
	}

	images, err := driver.ListImages()
	if err != nil {
		return err
	}

	for _, value := range images {
		cancel := f(value)
		if cancel {
			break
		}
	}

	return nil
}