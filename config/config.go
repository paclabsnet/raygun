/*
Copyright © 2024 PACLabs
*/
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const STANDARD_OPA_PORT uint16 = 8181

const DEFAULT_RAYGUN_EXTENSION = ".raygun"

// const DEFAULT_RAYSUITE_EXTENSION = ".raysuite"

var Debug bool = false
var Verbose bool = false
var Warning bool = true
var Normal bool = true
var Error bool = true
var StopOnFailure bool = false
var SkipOnParseError bool = false
var SkipOnNetworkError bool = false
var RaygunExtension = DEFAULT_RAYGUN_EXTENSION

var ReportFormat = "text"

var OpaPort = STANDARD_OPA_PORT
var OpaExecutablePath = FindOpaExecutable("opa")
var OpaBundlePath = "bundle.tar.gz"
var OpaLogPath = filepath.FromSlash(fmt.Sprintf("%s/raygun_opa.log", os.Getenv("TMP")))

// performance
var PerformanceMetrics bool = false

// core concepts in a testing program
const PASS = "pass"
const FAIL = "fail"
const SKIP = "skip"

/**********************************************
 *  These are probably superfluous at this point  (2024-03-22)
 */
func SetDebug(v bool) {
	Debug = v
}

func SetVerbose(v bool) {
	Verbose = v
}

func SetNormal(v bool) {
	Normal = v
}

func SetWarning(v bool) {
	Warning = v
}

func SetError(v bool) {
	Error = v
}

func SetStopOnFailure(v bool) {
	StopOnFailure = v
}

func SetSkipOnParseError(v bool) {
	SkipOnParseError = v
}

func SetSkipOnNetworkError(v bool) {
	SkipOnNetworkError = v
}

func SetOpaExecutablePath(path string) {
	OpaExecutablePath = path
}

/*
 * end superfluous section
 ***********************************************/

func FindOpaExecutable(defaultOpa string) string {

	if tmp, found := os.LookupEnv("RAYGUN_OPA_EXEC"); found {
		return tmp
	}

	return defaultOpa
}
