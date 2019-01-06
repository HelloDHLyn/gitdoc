package gitdoc

import (
	"fmt"
	"testing"
)

var testStringOrDefaultCase = []struct {
	Value        string
	DefaultValue string
	Expected     string
}{
	{"", "Default", "Default"},
	{"A", "B", "A"},
	{"한글", "우리말", "한글"},
}

func TestStringOrDefault(t *testing.T) {
	for _, c := range testStringOrDefaultCase {
		if result := stringOrDefault(c.Value, c.DefaultValue); result != c.Expected {
			fmt.Printf("stringOrDefault return invalid value: %s (expected: %s).\n", result, c.Expected)
			t.Fail()
			return
		}
	}
}
