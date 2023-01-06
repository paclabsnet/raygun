package opa

import (
	"raygun/config"
	"raygun/log"
)

type OpaConfig struct {
	opa_port            uint16
	opa_executable_path string
	opa_data            string
	rego_source_files   []string
}

/*
 we PUT the opa_data into OPA at path /v1/data

 we launch OPA specifying the .rego file found at {rego_package_path}

*/

func DefineRuntime(rego_source_files []string, opa_data *string) *OpaConfig {

	o := &OpaConfig{
		opa_port:          config.DefaultOpaPort,
		rego_source_files: rego_source_files,
	}

	if config.DefaultOpaExecutablePath != nil {
		o.opa_executable_path = *config.DefaultOpaExecutablePath
	} else {
		o.opa_executable_path = "opa"
	}

	if opa_data != nil {
		o.opa_data = *opa_data
	}

	return o
}

func (opa *OpaConfig) Start() error {

	log.Warning("OPA Start not implemented, faking for now to test the rest of the system")

	return nil

}

func (opa *OpaConfig) Stop() error {

	log.Warning("OPA Stop not implemented, faking for now to test the rest of the system")
	return nil
}
