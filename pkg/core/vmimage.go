package core

// VMImage describes a template from which VMHosts are created. A VMDriver
// is expected to maintain a cache of VMImages locally. A VMImage may be
// downloaded from a driver-specific source using the Fetch method, or
// added to the cache from a local file using the FromFile method.
type VMImage interface {
	K8sVersion() string
	Status() string

	Fetch() error
	FromFile(filepath string) error
	PurgeLocal() error
}
