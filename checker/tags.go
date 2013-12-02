package checker

// IsValidHTMLTagName returns true if the argument *can be* a valid HTML 5 tag
// name.
//
// 		Tags contain a tag name, giving the element's name. HTML elements all
// 		have names that only use alphanumeric ASCII characters. In the HTML
// 		syntax, tag names, even those for foreign elements, may be written with
// 		any mix of lower- and uppercase letters that, when converted to all-
// 		lowercase, matches the element's tag name; tag names are case-
// 		insensitive.
//
// This function checks the structural syntax of the argument (i.e. alphanumeric
// ASCII characters). It does not check if the argument is a pre-defined HTML5
// tag name. Use IsHTMLTagName to see if it is a pre-defined name.
//
func IsValidHTMLTagName(name string) bool {

	for _, c := range name {

		var isNum bool
		var isLower bool
		var isUpper bool

		isNum = c >= '0' && c <= '9'
		isLower = c >= 'a' && c <= 'z'
		isUpper = c >= 'A' && c <= 'Z'
		if !(isNum || isLower || isUpper) {
			return false
		}
	}
	return true
}
