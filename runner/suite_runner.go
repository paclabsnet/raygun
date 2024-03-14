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

		if len(result.Failed) > 0 && config.StopOnFailure {
			log.Debug("Detected test failures, and StopOnFailure is true, so we're aborting")
			break
		}

	}

	suiteRunner.StopOpa()

	return results, nil
}

/*
 *  Execute a single suite of tests. If the OPA configuration is different than the last
 *  OPA configuration (different executable, different bundle, different log file), then this
 *  function will also stop OPA (if necessary) and start OPA with the new configuration
 */
func (suiteRunner *SuiteRunner) ExecuteSuite(suite types.TestSuite) (types.TestSuiteResult, error) {

	results := types.TestSuiteResult{}

	results.Source = suite

	if suiteRunner.DifferentOpaConfigurationThanLast(suite) {

		suiteRunner.StopOpa()

		err := suiteRunner.StartOpa(suite)

		if err != nil {
			return results, err
		}
	}

	/*
	 *   for each test, we POST data to OPA at the test-specified location, and
	 *   compare the results to our expected results
	 *
	 *
	 */
	for _, test := range suite.Tests {

		testRunner := NewTestRunner(test)

		response, network_err := testRunner.Post()

		testResult := types.TestResult{Source: test}
		var eval_err error = nil

		if network_err != nil {
			if config.SkipOnNetworkError {
				testResult.Status = config.SKIP
			} else {
				log.Error("Failed to POST data to OPA: %s", network_err.Error())
				return results, network_err
			}

		} else {

			testResult, eval_err = testRunner.Evaluate(response)

			/*
			 *  This shouldn't happen, so its a fairly serious problem
			 */
			if eval_err != nil {
				log.Error("Failed to evaluate response from OPA: %s", eval_err.Error())
				return results, eval_err
			}
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

		if len(results.Failed) > 0 && config.StopOnFailure {
			log.Debug("Test failure detected and StopOnFailure is true, aborting...")
			break
		}

	}

	return results, nil

}

/*
 *  If there's an OPA process ID, this will stop it.
 */
func (suiteRunner *SuiteRunner) StopOpa() {

	if suiteRunner.OpaRunner != nil {

		/*
		 *  if this is the first time through the suite list, there won't be a Last Suite
		 *  to use for debug reporting.
		 */
		if suiteRunner.LastSuite != nil {
			log.Debug("Stopping OPA with config: %v", suiteRunner.LastSuite.Opa)
		} else {
			log.Debug("Stopping OPA")
		}
		suiteRunner.OpaRunner.Stop()

		suiteRunner.OpaRunner = nil
	} else {
		log.Debug("StopOpa() - OpaRunner wasn't found, skipping")
	}
}

/*
 *   Starts up OPA with the specified executable, and includes the appropriate bundle file
 *   in the command line arguments.
 *
 *   We handle the OPA output stderr as a separate file, so it doesn't clutter the test
 *   output
 */
func (suiteRunner *SuiteRunner) StartOpa(suite types.TestSuite) error {

	log.Debug("Starting OPA with config: %v", suite.Opa)

	opa_runner := opa.NewOpaRunner(suite.Opa)

	suiteRunner.OpaRunner = &opa_runner

	err := suiteRunner.OpaRunner.Start()

	return err
}

/*
 *  Compares various elements of the OPA configuration, to determine if we need to start a new
 *  copy of OPA with a different bundle, log file or executable
 */
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
