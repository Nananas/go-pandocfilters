/*

Pandoc filter to convert all regular text to uppercase.
Code, link URLs, etc. are not affected.

*/

package main

import (
	"strings"

	. "github.com/nananas/Pandocfilters/pandocfilters"
)

func main() {
	ToJSONFilter(caps)
}

func caps(k string, v Any, format string, meta string) Any {
	if k == "Str" {
		return Str(strings.ToUpper(AsString(v)))
	}

	return nil
}
