package checker

import (
	"fmt"
	"testing"
)

const UnicodePOI = '\u2318'

func TestIsValidAttributeName(t *testing.T) {

	valid := []string{
		"name",
		"type",
		"data-valid",
		string(UnicodePOI),
	}
	casesShouldBeTrue(t, valid, IsValidAttributeName,
		"Expected attribute %#v to be valid, got invalid.")

	invalid := []string{
		"",                 // "must consist of one or more characters"
		" ",                // "other than the space characters,"
		"\t",               // "other than the space characters,"
		"\n",               // "other than the space characters,"
		"\r\n",             // "other than the space characters,"
		"\f",               // "other than the space characters,"
		"this is invalid",  // "other than the space characters,"
		"not-\u0000-valid", // "U+0000 NULL, "
		"i-am-\"happy\"",   // "U+0022 QUOTATION MARK ("), "
		"i-am-'happy'",     // "U+0027 APOSTROPHE ('), "
		"<input>",          // "">" (U+003E),"
		"either/or",        // ""/" (U+002F), "
		"this /",           // "other than the space characters,"
		"mine=yours",       // ""=" (U+003D) characters, "
		"\a",               // "the control characters, "
		"\uffff",           // "and any characters that are not defined by Unicode"
	}
	casesShouldBeFalse(t, invalid, IsValidAttributeName,
		"Expected attribute %#v to be invalid, got valid.")
}

func ExampleIsValidAttributeName() {
	fmt.Println(IsValidAttributeName("example"))
	fmt.Println(IsValidAttributeName("not valid"))
	// Output:
	// true
	// false
}

// BenchmarkIsValidAttributeName_valid 50000000          30.2 ns/op         0 B/op        0 allocs/op
func BenchmarkIsValidAttributeName_valid(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsValidAttributeValue("data-test")
	}
}

// BenchmarkIsValidAttributeName_invalid 50000000          28.9 ns/op         0 B/op        0 allocs/op
func BenchmarkIsValidAttributeName_invalid(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsValidAttributeValue("data test")
	}
}

func TestIsValidAttributeValue(t *testing.T) {

	// Attribute values are a mixture of text and character references,
	// except with the additional restriction that the text cannot contain
	// an ambiguous ampersand.

	valid := []string{
		"data",
		"",
		"some other &#xBEEF;",
		"And &#1010; decimal",
		"And &amp; is fine",
		string(UnicodePOI),
	}
	casesShouldBeTrue(t, valid, IsValidAttributeValue,
		"Expected attribute value %#v to be valid, but got invalid.")

	invalid := []string{
		"&funky;",
	}
	casesShouldBeFalse(t, invalid, IsValidAttributeValue,
		"Expected attribute value %#v to be invalid, but got valid.")
}

//  BenchmarkIsValidAttributeValue	  200000	         12224 ns/op
//  BenchmarkIsValidAttributeValue_valid    500000        3316 ns/op       214 B/op        3 allocs/op # baseline
//  BenchmarkIsValidAttributeValue_valid  20000000         123 ns/op         0 B/op        0 allocs/op # optimized HasAmbiguousAmpersand
func BenchmarkIsValidAttributeValue_valid(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		IsValidAttributeValue("and &amp; is fine")
	}
}

// BenchmarkIsValidAttributeValue_invalid    500000        3888   ns/op       214 B/op        3 allocs/op # baseline
// BenchmarkIsValidAttributeValue_invalid  10000000         167   ns/op         0 B/op        0 allocs/op # optimized HasAmbiguousAmpersand
// BenchmarkIsValidAttributeValue_invalid  20000000          83.5 ns/op         0 B/op        0 allocs/op # Go 1.2
func BenchmarkIsValidAttributeValue_invalid(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		IsValidAttributeValue("and &funk; is fine") // ambiguous ampersand
	}
}

func TestIsValidAttributeValueUnquoted(t *testing.T) {
	valid := []string{
		"yes",
		"true",
		"email",
		"7",
		"wut?",
		string(UnicodePOI),
	}
	casesShouldBeTrue(t, valid, IsValidAttributeValueUnquoted,
		"Expected unquoted attribute value %#v to be valid, but got invalid.")

	invalid := []string{
		"",         // "must not be the empty string"
		"not this", // "must not contain any literal space characters"
		`"wut?"`,
	}
	casesShouldBeFalse(t, invalid, IsValidAttributeValueUnquoted,
		"Expected unquoted attribute value %#v to be invalid, but got valid.")
}

// BenchmarkIsValidAttributeValueUnquoted   5000000         467 ns/op         0 B/op        0 allocs/op # baseline
func BenchmarkIsValidAttributeValueUnquoted(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsValidAttributeValueUnquoted("email")
	}
}

func TestIsValidAttributeValueSingleQuoted(t *testing.T) {
	valid := []string{
		"this is valid",
		string(UnicodePOI),
	}
	casesShouldBeTrue(t, valid, IsValidAttributeValueSingleQuoted,
		"Expected single-quoted attribute value %#v to be valid, but got invalid.")

	invalid := []string{
		"this 'is not' valid",
	}
	casesShouldBeFalse(t, invalid, IsValidAttributeValueSingleQuoted,
		"Expected single-quoted attribute value %#v to be invalid, but got valid.")
}

func ExampleIsValidAttributeValueSingleQuoted() {
	fmt.Println(IsValidAttributeValueSingleQuoted("yes"))
	fmt.Println(IsValidAttributeValueSingleQuoted("'no'"))
	// Output:
	// true
	// false
}

// BenchmarkIsValidAttributeValueSingleQuoted  20000000          78.8 ns/op         0 B/op        0 allocs/op # baseline, ContainsAny
// BenchmarkIsValidAttributeValueSingleQuoted  50000000          31.4 ns/op         0 B/op        0 allocs/op # IndexRune
func BenchmarkIsValidAttributeValueSingleQuoted(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsValidAttributeValueSingleQuoted("yes")
	}
}

func TestIsValidAttributeValueDoubleQuoted(t *testing.T) {
	valid := []string{
		"this is valid",
		string(UnicodePOI),
	}
	casesShouldBeTrue(t, valid, IsValidAttributeValueDoubleQuoted,
		"Expected double-quoted attribute value %#v to be valid, but got invalid.")

	invalid := []string{
		"this \"is not\" valid",
	}
	casesShouldBeFalse(t, invalid, IsValidAttributeValueDoubleQuoted,
		"Expected double-quoted attribute value %#v to be invalid, but got valid.")
}

func ExampleIsValidAttributeValueDoubleQuoted() {
	fmt.Println(IsValidAttributeValueDoubleQuoted("yes"))
	fmt.Println(IsValidAttributeValueDoubleQuoted(`"no"`))
	// Output:
	// true
	// false
}

// BenchmarkIsValidAttributeValueDoubleQuoted  10000000         237   ns/op         0 B/op        0 allocs/op # baseline, ContainsAny
// BenchmarkIsValidAttributeValueDoubleQuoted  50000000          67.8 ns/op         0 B/op        0 allocs/op # IndexRune
func BenchmarkIsValidAttributeValueDoubleQuoted(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = IsValidAttributeValueDoubleQuoted("how are you?")
	}
}

func TestHasAmbiguousAmpersand(t *testing.T) {

	ambig := []string{
		"This &could; be",
		"&amp; &this;",
		"&a;",
		"&9;",
		"what; &about; this;&?",
		"&this &is; also;",
		"\u2318 &poi; \u2318",
	}
	casesShouldBeTrue(t, ambig, HasAmbiguousAmpersand,
		"Expected HasAmbiguousAmpersand(%q) to be true, go false")

	notAmbig := []string{
		"",
		"this is not ambiguous",
		"this & ain't ambiguous",
		"&#this",
		"neither &; this",
		"what; about& this?",
		"some; &and ;",
		"This &amp; that.",
		"Put &lt;input/>",
		"this & that;.",
		"This & that; &amp; others",
		"&this &amp; that;",
		"&\u2318;",
		"&\uFFFD;",
	}

	casesShouldBeFalse(t, notAmbig, HasAmbiguousAmpersand,
		"Expected HasAmbiguousAmpersand(%q) to be false, got true")
}

// BenchmarkHasAmbiguousAmpersand_false_with_amp   1000000        1925   ns/op       198 B/op        2 allocs/op # baseline
// BenchmarkHasAmbiguousAmpersand_false_with_amp  20000000          90.7 ns/op         0 B/op        0 allocs/op # manual loop
// BenchmarkHasAmbiguousAmpersand_false_with_amp 20000000           60.6 ns/op         0 B/op        0 allocs/op # With scanner
func BenchmarkHasAmbiguousAmpersand_false_with_amp(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		HasAmbiguousAmpersand("&amp;") // valid ref, so not ambiguous
	}
}

// BenchmarkHasAmbiguousAmpersand_true   1000000        1943  ns/op       198 B/op        2 allocs/op # baseline
// BenchmarkHasAmbiguousAmpersand_true  20000000         109  ns/op         0 B/op        0 allocs/op # manual loop
// BenchmarkHasAmbiguousAmpersand_true 20000000          65.5 ns/op         0 B/op        0 allocs/op # With scanner
func BenchmarkHasAmbiguousAmpersand_true(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		HasAmbiguousAmpersand("&yep;")
	}
}

// BenchmarkHasAmbiguousAmpersand_false_without_amp 100000000         18.1 ns/op         0 B/op        0 allocs/op # baseline
// BenchmarkHasAmbiguousAmpersand_false_without_amp  50000000         23.5 ns/op         0 B/op        0 allocs/op # With scanner
func BenchmarkHasAmbiguousAmpersand_false_without_amp(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		HasAmbiguousAmpersand("nope")
	}
}

func ExampleHasAmbiguousAmpersand() {
	fmt.Println(HasAmbiguousAmpersand("this &amp; that"))
	fmt.Println(HasAmbiguousAmpersand("this & that"))
	fmt.Println(HasAmbiguousAmpersand("nothing here"))
	fmt.Println(HasAmbiguousAmpersand("Can &what; be one?"))
	// Output:
	// false
	// false
	// false
	// true
}

func TestisUnicodeNonCharacter(t *testing.T) {

	// From Wikipedia:
	//
	//     Certain noncharacter code points are guaranteed never to be used for
	//     encoding characters, although applications may make use of these code
	//     points internally if they wish. There are sixty-six noncharacters:
	//     U+FDD0..U+FDEF and any code point ending in the value FFFE or FFFF
	//     (i.e. U+FFFE, U+FFFF, U+1FFFE, U+1FFFF, ... U+10FFFE, U+10FFFF). The
	//     set of noncharacters is stable, and no new noncharacters will ever be
	//     defined.[13]
	//
	// https://en.wikipedia.org/wiki/Unicode

	chars := []rune{
		'a',
		'\n',
		'\u1234', // U+1234 = &#4660;
		'喜',
	}
	runeCasesShouldBeFalse(t, chars, isUnicodeNonCharacter,
		"Expected rune %#v to be a valid Unicode character, but got true from isUnicodeNonCharacter")

	nonchars := []rune{
		'\uFDD0',
		'\uFDD1',
		'\uFDEF',
		'\uFDEE',
		'\uFFFE',
		'\uFFFF',
	}
	runeCasesShouldBeTrue(t, nonchars, isUnicodeNonCharacter,
		"Expected rune %#v to be a Unicode non-character, but got false from isUnicodeNonCharacter")
}

func runeCasesShouldBeTrue(t *testing.T, cases []rune, test func(rune) bool, pattern string) {
	for _, arg := range cases {
		if test(arg) != true {
			t.Errorf(pattern, arg)
		}
	}
}

func runeCasesShouldBeFalse(t *testing.T, cases []rune, test func(rune) bool, pattern string) {
	for _, arg := range cases {
		if test(arg) != false {
			t.Errorf(pattern, arg)
		}
	}
}

func TestNamedReferenceScanner_Next(t *testing.T) {

	var cases = []struct {
		Scanner *NamedReferenceScanner
		ExpRef  string
		ExpIdx  int
	}{
		{&NamedReferenceScanner{"value", -1}, "", -1},
		{&NamedReferenceScanner{"&value;", -1}, "value", 0},
		{&NamedReferenceScanner{"&value;", 6}, "", -1},
		{&NamedReferenceScanner{"a &value;", -1}, "value", 2},
		{&NamedReferenceScanner{"a &value;", 9}, "", -1},
		{&NamedReferenceScanner{"a &value;", 3}, "", -1},
		{&NamedReferenceScanner{";&", -1}, "", -1},
		{&NamedReferenceScanner{"a &;", -1}, "", -1},
		{&NamedReferenceScanner{"&not \u2318 this;", -1}, "", -1},
		{&NamedReferenceScanner{"&but \u2318 &this;", -1}, "this", 9},
		{&NamedReferenceScanner{"&but; &this;", -1}, "but", 0},
		{&NamedReferenceScanner{"&but; &this;", 4}, "this", 6},
	}

	for _, c := range cases {
		actRef, actIdx := c.Scanner.Next()
		// fmt.Println(c.Scanner)
		if actRef != c.ExpRef || actIdx != c.ExpIdx {
			t.Errorf("Expecting %#v.Next() to be %q, %d, but got %q, %d.",
				c.Scanner, c.ExpRef, c.ExpIdx, actRef, actIdx)
		}
	}
}
