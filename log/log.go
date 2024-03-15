/*
Copyright © 2024 PACLabs
*/
package log

/*
 *  Simple logging, with printf formatting directives baked in.
 */

import (
	"fmt"
	"os"
	"raygun/config"
)

func Verbose(format string, a ...any) {
	if config.Verbose || config.Debug {
		fmt.Printf(format+"\n", a...)
	}
}

func Debug(format string, a ...any) {
	if config.Debug {
		debug_msg := fmt.Sprintf(format, a...)
		fmt.Printf("DEBUG: %s\n", debug_msg)
	}
}

func Normal(format string, a ...any) {
	if config.Normal {
		fmt.Printf(format+"\n", a...)
	}
}

func Warning(format string, a ...any) {
	if config.Warning {
		warn_msg := fmt.Sprintf(format, a...)
		fmt.Printf("WARN : %s\n", warn_msg)
	}
}

func Error(format string, a ...any) {
	if config.Error {
		err_msg := fmt.Sprintf(format, a...)
		fmt.Printf("ERROR: %s\n", err_msg)
	}
}

/*
 *  Unlike the others, fatal the program after printing the fatal error
 */
func Fatal(format string, a ...any) {
	err_msg := fmt.Sprintf(format, a...)
	fmt.Printf("FATAL: %s\n", err_msg)
	os.Exit(-1)
}
