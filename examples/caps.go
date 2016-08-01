/*

Pandoc filter to convert all regular text to uppercase.
Code, link URLs, etc. are not affected.

*/

package main

import (
	"strings"

	. "github.com/nananas/go-pandocfilters/pandocfilters"
)

func main() {
	ToJSONFilter(caps)
}

func caps(k string, v Any, format string, meta Node) Any {
	if k == "Str" {
		return Str(strings.ToUpper(AsString(v)))
	}

	return nil
}
