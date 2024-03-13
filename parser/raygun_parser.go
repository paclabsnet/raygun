package parser

import (
	"errors"
	"fmt"
	"os"
	"raygun/config"
	"raygun/log"
	"raygun/types"
	"raygun/util"

	"gopkg.in/yaml.v3"
)

type RaygunParser struct {
	SkipOnParseError bool
}

/*
  To execute an OPA query, we POST to /v1/data/{OpaRuleUrlPath}
*/

func New(skip_on_parse_error bool) *RaygunParser {
	p := &RaygunParser{SkipOnParseError: skip_on_parse_error}

	return p
}

func (parser *RaygunParser) Parse(raygun_file_list []string) ([]types.TestSuite, error) {

	suite_list := make([]types.TestSuite, 0)

	for _, file := range raygun_file_list {

		suite, err := parser.ParseSuiteFile(file)

		if err != nil {
			if !parser.SkipOnParseError {
				log.Fatal("Parse error on suite file: %s [%v]", file, err)
			} else {
				log.Warning("Parse error on suite file: %s -> %v .. skipping", file, err)
			}
		}

		suite_list = append(suite_list, suite)

	}

	return suite_list, nil
}

func (parser *RaygunParser) ParseSuiteFile(filepath string) (types.TestSuite, error) {

	/*
		The Suite File is a .raygun file, but it
	*/

	log.Debug("RaygunParser: ParseSuiteFile: filepath: %s", filepath)

	// data, err := util.ReadFile(filepath)

	// if err != nil {
	// 	return nil, err
	// }

	tree := make(map[string]interface{})

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		if parser.SkipOnParseError {
			log.Warning("Parse Error on: %s [%v] . skipping...", filepath, err)
		} else {
			log.Fatal("Parse Error on %s [%v]", filepath, err)
		}
	}

	err = yaml.Unmarshal(yamlFile, tree)
	if err != nil {
		if parser.SkipOnParseError {
			log.Warning("Parse Error on: %s [%v] . skipping...", filepath, err)
		} else {
			log.Fatal("Parse Error on %s [%v]", filepath, err)
		}
	}

	log.Normal(" YAML File %s as Map:\n%v", filepath, tree)

	suite := types.TestSuite{}

	suite.Tests = make([]types.TestRecord, 0)

	err = parser.yamlToSuite(&suite, tree)
	if err != nil {
		if parser.SkipOnParseError {
			log.Warning("Parse Error on: %s [%v] . skipping...", filepath, err)
		} else {
			log.Fatal("Parse Error on %s [%v]", filepath, err)
		}
	}

	return suite, nil
}

/*
 *  process the tree of yaml elements and build out the test suite and test cases
 */
func (p RaygunParser) yamlToSuite(suite *types.TestSuite, tree map[string]interface{}) error {

	sorted_keys := util.SortMapKeys(tree)

	suite.Opa.OpaPort = config.OpaPort
	suite.Opa.OpaPath = config.OpaExecutablePath
	suite.Opa.BundlePath = config.OpaBundlePath
	suite.Opa.LogPath = config.OpaLogPath

	for _, k := range sorted_keys {

		v := tree[k]

		log.Debug("raygun file: key: %s has value: %v", k, v)

		switch k {
		case "opa":
			err := p.yamlToOpaConfig(suite, v.(map[string]interface{}))
			if err != nil {
				return err
			}
		case "description":
			suite.Description = v.(string)
		case "suite":
			suite.Name = v.(string)
		case "tests":
			err := p.yamlToTestArray(suite, v.([]interface{}))
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("suite parser: unknown key: %s", k)
		}

	}

	if suite.Name == "" {
		return errors.New("suite name (the suite: element) is required")
	}

	if len(suite.Tests) == 0 {
		return fmt.Errorf("suite %s has no tests", suite.Name)
	}

	return nil

}

/*
 *  logic to handle the OPA specification
 */
func (p RaygunParser) yamlToOpaConfig(suite *types.TestSuite, tree map[string]interface{}) error {

	sorted_keys := util.SortMapKeys(tree)

	for _, k := range sorted_keys {

		v := tree[k]

		log.Debug("raygun file opa configuration: key: %s has value: %v", k, v)

		switch k {
		case "path":
			suite.Opa.OpaPath = v.(string)
		case "bundle-path":
			suite.Opa.BundlePath = v.(string)
		default:
			return fmt.Errorf("unknown/unsupported opa config key: %s", k)
		}
	}

	return nil
}

/*
 * logic to handle the test array
 */
func (p RaygunParser) yamlToTestArray(suite *types.TestSuite, test_array []interface{}) error {

	for _, tree := range test_array {
		err := p.yamlToTest(suite, tree.(map[string]interface{}))

		if err != nil {
			return err
		}
	}

	return nil

}

func (p RaygunParser) yamlToTest(suite *types.TestSuite, tree map[string]interface{}) error {

	log.Debug("test info: %v", tree)

	test := types.TestRecord{}

	sorted_keys := util.SortMapKeys(tree)

	for _, k := range sorted_keys {

		v := tree[k]

		switch k {
		case "name":
			test.Name = v.(string)
		case "description":
			test.Description = v.(string)
		case "decision-path":
			test.DecisionPath = v.(string)
		case "expects":
			p.yamlToExpectations(&test, v.(map[string]interface{}))
		case "input":
			p.yamlToInputJSON(&test, v.(map[string]interface{}))
		default:
			return fmt.Errorf("test parser: unknown/unsupported key %s", k)
		}

	}

	suite.Tests = append(suite.Tests, test)

	return nil
}

/*
 * Process the test expectations
 */
func (p RaygunParser) yamlToExpectations(test *types.TestRecord, tree map[string]interface{}) error {

	sorted_keys := util.SortMapKeys(tree)

	for _, k := range sorted_keys {

		v := tree[k]

		switch k {
		case "type":
			test.Expects.ExpectationType = v.(string)
		case "target":
			test.Expects.Target = v.(string)
		default:
			return fmt.Errorf("unknown/unsupported 'expects' section key: %s", k)
		}
	}

	return nil
}

/*
 * Process the input - raw json or a filename
 */
func (p RaygunParser) yamlToInputJSON(test *types.TestRecord, tree map[string]interface{}) error {

	sorted_keys := util.SortMapKeys(tree)

	for _, k := range sorted_keys {

		v := tree[k]

		switch k {
		case "type":
			test.Input.InputType = v.(string)
		case "value":
			test.Input.Value = v.(string)
		default:
			return fmt.Errorf("unknown/unsupported 'input' section key: %s", k)
		}
	}

	return nil
}

// elements := strings.Split(*data, "#")

// log.Debug("RaygunParser: ParseSuiteFile: found elements: %d : %v\n", len(elements), elements)

// for count, element := range elements {

// 	log.Debug("RaygunParser: ParseSuiteFile: parsing element: %d -> %s", count, element)

// 	element = strings.TrimSpace(element)

// 	if len(element) == 0 {
// 		continue
// 	}

// 	header, body, err := util.SplitHeaderAndBody(element)

// 	if err != nil {
// 		return nil, err
// 	}

// 	switch *header {
// 	case "Name":
// 		suite.Name = util.Chomp(*body)
// 	case "Description":
// 		suite.SuiteDescription = strings.TrimSpace(*body)
// 	case "Data":
// 		tmp := strings.TrimSpace(*body)
// 		details.OpaData = &tmp
// 	case "Rule URL Path":
// 		tmp := strings.TrimSpace(*body)
// 		details.RuleUrlPath = &tmp
// 	case "Rego Source Files":
// 		details.RegoSourceFiles = util.Listify(*body)
// 		log.Debug("RayParser: ParseSuiteFile: RegoSourceFiles: %v", details.RegoSourceFiles)
// 	case "OPA Executable Path":
// 		tmp := strings.TrimSpace(*body)
// 		details.OpaExecutablePath = &tmp
// 	case "OPA Working Directory":
// 		tmp := strings.TrimSpace(*body)
// 		details.OpaWorkingDirectory = &tmp
// 	case "With Test Files":
// 		file_list := make([]string, 0)
// 		glob_list := util.Listify(*body)

// 		log.Debug("RayParser: ParseSuiteFile: glob_list: %v", glob_list)

// 		for _, filename_glob := range glob_list {

// 			fullPath := filepath.Join(path, filename_glob)

// 			if strings.Contains(fullPath, "*") {

// 				log.Debug("RayParser: ParseSuiteFile: matching for file glob: %s", fullPath)

// 				matches, err := filepath.Glob(fullPath)

// 				if err != nil {
// 					return nil, err
// 				}

// 				file_list = append(file_list, matches...)

// 			} else {
// 				file_list = append(file_list, fullPath)
// 			}

// 		}

// 		details.RaygunTestFiles = file_list

// 		log.Debug("RayParswer: ParseSuiteFile: RaygunTestFiles: %v", details.RaygunTestFiles)
// 	}

// }

// // if there are no test files, we'll try again with a simple glob of the entire working
// // directory for the suite
// if len(details.RaygunTestFiles) == 0 {
// 	details.RaygunTestFiles, err = filepath.Glob(filepath.Join(path, "*."+config.RaygunExtension))

// 	if err != nil {
// 		return nil, err
// 	}
// }

// // if there are still no test files, we can't really do anything
// if len(details.RaygunTestFiles) == 0 {
// 	return nil, errors.New("no_test_files_found")
// }

// log.Debug("RayParser: parse_suite: Successfully parsed suite file: %+v", details)
// if details.OpaData != nil {
// 	log.Debug("RayParser: parse_suite: OpaData is: %s", *details.OpaData)
// }

// func (parser *RayParser) ParseTestFile(fullPath string, suite *types.SuiteDetails) (*types.TestDetails, error) {
// 	log.Debug("RayParser: ParseTestFile: filename: %s", fullPath)

// 	data, err := util.ReadFile(filepath.Join(fullPath))

// 	if err != nil {
// 		return nil, err
// 	}

// 	elements := strings.Split(*data, "#")

// 	log.Debug("RayParser: ParseTestFile: found elements: %d : %v\n", len(elements), elements)

// 	testDetails := &types.TestDetails{}

// 	for count, element := range elements {

// 		log.Debug("RayParser: ParseTestFile: parsing element: %d -> %s", count, element)

// 		element = strings.TrimSpace(element)

// 		if len(element) == 0 {
// 			continue
// 		}

// 		header, body, err := util.SplitHeaderAndBody(element)

// 		if err != nil {
// 			return nil, err
// 		}

// 		switch *header {
// 		case "Name":
// 			testDetails.Name = util.Chomp(*body)
// 		case "Description":
// 			testDetails.Description = strings.TrimSpace(*body)
// 		case "Data":
// 			tmp := strings.TrimSpace(*body)
// 			testDetails.Data = &tmp
// 		case "Expects":
// 			expectations, err := buildExpectations(*body)
// 			if err != nil {
// 				log.Error("RayParser: ParseTestFile: unable to process expectations: %s", *body)
// 				return nil, err
// 			}
// 			testDetails.Expects = expectations
// 		case "Input":
// 			testDetails.Input = strings.TrimSpace(*body)
// 		case "Rule URL Path":
// 			tmp := strings.TrimSpace(*body)
// 			testDetails.RuleUrlPath = &tmp
// 		}
// 	}

// 	log.Debug("RayParser: parse_suite: Successfully parsed test file: %+v", testDetails)
// 	if testDetails.Expects != nil {
// 		log.Debug("RayParser: parse_suite: Expectations: %+v", *testDetails.Expects)
// 		if testDetails.Expects.Json != nil {
// 			log.Debug("RayParser: parse_suite: Expects JSON: %+v", *testDetails.Expects.Json)
// 		}
// 	}
// 	if testDetails.RuleUrlPath != nil {
// 		log.Debug("RayParser: parse_suite: Rule URL Path: %+v", *testDetails.RuleUrlPath)
// 	}

// 	return testDetails, nil
// }

// // partial implementation for now
// func buildExpectations(body string) (*types.TestExpectations, error) {

// 	log.Warning("RayParser: buildExpectations: partial implementation (JSON only)")

// 	var te = &types.TestExpectations{
// 		Format: "json",
// 		Json:   &body,
// 	}

// 	return te, nil
// }
