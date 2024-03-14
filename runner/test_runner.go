package runner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"raygun/config"
	"raygun/log"
	"raygun/types"
	"strings"
)

type TestRunner struct {
	Source types.TestRecord
}

func NewTestRunner(test types.TestRecord) TestRunner {
	testRunner := TestRunner{Source: test}

	return testRunner
}

func (tr TestRunner) Post() (string, error) {

	postUrl := fmt.Sprintf("http://localhost:%d%s", config.OpaPort, tr.Source.DecisionPath)

	bodyString := fmt.Sprintf("{\"input\":%s}", tr.Source.Input.Value)

	bodyBytes := []byte(bodyString)

	response, err := http.Post(postUrl, "application/json", bytes.NewReader(bodyBytes))

	if err != nil {
		log.Error("Attempted to complete POST to %s with payload %s -> %s", postUrl, bodyString, err.Error())
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

func (tr TestRunner) Evaluate(response string) (types.TestResult, error) {

	result := types.TestResult{}

	result.Source = tr.Source

	switch tr.Source.Expects.ExpectationType {
	case "substring":
		actual := strings.ReplaceAll(response, " ", "")

		result.Actual = actual

		expected := strings.ReplaceAll(tr.Source.Expects.Target, " ", "")
		if strings.Contains(actual, expected) {
			result.Status = config.PASS
		} else {
			result.Status = config.FAIL
		}

	default:
		log.Fatal("Unsupported ExpectationType for %s -> %s", tr.Source, tr.Source.Expects.ExpectationType)
	}

	return result, nil
}
