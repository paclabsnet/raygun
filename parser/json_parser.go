/*
Copyright © 2025 PACLabs
*/
package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"raygun/types"
)

/*
 *   Parses JSON files.  Primarily for backtest
 */

type JsonParser struct {
}

/*
  To execute an OPA query, we POST to /v1/data/{OpaRuleUrlPath}
*/

func NewJsonParser() *JsonParser {
	p := &JsonParser{}

	return p
}

func (parser *JsonParser) Parse(json_filename string) ([]types.TestSuite, error) {

	suite_list := make([]types.TestSuite, 0)

	data, err := os.ReadFile(json_filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	suite := CreateEmptySuite(json_filename)

	var decision_list []types.Decision
	err = json.Unmarshal(data, &decision_list)
	if err != nil {
		return nil, err
	}

	// convert decision records into tests and add them to the suite
	for _, decision := range decision_list {

		test := types.TestRecord{Suite: suite}

		test.DecisionPath = decision.Path
		test.Name = decision.DecisionId
		test.Input = createTestInput("inline", string(decision.Input))
		test.ExpectData = createTestExpectation(decision.Result)
		test.Description = fmt.Sprintf("%s (%s) -> %s", test.DecisionPath, test.Input.Value, test.ExpectData[0].Target)

		suite.Tests = append(suite.Tests, test)
	}

	suite_list = append(suite_list, suite)

	return suite_list, nil
}

func createTestExpectation(result string) []types.TestExpectation {

	var expected_result string = fmt.Sprintf("\"result\":%s", result)

	arr := []types.TestExpectation{
		{ExpectationType: "substring",
			Target: expected_result}}

	return arr
}

func createTestInput(inputType string, value string) types.TestInput {
	ti := types.TestInput{InputType: inputType, Value: value}
	return ti
}
