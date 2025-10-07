package example4

import rego.v1


default allow := false

deny contains "not ray" if {
  input.name != "ray"
} 

deny contains "not delegated to bob" if {
  input.delegate != "bob"
}

allow := true if {
 count(deny) == 0
}



