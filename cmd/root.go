/*
Copyright © 2023 John Brothers <johnbr@paclabs.net>
*/
package cmd

import (
	"os"
	"raygun/config"

	"github.com/spf13/cobra"
)

var verbose bool
var debug bool

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
	rootCmd.PersistentFlags().StringVar(&config.OpaExecutablePath, "opa-exec", "", "Location of the OPA executable")
}
