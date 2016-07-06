/*

Pandoc filter that causes emphasis to be rendered using
the custom macro '\myemph{...}' rather than '\emph{...}'
in latex.  Other output formats are unaffected.

*/

package main

import . "github.com/nananas/Pandocfilters/pandocfilters"

func main() {
	ToJSONFilter(myemph)
}

func myemph(k string, v Any, f string, meta string) Any {
	// fmt.Println(f)
	// fmt.Println(v)
	// return nil
	if k == "Emph" && f == "latex" {
		// return append(NewList(latex("\\myemph{")), AsList(v)...)
		l := NewList(latex("\\myemph{"))
		l = append(l, AsList(v)...)
		return append(l, latex("}"))

		// latex("}"))
	}

	return nil
}

func latex(s string) Node {
	return RawInline("latex", s)
}
