package vboxdriver

import (
	"fmt"
	"path"

	"github.com/rajch/kutti/internal/pkg/fileutils"
	"github.com/rajch/kutti/pkg/core"
)

// VBoxVMImage implements the VMImage interface for VirtualBox.
type VBoxVMImage struct {
	ImageK8sVersion string
	ImageChecksum   string
	ImageSourceURL  string
	ImageStatus     string
}

// K8sVersion returns the version of Kubernetes present in the image.
func (i *VBoxVMImage) K8sVersion() string {
	return i.ImageK8sVersion
}

// Status returns the status of the image.
// Status can be Available, meaning the image exists in the local cache and can
// be used to create hosts, or Unavailable, meaning it has to be downloaded using
// Fetch.
func (i *VBoxVMImage) Status() string {
	return i.ImageStatus
}

// Fetch fetches the image from its source URL.
func (i *VBoxVMImage) Fetch() error {
	cachedir, _ := core.CacheDir()
	tempfilename := fmt.Sprintf("kutti-k8s-%s.ovadownload", i.ImageK8sVersion)
	tempfilepath := path.Join(cachedir, tempfilename)

	// Download file
	err := fileutils.DownloadFile(i.ImageSourceURL, tempfilepath)
	if err != nil {
		return err
	}
	defer fileutils.RemoveFile(tempfilepath)

	// Add
	return i.FromFile(tempfilepath)
}

// FromFile verifies an image file on a local path, and if valid, copies it to the cache.
func (i *VBoxVMImage) FromFile(filepath string) error {
	err := addfromfile(i.ImageK8sVersion, filepath, i.ImageChecksum)
	if err != nil {
		return err
	}

	i.ImageStatus = "Available"
	return imageconfigmanager.Save()
}

// PurgeLocal removes the local cached copy of an image.
func (i *VBoxVMImage) PurgeLocal() error {
	if i.ImageStatus == "Available" {
		err := removefile(i.K8sVersion())
		if err == nil {
			i.ImageStatus = "Unavailable"

			return imageconfigmanager.Save()
		}
		return err
	}

	return nil
}
