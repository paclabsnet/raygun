/*
Copyright © 2025 PACLabs
*/
package cmd

import (
	"os"
	"raygun/config"
	"raygun/log"
	"raygun/parser"
	"raygun/report"
	"raygun/runner"

	"github.com/spf13/cobra"
)

/*
   Backtest allows the user to test an OPA loaded with a particular bundle against
   a JSON array of input / result pairs.

   Typically, this would be a way to test a bundle against a set of historical decision
   logs, but one could imagine other purposes, perhaps test generation, etc.

   The key steps

   1. OPA is up and running with the bundle we want
   2. We parse a JSON file containing an array of decision log entries. These entries
      may have some fields omitted for brevity
   3. Each record in that array includes at least an id, a path, an input document and a result
      These records are the tests!
   4. For each record
      - we call OPA at the specified path with the specified input
	  - we compare the result to the one we wanted
	  - we pass or fail tests as appropriate
   5. We generate a report
   6. We generate a non-zero exit code if there are any test failures
*/

var backtestCmd = &cobra.Command{
	Use:   "backtest <decision file>",
	Short: "Test the bundle against the historical decision records found in the decision file",
	Long:  `Test the bundle against the historical decision records found in the decision JSON file`,
	RunE: func(cmd *cobra.Command, args []string) error {

		config.Debug = debug
		config.Verbose = verbose
		config.Resolver = resolver

		/*
		 *  Find all of the directories and/or files specified on the command line.
		 *  If nothing is specified, add the current directory
		 */
		var entities = make([]string, 0)

		if len(args) < 1 {
			entities = append(entities, config.DecisionFile)
		} else {
			// we ignore everything after the first parameter for now
			entities = append(entities, args[0])
		}

		log.Verbose("Parsing decision log file: %v", entities)

		parser := parser.NewJsonParser()

		test_suite_list, err := parser.Parse(entities[0])

		if err != nil {
			log.Error("Unable to parse test files: %v", err)
			return err
		}

		if config.Verbose {
			log.Verbose("Backtests: ")
			for _, test := range test_suite_list[0].Tests {
				log.Verbose("   %s", test.Description)
			}
		}

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
	rootCmd.AddCommand(backtestCmd)
}
