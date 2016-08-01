package pandocfilters

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

	AST, format, meta := Load(FromStdin())

	resultAST := Reduce(actions, AST, format, meta)
	ToStdout(Save(resultAST))
}

func Reduce(actions []Action, AST List, format string, meta Node) Any {
	return reduce(
		func(x Any, action Action) Any {
			return Walk(x, action, format, meta)
			// return Walk(x, action, format, NewNode())
		},
		actions, AST)

}

func FromStdin() []byte {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func ToStdout(input []byte) {
	os.Stdout.Write(input)
}

// Returns the AST, the format and the meta from the []byte input
func Load(bytes []byte) (List, string, Node) {

	var j interface{}
	if err := json.Unmarshal(bytes, &j); err != nil {
		log.Println(err)
	}

	format := ""
	if len(os.Args) > 1 {
		format = os.Args[1]
	}

	AST := convert(j)

	meta := AsNode(AsNode(AsList(AST)[0])["unMeta"])

	return AsList(AST), format, meta
}

func Save(AST Any) []byte {

	encoded, err := json.Marshal(AST)
	if err != nil {
		log.Println(err)
	}

	return encoded

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

func Empty() List {
	return NewList()
}

// Returns an attribute list, constructed from the
// dictionary attrs.
func AttributeList(attributes Node) List {
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

// Walks the tree x and returns concatenated string content,
// leaving out all formatting.
func Stringify(x Any) string {
	result := []string{}

	g := func(key string, val Any, format string, meta Node) Any {
		log.Println(key)
		switch key {
		case "Str", "MetaString":
			result = append(result, AsString(val))
		case "Code", "Math":
			result = append(result, AsString(AsList(val)[1]))
		case "LineBreak":
			result = append(result, " ")
		case "Space":
			result = append(result, " ")
		}

		return nil
	}

	log.Println(x)
	Walk(x, g, "", NewNode())

	log.Println(result)

	return strings.Join(result, "")
}

// Retrieves the metadata variable 'name' from the 'meta' dict.
func GetMeta(meta Node, name string) Any {

	// fmt.Println("META:", meta)
	// fmt.Println("NAME:", name)
	if d, ok := meta[name]; ok {

		data := AsNode(d)

		if data["t"] == "MetaString" || data["t"] == "MetaBool" {
			return AsString(data["c"])
		} else if data["t"] == "MetaInlines" {
			if len(AsList(data["c"])) == 1 {
				return Stringify(data["c"])
			}
		} else if data["t"] == "MetaList" {
			l := NewList()
			for _, v := range AsList(data["c"]) {
				l = append(l, Stringify(AsNode(v)["c"]))
			}

			return l
		}

	}

	log.Fatal("Could not understand metadata variable ", name)
	return nil
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
