/*
Copyright © 2025 PACLabs
*/

package parser

import (
	"path/filepath"
	"raygun/config"
	"raygun/types"
)

func CreateEmptySuite(source string) types.TestSuite {

	suite := types.TestSuite{Directory: filepath.Dir(source)}
	suite.Tests = make([]types.TestRecord, 0)
	suite.Opa.OpaPort = config.OpaPort
	suite.Opa.OpaPath = config.OpaExecutablePath
	suite.Opa.BundlePath = config.OpaBundleUrl
	suite.Opa.LogPath = config.OpaLogPath
	suite.Opa.BundleUrl = config.OpaBundleUrl
	suite.Opa.EndpointUrl = config.OpaEndpointUrl

	return suite
}
