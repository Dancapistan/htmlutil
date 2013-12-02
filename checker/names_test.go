package checker

import (
	"fmt"
	"testing"
)

func TestIsHTMLTagName(t *testing.T) {
	assert(t, IsHTMLTagName("TEXTAREA"), "Expected \"TEXTAREA\" to be a valid HTML Tag name, but got false.")
	assert(t, IsHTMLTagName("canvas"), "Expected \"canvas\" to be a valid HTML Tag name, but got false.")
	assert(t, IsHTMLTagName("strong"), "Expected \"strong\" to be a valid HTML Tag name, but got false.")
	refute(t, IsHTMLTagName("Tuesday"), "Expected \"Tuesday\" to NOT be a valid HTML Tag name, but got true.")
	refute(t, IsHTMLTagName("text\narea"), "Expected \"text\\narea\" to NOT be a valid HTML Tag name, but got true.")
}

func ExampleIsHTMLTagName() {
	fmt.Println(IsHTMLTagName("strong"))
	fmt.Println(IsHTMLTagName("tuesday"))
	// Output:
	// true
	// false
}

// BenchmarkIsHTMLTagName_true  5000000         347   ns/op        16 B/op        2 allocs/op
// BenchmarkIsHTMLTagName_true 50000000          26.2 ns/op         0 B/op        0 allocs/op
func BenchmarkIsHTMLTagName_true(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		IsHTMLTagName("canvas")
	}
}

// BenchmarkIsHTMLTagName_false   5000000         382   ns/op        16 B/op        2 allocs/op
// BenchmarkIsHTMLTagName_false  50000000          51.0 ns/op         0 B/op        0 allocs/op
// BenchmarkIsHTMLTagName_false  50000000          39.3 ns/op         0 B/op        0 allocs/op # Go 1.2
func BenchmarkIsHTMLTagName_false(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		IsHTMLTagName("tuesday")
	}
}

func assert(t *testing.T, ok bool, message string) {
	if !ok {
		t.Error(message)
	}
}

func refute(t *testing.T, ok bool, message string) {
	assert(t, !ok, message)
}
