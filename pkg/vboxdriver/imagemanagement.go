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

//var images map[string]*VBoxVMImage

var (
	imageconfigmanager configfilemanager.ConfigManager
	imagedata          *imageconfigdata
)

type imageconfigdata struct {
	images map[string]*VBoxVMImage
}

func (icd *imageconfigdata) Serialize() ([]byte, error) {
	return json.Marshal(icd.images)
}

func (icd *imageconfigdata) Deserialize(data []byte) error {
	loaddata := make(map[string]*VBoxVMImage)
	err := json.Unmarshal(data, &loaddata)
	if err == nil {
		icd.images = loaddata
	}
	return err
}

func (icd *imageconfigdata) Setdefaults() {
	icd.images = defaultimages()
}

// func ensureimages() {
// 	if images == nil {
// 		loadimages()
// 	}
// }

// func saveimages() error {
// 	ensureimages()

// 	data, err := json.Marshal(images)
// 	if err != nil {
// 		return err
// 	}

// 	return configfilemanager.SaveConfigfile(imagesConfigFile, data)
// }

// func loadimages() error {
// 	data, notexist, err := configfilemanager.LoadConfigfile(imagesConfigFile)
// 	if notexist {
// 		images = defaultimages()
// 		return nil
// 	}

// 	err = json.Unmarshal(data, &images)
// 	if err != nil {
// 		images = defaultimages()
// 		return err
// 	}

// 	return nil
// }

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

func removefile(k8sversion string) error {
	filename, err := imagepathfromk8sversion(k8sversion)
	if err != nil {
		return err
	}

	return fileutils.RemoveFile(filename)
}

func fetchimagelist() error {
	// Download image list into temp directory
	confdir, _ := core.ConfigDir()
	tempfilename := "vboximagesnewlist.json"
	tempfilepath := path.Join(confdir, tempfilename)
	err := fileutils.DownloadFile(ImagesSourceURL, tempfilepath)
	if err != nil {
		return err
	}
	defer fileutils.RemoveFile(tempfilepath)

	// Load into object
	data, _, err := configfilemanager.LoadConfigfile(tempfilename)
	if err != nil {
		return err
	}

	tempobj := make(map[string]*VBoxVMImage)
	err = json.Unmarshal(data, &tempobj)
	if err != nil {
		return err
	}

	// Compare against current and update
	for key, newimage := range tempobj {
		oldimage := imagedata.images[key]
		if oldimage != nil &&
			newimage.ImageChecksum == oldimage.ImageChecksum &&
			newimage.ImageSourceURL == oldimage.ImageSourceURL &&
			oldimage.ImageStatus == "Available" {

			newimage.ImageStatus = "Available"
		}
	}

	// Make it current
	imagedata.images = tempobj

	// Save as local configuration
	imageconfigmanager.Save()

	return nil
}

func init() {
	imagedata = &imageconfigdata{}
	imageconfigmanager = configfilemanager.New(imagesConfigFile, imagedata)
}
