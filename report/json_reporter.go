/*
Copyright © 2024 PACLabs
*/
package report

import (
	"encoding/json"
	"raygun/config"
	"raygun/log"
	"raygun/types"
	"strings"
)

type JsonReporter struct {
	BaseReporter
}

func (tr JsonReporter) Generate(results types.CombinedResult) string {

	report := make(map[string]interface{}, 0)

	report["suite"] = generate_aggregate_suite_reports(results.ResultList)

	report["TEST_FAILURES"] = tr.TestFailuresExist(results)

	jsonBytes, err := json.Marshal(report)

	if err != nil {
		log.Error("Unable to convert report into JSON: %s", err.Error())
		return ""
	}

	return string(jsonBytes)

}

func generate_aggregate_suite_reports(list []types.TestSuiteResult) []interface{} {

	aggregate_suite_report := make([]interface{}, 0)

	for _, suite_result := range list {

		suite_report := make(map[string]interface{})

		suite_report["name"] = suite_result.Source.Name

		suite_report["SKIPPED"] = generate_aggregate_test_reports(suite_result.Skipped)
		suite_report["PASSED"] = generate_aggregate_test_reports(suite_result.Passed)
		suite_report["FAILED"] = generate_aggregate_test_reports(suite_result.Failed)

		aggregate_suite_report = append(aggregate_suite_report, suite_report)

	}

	return aggregate_suite_report
}

func generate_aggregate_test_reports(test_list []types.TestResult) []interface{} {

	aggregate_test_report := make([]interface{}, 0)

	for _, test_result := range test_list {
		report := make(map[string]interface{}, 0)

		report["name"] = test_result.Source.Name
		report["description"] = test_result.Source.Description
		if test_result.Status == config.FAIL {

			var comparison_type_array []string = make([]string, 0)
			var expected_value_array []string = make([]string, 0)

			report["actual"] = strings.TrimRight(test_result.Actual, "\r\n")

			for _, expectation := range test_result.Source.ExpectData {
				comparison_type_array = append(comparison_type_array, expectation.ExpectationType)
				expected_value_array = append(expected_value_array, expectation.Target)
			}

			report["comparison_type"] = comparison_type_array
			report["expected_value"] = expected_value_array

		}
		if config.PerformanceMetrics {
			report["durationMicroseconds"] = test_result.Duration.Microseconds()
		}

		aggregate_test_report = append(aggregate_test_report, report)
	}

	return aggregate_test_report

}
