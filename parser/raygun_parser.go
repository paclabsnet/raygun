/*
Copyright © 2025 PACLabs
*/
package parser

/*
 *   Parses the .raygun files (YAML) and validates that they have the information
 *   we need to create test suites and cases
 *
 *   Originally, I used a map structure, but the ability to include metadata in
 *   the types makes it much easier to parse, which makes changes much easier
 *   so except in a few places, I am parsing the yaml using type metastructure
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

	for _, raygun_filename := range raygun_file_list {

		data, err := os.ReadFile(raygun_filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		suite := CreateEmptySuite(raygun_filename)

		err = yaml.Unmarshal(data, &suite)

		var skip bool = false

		if err != nil {
			if !parser.SkipOnParseError {
				log.Fatal("Parse error on suite file: %s [%v]", raygun_filename, err)
				// this will exit, no need to handle this elegantly
			} else {
				log.Warning("Parse error on suite file: %s -> %v .. skipping", raygun_filename, err)
				skip = true
			}
		}

		log.Debug("RaygunParser.Parse: %v", suite)

		err = parser.parseExpectations(&suite)

		if err != nil {
			if !parser.SkipOnParseError {
				log.Fatal("Parse error on suite file: %s [%v]", raygun_filename, err)
			} else {
				log.Warning("Parse error on suite file: %s -> %v .. skipping", raygun_filename, err)
				skip = true
			}
		}

		if !skip {
			suite_list = append(suite_list, suite)
		}

	}

	return suite_list, nil
}

func (parser *RaygunParser) parseExpectations(suite *types.TestSuite) error {

	for i := range suite.Tests {

		// this was an interesting and annoying bug - I thought
		// for _, tr := range suite.Tests would give me the record
		// so I could pass it by reference to yamlToExpectationsMap
		//
		// but instead, it appeared to make a copy, so the parsing data was lost

		err := parser.yamlToExpectationsMap(&suite.Tests[i], suite.Tests[i].ExpectsMap)
		if err != nil {
			return err
		}

	}

	return nil
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
	suite.Opa.BundlePath = config.OpaBundleUrl
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
				return fmt.Errorf("invalid test name type: %v, expecting string", v)
			}
		case "description":
			if util.IsString(v) {
				test.Description = v.(string)
			} else {
				return fmt.Errorf("invalid test description type: %v, expecting string", v)
			}
		case "decision-path":
			if util.IsString(v) {
				test.DecisionPath = v.(string)
			} else {
				return fmt.Errorf("invalid test DecisionPath type: %v, expecting string", v)
			}
		case "expects":
			if util.IsArray(v) {
				err := p.yamlToExpectationsArray(&test, v.([]interface{}))
				if err != nil {
					return err
				}
			} else if util.IsMap(v) {
				err := p.yamlToExpectationsMap(&test, v.(map[string]interface{}))
				if err != nil {
					return err
				}
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
 * after writing this originally with explicit elements for type and target, I
 * realized that there was a nice implicit shorthand provided by allowing
 * for the type to be the key, and the value to be interpreted as the target
 *
 * But I don't want to break existing .raygun files so I want to support both
 */
func (p RaygunParser) yamlToExpectationsArray(test *types.TestRecord, tree_array []interface{}) error {

	for _, element := range tree_array {

		tree := element.(map[string]interface{})

		err := p.yamlToExpectationsMap(test, tree)

		if err != nil {
			return err
		}
	}

	return nil
}

/*
 * Process a single TestExpectations map
 */
func (p RaygunParser) yamlToExpectationsMap(test *types.TestRecord, tree map[string]interface{}) error {

	test.ExpectData = append(test.ExpectData, types.TestExpectation{})

	for _, k := range util.SortMapKeys(tree) {

		v := tree[k]

		switch k {
		case "type":
			if util.IsString(v) {

				test.ExpectData[len(test.ExpectData)-1].ExpectationType = v.(string)

			} else {
				return fmt.Errorf("invalid Expects.ExpectationType value: %v, expecting string", v)
			}
		case "target":
			if util.IsString(v) {

				test.ExpectData[len(test.ExpectData)-1].Target = v.(string)
			} else {
				return fmt.Errorf("invalid Expects.Target value: %v, expecting string", v)
			}
		case "substring":
			if util.IsString(v) {

				test.ExpectData[len(test.ExpectData)-1].ExpectationType = "substring"
				test.ExpectData[len(test.ExpectData)-1].Target = v.(string)

			} else {
				return fmt.Errorf("invalid substring value: %v, expecting string", v)
			}
		default:
			return fmt.Errorf("unknown/unsupported 'expects' section key: %s", k)
		}

	}

	// log.Debug("Test %s ExpectData array: %v", test.Name, test.ExpectData)

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
				return fmt.Errorf("invalid input type value %v, expecting string", v)
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
