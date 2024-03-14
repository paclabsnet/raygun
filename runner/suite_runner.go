package runner

import (
	"raygun/config"
	"raygun/log"
	"raygun/opa"
	"raygun/types"
)

type SuiteRunner struct {
	SuiteList []types.TestSuite
	LastSuite *types.TestSuite
	OpaRunner *opa.OpaRunner
}

func NewSuiteRunner(suite_list []types.TestSuite) SuiteRunner {

	suiteRunner := SuiteRunner{SuiteList: suite_list, LastSuite: nil}

	return suiteRunner
}

/*
	// Set the current working directory for the executable to "/tmp"
	cmd := exec.Command("my-executable")
	cmd.Dir = "/tmp"

	// Run the executable
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
*/

func (suiteRunner *SuiteRunner) Execute() (types.CombinedResult, error) {

	results := types.CombinedResult{}

	for i, suite := range suiteRunner.SuiteList {

		log.Debug("Execute: About To Execute Suite: %v", suite)

		result, err := suiteRunner.ExecuteSuite(suite)

		if err != nil {
			suiteRunner.StopOpa()
			return results, err
		}

		results.ResultList = append(results.ResultList, result)

		suiteRunner.LastSuite = &suiteRunner.SuiteList[i]

	}

	suiteRunner.StopOpa()

	return results, nil
}

func (suiteRunner *SuiteRunner) ExecuteSuite(suite types.TestSuite) (types.TestSuiteResult, error) {

	results := types.TestSuiteResult{}

	if suiteRunner.DifferentOpaConfigurationThanLast(suite) {

		suiteRunner.StopOpa()

		suiteRunner.StartOpa(suite)
	}

	//	log.Normal("Nothing is implemented yet, so test suite: %s cannot be run", suite)

	/*
		Now, it is time to figure out how to POST to the OPA url
	*/

	for _, test := range suite.Tests {

		testRunner := NewTestRunner(test)

		response, err := testRunner.Post()

		if err != nil {
			log.Error("Failed to POST data to OPA: %s", err.Error())
			return results, err
		}

		testResult, err := testRunner.Evaluate(response)

		if err != nil {
			log.Error("Failed to evaluate response from OPA: %s", err.Error())
			return results, err
		}

		switch testResult.Status {
		case config.PASS:
			results.Passed = append(results.Passed, testResult)
		case config.FAIL:
			results.Failed = append(results.Failed, testResult)
		case config.SKIP:
			results.Skipped = append(results.Skipped, testResult)
		default:
			log.Fatal("Unknown testResult Status for test %s : %s", testResult.Source, testResult.Status)
		}

		log.Debug("Expectations: type: %s, value: %s", test.Expects.ExpectationType, test.Expects.Target)

		/*
			request, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(body))

			if err != nil {
				log.Error("attempted to create POST request to %s with data %v -> %s", postUrl, body, err.Error())
				return results, err
			}

			request.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			response, err := client.Do(request)


			defer response.Body.Close()
		*/

	}

	return results, nil

}

func (suiteRunner *SuiteRunner) StopOpa() {

	if suiteRunner.OpaRunner != nil {

		if suiteRunner.LastSuite != nil {
			log.Debug("Stopping OPA with config: %v", suiteRunner.LastSuite.Opa)
		} else {
			log.Debug("Stopping OPA")
		}
		suiteRunner.OpaRunner.Stop()

		suiteRunner.OpaRunner = nil
	} else {
		log.Debug("Attempted to stop OPA, but no runner was found")
	}
}

func (suiteRunner *SuiteRunner) StartOpa(suite types.TestSuite) {

	log.Debug("Starting OPA with config: %v", suite.Opa)

	opa_runner := opa.NewOpaRunner(suite.Opa)

	suiteRunner.OpaRunner = &opa_runner

	suiteRunner.OpaRunner.Start()
}

func (suiteRunner SuiteRunner) DifferentOpaConfigurationThanLast(suite types.TestSuite) bool {

	if suiteRunner.LastSuite == nil {
		log.Debug("DifferentOpaConfigurationThanLast: No previous opa configuration, so this is definitely new")
		return true
	}

	if suiteRunner.LastSuite.Opa.BundlePath != suite.Opa.BundlePath {
		log.Debug("DifferentOpaConfigurationThanLast: Last suite bundlePath: %s is not the same as the new one %s", suiteRunner.LastSuite.Opa.BundlePath, suite.Opa.BundlePath)
		return true
	}

	if suiteRunner.LastSuite.Opa.LogPath != suite.Opa.LogPath {
		log.Debug("DifferentOpaConfigurationThanLast: Last Suite logPath: %s is different from the new logpath: %s", suiteRunner.LastSuite.Opa.LogPath, suite.Opa.LogPath)
		return true
	}

	if suiteRunner.LastSuite.Opa.OpaPath != suite.Opa.OpaPath {
		log.Debug("DifferentOpaConfigurationThanLast: Last Suite opaPath: %s is different from the new opaPath: %s", suiteRunner.LastSuite.Opa.OpaPath, suite.Opa.OpaPath)
		return true
	}

	log.Debug("DifferentOpaConfigurationThanLast: No differences in Opa Config detected: %v vs %v", suiteRunner.LastSuite.Opa, suite.Opa)
	return false
}

// parser := ray_parser.New()

// suite, err := parser.ParseSuiteFile(suitePath, suiteFilename)

// if err != nil {
// 	log.Error("suite_runner: unable to process suite filename: %s/%s", suitePath, suiteFilename)
// 	return err
// }

// if len(suite.RaygunTestFiles) == 0 {
// 	log.Error("no test files found\n")
// 	return nil
// }
// log.Verbose("Using test files: %v\n", suite.RaygunTestFiles)

// log.Verbose("launching OPA with the appropriate rules and data\n")

// opa := opa.DefineRuntime(suite.RegoSourceFiles, suite.OpaData)

// err = opa.Start()

// if err != nil {
// 	return err
// }

// log.Verbose("iterating over the test files and executing them\n")

// for _, test_file := range suite.RaygunTestFiles {

// 	testDetails, err := parser.ParseTestFile(test_file, suite)

// 	if err != nil {
// 		return err
// 	}

// 	runner, err := build_test_runner(testDetails, opa)

// 	if err != nil {
// 		return err
// 	}

// 	outcome, failure_reason, err := runner.ExecuteTest()

// 	if err != nil {
// 		return err
// 	}

// 	if outcome == "PASS" {
// 		log.Normal("Test %s : PASS\n", testDetails.Name)
// 	} else {
// 		log.Normal("test %s : FAIL: %s\n", testDetails.Name, failure_reason)
// 	}

// }

// opa.Stop()

// return nil

// func build_test_runner(details *types.TestRecord, opaConfig *types.OpaConfig) (*TestRunner, error) {
// 	return nil, errors.New("not_implemented_suite_runner_build_test_runner")
// }

/* func run_test(test *types.TestDetails, opa_context map[string]interface{}) (string, string, error) {

	return "FAIL", "Not implemented yet", nil

}
*/
/* func parse_test_file(test_file os.FileInfo, suite_context map[string]interface{}) (*types.TestDetails, error) {

	td := types.NewTestDetails(test_file.Name())

	return td, errors.New("not_implemented_suite_runner_parse_test_file")
}
*/
/* func launch_opa(suite_context map[string]interface{}) (map[string]interface{}, error) {
	opa_ctx := make(map[string]interface{})

	return opa_ctx, nil
}
*/
/* func parse_suite_file(suite string) (map[string]interface{}, []os.FileInfo, error) {
	suite_ctx := make(map[string]interface{})

	test_file_array := make([]os.FileInfo, 0)

	return suite_ctx, test_file_array, nil
}
*/
