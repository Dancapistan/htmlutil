// Package checker provides some utility functions for checking the validity of
// HTML5 tags, attribute name, and attribute values.
package checker

import (
	"strings"
	"unicode"
)

// Unicode characters
const (
	UnicodeAmpersand  = '\u0026'
	UnicodeSemicolon  = '\u003B'
	UnicodeQuoteMark  = '\u0022'
	UnicodeApostrophe = '\u0027'
)

// SpaceCharacters are space, tab, linefeed, formfeed, and carriagereturn.
//
const SpaceCharacters string = "\u0020\u0009\u000A\u000C\u000D"

// ControlCharacters are the non-SpaceCharacters ASCII control characters
// NUL, BEL, BS, VT, ESC, and DEL.
//
const ControlCharacters string = "\u0000\u0007\u0008\u000B\u001B\u007F"

// InvalidAttributeNameCharacters are the characters not valid in an attribute's
// name, minus the SpaceCharacters and the ControlCharacters.
//
const InvalidAttributeNameCharacters string = "\u0022\u0027\u003E\u002F\u003D"

// InvalidAttributeValueUnquotedCharacters are the characters not valid in an
// unquoted attribute's value.
//
const InvalidAttributeValueUnquotedCharacters string = "\u0022\u0027\u003C\u003D\u003E\u0060"

// IsValidAttributeName returns true if the argument is a HTML5-valid attribute
// name, as defined here:
// http://www.w3.org/TR/html5/syntax.html#attributes-0
//
//     Attribute names must consist of one or more characters other than the
//     space characters, U+0000 NULL, U+0022 QUOTATION MARK ("), U+0027
//     APOSTROPHE ('), ">" (U+003E), "/" (U+002F), and "=" (U+003D) characters,
//     the control characters, and any characters that are not defined by
//     Unicode.
//
// This is merely a syntax check. This function makes no judgments about the
// semantic validity of the argument.
//
func IsValidAttributeName(name string) bool {

	var invalid bool

	invalid = len(name) == 0 ||
		strings.ContainsAny(name, SpaceCharacters) ||
		strings.ContainsAny(name, ControlCharacters) ||
		strings.ContainsAny(name, InvalidAttributeNameCharacters) ||
		strings.IndexFunc(name, isUnicodeNonCharacter) > -1

	return !invalid
}

func isUnicodeNonCharacter(char rune) bool {
	var table *unicode.RangeTable
	table = unicode.Noncharacter_Code_Point
	return unicode.Is(table, char)
}

// IsValidAttributeValue returns true if the argument is a valid attribute
// value, with the caveat that additional rules apply to unquoted, single-
// quoted, and double-quoted attribute values.
//
//     Attribute values are a mixture of text and character references,
//     except with the additional restriction that the text cannot contain
//     an ambiguous ampersand.
//
// Definition: http://www.w3.org/TR/html5/syntax.html#attributes-0
//
// Note: It is almost certainly a bad idea to call this function directly. Use
// one of IsValidAttributeValueUnquoted, IsValidAttributeValueSingleQuoted, or
// IsValidAttributeValueDoubleQuoted instead.
//
func IsValidAttributeValue(val string) bool {
	return !HasAmbiguousAmpersand(val)
}

// IsValidAttributeValueUnquoted return true if the argument is a valid
// unquoted attribute value. For example, "email" in <input type=email>.
//
//     The attribute name, followed by zero or more space characters, followed
//     by a single U+003D EQUALS SIGN character, followed by zero or more space
//     characters, followed by the attribute value, which, in addition to the
//     requirements given above for attribute values, must not contain any
//     literal space characters, any U+0022 QUOTATION MARK characters ("),
//     U+0027 APOSTROPHE characters ('), "=" (U+003D) characters, "<" (U+003C)
//     characters, ">" (U+003E) characters, or "`" (U+0060) characters, and
//     must not be the empty string
//
// From http://www.w3.org/TR/html5/syntax.html#attributes-0
//
func IsValidAttributeValueUnquoted(val string) bool {

	// "must not be the empty string"
	if len(val) == 0 {
		return false
	}

	// "must not contain any literal space characters"
	if strings.ContainsAny(val, SpaceCharacters) {
		return false
	}

	if strings.ContainsAny(val, InvalidAttributeValueUnquotedCharacters) {
		return false
	}

	// "in addition to the requirements given above for attribute values"
	return IsValidAttributeValue(val)
}

// IsValidAttributeValueSingleQuoted return true if the argument is a valid
// single-quoted attribute value. Note, the argument must not contain the
// single quotes.
//
//     the attribute value, which, in addition to the requirements given above
//     for attribute values, must not contain any literal "'" (U+0027)
//     characters
//
// From http://www.w3.org/TR/html5/syntax.html#attributes-0
func IsValidAttributeValueSingleQuoted(val string) bool {

	// must not contain any literal "'"
	if strings.IndexRune(val, UnicodeApostrophe) > -1 {
		return false
	}

	// "in addition to the requirements given above for attribute values"
	return IsValidAttributeValue(val)
}

// IsValidAttributeValueDoubleQuoted return true if the argument is a valid
// double-quoted attribute value. Note, the argument must not contain the
// double quotes.
//
//     the attribute value, which, in addition to the requirements given above
//     for attribute values, must not contain any literal """ (U+0022)
//     characters
//
// From http://www.w3.org/TR/html5/syntax.html#attributes-0
func IsValidAttributeValueDoubleQuoted(val string) bool {

	// must not contain any literal "
	if strings.IndexRune(val, UnicodeQuoteMark) > -1 {
		return false
	}

	// "in addition to the requirements given above for attribute values"
	return IsValidAttributeValue(val)
}

// HasAmbiguousAmpersand returns true if the argument contains a substring that
// is an ambiguous ampersand.
//
//     An ambiguous ampersand is a U+0026 AMPERSAND character (&) that is
//     followed by one or more alphanumeric ASCII characters, followed by
//     a ";" (U+003B) character, where these characters do not match any
//     of the names given in the named character references section.
//
// It is ambiguous if it looks like a named character reference but is NOT one:
// "&ambiguous;" is ambiguous, but "&amp;" is not because "&amp;" is a valid
// reference. See also IsCharacterReferenceName and IsCharacterReference.
//
func HasAmbiguousAmpersand(val string) bool {

	scanner := NamedReferenceScanner{val, -1}

	for {
		name, idx := (&scanner).Next()
		if idx == -1 {
			return false
		}

		if !IsCharacterReferenceName(name) {
			return true
		}
	}
}

// INTERNAL USE ONLY. NO API GUARANTEES.
//
// NamedReferenceScanner is a utility for scanning through strings and looking
// for named character references.
//
// The Next method only checks for the structure: an ampersand, one or more
// alphanumeric values, and a semicolon. Use IsCharacterReferenceName to see if
// the returned values are actually valid character reference names.
//
type NamedReferenceScanner struct {
	Value     string // The string to be scanned.
	LastIndex int    // The stopping point where Next last finished, or -1 to start from the beginning.
}

// NewNamedReferenceScanner creates a new scanner with the given value.
//
func NewNamedReferenceScanner(value string) *NamedReferenceScanner {
	return &NamedReferenceScanner{value, -1}
}

// Next returns the next named character reference (just the alphanumeric part,
// skipping the ampersand and semicolon) and the byte index of the leading
// ampersand.
//
// Or empty string and -1 if no named character references could be found.
//
func (scanner *NamedReferenceScanner) Next() (name string, ampIndex int) {

	length := len(scanner.Value)
	first := scanner.LastIndex + 1

	ampIndex = -1

	// Loop through the characters until we find an ampersand.
	//
	// TODO: This loop assumes one-byte wide characters. Test with multi-byte
	// characters.

	for i := first; i < length; i++ {
		cur := scanner.Value[i]

		if cur == UnicodeAmpersand {

			// There must be at least one character between the ampersand and
			// the semicolon.
			if i+1 < length && scanner.Value[i+1] == UnicodeSemicolon {
				continue
			}

			ampIndex = i

			// Loop through all the characters after the ampersand. The
			// characters will all be alphanumeric until the semicolon, in which
			// case we found a character reference name to be returned. Or there
			// will be a non-alphanumeric value, in which case we break and
			// continue searching for more ampersands.

			for j := i + 1; j < length; j++ {
				cur2 := scanner.Value[j]

				if cur2 == UnicodeSemicolon {
					name = scanner.Value[i+1 : j]
					scanner.LastIndex = j
					return
				}

				isLower := cur2 >= 'a' && cur2 <= 'z'
				isUpper := cur2 >= 'A' && cur2 <= 'Z'
				isNumber := cur2 >= '0' && cur2 <= '9'

				if !(isLower || isUpper || isNumber) {
					break
				}
			}
		}
	}

	// Didn't find anything.

	scanner.LastIndex = length
	return "", -1
}

// Reset resets the scanner to the beginning of the Value string.
//
func (scanner *NamedReferenceScanner) Reset() {
	scanner.LastIndex = -1
}
