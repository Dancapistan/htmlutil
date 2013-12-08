package escaper

import (
	"bytes"
	"github.com/Dancapistan/htmlutil/checker"
	"strings"
)

// Common HTML entities. Stolen from template/funcs.go.
const (
	htmlQuot = "&#34;"
	htmlAmp  = "&amp;"
)

const (
	unicodeDoubleQuote = '\u0022'
	unicodeAmpersand   = '\u0026'
	unicodeSemicolon   = '\u003B'
)

const doubleQuoteStr = string(unicodeDoubleQuote)

var doubleQuoteByte = []byte(doubleQuoteStr)
var htmlQuotByte = []byte(htmlQuot)

// EscapeAttributeValueDoubleQuoted returns the argument with double quotes
// escaped and with ambiguous ampersands escaped.
//
func EscapeAttributeValueDoubleQuoted(val string) string {

	idxAmp := strings.IndexRune(val, unicodeAmpersand)
	idxQuo := strings.IndexRune(val, unicodeDoubleQuote)

	// Heuristic: If the argument doesn't contain a quote, or an ampersand, then
	// it is most likely fine unescaped.

	if idxAmp == -1 && idxQuo == -1 {
		return val
	}

	var b []byte

	// If we have an ampersand, it *may* contain an ambiguous ampersand that
	// needs escaping. An ambiguous ampersand is of the form &xxx; (where xxx is
	// one or more alphanumeric characters). Therefore, we also need a semicolon
	// to appear after the ampersand for it to be a possible ambiguous
	// ampersand.

	idxSemi := strings.IndexRune(val, unicodeSemicolon)
	if idxAmp != -1 && idxAmp < idxSemi {
		b = escapeAmbiguousAmpersandsBuffer(val)
	}

	// Escape double quotes characters.

	if idxQuo != -1 {

		if b == nil {
			return strings.Replace(val, doubleQuoteStr, htmlQuot, -1)
		} else {
			b := bytes.Replace(b, doubleQuoteByte, htmlQuotByte, -1)
			return string(b)
		}

	} else {

		if b == nil {
			return val
		} else {
			return string(b)
		}
	}
}

// EscapeAmbiguousAmpersands returns a copy of the argument with ambiguous
// ampersands escaped with &amp;.
//
func EscapeAmbiguousAmpersands(val string) string {

	length := len(val)
	if length < 3 {
		return val
	}

	ampIdx := strings.IndexRune(val, unicodeAmpersand)

	// No ampersands? Nothing to do.
	if ampIdx == -1 || ampIdx >= length-2 {
		return val
	}

	b := escapeAmbiguousAmpersandsBuffer(val)
	if b != nil {
		return string(b)
	}

	return val
}

func escapeAmbiguousAmpersandsBuffer(val string) []byte {

	var scanner = checker.NamedReferenceScanner{val, -1}

	// Count how many ambiguous ampersands are actually in the string.

	var count int
	// var indexes [5]int // cache first 5 ambiguous ampersand indexes
	for {
		name, index := scanner.Next()
		if index == -1 {
			break
		}
		if !checker.IsCharacterReferenceName(name) {
			// if count < len(indexes) {
			// indexes[count] = index
			// }
			count++
		}
	}
	scanner.Reset()

	// If there are no ambiguous ampersands, there is nothing to be escaped.

	if count == 0 {
		return nil
	}

	// We have at least one ambiguous ampersand to be escaped. Calculate the
	// final width.
	//
	finalLength := len(val) - count + (count * len(htmlAmp))

	// Create the bytes buffer that we'll copy the final escaped string into.

	var b = make([]byte, finalLength)
	var dest int // Current write location relative to b.
	var src int  // Current read location relative to val.

	// We cached the first len(indexes) of ambiguous ampersands. Escape those
	// first.

	// var numCached = count
	// if count >= len(indexes) {
	// 	numCached = len(indexes)
	// }

	// for i := 0; i < numCached; i++ {
	// 	ampIndex := indexes[i]
	// 	dest += copy(b[dest:], val[src:ampIndex])
	// 	dest += copy(b[dest:], htmlAmp)
	// 	src = ampIndex + 1
	// }

	// if numCached == count {
	// 	copy(b[dest:], val[src:])
	// 	return b
	// }

	// If we get here, then the input has even more ambiguous ampersands than we
	// cached, so we need to re-scan for those and escape them, too.

	// scanner.LastIndex = indexes[len(indexes)-1] + 1
	for {
		name, index := scanner.Next()

		// If we're past the last possible ambiguous ampersand, then copy in the
		// remaining data from `val`.

		if index == -1 {
			copy(b[dest:], val[src:])
			break
		}

		// If we're at an ambiguous ampersand (i.e. if name is not a valid
		// character reference), then copy in the data from `val` from where we
		// left off up to but not including the ampersand. Then copy in the
		// escaped version of the ampersand.

		if !checker.IsCharacterReferenceName(name) {
			dest += copy(b[dest:], val[src:index])
			dest += copy(b[dest:], htmlAmp)
			src = index + 1 // skip the ampersand.
		}
	}

	return b
}
