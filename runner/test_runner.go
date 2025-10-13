/*
Copyright © 2024 PACLabs
*/
package runner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"raygun/config"
	"raygun/jwt"
	"raygun/log"
	"raygun/types"
	"raygun/util"
	"strings"
)

type TestRunner struct {
	Source types.TestRecord
}

func NewTestRunner(test types.TestRecord) TestRunner {
	testRunner := TestRunner{Source: test}

	return testRunner
}

var jwtBuilder jwt.JWTBuilder = jwt.NewJWTBuilder()

func (tr TestRunner) Post() (string, error) {

	//	postUrl := fmt.Sprintf("http://localhost:%d%s", config.OpaPort, tr.Source.DecisionPath)
	postUrl := fmt.Sprintf("%s%s", tr.Source.Suite.Opa.GetAgentUrl(), tr.Source.DecisionPath)

	log.Debug("postUrl: %s", postUrl)

	preExpansionInput := ""

	switch tr.Source.Input.InputType {
	case "inline":

		// read the JSON data directly from the .raygun file
		preExpansionInput = optionally_add_input_key(tr.Source.Input.Value)
	case "json-file":

		// read the JSON data from a file
		log.Debug("Suite Directory: %s , filename: %s", tr.Source.Suite.Directory, tr.Source.Input.Value)
		tmp, err := util.ReadFile(tr.Source.Suite.Directory, tr.Source.Input.Value)
		if err != nil {
			return "", err
		}

		preExpansionInput = optionally_add_input_key(tmp)

	default:
		return "", fmt.Errorf("unsupported input type: %s", tr.Source.Input.InputType)
	}

	// we only process the JWT data if there's anything present to process, otherwise
	// it's safe to ignore
	if tr.Source.Jwt.Active || (len(tr.Source.Jwt.Claims.Custom) > 0) {

		jwt_string, err := jwtBuilder.Generate(tr.Source.Suite, tr.Source.Jwt)

		if err != nil {
			return "", err
		}

		log.Debug("Test: %s Generated JWT: %s", tr.Source.Name, jwt_string)

		// by convention, we'll put the JWT into a property named
		// RAYGUN_GENERATED_JWT
		// so the user can just use ${RAYGUN_GENERATED_JWT} in their input
		// document
		if jwt_string != "" {
			config.Resolver.AddProperty("RAYGUN_GENERATED_JWT", jwt_string)
		}
	}

	// substitute any ${} tokens in the input with their appropriate values
	// which are pulled either from properties or from the environment
	bodyString := config.Resolver.ExpandProperties(preExpansionInput)

	return _post(postUrl, bodyString)

}

/*
 *  This is the most meaningful step of the entire process - does the response from OPA
 *  match the expectations defined in the test case
 */
func (tr TestRunner) Evaluate(response string) (types.TestResult, error) {

	result := types.TestResult{}

	result.Source = tr.Source

	log.Debug("Expectations: %v . Actual: %s", tr.Source.ExpectData, response)

	for _, expected := range tr.Source.ExpectData {

		if result.Status != config.FAIL {

			switch expected.ExpectationType {
			case "substring":
				compressed_actual := util.RemoveAllWhitespace(response)
				result.Actual = response

				expected := util.RemoveAllWhitespace(expected.Target)
				if strings.Contains(compressed_actual, expected) {
					result.Status = config.PASS
				} else {
					result.Status = config.FAIL
				}

			default:
				log.Fatal("Unsupported ExpectationType for %s -> %s", tr.Source, expected.ExpectationType)
			}
		}
	}

	return result, nil
}

/*
 *  Keeping it really simple until we know we need something more sophisticated
 */
func optionally_add_input_key(json string) string {

	no_whitespace := util.RemoveAllWhitespace(json)

	if strings.HasPrefix(no_whitespace, "{\"input\"") {
		return json
	}

	return fmt.Sprintf("{\"input\":%s}", json)

}

/*
 *  the core implementation of the http post and returning the response
 */
func _post(url string, body string) (string, error) {

	log.Debug("Request URL: %s", url)
	log.Debug("Request Content: \n%s", body)

	bodyBytes := []byte(body)

	response, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))

	if err != nil {
		log.Error("Attempted to complete POST to %s with payload %s -> %s", url, body, err.Error())
		return "", err
	}

	defer response.Body.Close()

	builderBuffer := new(strings.Builder)

	_, err = io.Copy(builderBuffer, response.Body)

	if err != nil {
		log.Error("Error reading body of response: %s", err.Error())
		return "", err
	}

	log.Debug("Response Content: %s", builderBuffer.String())

	return builderBuffer.String(), nil

}
