package debuglog

import (
	"log"
)

var debuggingEnabled bool

// On turns debug logging on
func On() {
	debuggingEnabled = true
	log.Print("Debug Logging Enabled")
}

// Off turns debug logging off
func Off() {
	debuggingEnabled = false
	log.Print("Debug Logging Disabled")
}

// Print writes msg to log file if debugging enabled
func Print(msg string) {
	if debuggingEnabled {
		caller := ""
		// _, file, no, ok := runtime.Caller(1)
		// if ok {
		// 	caller = fmt.Sprintf("[%s #%d] ", file, no)
		// }
		log.Printf(" * Debug %s:: %s", caller, msg)
	}
}

// IsEnabled reports current state
func IsEnabled() bool {
	return debuggingEnabled
}
