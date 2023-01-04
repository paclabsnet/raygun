# raygun

Raygun is a program for testing Rego rules. 

1. You specify a directory tree for .rego source code
2. You specify the data JSON to include in OPA
3. You specify a specific input JSON
4. You specify the expected output JSON

Using that information, we launch OPA with the .rego code, set up the data section, invoke OPA via the RESTful API and compare the actual 
output to the expected output. This represents a test.

Since many tests will likely use the same .rego source, and the same data, we provide a test suite where we launch an OPA instance with data and .rego files, and then run a bunch of tests against the same OPA / Rego / Data infrastructure.  This will save time


