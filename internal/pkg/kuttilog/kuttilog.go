package kuttilog

import "log"

var loglevel int = 2

// Setloglevel sets the current log level.
func Setloglevel(newlevel int) {
	loglevel = newlevel
}

// Loglevel returns the current log level.
func Loglevel() int {
	return loglevel
}

// V returns true if the specified level is between 0 and the current
// log level, false otherwise.
func V(level int) bool {
	return level <= loglevel && level >= 0
}

// Print prints to the log, if the level is <= current log level.
// Arguments are handled in the manner of fmt.Print.
func Print(level int, v ...interface{}) {
	if V(level) {
		log.Print(v...)
	}
}

// Printf prints to the log, if the level is <= current log level.
// Arguments are handled in the manner of fmt.Printf.
func Printf(level int, format string, v ...interface{}) {
	if V(level) {
		log.Printf(format, v...)
	}
}

// Println prints to the log, if the level is <= current log level.
// Arguments are handled in the manner of fmt.Println.
func Println(level int, v ...interface{}) {
	if V(level) {
		log.Println(v...)
	}
}

// SetPrefix sets a log prefix
func SetPrefix(prefix string) {
	log.SetPrefix(prefix)
}

func init() {
	log.SetFlags(0)

}
