/*
Copyright © 2024 PACLabs
*/
package util

import "testing"

func TestListify_Empty(t *testing.T) {

	s := ""

	results := Listify(s)

	if len(results) != 0 {
		t.Errorf("got non-empty result for empty string")
	}

}

func TestListify_OneLine(t *testing.T) {

	s := "one line"

	results := Listify(s)

	if len(results) != 1 {
		t.Errorf("got %d lines, expected 1", len(results))
	}

}

func TestListify_OneLineWithNewLines(t *testing.T) {

	s := "\none line\n"

	results := Listify(s)

	if len(results) != 1 {
		t.Errorf("got %d lines, expected 1", len(results))
	}

}

func TestListify_ThreeLineWithNewLines(t *testing.T) {

	s := `
     
	one
	   two
	three
            
	`

	results := Listify(s)

	if len(results) != 3 {
		t.Errorf("got %d lines, expected 3", len(results))
	}

}
