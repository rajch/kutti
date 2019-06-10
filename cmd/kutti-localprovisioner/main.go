package main

import (
	"flag"
	"os"

	"github.com/rajch/kutti/pkg/localprovisioner"
	glog "k8s.io/klog"
)

func main() {
	glog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "INFO")
	flag.Parse()

	// Fetch and sanity-check nodename and rootpath
	//   from the environment variables KUTTI_NODE_NAME
	//   and KUTTI_ROOT_PATH respectively.
	nodename := os.Getenv("KUTTI_NODE_NAME")
	if nodename == "" {
		glog.Exit("Could not fetch node name from variable KUTTI_NODE_NAME. Cannot continue.")
	}

	rootpath := os.Getenv("KUTTI_ROOT_PATH")
	if rootpath == "" {
		glog.Exit("Could not fetch root path from variable KUTTI_ROOT_PATH. Cannot continue.")
	}

	if stat, err := os.Stat(rootpath); !(err == nil && stat.IsDir()) {
		glog.Exitf("Root path %s does not exist, or is not a directory.", rootpath)
	}

	// Start
	glog.Info("Starting local provisioner...")
	err := localprovisioner.RunProvisioner(nodename, rootpath)
	if err != nil {
		glog.Fatalln(err)
	} else {
		glog.Exit("Kutti Local provisioner stopped by itself. Something is not right.")
	}
}
