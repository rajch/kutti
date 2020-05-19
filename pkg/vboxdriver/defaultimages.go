package vboxdriver

// This contains a list of images available at the release of a version of
// vboxdriver
func defaultimages() map[string]*VBoxVMImage {
	return map[string]*VBoxVMImage{
		"1.16": &VBoxVMImage{
			ImageK8sVersion: "1.16",
			ImageChecksum:   "495033e926241e4fed842acf7a892cb2c65c5f6cf7803d8a01af259f7162e07c",
			ImageSourceURL:  "https://github.com/rajch/kutti-images/releases/download/v0.1.13-beta/kutti-0.1.13-k8s-1.16.ova",
			ImageStatus:     "Unavailable",
		},
	}
}

// ImagesSourceURL is the location where the master list of images can be found
var ImagesSourceURL string = "https://github.com/rajch/kutti-images/releases/download/v0.1.13-beta/kutti-images.json"
