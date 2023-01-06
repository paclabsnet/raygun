package types

type SuiteDetails struct {
	OpaPort             uint16
	OpaExecutablePath   *string
	OpaWorkingDirectory *string
	RuleUrlPath         *string
	RegoSourceFiles     []string
	RaygunTestFiles     []string
	OpaData             *string
	SuiteName           string
	SuiteDescription    string
}

type TestDetails struct {
	Name        string
	Description string
	Expects     *TestExpectations
	Input       string
	RuleUrlPath *string
	Data        *string
}

type TestExpectations struct {
	Format   string
	Json     *string
	JsonPath []string
}

func NewTestDetails(name string) *TestDetails {

	td := &TestDetails{
		Name: name,
	}

	return td

}
