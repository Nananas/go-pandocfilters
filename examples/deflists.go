/*

Pandoc filter to convert definition lists to bullet
lists with the defined terms in strong emphasis (for
compatibility with standard markdown).

*/

package main

import . "github.com/nananas/Pandocfilters/pandocfilters"

func main() {
	ToJSONFilter(deflists)
}

func deflists(key string, val Any, format string, meta string) Any {
	if key == "DefinitionList" {
		array := Empty()
		value := val.(List)

		for _, v := range value {
			va := v.(List)
			array = append(array, tobullet(va[0], va[1]))
		}
		return BulletList(array)
	}
	return nil
}

func tobullet(term Any, def Any) Any {

	para := Empty()
	para = append(para, Strong(term))
	paraBlock := Empty()
	paraBlock = append(paraBlock, Para(para))

	defs := def.(List)
	for _, d := range defs {
		dd := d.(List)
		for _, b := range dd {
			paraBlock = append(paraBlock, b)
		}
	}
	return paraBlock
}
