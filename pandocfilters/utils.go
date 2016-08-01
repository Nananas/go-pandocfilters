package pandocfilters

import (
	"fmt"
	"log"
	"reflect"
)

func Contains(slice List, value string) bool {
	ok := false
	for _, v := range slice {
		if v == value {
			ok = true
		}
	}

	return ok
}

func Deepcopy(input Any) Any {
	switch t := input.(type) {
	case List:
		newlist := NewList()
		for _, v := range t {
			newlist = append(newlist, Deepcopy(v))
		}

		return newlist
	case Node:
		newnode := NewNode()
		for k, v := range t {
			newnode[k] = Deepcopy(v)
		}

		return newnode
	default:
		return t

	}

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

func reduce(function func(c Any, n Action) Any, seq []Action, init Any) Any {

	current := init
	for _, next := range seq {
		current = function(current, next)
	}

	return current

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
