/*
Copyright © 2024 PACLabs
*/
package report

/*
 *  This defines the generic interface for test reports, as well as some
 *  common methods that can be re-ussed by the implementations
 */

import (
	"raygun/log"
	"raygun/types"
)

type TestReporter interface {
	Generate(types.CombinedResult) string
	TestFailuresExist(types.CombinedResult) bool
}

func Build(outputFormat string) TestReporter {

	switch outputFormat {
	case "text":
		return TextReporter{}
	case "json":
		return JsonReporter{}
	default:
		log.Warning("Unrecognized output format %s, using text", outputFormat)
		return TextReporter{}
	}

}

type BaseReporter struct {
}

// we can save a few microns by making this a flag on the suite and on
// the combined result, but that's not a priority
func (br BaseReporter) TestFailuresExist(results types.CombinedResult) bool {

	for _, suite_result := range results.ResultList {
		if len(suite_result.Failed) > 0 {
			return true
		}
	}

	return false
}
