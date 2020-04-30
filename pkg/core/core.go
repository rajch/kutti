package core

var (
	drivers map[string]VMDriver
)

// SSHClient defines a simple SSH client
type SSHClient interface {
	RunWithResults(address string, command string) (string, error)
}

func ensuredrivers() {
	if drivers == nil {
		drivers = make(map[string]VMDriver)
	}
}

// RegisterDriver registers a VMdriver with a name to core.
// If a driver with the specified name already exists, it is replaced.
func RegisterDriver(name string, d VMDriver) {
	ensuredrivers()

	if d != nil {
		drivers[name] = d
	}
}

// GetDriver returns a VMDriver corresponding to the name.
// If there is no driver registered against the name, nil is returned.
func GetDriver(name string) (VMDriver, bool) {
	result, ok := drivers[name]
	return result, ok
}

// ForEachDriver iterates over VM drivers.
// The callback function can return false to stop the iteration.
func ForEachDriver(f func(VMDriver) bool) {
	for _, driver := range drivers {
		cancel := f(driver)
		if cancel {
			break
		}
	}
}

func init() {
	ensuredrivers()
}
