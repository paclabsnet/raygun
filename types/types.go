/*
Copyright © 2024 PACLabs
*/
package types

/*
 *  all of the test-related structs we need for this project
 */

import (
	"encoding/json"
	"fmt"
	"raygun/opa"
	"time"
)

type TestSuite struct {
	Opa         opa.OpaConfig `yaml:"opa"`
	Name        string        `yaml:"name"`
	Description string        `yaml:"description,omitempty"`
	Directory   string        `yaml:"directory"`
	Jwt         TestJwt       `yaml:"jwt,omitempty"`
	Tests       []TestRecord  `yaml:"tests"`
}

func (suite TestSuite) String() string {

	return fmt.Sprintf("Suite: %s with %d Tests.\n  OPA config: %v\n  JWT config: %v\n", suite.Name, len(suite.Tests), suite.Opa.String(), suite.Jwt)
}

type TestRecord struct {
	Suite        TestSuite              `yaml:"suite"`
	Name         string                 `yaml:"name"`
	Skip         bool                   `yaml:"skip,omitempty"`
	Description  string                 `yaml:"description,omitempty"`
	ExpectsMap   map[string]interface{} `yaml:"expects"`
	Input        TestInput              `yaml:"input"`
	DecisionPath string                 `yaml:"decision-path"` // the path part of the URL to use to call opa
	Jwt          TestJwt                `yaml:"jwt,omitempty"` // the structure containing the parts of the JWT
	ExpectData   []TestExpectation      // we parse ExpectsMap to create this
}

func (tr TestRecord) String() string {

	return fmt.Sprintf("Test: %s (%s)", tr.Name, tr.Description)
}

type CombinedResult struct {
	ResultList []TestSuiteResult
}

type TestSuiteResult struct {
	Source  TestSuite
	Failed  []TestResult
	Passed  []TestResult
	Skipped []TestResult
}

func (tsr TestSuiteResult) String() string {

	return fmt.Sprintf("Suite Results: %s - Passed: %d, Failed: %d, Skipped: %d",
		tsr.Source.Name, len(tsr.Passed), len(tsr.Failed), len(tsr.Skipped))

}

type TestResult struct {
	Source   TestRecord
	Actual   string
	Status   string // fail, pass, skip
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

func (tr TestResult) String() string {

	return fmt.Sprintf("TestResult: %s - status: %s", tr.Source.Name, tr.Status)
}

type TestExpectation struct {
	ExpectationType string // exact, substring, jsonpath
	Target          string
}

func (te TestExpectation) String() string {

	return fmt.Sprintf("TestExpectation: Type: %s  - Target: %s", te.ExpectationType, te.Target)

}

type TestInput struct {
	InputType string `yaml:"type"` // inline, filepath
	Value     string `yaml:"value"`
}

func (ti TestInput) String() string {

	if ti.InputType == "filepath" {
		return fmt.Sprintf("TestInput File: %s", ti.Value)
	}

	return fmt.Sprintf("TestInput: %s...", ti.Value[:20])
}

type TestJwt struct {
	Algorithm  string       `yaml:"algorithm"`
	Secret     string       `yaml:"secret,omitempty"`
	PrivateKey string       `yaml:"private_key,omitempty"`
	PublicKey  string       `yaml:"public_key,omitempty"`
	Claims     ClaimsConfig `yaml:"claims"`
	Active     bool         `yaml:"active,omitempty"`
}

type ClaimsConfig struct {
	Issuer     string                 `yaml:"iss,omitempty"`
	Subject    string                 `yaml:"sub,omitempty"`
	Audience   interface{}            `yaml:"aud,omitempty"`
	Lifetime   string                 `yaml:"lifetime"`
	IncludeIat bool                   `yaml:"include-iat,omitempty"`
	JWTID      string                 `yaml:"jti,omitempty"`
	Custom     map[string]interface{} `yaml:"custom,omitempty"`
}

type Decision struct {
	DecisionId string          `json:"id"`
	Path       string          `json:"path"`
	Input      json.RawMessage `json:"input"`
	Result     string          `json:"result"`
}
