/*

Pandoc filter that causes everything between
'<!-- BEGIN COMMENT -->' and '<!-- END COMMENT -->'
to be ignored.  The comment lines must appear on
lines by themselves, with blank lines surrounding
them.

*/

package main

import (
	"regexp"

	. "github.com/nananas/Pandocfilters/pandocfilters"
)

var incomment = false

func main() {
	ToJSONFilter(comment)
}

func comment(k string, v Any, format string, meta string) Any {

	if k == "RawBlock" {
		// format := v.([]interface{})[0].(string)
		// s := v.([]interface{})[1].(string)
		format := AsString(AsList(v)[0])
		s := AsString(AsList(v)[1])

		if format == "html" {
			if m, err := regexp.MatchString("<!-- BEGIN COMMENT -->", s); err == nil && m {
				incomment = true
				return Empty()
			} else if m, err := regexp.MatchString("<!-- END COMMENT -->", s); err == nil && m {
				incomment = false
				return Empty()
			}
		}
	}

	if incomment {
		return Empty()
	}

	return nil
}
