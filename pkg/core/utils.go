package core

import (
	"os"
	"path"
)

// RegisterDriver registers a driver with a name to core
func RegisterDriver(name string, d VMDriver) {
	if d != nil {
		drivers[name] = d
	}
}

// GetDriver returns a VMDriver corresponding to the name
func GetDriver(name string) (VMDriver, bool) {
	result, ok := drivers[name]
	return result, ok
	// newdriver, ok := drivers[name]
	// if ok {
	// 	return newdriver, nil
	// }

	// f, ok := driverfuncs[name]
	// if !ok {
	// 	return nil, fmt.Errorf("Driver '%s' not found", name)
	// }

	// newdriver, err := f()
	// if err != nil {
	// 	return nil, err
	// }

	// return newdriver, nil
}

// ForEachDriver iterates over drivers
func ForEachDriver(f func(VMDriver) bool) {
	for _, driver := range drivers {
		cancel := f(driver)
		if cancel {
			break
		}
	}
}

func ensureDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
	}
	return err
}

// CacheDir returns the location where the kutti cache should reside
func CacheDir() (result string, err error) {
	result, err = os.UserCacheDir()
	if err != nil {
		return
	}

	result = path.Join(result, "kutti")
	err = ensureDirectory(result)

	if err != nil {
		result = ""
	}
	return
}

// ConfigDir returns the location where the kutti config files reside
func ConfigDir() (result string, err error) {
	result, err = os.UserConfigDir()
	if err != nil {
		return
	}

	result = path.Join(result, "kutti")
	err = ensureDirectory(result)

	if err != nil {
		result = ""
	}

	return
}
