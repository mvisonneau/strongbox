package rand

import (
	"testing"
)

func TestString(t *testing.T) {
	string1 := String(10)
	string2 := String(10)

	if string1 == string2 {
		t.Fatalf("Expected strings to be differents")
	}

	if len(string1) != 10 {
		t.Fatalf("Expected length of string1 to equal 10, got '%v'", len(string1))
	}
}
