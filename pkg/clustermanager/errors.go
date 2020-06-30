package clustermanager

import "errors"

var (
	errInvalidName             = errors.New("invalid name. Valid names are up to 10 characters long, must start with a lowercase letter, and may contain lowercase letters and digits only")
	errClusterExists           = errors.New("cluster already exists")
	errClusterDoesNotExist     = errors.New("cluster does not exist")
	errClusterNotEmpty         = errors.New("cluster is not empty")
	errDriverDoesNotExist      = errors.New("driver does not exist")
	errImageNotAvailable       = errors.New("image not available")
	errNodeNotFound            = errors.New("node not found")
	errNodeIsRunning           = errors.New("node is running")
	errNodeCannotStart         = errors.New("cannot start node")
	errNodeCannotStop          = errors.New("node not started. Cannot stop node")
	errPortForwardNotSupported = errors.New("port forwarding not supported")
	errPortAlreadyUsed         = errors.New("port already used")
)
