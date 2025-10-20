/*
Copyright © 2024 PACLabs
*/
package cmd

import (
	"os"
	"raygun/config"

	"github.com/spf13/cobra"
)

var verbose bool
var debug bool

// 0.1 - basic core stuff
// 0.2 - adding support jwt and remote bundles
// 0.3 - adding the backtest command
// 0.3.1 - added -b for bundle, even though it isn't necessary
//   - also discovered you can 'fake' http request attributes and headers
//     by properly adjusting the input documents. see sample/example7
var Version = "v0.3.2"

var resolver = config.NewPropertyResolver()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "raygun",
	Short: "A tool for testing .rego logic",
	Long:  `A tool for testing Rego rule logic, by executing a series of tests against it and verifying we get back the expected results`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// Parse -D flags before cobra processes the flags
	// We need to filter them out from os.Args
	filteredArgs := resolver.ParseFlags(os.Args[1:])

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.raygun.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug")
	rootCmd.PersistentFlags().BoolVar(&config.StopOnFailure, "stop-on-failure", false, "Stop immediately at the first failed test")
	rootCmd.PersistentFlags().BoolVar(&config.SkipOnNetworkError, "skip-on-network-failure", false, "skip tests that fail because of OPA connectivity issues")
	rootCmd.PersistentFlags().BoolVar(&config.SkipOnParseError, "skip-on-parse-failure", false, "skip Raygun files that aren't valid YAML")

	//
	// flags related to OPA
	//
	rootCmd.PersistentFlags().StringVar(&config.OpaExecutablePath, "opa-exec",
		config.FindOpaExecutable("opa"),
		"OPA executable. Consider env var: RAYGUN_OPA_EXEC")
	rootCmd.PersistentFlags().StringVar(&config.OpaLogPath, "opa-log", config.OpaLogPath, "Location of the OPA log file")
	rootCmd.PersistentFlags().Uint16Var(&config.OpaPort, "opa-port", config.OpaPort, "The port upon which OPA is listening")

	// specify a remote server where the bundle can be found
	rootCmd.PersistentFlags().StringVarP(&config.OpaBundleUrl, "opa-bundle-url", "b", config.OpaBundleUrl, "URL of a hosted OPA bundle")

	// specify the endpoint of an OPA that is already running, so we don't have to start one
	rootCmd.PersistentFlags().StringVar(&config.OpaEndpointUrl, "opa-url", config.OpaEndpointUrl, "URL of an existing OPA that we will use for our tests")

	//
	// flags related to the decision logs
	//
	rootCmd.PersistentFlags().StringVar(&config.DecisionFile, "decision-file", config.DecisionFile, "Location of the file containing the json array of decision logs for backtest")

	// flags related to the output format
	rootCmd.PersistentFlags().StringVar(&config.ReportFormat, "report-format",
		config.ReportFormat, "Format of the test completion report (text, json)")

	// flags related to performance
	rootCmd.PersistentFlags().BoolVar(&config.PerformanceMetrics, "perf-metrics", false, "Measure the time required for each call & report")

	// Note: -D flags are handled by PropertyResolver before cobra sees them
	rootCmd.SetHelpTemplate(rootCmd.HelpTemplate() +
		"\nDynamic Properties:\n  -D KEY=VALUE    Define property (repeatable, takes precedence over env vars)\n")

	rootCmd.SetArgs(filteredArgs)
}
