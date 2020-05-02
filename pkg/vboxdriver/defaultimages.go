package vboxdriver

// This contains a list of images available at the release of a version of
// vboxdriver
func defaultimages() map[string]*VBoxVMImage {
	return map[string]*VBoxVMImage{
		"1.14": &VBoxVMImage{
			ImageK8sVersion: "1.14",
			ImageChecksum:   "d071b0f991e4c2ee6b0cd95c77c2ca6336e351717f724236f30b52a31002ff1a",
			ImageStatus:     "Unavailable",
		},
	}
}
