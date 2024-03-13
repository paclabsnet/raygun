package suite_runner

import (
	"raygun/log"
	"raygun/types"
)

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

func Run(suite types.TestSuite) types.TestSuiteResult {

	log.Verbose("Running test suite: %v\n", suite)

	results := types.TestSuiteResult{Source: suite}

	return results

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

}

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
