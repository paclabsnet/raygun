/*
Copyright © 2024 PACLabs
*/
package parser

/*
 *   Parses the .raygun files (YAML) and validates that they have the information
 *   we need to create test suites and cases
 *
 *   I elected for a map-based approach instead of using the automatic marshall/unmarshall
 *   of structures, because I wanted to be able to do more in-depth validation.  It also
 *
 */

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

func NewRaygunParser(skip_on_parse_error bool) *RaygunParser {
	p := &RaygunParser{SkipOnParseError: skip_on_parse_error}

	return p
}

/*
 * Struct Methods for RaygunParser
 */

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
		} else {
			suite_list = append(suite_list, suite)
		}

	}

	return suite_list, nil
}

/*
the .raygun file is YAML, so we'll use a YAML parser to process it and then work with the maps
and strings that come out of that process
*/
func (parser *RaygunParser) ParseSuiteFile(suite_file_path string) (types.TestSuite, error) {

	// make sure the directory is in the os-friendly format, so users can reference
	// bundles and input files relative to that directory, instead of relative to
	// where raygun is running.
	//
	suite := types.TestSuite{Directory: filepath.Dir(suite_file_path)}
	suite.Tests = make([]types.TestRecord, 0)

	tree := make(map[string]interface{})

	/*
	 *  Find and parse the YAML content of the .raygun file
	 */

	yamlFile, err := os.ReadFile(suite_file_path)
	if err != nil {
		if parser.SkipOnParseError {
			log.Warning("Parse Error on: %s [%v] . skipping...", suite_file_path, err)
			return suite, err
		} else {
			log.Fatal("Parse Error on %s [%v]", suite_file_path, err)
		}
	}

	err = yaml.Unmarshal(yamlFile, tree)
	if err != nil {
		if parser.SkipOnParseError {
			log.Warning("Parse Error on: %s [%v] . skipping...", suite_file_path, err)
			return suite, err
		} else {
			log.Fatal("Parse Error on %s [%v]", suite_file_path, err)
		}
	}

	/*
	 *  Now we can parse the tree of YAML nodes to create our Test Suite
	 */
	log.Debug("Filepath is: %s, Suite Directory: %s", suite_file_path, suite.Directory)

	err = parser.yamlToSuite(&suite, tree)
	if err != nil {
		if parser.SkipOnParseError {
			log.Warning("Parse Error on: %s [%v] . skipping...", suite_file_path, err)
			return suite, err
		} else {
			log.Fatal("Parse Error on %s [%v]", suite_file_path, err)
		}
	}

	return suite, nil
}

/*
 *  process the tree of yaml elements and build out the test suite and test cases
 */
func (p RaygunParser) yamlToSuite(suite *types.TestSuite, tree map[string]interface{}) error {

	// set up defaults, in case the .raygun file doesn't specify some or all of these
	// values
	suite.Opa.OpaPort = config.OpaPort
	suite.Opa.OpaPath = config.OpaExecutablePath
	suite.Opa.BundlePath = config.OpaBundlePath
	suite.Opa.LogPath = config.OpaLogPath

	//
	//  sorting the keys helps ensure they're in a consistent order from run to run
	//
	for _, k := range util.SortMapKeys(tree) {

		v := tree[k]

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

	// validate the suite

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

	for _, k := range util.SortMapKeys(tree) {

		v := tree[k]

		switch k {
		case "path":
			suite.Opa.OpaPath = v.(string)
		case "bundle-path":
			suite.Opa.BundlePath = filepath.Join(suite.Directory, v.(string))
			log.Debug("OPA BundlePath is: %s", suite.Opa.BundlePath)
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

/*
 *  Process the individual test subrecords
 */
func (p RaygunParser) yamlToTest(suite *types.TestSuite, tree map[string]interface{}) error {

	test := types.TestRecord{Suite: *suite}

	for _, k := range util.SortMapKeys(tree) {

		v := tree[k]

		switch k {
		case "name":
			if util.IsString(v) {
				test.Name = v.(string)
			} else {
				return fmt.Errorf("Invalid test name type: %v, expecting string", v)
			}
		case "description":
			if util.IsString(v) {
				test.Description = v.(string)
			} else {
				return fmt.Errorf("Invalid test description type: %v, expecting string", v)
			}
		case "decision-path":
			if util.IsString(v) {
				test.DecisionPath = v.(string)
			} else {
				return fmt.Errorf("Invalid test DecisionPath type: %v, expecting string", v)
			}
		case "expects":
			err := p.yamlToExpectations(&test, v.(map[string]interface{}))
			if err != nil {
				return err
			}
		case "input":
			err := p.yamlToInputJSON(&test, v.(map[string]interface{}))
			if err != nil {
				return err
			}
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

	for _, k := range util.SortMapKeys(tree) {

		v := tree[k]

		switch k {
		case "type":
			if util.IsString(v) {
				test.Expects.ExpectationType = v.(string)
			} else {
				return fmt.Errorf("invalid Expects.ExpectationType value: %v, expecting string", v)
			}
		case "target":
			if util.IsString(v) {
				test.Expects.Target = v.(string)
			} else {
				return fmt.Errorf("Invalid Expects.Target value: %v, expecting string", v)
			}
		default:
			return fmt.Errorf("unknown/unsupported 'expects' section key: %s", k)
		}
	}

	return nil
}

/*
 * Process where to find the input.
 *
 * type:
 *   inline - a JSON string embedded in the .raygun file
 *   json-file - a reference to an external JSON file to be read at test time.
 */
func (p RaygunParser) yamlToInputJSON(test *types.TestRecord, tree map[string]interface{}) error {

	for _, k := range util.SortMapKeys(tree) {

		v := tree[k]

		switch k {
		case "type":
			if util.IsString(v) {
				test.Input.InputType = v.(string)
			} else {
				return fmt.Errorf("Invalid input type value %v, expecting string", v)
			}
		case "value":
			if util.IsString(v) {
				test.Input.Value = v.(string)
			} else {
				return fmt.Errorf("invalid input Value: %v, expecting string", v)
			}
		default:
			return fmt.Errorf("unknown/unsupported 'input' section key: %s", k)
		}
	}

	return nil
}
