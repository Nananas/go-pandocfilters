package pandocfilters

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

// IS KEY ALWAYS A STRING
type Action func(k string, v Any, fmt string, meta string) Any
type Any interface{}
type List []Any
type Node map[string]Any

func ToJSONFilter(action Action) {
	ToJSONFilters([]Action{action})
}

func ToJSONFilters(actions []Action) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// fmt.Println("START")
	bytes, err := ioutil.ReadAll(os.Stdin)

	// fmt.Println(string(bytes))

	var j interface{}
	if err := json.Unmarshal(bytes, &j); err != nil {
		log.Println(err)
	}

	format := ""
	if len(os.Args) > 1 {
		format = os.Args[1]
	}

	c := convert(j)

	r := reduce(
		func(x Any, action Action) Any {
			return Walk(x, action, format, "")
		},
		actions, c)

	// fmt.Println(r)
	// fmt.Println("")
	encoded, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	}

	os.Stdout.Write(encoded)
}

func convert(input interface{}) Any {
	switch t := input.(type) {
	case []interface{}:
		list := make(List, 0)
		for _, e := range t {
			list = append(list, convert(e))
		}

		return list
	case map[string]interface{}:
		Node := make(Node, 0)
		tt, ok := t["t"]
		if !ok {
			for k, v := range t {
				Node[k] = v
			}
		} else {
			Node["t"] = tt
			Node["c"] = convert(t["c"])
		}
		return Node
	case string: // end of tree, leaf
		return string(input.(string))
	default:
		log.Fatal("this should not happen")
	}
	return nil
}

func Walk(x Any, action Action, format string, meta string) Any {
	switch v := x.(type) {
	case List:
		array := NewList()

		for _, item := range v {
			switch vv := item.(type) {
			case Node:
				if _, ex := vv["t"]; ex {
					result := action(vv["t"].(string), vv["c"], format, meta)

					if result == nil {
						array = append(array, Walk(vv, action, format, meta))
					} else {
						switch rt := result.(type) {
						case List:
							for _, z := range rt {
								array = append(array, Walk(z, action, format, meta))
							}
						default:
							array = append(array, Walk(rt, action, format, meta))
						}
					}
				} else {
					array = append(array, Walk(vv, action, format, meta))
				}
			default:
				array = append(array, Walk(vv, action, format, meta))
			}
		}

		return array

	// Node
	case Node:
		// fmt.Println("Node")
		// obj := make(map[string]interface{})
		obj := NewNode()

		for k, _ := range v {
			obj[k] = Walk(v[k], action, format, meta)
		}
		return obj
	default:
		// fmt.Println("Other?: ", x)
		return x

	}
	return x
}

func NewNode() Node {
	return make(Node, 0)
}

func NewList(args ...Any) List {
	l := make(List, 0)
	for _, e := range args {
		l = append(l, e)
	}

	return l
}
func reduce(function func(c Any, n Action) Any, seq []Action, init Any) Any {

	current := init
	for _, next := range seq {
		current = function(current, next)
	}

	return current

}

func Empty() List {
	return NewList()
}

func elt(eltType string, numargs int) EltFunc {
	return func(args ...Any) Node {
		var xs Any

		lenargs := len(args)
		if lenargs != numargs {
			log.Fatal("Error numargs != lenargs: ", numargs, lenargs)
		}
		if numargs == 0 {
			xs = Empty()
		} else if lenargs == 1 {
			xs = args[0]
		} else {
			xs = args
		}

		result := Node{
			"t": eltType,
			"c": xs,
		}

		return result

	}
}

type EltFunc func(...Any) Node

var (
	Plain          EltFunc
	Para           EltFunc
	CodeBlock      EltFunc
	RawBlock       EltFunc
	BlockQuote     EltFunc
	OrderedList    EltFunc
	BulletList     EltFunc
	DefinitionList EltFunc
	Header         EltFunc
	HorizontalRule EltFunc
	Table          EltFunc
	Div            EltFunc
	Null           EltFunc

	Str         EltFunc
	Emph        EltFunc
	Strong      EltFunc
	Strikeout   EltFunc
	Superscript EltFunc
	Subscript   EltFunc
	SmallCaps   EltFunc
	Quoted      EltFunc
	Cite        EltFunc
	Code        EltFunc
	Space       EltFunc
	LineBreak   EltFunc
	Math        EltFunc
	RawInline   EltFunc
	Link        EltFunc
	Image       EltFunc
	Note        EltFunc
	SoftBreak   EltFunc
	Span        EltFunc
)

func init() {
	Plain = elt("Plain", 1)
	Para = elt("Para", 1)
	CodeBlock = elt("CodeBlock", 2)
	RawBlock = elt("RawBlock", 2)
	BlockQuote = elt("BlockQuote", 1)
	OrderedList = elt("OrderedList", 2)
	BulletList = elt("BulletList", 1)
	DefinitionList = elt("DefinitionList", 1)
	Header = elt("Header", 3)
	HorizontalRule = elt("HorizontalRule", 0)
	Table = elt("Table", 5)
	Div = elt("Div", 2)
	Null = elt("Null", 0)

	Str = elt("Str", 1)
	Emph = elt("Emph", 1)
	Strong = elt("Strong", 1)
	Strikeout = elt("Strikeout", 1)
	Superscript = elt("Superscript", 1)
	Subscript = elt("Subscript", 1)
	SmallCaps = elt("SmallCaps", 1)
	Quoted = elt("Quoted", 2)
	Cite = elt("Cite", 2)
	Code = elt("Code", 2)
	Space = elt("Space", 0)
	LineBreak = elt("LineBreak", 0)
	Math = elt("Math", 2)
	RawInline = elt("RawInline", 2)
	Link = elt("Link", 3)
	Image = elt("Image", 3)
	Note = elt("Note", 1)
	SoftBreak = elt("SoftBreak", 0)
	Span = elt("Span", 2)

}

func AsString(input Any) string {
	switch t := input.(type) {
	case string:
		return t
	default:
		log.Fatal("Trying to convert to String, but not a string type (", reflect.TypeOf(input), ")")
	}

	return ""

}

func AsNode(input Any) Node {
	switch t := input.(type) {
	case Node:
		return t
	default:
		log.Fatal("Trying to convert to Node, but not a node type (", reflect.TypeOf(input), ")")
	}

	return nil

}

func AsList(input Any) List {
	switch t := input.(type) {
	case List:
		return t
	default:
		log.Fatal("Trying to convert to List, but not a list type (", reflect.TypeOf(input), ")")
	}

	return nil

}
