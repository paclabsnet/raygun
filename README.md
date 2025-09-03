# raygun

Current Version: 0.1.4

## What does this project do?
   Raygun is a tool for testing Rego policy against OPA servers in a way that resembles a "real-world" usage.  Specifically:

   1. It doesn't require any access to the Rego codebase - you need a bundle, and that's it.
   2. The people developing the tests don't need to know anything about Rego.  They just need to know how to create the test cases (YAML) and the test inputs and outputs (JSON)

   It is written in go, so it should be easy to port to various platforms


## Why is this project useful

   OPA has built in testing, but from our experience, it basically requires a very "white-box"
   approach to writing and maintaining test cases.

   Raygun is an attempt at a 'black-box' testing framework for policy, where the testers can
   create JSON to represent inputs, use some sort of pre-generated bundle for the policy code,
   and specify what they expect as the output from that policy.  They won't have to know anything
   about Rego, or about how OPA works.

   We have found this to be tremendously helpful in our own work, and thought it made sense
   to share it with the community.

## How do I get started

* build the raygun executable for your platform (go build)
* put the executable somewhere in your path
* make sure you have OPA somewhere on your path
* use OPA to build a bundle.tar.gz of Rego code & data
* create a .raygun test case, using the examples in sample/ as a starting point
   * determine the appropriate URL path for the policy you want to test
   * create the appropriate input json
   * identify what you expect the output to be and determine the substring of that output that indicates that the policy behaved the way you expected
* raygun execute ```testfile.raygun```

### Helpful flags

```-d``` or ```--debug``` for debug logs

```-v``` or ```--verbose``` for verbose (more detail for failures)

set the environment variable RAYGUN_OPA_EXEC if it isn't "opa"

```--report-format json``` if you want JSON output instead of a text report

```--stop-on-failure``` if you want the testing to stop at the first failed test

### Building/Installing

If you have an executable for raygun, put it somewhere on your path.

If you want to build the raygun executable:

```
go build
```

### Usage

```
raygun execute  <list of .raygun test files>
```

#### Example
```
raygun execute --verbose  sample/*/*.raygun
```

### Understanding the code

   execute.go (in cmd/) is the best place to start if you want to understand what this code does

### Areas that need help

* additional expectation capabilities
   * regex
   * jsonpath
   * compare key/value pairs from the response against expecations

### Running the tests

```
go test  raygun/util
```


## Goals of the Project

   * Providing a black-box testing apparatus for testing Rego Policy
   * Making it easy to embed this into a build pipeline


## Where can I get more help, if I need it?

   for simple questions/suggestions/etc: email info@paclabs.io

   if you'd like formal support: email sales@paclabs.io



## Contributing

Please read [CONTRIBUTING.md]() for details on our code of conduct, and the process for submitting pull requests to us.


## Authors

* **John Brothers** - *Initial work* - [johndbro1](https://github.com/johndbro1)

See also the list of [contributors](https://github.com/paclabsnet/raygun/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* the OPA team for inspiring me to write this, especially Anders Eknert
* Inspiration : all the test tools I've used in the past


## Future Features

1. multiple expectations
2. eopa support
3. simplified expectation format
4. jsonpath expectation
5. input templates (so you can pull in environment variables and other external configuration (date, time, things like that))

