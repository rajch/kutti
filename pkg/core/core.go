package core

// VMDriver describes common VM operations
type VMDriver interface {
	CreateNetwork() error
	CreateNode() error

	DeleteNetwork() error
	DeleteNode() error
}
