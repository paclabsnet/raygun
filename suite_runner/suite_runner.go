package suite_runner

import (
	"fmt"
	"os"
)

var verbose bool
var debug bool

type TestDetails struct {
	Name string
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

func Run(suite string) error {

	if verbose || debug {
		fmt.Printf("Running test suite: %s\n", suite)
	}

	suite_context, test_files, err := parse_suite_file(suite)

	if err != nil {
		return err
	}

	if verbose || debug {
		fmt.Printf("Using test files: %v\n", test_files)
	}

	if verbose || debug {
		fmt.Printf("launching OPA with the appropriate rules and data\n")
	}

	opa_context, err := launch_opa(suite_context)

	if err != nil {
		return err
	}

	if len(test_files) == 0 {
		fmt.Printf("ERROR: no test files found\n")
		return nil
	}

	if verbose || debug {
		fmt.Printf("iterating over the test files and executing them\n")
	}

	for _, test_file := range test_files {

		test_details, err := parse_test_file(test_file, suite_context)

		if err != nil {
			return err
		}

		outcome, failure_reason, err := run_test(test_details, opa_context)

		if err != nil {
			return err
		}

		if outcome == "PASS" {
			fmt.Printf("Test %s : PASS\n", test_details.Name)
		} else {
			fmt.Printf("test %s : FAIL: %s\n", test_details.Name, failure_reason)
		}

	}

	return nil

}

func run_test(test TestDetails, opa_context map[string]interface{}) (string, string, error) {

	return "FAIL", "Not implemented yet", nil

}

func parse_test_file(test_file os.FileInfo, suite_context map[string]interface{}) (TestDetails, error) {

	var td TestDetails

	td.Name = test_file.Name()

	return td, nil
}

func launch_opa(suite_context map[string]interface{}) (map[string]interface{}, error) {
	opa_ctx := make(map[string]interface{})

	return opa_ctx, nil
}

func parse_suite_file(suite string) (map[string]interface{}, []os.FileInfo, error) {
	suite_ctx := make(map[string]interface{})

	test_file_array := make([]os.FileInfo, 0)

	return suite_ctx, test_file_array, nil
}

func SetVerbose(v bool) {
	verbose = v
}

func SetDebug(d bool) {
	debug = d
}
