package checker

import (
	"testing"
)

func TestIsValidHTMLTagName(t *testing.T) {
	valid := []string{
		"a",
		"input",
		"html",
		"BODY",
		"pre",
	}
	casesShouldBeTrue(t, valid, IsValidHTMLTagName,
		"Expecting %q to be a valid HTML tag name, but got false.")

	invalid := []string{
		"x-tag",
		" mine \n",
		string(UnicodePOI),
	}
	casesShouldBeFalse(t, invalid, IsValidHTMLTagName,
		"Expecting %q to NOT be a valid HTML tag name, but got true.")
}

// BenchmarkIsValidHTMLTagName 20000000          80.8 ns/op         0 B/op        0 allocs/op
// BenchmarkIsValidHTMLTagName 50000000          53.3 ns/op         0 B/op        0 allocs/op # manual loop
func BenchmarkIsValidHTMLTagName(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsValidHTMLTagName("input")
	}
}
