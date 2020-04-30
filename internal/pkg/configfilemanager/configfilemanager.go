package configfilemanager

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/rajch/kutti/pkg/core"
)

func getconfigfilepath(configFileName string) (string, error) {
	configPath, err := core.ConfigDir()
	if err != nil {
		return "", err
	}

	datafilepath := path.Join(configPath, configFileName)
	return datafilepath, nil
}

// Save saves the specified data into the named file in the kutti config directory.
func Save(configfilename string, data []byte) error {
	datafilepath, err := getconfigfilepath(configfilename)
	if err != nil {
		return err
	}

	file, err := os.Create(datafilepath)
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Load loads data from the named file in the kutti config directory.
// If the named file does not exist, the second returned value is true
func Load(configfilename string) ([]byte, bool, error) {
	datafilepath, err := getconfigfilepath(configfilename)
	if err != nil {
		return nil, false, err
	}
	_, err = os.Stat(datafilepath)
	if os.IsNotExist(err) {
		return nil, true, err
	}

	if err != nil {
		return nil, false, err
	}

	data, err := ioutil.ReadFile(datafilepath)

	if err != nil {
		return nil, false, err
	}

	return data, false, nil
}
