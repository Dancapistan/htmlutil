package checker

import (
	"strings"
)

// IsValidHtml4IdValue returns true if the argument is a valid HTML 4 ID
// attribute value.
//
// Note: HTML 4 is more strict than HTML 5.
//
// Note: HTML 4 ID values are strict enough that a valid ID value will be valid
// also for a quoted or unquoted attribute value.
//
// Note: This function only checks the syntax. Additional restrictions on ID
// values, such as global uniqueness, also apply.
//
// From http://www.w3.org/TR/html401/types.html#type-name
//
//     ID and NAME tokens must begin with a letter ([A-Za-z]) and may be
//     followed by any number of letters, digits ([0-9]), hyphens ("-"),
//     underscores ("_"), colons (":"), and periods (".").
//
func IsValidHtml4IdValue(val string) bool {

	// must be one or more characters long

	length := len(val)
	if length == 0 {
		return false
	}

	// "must begin with a letter"

	first := val[0]
	isLower := first >= 'a' && first <= 'z'
	isUpper := first >= 'A' && first <= 'Z'
	if !isLower && !isUpper {
		return false
	}

	// "may be followed by any number of ..."

	for _, cur := range val {

		if !(cur >= 'a' && cur <= 'z') && // lowercase letter?
			!(cur >= 'A' && cur <= 'Z') && // uppercase letter?
			!(cur >= '0' && cur <= '9') && // digit?
			cur != ':' && // colon?
			cur != '_' && // underscore?
			cur != '-' && // hyphen?
			cur != '.' { // period?

			return false
		}
	}

	return true
}

// IsValidHtml5IdValue returns true if the argument is a valid HTML 5 ID value.
//
// Note: This is much more permissive than HTML 4 ID values and CSS3 ID values,
// so don't get too wacky.
//
// Note: Additional restrictions of HTML attribute values apply.
//
// From http://dev.w3.org/html5/markup/global-attributes.html#common.attrs.id
//
//     Any string, with the following restrictions:
//      - must be at least one character long
//      - must not contain any space characters
//
func IsValidHtml5IdValue(val string) bool {

	// "must be at least one character long"

	var invalid bool

	invalid = len(val) == 0 ||
		strings.ContainsAny(val, SpaceCharacters)

	return !invalid
}

// IsValidCss3IdValue returns true if the argument is a valid CSS 3 ID value.
//
// From http://www.w3.org/TR/css3-selectors/#id-selectors
//
//     An ID selector contains a "number sign" (U+0023, #) immediately followed
//     by the ID value, which must be an CSS identifiers.
//
// See also IsValidCss3Identifier.
//
func IsValidCss3IdValue(val string) bool {
	if len(val) < 2 {
		return false
	}

	return val[0] == '\u0023' && IsValidCss3Identifier(val[1:])
}

// IsValidCss3Identifier returns true if the argument is a valid CSS3
// identifier. Identifiers include element names, and the selector part of ID
// and class names.
//
// Prohibited characters can be included using escaped values, and this function
// takes escaped values into account.
//
// From http://www.w3.org/TR/CSS21/syndata.html#value-def-identifier
//
//     In CSS, identifiers ... can contain only:
//
//     - the characters [a-zA-Z0-9] and
//     - ISO 10646 characters U+00A0 and higher, plus
//     - the hyphen (-) and
//     - the underscore (_);
//     - they cannot start with a digit, two hyphens, or a hyphen followed by a
//       digit.
//
// Also, special characters can be escaped:
//
//     Any character (except a hexadecimal digit, linefeed, carriage return, or
//     form feed) can be escaped with a backslash to remove its special meaning.
//
// or
//
//     Third, backslash escapes allow authors to refer to characters they cannot
//     easily put in a document. In this case, the backslash is followed by at
//     most six hexadecimal digits (0..9A..F), which stand for the ISO 10646
//     ([ISO10646]) character with that number, which must not be zero.
//
//
// BUG(dr): IsValidCss3Identifier uses the CSS 2.1 spec, which the CSS 3 spec
// links to when it refers "identifiers". Is this the most up-to-date?
//
func IsValidCss3Identifier(val string) bool {

	if len(val) == 0 {
		return false
	}

	var first, second rune
	var wasSlash, inEscape bool
	var hexCount int

	for i, char := range val {
		if i == 0 {
			first = char
			// "they cannot start with a digit" TODO Alternative Unicode digits?
			if first >= '0' && first <= '9' {
				return false
			}
		}

		if i == 1 {
			second = char

			// "they cannot start with ... two hyphens, or a hyphen followed by
			// a digit."
			if first == '-' {
				if second == '-' || (second >= '0' && second <= '9') {
					return false
				}
			}
		}

		if hexCount > 6 {
			inEscape = false
		}

		// "can contain ... ISO 10646 characters U+00A0 and higher"

		if char >= '\u00A0' {
			inEscape = false
			wasSlash = false
			continue
		}

		// "can contain ... the hyphen and the underscore"

		if char == '-' || char == '_' {
			inEscape = false
			wasSlash = false
			continue
		}

		// "can contain ... the characters [a-zA-Z0-9]"

		if char >= 'a' && char <= 'z' {
			wasSlash = false
			if char > 'f' {
				inEscape = false
			} else {
				hexCount++
			}
			continue
		}

		// "can contain ... the characters [a-zA-Z0-9]"

		if char >= 'A' && char <= 'Z' {
			wasSlash = false
			if char > 'F' {
				inEscape = false
			} else {
				hexCount++
			}
			continue
		}

		// "can contain ... the characters [a-zA-Z0-9]"

		if char >= '0' && char <= '9' {
			wasSlash = false
			hexCount++
			continue
		}

		// "backslash escapes allow authors to refer to characters they cannot
		// easily put in a document"

		if char == '\\' {
			hexCount = 0
			wasSlash = true
			inEscape = true
			continue
		}

		// "the backslash is followed by at most six hexadecimal digits"
		//
		// "If a character in the range [0-9a-fA-F] follows the hexadecimal
		// number, the end of the number needs to be made clear ... with a space
		// (or other white space character) ... [or] by providing exactly 6
		// hexadecimal digits".

		if inEscape && (char == '\u0020' || char == '\u0009' || char == '\u000A' || char == '\u000D' || char == '\u000C') {
			inEscape = false
			wasSlash = false
			continue
		}

		if wasSlash {
			wasSlash = false
			inEscape = false
			continue
		}

		return false
	}

	return true
}
