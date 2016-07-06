/*

Pandoc filter to allow interpolation of metadata fields
into a document.  %{fields} will be replaced by the field's
value, assuming it is of the type MetaInlines or MetaString.

*/

package main

import (
	"log"
	"regexp"

	. "github.com/nananas/Pandocfilters/pandocfilters"
)

var pattern *regexp.Regexp

func main() {
	var err error
	pattern, err = regexp.Compile("%\\{(.*)\\}$")
	if err != nil {
		log.Fatal("Error in compiling regex: ", err)
	}

	ToJSONFilter(metavars)
}

func metavars(key string, value Any, format string, meta Node) Any {
	if key == "Str" {
		m := pattern.FindStringSubmatch(AsString(value))

		if len(m) > 0 {
			field := m[1]
			result := AsNode(meta[field])

			// fmt.Println(field)

			metatype, ok := result["t"]
			if ok {
				if metatype == "MetaInlines" {
					return Span(Attributes(Node{
						"class": "interpolated",
						"field": field,
					}), result["c"])
				} else if metatype == "MetaString" {
					return Str(result["c"])
				}
			}
		}

	}

	return nil
}
