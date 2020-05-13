package vboxdriver

// This contains a list of images available at the release of a version of
// vboxdriver
func defaultimages() map[string]*VBoxVMImage {
	return map[string]*VBoxVMImage{
		"1.18": &VBoxVMImage{
			ImageK8sVersion: "1.18",
			ImageChecksum:   "dbdfdfaa686143199887d605d6971074886189f68c5fbf081cae874bc9a56da8",
			ImageStatus:     "Unavailable",
		},
		"1.14": &VBoxVMImage{
			ImageK8sVersion: "1.14",
			ImageChecksum:   "d071b0f991e4c2ee6b0cd95c77c2ca6336e351717f724236f30b52a31002ff1a",
			ImageStatus:     "Unavailable",
		},
	}
}
