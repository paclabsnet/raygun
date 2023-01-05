package ray_parser

import (
	"errors"
	"raygun/types"
)

type RayParser struct {
}

/*
  To execute an OPA query, we POST to /v1/data/{OpaRuleUrlPath}
*/

func New() *RayParser {
	p := &RayParser{}

	return p
}

func (parser *RayParser) ParseSuiteFile(filename string) (*types.SuiteDetails, error) {

	return nil, errors.New("not_implemented")
}

func (parser *RayParser) ParseTestFile(filename string, suite *types.SuiteDetails) (*types.TestDetails, error) {
	return nil, errors.New("not_implemented")

}
