# raygun

## What does this project do?
   Raygun is a tool for testing Rego policy against OPA servers in a "real-life" facsimile. 

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

* Hat tip to anyone whose code was used
* Inspiration
* etc
