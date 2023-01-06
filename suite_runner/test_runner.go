package suite_runner

import "errors"

type TestRunner struct {
}

func (runner *TestRunner) ExecuteTest() (string, *string, error) {
	return "ERROR", nil, errors.New("not_implemented_test_runner_execute_test")
}
