package example3

import rego.v1


default allow := false

deny contains "not ray" if {
  input.name != "ray"
} 

allow := true if {
 count(deny) == 0
}


