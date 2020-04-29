package core

// VMImage describes a template from which VMHosts are created
type VMImage interface {
	K8sVersion() string
	Status() string

	Fetch() error
	Verify() bool
	FromFile(string) error
}
