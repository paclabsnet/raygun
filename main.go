/*
Copyright © 2022 John Brothers <johnbr@paclabs.net>
*/
package main

import (
	"raygun/cmd"
)

/*
how do I want to organize this thing?

we start with a command that specifies what to do (test)

we look for test suite filess in the specified place(s)

we parse the test suite files
	- each suite file parse produces suite metadata

- we iterate over test suites, using the suite metadata
	- we parse the individual test files specified in the suite metadata
		- this produces test details
	- we build the opa launch configuration, using a mix of defaults, cli flags and suite metadata
	- we launch OPA with the specified launch configuration
	- for each file:
		- we launch a test runner, with the test details and the opa launch configuration
			- the test runner sends the input from the test file to OPA
			- the test runner gets back the results from OPA and compare to expected and
				return our findings
		- we log the findings
	= we shut down OPA

*/

func main() {

	cmd.Execute()
}
