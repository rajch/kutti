package clustermanager

import (
	"regexp"
)

// IsValidName checks for the validity of a name.
// Valid names are up to 10 characters long, must start with a lowercase letter, and may contain lowercase letters and digits only.
func IsValidName(name string) bool {
	matched, _ := regexp.MatchString("^[a-z]([a-z0-9]{1,9})$", name)
	return matched
}

// // DefaultCluster returns the default cluster, or nil if none has been set
// func DefaultCluster() *Cluster {
// 	if config.DefaultClusterName == "" {
// 		return nil
// 	}

// 	result, ok := config.Clusters[config.DefaultClusterName]
// 	if !ok {
// 		ClearDefaultCluster()
// 	}
// 	return result
// }

// // SetDefaultCluster sets the default cluster name.
// // It returns an error if the cluster does not exist.
// func SetDefaultCluster(clustername string) error {
// 	_, ok := config.Clusters[clustername]
// 	if !ok {
// 		return errClusterDoesNotExist
// 	}

// 	config.DefaultClusterName = clustername
// 	return Save()
// }

// // ClearDefaultCluster clears the default cluster name.
// func ClearDefaultCluster() {
// 	config.DefaultClusterName = ""
// 	Save()
// }
