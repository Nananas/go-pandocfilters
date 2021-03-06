/*

Pandoc filter that causes emphasized text to be displayed
in ALL CAPS.

*/

package main

import (
	"strings"

	. "github.com/nananas/go-pandocfilters/pandocfilters"
)

func main() {
	ToJSONFilter(deemph)
}

func deemph(key string, val Any, format string, meta Node) Any {
	if key == "Emph" {
		return Walk(val, caps, format, meta)
	}
	return nil
}

// copied from caps.go
func caps(k string, v Any, format string, meta Node) Any {

	if k == "Str" {
		// fmt.Println(v)
		return Str(strings.ToUpper(AsString(v)))
	}
	return nil
}
