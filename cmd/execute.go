/*
Copyright © 2022 John Brothers <johnbr@paclabs.net>
*/
package cmd

import (
	"fmt"
	"raygun/config"
	"raygun/finder"
	"raygun/log"
	"raygun/parser"
	"raygun/runner"

	"github.com/spf13/cobra"
)

/*

	1. we parse each file to create the TestSuite.
		a. A TestSuite consists of:
		* Name
		* OPA startup information (executable and bundle file)
		* one or more TestRecords
		b. TestRecords contain:
		* Name
		* Invocation Path (i.e. the URL we use for calling OPA to get the decision)
		* Input JSON (embedded in .raygun file)
		* Expected Output, in one format for now:
			* substring: "allow": true       (this is a string match, stripping all whitespace)
			* eventually we might include JSONPath as well

	2. We iterate over the TestSuite records
		a.  We start OPA, if necessary.
			* is OPA already running?  Does the configuration match the one we already have?
			If so, do nothing.
			* if OPA is running with the wrong configuration, shut it down
			* start up OPA with the right configuration
		b. We iterate over the TestRecords
			* we call OPA on the specified Invocation Path, and pass in the specified Input JSON
			* We accept the response
			* We compare the response to our match
			* If the match fails, we generate a message including the test name, the expected & actual output
		c. After all the TestRecords have been run, we determine if all of the TestRecords have passed (or not)
			* If they all passed, we generate a message including the testsuite name, and the word PASS and the count of tests
			* If any failed, we generate a message including the testsuite name, the word FAIL and the failed and total test record count


	3. Once all of the TestSuites have been run, we determine if there were any failures
		a.  If there were any failures, we generate a message including "There were test failures:" and the number of failures
		b.  Otherwise, we generate a message including "All tests passed: " and the number of tests executed


	special cases:
	* OPA isn't running
	* honor --stop-on-fail and --skip-network-error  (in that order)
	* One of the test suites or test records does not have all of the required information
	* by default, stop immediately
	* honor --skip-parse-error



	Full set of command line arguments:
	--opa-exec  = location of opa executable  (skip for now)
	--opa-port  = port number to run OPA by default   (8181 for now)
	--opa-bundle-path = location of bundle file (includes filename)  (skip for now)
	--opa-log  = location of opa log file  (/tmp/raygun-opa-<DATE_TIME>.log for now)

	--stop-on-fail         (continue on fail is default)
	--skip-parse-error     (stop immediately on parse error is default)
	--skip-network-error   (mark the test failed by default)

*/

var executeCmd = &cobra.Command{
	Use:   fmt.Sprintf("execute <test directories or %s files>", config.RaygunExtension),
	Short: "Execute the .raygun files in the specified directory",
	Long:  `Execute the .raygun test cases specified via the command line directives`,
	RunE: func(cmd *cobra.Command, args []string) error {

		config.Debug = debug
		config.Verbose = verbose

		var entities = make([]string, 0)

		if len(args) < 1 {
			entities = append(entities, ".")
		} else {
			entities = append(entities, args...)
		}

		finder := finder.NewFinder(config.RaygunExtension)

		suite_files, err := finder.FindTargets(entities)

		if err != nil {
			log.Error("Error finding test suites: %v", err)
			return err
		}

		log.Verbose("Suite files: %v", suite_files)

		parser := parser.New(config.SkipOnParseError)

		test_suite_list, err := parser.Parse(suite_files)

		if err != nil {
			log.Error("Unable to parse test files: %v", err)
			return err
		}

		log.Verbose("Test Suite List: %v", test_suite_list)

		suiteRunner := runner.NewSuiteRunner(test_suite_list)

		results, err := suiteRunner.Execute()

		if err != nil {
			log.Error("Unable to execute test suite: %v", err)
			return err
		}

		log.Normal("Test Results: %v", results)

		return nil

	},
}

func init() {
	rootCmd.AddCommand(executeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
