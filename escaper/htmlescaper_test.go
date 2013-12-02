package escaper

import (
	"fmt"
	"testing"
)

func TestEscapeAmbiguousAmpersands(t *testing.T) {
	cases := map[string]string{
		"":                                     "",
		"&\u2601;":                             "&\u2601;",
		"test":                                 "test",
		"none&;":                               "none&;",
		"&amp;":                                "&amp;",
		"this &\u2318&that;.":                  "this &\u2318&amp;that;.",
		"&this;&lt;&this;&that;":               "&amp;this;&lt;&amp;this;&amp;that;",
		"a &tuesday;":                          "a &amp;tuesday;",
		"this &could; be &ambigous;.":          "this &amp;could; be &amp;ambigous;.",
		"no &amp; here":                        "no &amp; here",
		"test &a;&b;&c;&d;&e;&f;&g;&h; \u2318": "test &amp;a;&amp;b;&amp;c;&amp;d;&amp;e;&amp;f;&amp;g;&amp;h; \u2318",
	}
	checkTestCases(t, cases, EscapeAmbiguousAmpersands,
		"EscapeAmbiguousAmpersands")
}

// BenchmarkEscapeAmbiguousAmpersands_simple 10000000         106   ns/op         0 B/op        0 allocs/op
// BenchmarkEscapeAmbiguousAmpersands_simple 20000000          84.8 ns/op         0 B/op        0 allocs/op # Using scanner
// BenchmarkEscapeAmbiguousAmpersands_simple 20000000         112   ns/op         0 B/op        0 allocs/op # Switched from buffer to []byte
func BenchmarkEscapeAmbiguousAmpersands_simple(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = EscapeAmbiguousAmpersands("no &amp; here")
	}
}

// BenchmarkEscapeAmbiguousAmpersands_complex   1000000        1692 ns/op       237 B/op        5 allocs/op
// BenchmarkEscapeAmbiguousAmpersands_complex   1000000        1795 ns/op       302 B/op        6 allocs/op # Using scanner
// BenchmarkEscapeAmbiguousAmpersands_complex   1000000        1166 ns/op       220 B/op        4 allocs/op # Bigger initial buffer
// BenchmarkEscapeAmbiguousAmpersands_complex   1000000        1118 ns/op        98 B/op        2 allocs/op # Switched from buffer to []byte
// BenchmarkEscapeAmbiguousAmpersands_complex   2000000         878 ns/op        98 B/op        2 allocs/op # cache amp indexes
// BenchmarkEscapeAmbiguousAmpersands_complex   5000000         661 ns/op        96 B/op        2 allocs/op # Go 1.2
func BenchmarkEscapeAmbiguousAmpersands_complex(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = EscapeAmbiguousAmpersands("this &could; be &ambigous;.")
	}
}

func ExampleEscapeAmbiguousAmpersands() {
	fmt.Println(EscapeAmbiguousAmpersands("Who &writes; like this?"))
	// Output:
	// Who &amp;writes; like this?
}

func TestEscapeAttributeValueDoubleQuoted(t *testing.T) {

	cases := map[string]string{
		"":             "",
		"a":            "a",
		"\n":           "\n",
		"ok fine":      "ok fine",
		"'fine'":       "'fine'",
		"none&;":       "none&;",
		"&dan;":        "&amp;dan;",              // ambiguous ampersand
		"&dan;\"xxx\"": "&amp;dan;&#34;xxx&#34;", // ambiguous ampersand and double quote
		"\"\u2318\"":   "&#34;\u2318&#34;",       // double quotes
	}

	checkTestCases(t, cases, EscapeAttributeValueDoubleQuoted,
		"EscapeAttributeValueDoubleQuoted")
}

func ExampleEscapeAttributeValueDoubleQuoted() {
	title := EscapeAttributeValueDoubleQuoted(`My name is "Franklin".`)
	fmt.Printf("title=%q", title)
	// Output:
	// title="My name is &#34;Franklin&#34;."
}

//   BenchmarkEscapeAttributeValueDoubleQuoted_none 20000000         109 ns/op         0 B/op        0 allocs/op
func BenchmarkEscapeAttributeValueDoubleQuoted_none(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = EscapeAttributeValueDoubleQuoted("nothing to be escaped")
	}
}

// BenchmarkEscapeAttributeValueDoubleQuoted_quote   5000000        677 ns/op        65 B/op        2 allocs/op
// BenchmarkEscapeAttributeValueDoubleQuoted_quote  5000000         436 ns/op        64 B/op        2 allocs/op # Go 1.2
func BenchmarkEscapeAttributeValueDoubleQuoted_quote(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = EscapeAttributeValueDoubleQuoted(`My name is "Franklin".`)
	}
}

// BenchmarkEscapeAttributeValueDoubleQuoted_both  1000000        1778 ns/op       269 B/op        5 allocs/op
// BenchmarkEscapeAttributeValueDoubleQuoted_both  1000000        1486 ns/op       130 B/op        3 allocs/op # improved escapeAmbiguousAmpersandsBuffer
// BenchmarkEscapeAttributeValueDoubleQuoted_both   2000000        955 ns/op       128 B/op        3 allocs/op # Go 1.2
func BenchmarkEscapeAttributeValueDoubleQuoted_both(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = EscapeAttributeValueDoubleQuoted(`An "&ambiguous;" ampersand.`)
	}
}

func checkTestCases(t *testing.T, cases map[string]string, mutator func(string) string, testLabel string) {
	for input, expected := range cases {
		output := mutator(input)
		if output != expected {
			t.Errorf(
				"Expectation failed. For %s, was expecting %q => %q, but got %q",
				testLabel, input, expected, output)
		}
	}
}
