package localprovisioner

import (
	"os"
	"path"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	glog "k8s.io/klog"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

const (
	provisionerName = "rajware.net/kutti-local-provisioner"
)

type kuttiLocalProvisioner struct {
	nodeName string // The hostname/nodename where the provisioner runs
	rootPath string // The directory under which volume directories will be created
}

func (p *kuttiLocalProvisioner) Provision(options controller.ProvisionOptions) (*v1.PersistentVolume, error) {
	defer glog.Flush()

	glog.Infof("Request received for creating a PV called %s.\n", options.PVName)

	// Create a directory
	newvolumepath := path.Join(p.rootPath, options.PVName)
	if err := os.MkdirAll(newvolumepath, 0755); err != nil {
		return nil, err
	}

	// Explicitly chmod created dir, so we know mode is set to 0777 regardless of umask
	if err := os.Chmod(newvolumepath, 0755); err != nil {
		return nil, err
	}

	glog.Infof("Directory called %s created. Now provisioning volume...\n", options.PVName)
	// Create the PersistentVolume object with node affinity
	pv := &v1.PersistentVolume{
		ObjectMeta: meta.ObjectMeta{
			Name: options.PVName,
			Annotations: map[string]string{
				"rajware.net/provisionedBy": provisionerName,
			},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: *options.StorageClass.ReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceStorage: options.PVC.Spec.Resources.Requests[v1.ResourceStorage],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				Local: &v1.LocalVolumeSource{
					Path: newvolumepath,
				},
			},
			NodeAffinity: &v1.VolumeNodeAffinity{
				Required: &v1.NodeSelector{
					NodeSelectorTerms: []v1.NodeSelectorTerm{
						{
							MatchExpressions: []v1.NodeSelectorRequirement{
								{
									Key:      "kubernetes.io/hostname",
									Operator: "In",
									Values:   []string{p.nodeName},
								},
							},
						},
					},
				},
			},
		},
	}

	glog.Infoln("Volume object created.")

	return pv, nil
}

func (p *kuttiLocalProvisioner) Delete(volume *v1.PersistentVolume) error {
	defer glog.Flush()

	glog.Infof("Request received for deleting a PV called %s.\n", volume.Name)

	// Sanity check the volume before removing underlying storage
	ann, ok := volume.Annotations["rajware.net/provisionedBy"]
	if !ok || ann != provisionerName {
		glog.Errorf("The persistent volume %s has not been created by Kutti Local Provisioner.\n", volume.Name)
		return errors.New("This persistent volume has not been created by Kutti Local Provisioner")
	}

	// Remove underlying storage
	volumepath := path.Join(p.rootPath, volume.Name)
	if err := os.RemoveAll(volumepath); err != nil {
		glog.Errorf("Problem removing PV source directory %s.\n", volumepath)
		return errors.Wrap(err, "problem removing PV source directory "+volumepath)
	}

	glog.Infof("Underlying storage for PV %s deleted.\n", volume.Name)
	return nil
}

var _ controller.Provisioner = &kuttiLocalProvisioner{}

// NewKuttiLocalProvisioner returns a new Provisioner using Local volume driver
func NewKuttiLocalProvisioner(nodename string, rootpath string) controller.Provisioner {
	return &kuttiLocalProvisioner{
		nodeName: nodename,
		rootPath: rootpath,
	}
}

// RunProvisioner creates and runs a provision controller.
func RunProvisioner(nodename string, rootpath string) error {
	glog.Infoln("Getting kube config and client...")

	// Get config
	// We will always be running in-cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Errorln("Could not fetch kube config from cluster.")
		return errors.Wrap(err, "Could not fetch kube config from cluster.")
	}

	// Create a client from config
	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorln("Could not create kube client from in-cluster config.")
		return errors.Wrap(err, "Could not create kube client from in-cluster config.")
	}

	// Get server verion from client
	//   because the provision controller needs it
	kubeVersion, err := kubeclient.Discovery().ServerVersion()
	if err != nil {
		glog.Errorln("Could not fetch server version.")
		return errors.Wrap(err, "Could not fetch server version.")
	}

	localprovisioner := NewKuttiLocalProvisioner(nodename, rootpath)

	// Create a controller with provisioner
	pc := controller.NewProvisionController(
		kubeclient,
		provisionerName,
		localprovisioner,
		kubeVersion.GitVersion,
	)

	glog.Infof("Controller created. Details:\n%+v\nRun commencing...", pc)
	glog.Flush()

	pc.Run(wait.NeverStop)

	return nil
}
