/*
Copyright © 2024 PACLabs
*/
package util

import (
	"raygun/types"
	"testing"
)

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

func TestLast_Simple(t *testing.T) {

	var tmp []string = []string{"one", "two", "three"}

	last, _ := Last(tmp)

	if last != "three" {
		t.Errorf("expected 'three', got: %s", last)
	}

}

func TestLast_Struct(t *testing.T) {

	var tmp []types.TestExpectation = []types.TestExpectation{{Target: "blah"}}

	last, _ := Last(tmp)

	if last.Target != "blah" {
		t.Errorf("expected 'blah', got: %v", last)
	}

	last = tmp[len(tmp)-1]

	tmp[len(tmp)-1].ExpectationType = "substring"

	if tmp[0].ExpectationType != "substring" {
		t.Errorf("expected to be able to set the expectationType of the last element, but can't: %v", tmp[0])
	}

}
