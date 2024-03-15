/*
Copyright © 2024 - PACLabs
*/
package main

import (
	"raygun/cmd"
)

/*

   Raygun is a tool for testing Rego policy against OPA servers in a "real-life" facsimile

   OPA has built in testing, but from our experience, it basically requires a very "white-box"
   approach to writing and maintaining test cases.

   Raygun is an attempt at a 'black-box' testing framework for policy, where the testers can
   create JSON to represent inputs, use some sort of pre-generated bundle for the policy code,
   and specify what they expect as the output from that policy.

   We have found this to be tremendously helpful in our own work, and thought it made sense
   to share it with the community.

   Usage:

      raygun execute  <list of .raygun test files>

   The Execute command (in cmd/execute) is the best place to start

*/

func main() {

	cmd.Execute()
}
