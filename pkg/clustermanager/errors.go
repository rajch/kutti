package clustermanager

import "errors"

var (
	errInvalidName             = errors.New("invalid name. Valid names are up to 10 characters long, must start with a lowercase letter, and may contain lowercase letters and digits only")
	errClusterExists           = errors.New("cluster already exists")
	errClusterDoesNotExist     = errors.New("cluster does not exist")
	errClusterNotEmpty         = errors.New("cluster is not empty")
	errDriverDoesNotExist      = errors.New("driver does not exist")
	errImageNotAvailable       = errors.New("image not available")
	errNodeExists              = errors.New("node already exists")
	errNodeNotFound            = errors.New("node not found")
	errNodeIsRunning           = errors.New("node is running")
	errNodeCannotStart         = errors.New("cannot start node")
	errNodeCannotStop          = errors.New("node not started. Cannot stop node")
	errPortForwardNotSupported = errors.New("port forwarding not supported")
	errPortNotForwarded        = errors.New("port not forwarded")
	errPortCannotUnmap         = errors.New("the SSH port cannot be unmapped")
	errPortNodePortInvalid     = errors.New("node port is invalid")
	errPortNodePortInUse       = errors.New("node port has already been forwarded")
	errPortHostPortInvalid     = errors.New("host port is invalid")
	errPortHostPortAlreadyUsed = errors.New("port already used")
)
