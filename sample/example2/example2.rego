package example2

import rego.v1


default allow := false


deny contains "not not-ray" if {
	input.name != "not-ray"
}


allow := true if {
	count(deny) == 0
}


