package opa

import (
	"raygun/log"
)

type OpaConfig struct {
	OpaPort    uint16
	OpaPath    string
	BundlePath string
	LogPath    string
}

/*

 we launch OPA specifying the bundle.tar.gz file to use as the source for .rego and the data section

*/

func DefineRuntime(opa_path string, opa_port uint16, bundle_path string, log_path string) *OpaConfig {

	o := &OpaConfig{
		OpaPort:    opa_port,
		OpaPath:    opa_path,
		BundlePath: bundle_path,
		LogPath:    log_path,
	}

	return o
}

func (opa *OpaConfig) Start() error {

	log.Warning("Opa Start() not implemented")
	return nil

}

func (opa *OpaConfig) Stop() error {

	log.Warning("OPA Stop not implemented, faking for now to test the rest of the system")
	return nil
}
