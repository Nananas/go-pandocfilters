package pandocfilters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

// IS KEY ALWAYS A STRING
// Action function type used as Filter
//
// Key is the type of the pandoc object as string(e.g. 'Str', 'Para'),
// Value is the contents of the object as <Any>(e.g. a string for 'Str',
// a list of inline elements for 'Para', typed as <Any>),
// Format is the target output format as string (which will be taken
// for the first command line argument if present),
// and Meta is the document's metadata as <Node>.
type Action func(key string, value Any, format string, meta Node) Any
type Any interface{}
type List []Any
type Node map[string]Any

func ToJSONFilter(action Action) {
	ToJSONFilters([]Action{action})
}

// Converts a list of actions <Action> into a filter that reads a JSON-formatted
// pandoc document from stdin, transforms it by walking the tree
// with the actions, and returns a new JSON-formatted pandoc document
// to stdout.
//
// The argument is a list of functions of type <Action>,
// If the function returns None, the object to which it applies
// will remain unchanged.  If it returns an object, the object will
// be replaced.    If it returns a list, the list will be spliced in to
// the list to which the target object belongs.    (So, returning an
// empty list deletes the object.)
//
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
	// fmt.Println(m)
	// fmt.Println(AsNode(m)["author"])
	meta := AsNode(AsNode(AsList(c)[0])["unMeta"])

	r := reduce(
		func(x Any, action Action) Any {
			return Walk(x, action, format, meta)
			// return Walk(x, action, format, NewNode())
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
	case List:
		fmt.Println("Already converted...")
	case []interface{}: // List
		list := make(List, 0)
		for _, e := range t {
			list = append(list, convert(e))
		}

		return list
	case Node:
		fmt.Println("Already converted...")
	case map[string]interface{}: // Node
		Node := make(Node, 0)
		tt, ok := t["t"]
		if !ok {
			for k, v := range t {
				Node[k] = convert(v)
			}
		} else {
			Node["t"] = tt
			Node["c"] = convert(t["c"])
		}
		return Node
	case string: // end of tree, leaf
		return string(input.(string))
	case float64:
		return float64(input.(float64))
	default:
		log.Println(input, ": is :", reflect.TypeOf(input))
		log.Fatal("[FATAL] This should not happen")
	}
	return nil
}

// Walk a tree, applying an action to every object.
// Returns a modified tree.
func Walk(x Any, action Action, format string, meta Node) Any {
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

func NewListUnpack(args ...Any) List {
	l := make(List, 0)
	for _, e := range args {
		switch t := e.(type) {
		case List:
			for _, ee := range t {
				l = append(l, ee)
			}
		default:
			l = append(l, t)
		}
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

// Returns an attribute list, constructed from the
// dictionary attrs.
func Attributes(attributes Node) List {
	ident, ok := attributes["id"]
	if !ok {
		ident = ""
	}

	classes, ok := attributes["classes"]
	if !ok {
		classes = []string{}
	}

	keyvals := NewList()
	for k, v := range attributes {
		if k != "classes" && k != "id" {
			l := NewList(k, v)
			keyvals = append(keyvals, l)
		}
	}

	result := NewList(ident, classes, keyvals)
	return result
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
		log.Println("Trying to convert to String, but not a string type (", reflect.TypeOf(input), ") for object: ", input)
		panic("")
	}

	return ""

}

func AsNode(input Any) Node {
	switch t := input.(type) {
	case Node:
		return t
	default:
		log.Println("Trying to convert to Node, but not a node type (", reflect.TypeOf(input), ") for object: ", input)
		panic("")
	}

	return nil

}

func AsList(input Any) List {
	switch t := input.(type) {
	case List:
		return t
	default:
		log.Println("Trying to convert to List, but not a list type (", reflect.TypeOf(input), ") for object: ", input)
		panic("")
	}

	return nil

}
