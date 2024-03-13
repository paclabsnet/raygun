package config

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
var OpaPort = STANDARD_OPA_PORT
var OpaExecutablePath = "opa"
var OpaBundlePath = "bundle.tar.gz"
var OpaLogPath = "/tmp/raygun_opa.log"

//var RaysuiteExtension = DEFAULT_RAYSUITE_EXTENSION

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
