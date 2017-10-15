package bulletin

import (
	"fmt"
	"time"
)

// for tons of debug output (see also WorkerVerbose)
var Verbose bool = true

func p(format string, a ...interface{}) {
	if Verbose {
		TSPrintf(format, a...)
	}
}

// time-stamped printf
func TSPrintf(format string, a ...interface{}) {
	fmt.Printf("\n%s ", ts())
	fmt.Printf(format+"\n", a...)
}

// get timestamp for logging purposes
func ts() string {
	return time.Now().Format("2006-01-02 15:04:05.999 -0700 MST")
}
