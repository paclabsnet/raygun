package types

import (
	"fmt"
	"raygun/opa"
)

type TestSuite struct {
	Opa         opa.OpaConfig
	Name        string
	Description string
	Tests       []TestRecord
}

func (suite TestSuite) String() string {

	return fmt.Sprintf("Suite: %s with %d Tests", suite.Name, len(suite.Tests))
}

type TestRecord struct {
	Suite        TestSuite
	Name         string
	Description  string
	Expects      TestExpectation
	Input        TestInput
	DecisionPath string // the path part of the URL to use to call opa
}

func (tr TestRecord) String() string {

	return fmt.Sprintf("Test: %s (%s)", tr.Name, tr.Description)
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
	Expected TestExpectation
	Actual   string
	Status   string // fail, pass, skip
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
	InputType string // inline, filepath
	Value     string
}

func (ti TestInput) String() string {

	if ti.InputType == "filepath" {
		return fmt.Sprintf("TestInput File: %s", ti.Value)
	}

	return fmt.Sprintf("TestInput: %s...", ti.Value[:20])
}
