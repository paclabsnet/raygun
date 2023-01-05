package config

const STANDARD_OPA_PORT uint16 = 8181

const DEFAULT_RAYGUN_EXTENSION = ".raygun"
const DEFAULT_RAYSUITE_EXTENSION = ".raysuite"

var Debug bool = false
var Verbose bool = false
var Warning bool = true
var Normal bool = true
var Error bool = true
var StopOnTestFail bool = false
var ContinueOnSuiteError bool = false
var DefaultOpaPort = STANDARD_OPA_PORT
var DefaultOpaExecutablePath *string = nil
var RaygunExtension = DEFAULT_RAYGUN_EXTENSION
var RaysuiteExtension = DEFAULT_RAYSUITE_EXTENSION

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

func SetStopOnTestFail(v bool) {
	StopOnTestFail = v
}

func SetContinueOnSuiteError(v bool) {
	ContinueOnSuiteError = v
}

func SetDefaultOpaPort(port uint16) {
	DefaultOpaPort = port
}

func SetDefaultOpaExecutablePath(path string) {
	DefaultOpaExecutablePath = &path
}
