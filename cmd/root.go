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

var Version = "v0.1.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "raygun",
	Short: "A tool for testing .rego logic",
	Long: `A tool for testing Rego rule logic, by executing a series of tests 
	against it and verifying we get back the expected results`,
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

	rootCmd.PersistentFlags().StringVar(&config.OpaExecutablePath, "opa-exec",
		config.FindOpaExecutable("opa"),
		"OPA executable. Consider env var: RAYGUN_OPA_EXEC")
	rootCmd.PersistentFlags().StringVar(&config.OpaLogPath, "opa-log", config.OpaLogPath, "Location of the OPA log file")

	rootCmd.PersistentFlags().StringVar(&config.ReportFormat, "report-format",
		config.ReportFormat, "Format of the test completion report (text, json)")

}
