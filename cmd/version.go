/*
Copyright © 2024 PACLabs
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "the version of raygun",
	Long:  `displays the version of raygun`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

}
