package core

import (
	"os"
	"path"
)

func ensureDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
	}
	return err
}

// CacheDir returns the location where the kutti cache should reside.
// The kutti cache contains images.
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

// ConfigDir returns the location where the kutti config files reside.
// The kutti config files include driver-specific config files and image lists.
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
