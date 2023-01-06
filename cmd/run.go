/*
Copyright © 2022 John Brothers <johnbr@paclabs.net>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"raygun/config"
	"raygun/log"
	"raygun/suite_runner"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   fmt.Sprintf("run  <test directories or %s files>", config.RaysuiteExtension),
	Short: "Execute the .raygun and .raysuite files in the specified directory",
	Long: `Execute the .raygun test cases and .raysuite test suites specified
	via the command line directives`,
	RunE: func(cmd *cobra.Command, args []string) error {

		config.Debug = debug
		config.Verbose = verbose

		var entities = make([]string, 0)

		if len(args) < 1 {
			entities = append(entities, ".")
		} else {
			entities = append(entities, args...)
		}

		test_files, test_suites, err := findRaygunFiles(entities)

		if err != nil {
			return err
		}

		if len(test_suites) > 0 {
			log.Verbose("Test Suites to execute: %v\n", test_suites)

			for _, suite := range test_suites {

				path := filepath.Dir(suite)
				filename := filepath.Base(suite)

				err := suite_runner.Run(path, filename)
				if err != nil {
					log.Error("Error executing suite: %s", suite)
					return err
				}
			}

		} else {
			log.Normal("Not Implemented: Test Files to execute: %v\n", test_files)
			log.Warning("WARNING: ignoring test files for the moment, only implementing test suites\n")
		}

		return nil

	},
}

func findRaygunFiles(entities []string) ([]string, []string, error) {

	directories := make([]string, 0)
	test_files := make([]string, 0)
	test_suites := make([]string, 0)

	for _, entity := range entities {

		file_info, err := os.Stat(entity)

		if err != nil {
			log.Error("unable to find file/directory information about: %s", entity)
			return nil, nil, err
		}

		if file_info.IsDir() {
			directories = append(directories, entity)
		} else {

			if isRaygunFile(entity) {
				test_files = append(test_files, entity)
			} else if isRaysuiteFile(entity) {
				test_suites = append(test_suites, entity)
			} else {
				log.Debug("Ignoring non-raygun file: %s", entity)
			}

		}

		log.Verbose("Directories in which to look for %s files: %v", config.RaysuiteExtension, directories)

		for _, dir := range directories {

			files_in_dir, err := os.ReadDir(dir)

			if err != nil {
				log.Error("Unable to open directory %s to look for test files", dir)
				return nil, nil, err
			}

			for _, file_info := range files_in_dir {

				if !file_info.IsDir() {

					if isRaygunFile(file_info.Name()) {
						test_files = append(test_files, fmt.Sprintf("%s/%s", dir, file_info.Name()))
					} else if isRaysuiteFile(file_info.Name()) {
						test_suites = append(test_suites, fmt.Sprintf("%s/%s", dir, file_info.Name()))
					}
				}

			}

		}

	}

	log.Debug("run: test_files: %v", test_files)
	log.Debug("run: test_suites: %v", test_suites)

	return test_files, test_suites, nil
}

func isRaysuiteFile(entity string) bool {
	return filepath.Ext(entity) == config.RaysuiteExtension
}

func isRaygunFile(entity string) bool {
	return filepath.Ext(entity) == config.RaygunExtension
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
