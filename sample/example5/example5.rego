package example5

import rego.v1


default allow := false

deny contains "not ray" if {
  lower(input.name) != "ray"
} 

deny contains "bad method" if {
  input.method != "GET"
}

allow := true if {
 count(deny) == 0
}


