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

// Status returns the local availability of the version
func (v *Version) Status() string {
	return v.image.Status()
}
