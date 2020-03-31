// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Log with stack trace.
package log

import (
	glog "log"
	"os"
	"runtime/debug"
)

/// Is equivalent to log.Fatalln with stack trace
func Fatal(v ...interface{}) {
	glog.Println(v...)
	debug.PrintStack()
	os.Exit(1)
}

/// Is equivalent to log.Fatalf with stack trace
func Fatalf(format string, v ...interface{}) {
	glog.Printf(format, v...)
	debug.PrintStack()
	os.Exit(1)
}

/// Is equivalent to log.Println with stack trace
func Print(v ...interface{}) {
	glog.Println(v...)
	debug.PrintStack()
}

/// Is equivalent to log.Printf with stack trace
func Printf(format string, v ...interface{}) {
	glog.Printf(format, v...)
	debug.PrintStack()
}

/// Is equivalent to log.Panicln with stack trace
func Panic(v ...interface{}) {
	glog.Panicln(v...)
}

/// Is equivalent to log.Panicf with stack trace
func Panicf(format string, v ...interface{}) {
	glog.Panicf(format, v...)
}
