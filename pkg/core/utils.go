package core

import (
	"fmt"
	"os"
	"path"
)

// RegisterDriver registers a driver with a name to core
func RegisterDriver(name string, f func() (VMDriver, error)) {
	drivers[name] = f
}

// NewDriver returns a VMDriver corresponding to the name
func NewDriver(name string) (VMDriver, error) {
	f, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("Driver '%s' not found", name)
	}

	newdriver, err := f()
	if err != nil {
		return nil, err
	}

	return newdriver, nil
}

// CacheDir returns the location where the kutti cache should reside
func CacheDir() (result string, err error) {
	result, err = os.UserCacheDir()
	if err != nil {
		return
	}

	result = path.Join(result, "kutti")
	_, err = os.Stat(result)
	if os.IsNotExist(err) {
		err = os.Mkdir(result, 0755)
		if err != nil {
			result = ""
		}
	}

	return
}
