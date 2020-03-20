package core

var (
	driverfuncs map[string]func() (VMDriver, error)
	drivers     map[string]VMDriver
)
