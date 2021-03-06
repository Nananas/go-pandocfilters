/*

Pandoc filter to convert divs with class="theorem" to LaTeX
theorem environments in LaTeX output, and to numbered theorems
in HTML output.

*/

package main

import (
	"strconv"

	. "github.com/nananas/go-pandocfilters/pandocfilters"
)

var theoremcount int

func main() {
	theoremcount = 0
	// Load()
	ToJSONFilter(theorems)
}

func theorems(key string, value Any, format string, meta Node) Any {
	if key == "Div" {
		attributes := AsList(AsList(value)[0])
		ident := AsString(attributes[0])
		classes := AsList(attributes[1])
		contents := AsList(AsList(value)[1])

		label := ""

		if Contains(classes, "theorem") {
			if format == "latex" {
				if ident != "" {
					label = "\\label{" + ident + "}"
				}

				return NewListUnpack(
					latex("\\begin{theorem}"+label),
					contents,
					latex("\\end{theorem}"),
				)
			} else if format == "html" || format == "html5" {
				theoremcount++

				newcontent := NewListUnpack(
					html("<dt>Theorem "+strconv.Itoa(theoremcount)+"</dt>"),
					html("<dd>"),
					contents,
					html("</dd>\n</dl>"),
				)

				return Div(attributes, newcontent)
			}
		}

		// if
	}

	return nil
}

func latex(x string) Node {
	return RawBlock("latex", x)
}

func html(x string) Node {
	return RawBlock("html", x)
}
