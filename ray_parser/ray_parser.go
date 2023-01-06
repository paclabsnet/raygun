package ray_parser

import (
	"errors"
	"path/filepath"
	"raygun/config"
	"raygun/log"
	"raygun/types"
	"raygun/util"
	"strings"
)

type RayParser struct {
}

/*
  To execute an OPA query, we POST to /v1/data/{OpaRuleUrlPath}
*/

func New() *RayParser {
	p := &RayParser{}

	return p
}

func (parser *RayParser) ParseSuiteFile(path string, filename string) (*types.SuiteDetails, error) {

	/*
		Naive implementation:  pull the entire file into a string,
		split on #

		for each sub-element, everything up to the first newline will
		be the directive

		everything after the first newline will be the data
	*/

	log.Debug("RayParser: ParseSuiteFile: path: %s; filename: %s", path, filename)

	data, err := util.ReadFile(filepath.Join(path, filename))

	if err != nil {
		return nil, err
	}

	elements := strings.Split(*data, "#")

	log.Debug("RayParser: ParseSuiteFile: found elements: %d : %v\n", len(elements), elements)

	details := &types.SuiteDetails{}

	for count, element := range elements {

		log.Debug("RayParser: ParseSuiteFile: parsing element: %d -> %s", count, element)

		element = strings.TrimSpace(element)

		if len(element) == 0 {
			continue
		}

		header, body, err := util.SplitHeaderAndBody(element)

		if err != nil {
			return nil, err
		}

		switch *header {
		case "Name":
			details.SuiteName = util.Chomp(*body)
		case "Description":
			details.SuiteDescription = strings.TrimSpace(*body)
		case "Data":
			tmp := strings.TrimSpace(*body)
			details.OpaData = &tmp
		case "Rule URL Path":
			tmp := strings.TrimSpace(*body)
			details.RuleUrlPath = &tmp
		case "Rego Source Files":
			details.RegoSourceFiles = util.Listify(*body)
			log.Debug("RayParser: ParseSuiteFile: RegoSourceFiles: %v", details.RegoSourceFiles)
		case "OPA Executable Path":
			tmp := strings.TrimSpace(*body)
			details.OpaExecutablePath = &tmp
		case "OPA Working Directory":
			tmp := strings.TrimSpace(*body)
			details.OpaWorkingDirectory = &tmp
		case "With Test Files":
			file_list := make([]string, 0)
			glob_list := util.Listify(*body)

			log.Debug("RayParser: ParseSuiteFile: glob_list: %v", glob_list)

			for _, filename_glob := range glob_list {

				fullPath := filepath.Join(path, filename_glob)

				if strings.Contains(fullPath, "*") {

					log.Debug("RayParser: ParseSuiteFile: matching for file glob: %s", fullPath)

					matches, err := filepath.Glob(fullPath)

					if err != nil {
						return nil, err
					}

					file_list = append(file_list, matches...)

				} else {
					file_list = append(file_list, fullPath)
				}

			}

			details.RaygunTestFiles = file_list

			log.Debug("RayParswer: ParseSuiteFile: RaygunTestFiles: %v", details.RaygunTestFiles)
		}

	}

	// if there are no test files, we'll try again with a simple glob of the entire working
	// directory for the suite
	if len(details.RaygunTestFiles) == 0 {
		details.RaygunTestFiles, err = filepath.Glob(filepath.Join(path, "*."+config.RaygunExtension))

		if err != nil {
			return nil, err
		}
	}

	// if there are still no test files, we can't really do anything
	if len(details.RaygunTestFiles) == 0 {
		return nil, errors.New("no_test_files_found")
	}

	log.Debug("RayParser: parse_suite: Successfully parsed suite file: %+v", details)
	if details.OpaData != nil {
		log.Debug("RayParser: parse_suite: OpaData is: %s", *details.OpaData)
	}

	return details, nil
}

func (parser *RayParser) ParseTestFile(fullPath string, suite *types.SuiteDetails) (*types.TestDetails, error) {
	log.Debug("RayParser: ParseTestFile: filename: %s", fullPath)

	data, err := util.ReadFile(filepath.Join(fullPath))

	if err != nil {
		return nil, err
	}

	elements := strings.Split(*data, "#")

	log.Debug("RayParser: ParseTestFile: found elements: %d : %v\n", len(elements), elements)

	testDetails := &types.TestDetails{}

	for count, element := range elements {

		log.Debug("RayParser: ParseTestFile: parsing element: %d -> %s", count, element)

		element = strings.TrimSpace(element)

		if len(element) == 0 {
			continue
		}

		header, body, err := util.SplitHeaderAndBody(element)

		if err != nil {
			return nil, err
		}

		switch *header {
		case "Name":
			testDetails.Name = util.Chomp(*body)
		case "Description":
			testDetails.Description = strings.TrimSpace(*body)
		case "Data":
			tmp := strings.TrimSpace(*body)
			testDetails.Data = &tmp
		case "Expects":
			expectations, err := buildExpectations(*body)
			if err != nil {
				log.Error("RayParser: ParseTestFile: unable to process expectations: %s", *body)
				return nil, err
			}
			testDetails.Expects = expectations
		case "Input":
			testDetails.Input = strings.TrimSpace(*body)
		case "Rule URL Path":
			tmp := strings.TrimSpace(*body)
			testDetails.RuleUrlPath = &tmp
		}
	}

	log.Debug("RayParser: parse_suite: Successfully parsed test file: %+v", testDetails)
	if testDetails.Expects != nil {
		log.Debug("RayParser: parse_suite: Expectations: %+v", *testDetails.Expects)
		if testDetails.Expects.Json != nil {
			log.Debug("RayParser: parse_suite: Expects JSON: %+v", *testDetails.Expects.Json)
		}
	}
	if testDetails.RuleUrlPath != nil {
		log.Debug("RayParser: parse_suite: Rule URL Path: %+v", *testDetails.RuleUrlPath)
	}

	return testDetails, nil
}

// partial implementation for now
func buildExpectations(body string) (*types.TestExpectations, error) {

	log.Warning("RayParser: buildExpectations: partial implementation (JSON only)")

	var te = &types.TestExpectations{
		Format: "json",
		Json:   &body,
	}

	return te, nil
}
