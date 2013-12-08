[![Build Status](https://drone.io/github.com/Dancapistan/htmlutil/status.png)](https://drone.io/github.com/Dancapistan/htmlutil/latest)

htmlutil
========

Some HTML utilities for Go.

This repo provides two packages: `checker` and `escaper`. These are works in progress so far, and I intend to add more functions to each package as needed. Pull requests welcome.

checker
=======

    import "github.com/Dancapistan/htmlutil/checker"

[GoDoc for package checker](http://godoc.org/github.com/Dancapistan/htmlutil/checker)

This package provides a number of utility functions to verify that a given input conforms to the HTML 5 spec for a given role.

For example, the function `checker.IsValidAttributeName` returns true if the argument can be used as an HTML tag's attribute name, like "type" in `<input type=email>`:

    name := "type"
    ok := checker.IsValidAttributeName(name)
    // ok == true

Please see the package's GoDoc page for a complete list of functions and fairly comprehensive examples.

escaper
=======

    import "github.com/Dancapistan/htmlutil/escaper"

[GoDoc for package escaper](http://godoc.org/github.com/Dancapistan/htmlutil/escaper)

This package is more work-in-progress. It provides some functions for escaping arguments to be safe in particular HTML contexts.

For example, `escaper.EscapeAttributeValueDoubleQuoted` makes sure the argument is safe to use as an HTML attribute value when double quoted. In this case, if the argument contains a double quote character, that character is converted into a safe HTML entity.



[GoDoc](http://godoc.org/github.com/Dancapistan/htmlutil)
