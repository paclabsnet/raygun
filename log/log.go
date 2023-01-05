package log

import (
	"fmt"
	"raygun/config"
)

// var debug bool = false
// var verbose bool = false
// var normal bool = true
// var warning bool = true

// func SetDebug(v bool) {
// 	debug = v
// }

// func SetVerbose(v bool) {
// 	verbose = v
// }

// func SetWarning(v bool) {
// 	warning = v
// }

// func SetNormal(v bool) {
// 	normal = v
// }

func Verbose(format string, a ...any) {
	if config.Verbose {
		fmt.Printf(format, a...)
	}
}

func Debug(format string, a ...any) {
	if config.Debug {
		debug_msg := fmt.Sprintf(format, a...)
		fmt.Printf("DEBUG: %s", debug_msg)
	}
}

func Normal(format string, a ...any) {
	if config.Normal {
		fmt.Printf(format, a...)
	}
}

func Warning(format string, a ...any) {
	if config.Warning {
		warn_msg := fmt.Sprintf(format, a...)
		fmt.Printf("WARN : %s", warn_msg)
	}
}

func Error(format string, a ...any) {
	if config.Error {
		err_msg := fmt.Sprintf(format, a...)
		fmt.Printf("ERROR: %s", err_msg)
	}
}
