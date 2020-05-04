package vboxdriver

import (
	"encoding/json"
	"errors"
	"fmt"

	"path"

	"github.com/rajch/kutti/internal/pkg/configfilemanager"
	"github.com/rajch/kutti/internal/pkg/fileutils"
	"github.com/rajch/kutti/pkg/core"
)

const imagesConfigFile = "vboximages.json"

var images map[string]*VBoxVMImage

func ensureimages() {
	if images == nil {
		loadimages()
	}
}

func saveimages() error {
	ensureimages()

	data, err := json.Marshal(images)
	if err != nil {
		return err
	}

	return configfilemanager.Save(imagesConfigFile, data)
}

func loadimages() error {
	data, notexist, err := configfilemanager.Load(imagesConfigFile)
	if notexist {
		images = defaultimages()
		return nil
	}

	err = json.Unmarshal(data, &images)
	if err != nil {
		images = defaultimages()
		return err
	}

	return nil
}

func imagenamefromk8sversion(k8sversion string) string {
	return "kutti-" + k8sversion + ".ova"
}

func imagepathfromk8sversion(k8sversion string) (string, error) {
	cachedir, err := core.CacheDir()
	if err != nil {
		return "", fmt.Errorf("Could not retrieve CacheDir: %v", err)
	}

	result := path.Join(cachedir, imagenamefromk8sversion(k8sversion))
	return result, nil
}

func addfromfile(k8sversion string, filepath string, checksum string) error {
	filechecksum, err := fileutils.ChecksumFile(filepath)
	if err != nil {
		return err
	}

	if filechecksum != checksum {
		return errors.New("file  is not valid")
	}

	localfilepath, err := imagepathfromk8sversion(k8sversion)
	if err != nil {
		return err
	}

	err = fileutils.CopyFile(filepath, localfilepath, 1000, true)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	loadimages()
}
