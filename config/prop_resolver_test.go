package config

import (
	"os"
	"strings"
	"testing"
)

func TestEnvironment_Substitution(t *testing.T) {

	s := "This is the input string: ${ABC}"

	os.Setenv("ABC", "123")

	resolver := NewPropertyResolver()

	result := resolver.ExpandProperties(s)

	if !strings.Contains(result, "123") {
		t.Errorf("Expected environment substitution. Expected: 123, got: %s", result)
	}
}

func TestProperty_Substitution(t *testing.T) {

	s := "This is the input string: ${ABC}"

	os.Setenv("ABC", "123")

	resolver := NewPropertyResolver()
	resolver.AddProperty("ABC", "456")

	result := resolver.ExpandProperties(s)

	if !strings.Contains(result, "456") {
		t.Errorf("Expected property substitution. Expected: 456, got: %s", result)
	}
}

func TestNo_Substitution(t *testing.T) {

	s := "This is the input string: ${NOOP}"

	os.Setenv("ABCD", "123")

	resolver := NewPropertyResolver()
	resolver.AddProperty("ABCD", "456")

	result := resolver.ExpandProperties(s)

	if !strings.Contains(result, "${NOOP}") {
		t.Errorf("Expected no property substitution. Expected: ${NOOP}, got: %s", result)
	}
}

func TestMultiple_Substitution(t *testing.T) {

	s := "This ${ABCD} is the input string: ${NOOP}"

	os.Setenv("ABCD", "123")

	resolver := NewPropertyResolver()
	resolver.AddProperty("ABCD", "456")

	result := resolver.ExpandProperties(s)

	if !strings.Contains(result, "${NOOP}") {
		t.Errorf("Expected no property substitution. Expected: ${NOOP}, got: %s", result)
	}

	if !strings.Contains(result, "456") {
		t.Errorf("Expected property substitution. Expected: 456, got: %s", result)
	}
}
