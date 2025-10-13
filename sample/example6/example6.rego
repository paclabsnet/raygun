package example6

import rego.v1


default allow := false

deny contains "not ray" if {
  lower(input.name) != "ray"
} 

deny contains "bad method" if {
  not input.method in ["GET","PUT"]
}

allow := true if {
 count(deny) == 0
}


