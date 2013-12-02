package checker

import (
	"fmt"
	"testing"
)

func TestIsValidHtml4IdValue(t *testing.T) {
	var cases = map[string]bool{
		"":       false,
		"a b":    false,
		"1abc":   false,
		"\u2601": false,
		"abc":    true,
		"a":      true,
		"a9":     true,
		"A-":     true,
		"a_b":    true,
		"a:b":    true,
		"a.b.c.": true,
	}

	for input, expected := range cases {
		actual := IsValidHtml4IdValue(input)
		if expected != actual {
			t.Errorf("Expecting IsValidHtml4IdValue(%q) to be %v, got %v.\n",
				input, expected, actual)
		}
	}
}

// BenchmarkIsValidHtml4IdValue  20000000          81.2 ns/op         0 B/op        0 allocs/op # baseline
func BenchmarkIsValidHtml4IdValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsValidHtml4IdValue("introduction")
	}
}

func ExampleIsValidHtml4IdValue() {
	fmt.Println(IsValidHtml4IdValue("introduction"))
	fmt.Println(IsValidHtml4IdValue("last remarks"))
	// Output:
	// true
	// false
}

func TestIsValidHtml5IdValue(t *testing.T) {
	var cases = map[string]bool{
		"":       false, // "must be at least one character long"
		"a b":    false, // "must not contain any space characters"
		"1abc":   true,
		"\u2601": true,
		"abc":    true,
		"a":      true,
		"a9":     true,
		"A-":     true,
		"a_b":    true,
		"a:b":    true,
		"a.b.c.": true,
	}

	for input, expected := range cases {
		actual := IsValidHtml5IdValue(input)
		if expected != actual {
			t.Errorf("Expecting IsValidHtml5IdValue(%q) to be %v, got %v.\n",
				input, expected, actual)
		}
	}
}

// BenchmarkIsValidHtml5IdValue   5000000         489 ns/op         0 B/op        0 allocs/op
//
// TODO: Investigate why this is so much slower than IsValidAttributeName.
//
func BenchmarkIsValidHtml5IdValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsValidHtml5IdValue("introduction")
	}
}

func ExampleIsValidHtml5IdValue() {
	fmt.Println(IsValidHtml5IdValue("introduction"))
	fmt.Println(IsValidHtml5IdValue("last remarks"))
	// Output:
	// true
	// false
}

func TestIsValidCss3IdValue(t *testing.T) {
	var cases = map[string]bool{
		"":   false,
		"a":  false, // starts with #
		"#a": true,
	}

	for input, expected := range cases {
		actual := IsValidCss3IdValue(input)
		if expected != actual {
			t.Errorf("Expecting IsValidCss3IdValue(%q) to be %v, got %v.\n",
				input, expected, actual)
		}
	}
}

func TestIsValidCss3Identifier(t *testing.T) {
	var cases = map[string]bool{
		"":          false,
		"1a":        false, // cannot start with "a digit".
		"--lua":     false, // cannot start with "two hyphens".
		"-1a":       false, // cannot start with "a hyphen followed by a digit".
		"a#":        false,
		"a{b}":      false,
		"a|b":       false,
		"!":         false,
		"a b":       false,
		"2":         false,
		`\32`:       true, // "\32 is allowed at the start of a class name, even though "2" is not"
		"abc123DEF": true, // can contain "the characters [a-zA-Z0-9] and".
		"a\u00A0b":  true, // can contain "ISO 10646 characters U+00A0 [&nbsp;] and higher". (NOTE: ISO 10646 code points map to Unicode code points.)
		"¡":         true, // U+00A1
		"©2005":     true, // U+00A9
		"a\u2601b":  true, // "and higher"
		"a-b":       true, // can contain "the hyphen".
		"a_b":       true, // can contain "the underscore".
		`B\&W\?`:    true, // illegal characters escaped
		`B\26 W\3F`: true,
		`a\ b`:      true,
		`te\st`:     true, // "The identifier 'te\st' is exactly the same identifier as 'test'."
	}

	// From http://www.w3.org/TR/CSS21/syndata.html#value-def-identifier :

	// If a character in the range [0-9a-fA-F] follows the hexadecimal number,
	// the end of the number needs to be made clear. There are two ways to do
	// that:
	//
	//   - with a space (or other white space character): "\26 B" ("&B"). In
	//     this case, user agents should treat a "CR/LF" pair (U+000D/U+000A) as
	//     a single white space character.
	//
	//   - by providing exactly 6 hexadecimal digits: "\000026B" ("&B")
	//
	// In fact, these two methods may be combined. Only one white space
	// character is ignored after a hexadecimal escape. Note that this means
	// that a "real" space after the escape sequence must be doubled.
	//
	// The identifier "te\st" is exactly the same identifier as "test"

	for input, expected := range cases {
		actual := IsValidCss3Identifier(input)
		if expected != actual {
			t.Errorf("Expecting IsValidCss3Identifier(%q) to be %v, got %v.\n",
				input, expected, actual)
		}
	}
}

func ExampleIsValidCss3Identifier() {
	fmt.Println(IsValidCss3Identifier("hullo"))
	// Output:
	// true
}

// BenchmarkIsValidCss3Identifier  20000000         116   ns/op         0 B/op        0 allocs/op
// BenchmarkIsValidCss3Identifier  20000000          81.6 ns/op         0 B/op        0 allocs/op # Go 1.2
func BenchmarkIsValidCss3Identifier(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsValidCss3Identifier("wrapper2")
	}
}
