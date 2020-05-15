package clustermanager

import (
	"github.com/rajch/kutti/pkg/core"
)

// Version is a Kubernetes version that may be available to create a cluster
type Version struct {
	image core.VMImage
}

// K8sversion returns the Kubernetes version string
func (v *Version) K8sversion() string {
	return v.image.K8sVersion()
}

// Status returns the local availability of the version image
func (v *Version) Status() string {
	return v.image.Status()
}

// FromFile imports a version image from the specified file
func (v *Version) FromFile(filename string) error {
	return v.image.FromFile(filename)
}

// PurgeLocal removes the local cached copy of a version image
func (v *Version) PurgeLocal() error {
	return v.image.PurgeLocal()
}
