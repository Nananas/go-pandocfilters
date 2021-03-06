/*

Pandoc filter that causes emphasis to be rendered using
the custom macro '\myemph{...}' rather than '\emph{...}'
in latex.  Other output formats are unaffected.

*/

package main

import . "github.com/nananas/go-pandocfilters/pandocfilters"

func main() {
	ToJSONFilter(myemph)
}

func myemph(k string, v Any, f string, meta Node) Any {
	if k == "Emph" && f == "latex" {
		return NewListUnpack(
			latex("\\myemph{"),
			AsList(v),
			latex("}"),
		)
		// l := NewList(latex("\\myemph{"))
		// l = append(l, AsList(v)...)
		// return append(l, latex("}"))
	}

	return nil
}

func latex(s string) Node {
	return RawInline("latex", s)
}
