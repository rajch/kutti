package vboxdriver

import "github.com/rajch/kutti/internal/pkg/fileutils"

// VBoxVMImage implements the VMImage interface for VirtualBox
type VBoxVMImage struct {
	ImageK8sVersion string
	ImageChecksum   string
	ImageStatus     string
}

// Name returns the name of the image
func (i *VBoxVMImage) Name() string {
	return imagenamefromk8sversion(i.ImageK8sVersion)
}

// K8sVersion returns the version of Kubernetes present in the image
func (i *VBoxVMImage) K8sVersion() string {
	return i.ImageK8sVersion
}

// Checksum returns the SHA256 checksum of the image
func (i *VBoxVMImage) Checksum() string {
	return i.ImageChecksum
}

// Status returns the status of the image
func (i *VBoxVMImage) Status() string {
	return i.ImageStatus
}

// Fetch fetches the image from wherever
func (i *VBoxVMImage) Fetch() error {
	panic("not implemented") // TODO: Implement
}

// Verify verifies the checksum
func (i *VBoxVMImage) Verify() bool {
	// Check status
	if i.ImageStatus != "Available" {
		return false
	}
	// Check for file existence
	imagepath, err := imagepathfromk8sversion(i.ImageK8sVersion)
	if err != nil {
		return false
	}

	// Check for checksum
	checksum, err := fileutils.ChecksumFile(imagepath)
	if err != nil {
		return false
	}

	if checksum != i.ImageChecksum {
		return false
	}

	return true
}

// FromFile verfies an image file, and if valid, copies it to the cache.
func (i *VBoxVMImage) FromFile(filepath string) error {
	err := addfromfile(i.ImageK8sVersion, filepath, i.ImageChecksum)
	if err != nil {
		return err
	}

	i.ImageStatus = "Available"
	return saveimages()
}
