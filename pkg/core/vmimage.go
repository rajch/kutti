package core

// VMImage describes a template from which VMHosts are created
type VMImage interface {
	K8sVersion() string
	Status() string

	Fetch() error
	FromFile(string) error
	PurgeLocal() error
}
