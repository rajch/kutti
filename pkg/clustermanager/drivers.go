package clustermanager

import "github.com/rajch/kutti/pkg/core"

// ForEachDriver iterates over drivers
func ForEachDriver(f func(*Driver) bool) {
	core.ForEachDriver(func(vd core.VMDriver) bool {
		d := &Driver{vmdriver: vd}
		return f(d)
	})
}

// GetDriver gets the specified driver OR an error
func GetDriver(drivername string) (*Driver, bool) {
	vd, ok := core.GetDriver(drivername)
	if ok {
		return &Driver{vmdriver: vd}, ok
	}
	return nil, ok
}
