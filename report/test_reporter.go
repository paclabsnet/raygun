package report

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

func (br BaseReporter) TestFailuresExist(results types.CombinedResult) bool {

	for _, suite_result := range results.ResultList {
		if len(suite_result.Failed) > 0 {
			return true
		}
	}

	return false
}
