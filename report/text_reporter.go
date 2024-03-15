/*
Copyright © 2024 PACLabs
*/
package report

/*
 *  Generate a simple text output that looks familiar to
 */

import (
	"fmt"
	"raygun/config"
	"raygun/types"
	"strings"
)

type TextReporter struct {
	BaseReporter
}

func (tr TextReporter) Generate(results types.CombinedResult) string {

	var sb strings.Builder

	sb.WriteString("Test Results:\n")

	test_failures := false

	for _, suite_result := range results.ResultList {

		sb.WriteString(fmt.Sprintf("   Suite: %s :\n", suite_result.Source.Name))

		if config.Verbose {
			sb.WriteString("      OPA Configuration:\n")
			sb.WriteString(fmt.Sprintf("         OPA Output Log: %s\n", suite_result.Source.Opa.LogPath))
			sb.WriteString(fmt.Sprintf("         Using OPA Bundle: %s\n", suite_result.Source.Opa.BundlePath))
		}

		for _, test_result := range suite_result.Skipped {
			sb.WriteString(fmt.Sprintf("      SKIPPED: %s\n", test_result.Source.Name))
			if config.Verbose {
				sb.WriteString(fmt.Sprintf("        - %s\n", test_result.Source.Description))
			}
		}
		for _, test_result := range suite_result.Passed {
			sb.WriteString(fmt.Sprintf("      PASSED: %s\n", test_result.Source.Name))
			if config.Verbose {
				sb.WriteString(fmt.Sprintf("        - %s\n", test_result.Source.Description))
			}
		}

		for _, test_result := range suite_result.Failed {
			test_failures = true
			sb.WriteString("\n")
			sb.WriteString(fmt.Sprintf("      FAILED: %s\n", test_result.Source.Name))
			if config.Verbose {
				sb.WriteString(fmt.Sprintf("        - %s\n", test_result.Source.Description))
			}

			if config.Verbose {

				// the TrimRight at the end is to make sure we don't have a dangling ] on a single line
				sb.WriteString(fmt.Sprintf("        Comparison: %s. Expected:[%s] Actual: [%s]\n",
					test_result.Source.Expects.ExpectationType,
					test_result.Source.Expects.Target,
					strings.TrimRight(test_result.Actual, "\r\n")))

			}
			sb.WriteString("\n")
		}

	}

	if test_failures {
		sb.WriteString("\n")
		sb.WriteString("WARNING: There are test failures\n")
	}

	return sb.String()
}
