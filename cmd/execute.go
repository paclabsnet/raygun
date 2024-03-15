/*
Copyright © 2024 PACLabs
*/
package cmd

import (
	"fmt"
	"os"
	"raygun/config"
	"raygun/finder"
	"raygun/log"
	"raygun/parser"
	"raygun/report"
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
		* Description
		* Decision Path (i.e. the URL we use for calling OPA to get the decision)
		* Input JSON (embedded in .raygun file or found in external files)
		* Expected Output, in one format for now:
			* substring match
			* eventually we might include JSONPath as well

	2. We iterate over the TestSuite records
		a.  We start OPA, if necessary.
			* is OPA already running?  Does the configuration match the one we already have?
			If so, do nothing.
			* if OPA is running with the wrong configuration, shut it down
			* start up OPA with the right configuration
		b. We iterate over the TestRecords
			* we call OPA at the specified Decision Path, and pass in the specified Input JSON
			* We accept the response
			* We take the response and substring and strip out all spaces
			* We compare the stripped response to our stripped substring to see if there's a match
			* We track the expected result, actual result and outcome for the report
		c. After all the TestRecords have been processed, we create a report for the Test Suite.

	4. Once all the TestSuites have been run, we generate a report  (text or json)
	        * For each Test Suite
				* We list all the tests that were skipped
				* We list all the tests that passed
				* We list all the tests that failed
					- in verbose mode, we'll print out more details, including the expected and actual
			* If there were test failures, we return with an error code




	Full set of command line arguments:
	--opa-exec  = location of opa executable . If not specified, look up RAYGUN_OPA_EXEC
	--opa-bundle-path = location of bundle file (includes filename)
	--opa-log  = location of opa log file  ($TMP/raygun-opa.log for now)

	--stop-on-failure       (continue on fail is default)
	--skip-on-parse-error   (stop immediately on parse error is default)
	--skip-on-network-error (mark the test failed by default)

	--report-format - text (default) or json

    -d --debug
	-v --verbose

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

		/*
		 *  Find the raygun files amidst the files and directories specified on the command line
		 */
		finder := finder.NewFinder(config.RaygunExtension)

		suite_files, err := finder.FindTargets(entities)

		if err != nil {
			log.Error("Error finding test suites: %v", err)
			return err
		}

		if len(suite_files) == 0 {
			log.Warning("No .raygun files found in specified location(s)")
			os.Exit(2)
		}

		/*
		 *  Parse the .raygun files that we found in the previous step
		 */

		log.Verbose("Parsing Raygun files: %v", suite_files)

		parser := parser.NewRaygunParser(config.SkipOnParseError)

		test_suite_list, err := parser.Parse(suite_files)

		if err != nil {
			log.Error("Unable to parse test files: %v", err)
			return err
		}

		//
		// Show what test suites we're about to test
		//
		if config.Verbose {
			log.Verbose("Test Suites: ")
			for _, suite := range test_suite_list {
				log.Verbose("   %s : %d tests", suite.Name, len(suite.Tests))
			}
		}

		/*
			create and execute the test-suite runner. It returns the combined
			results from all of the test suites.

			if config.StopOnFailure is true, this will have a partial set of results
		*/

		suiteRunner := runner.NewSuiteRunner(test_suite_list)

		results, err := suiteRunner.Execute()

		if err != nil {
			log.Error("Unable to execute test suite: %v", err)
			return err
		}

		/*
		 *  Generate an output report from the test results
		 */
		reporter := report.Build(config.ReportFormat)

		output := reporter.Generate(results)

		log.Normal(output)

		/*
		 * Fail with an error code, so build tools can detect it
		 */
		if reporter.TestFailuresExist(results) {
			os.Exit(1)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(executeCmd)

}
