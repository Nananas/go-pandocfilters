package examples

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"testing"

	pf "github.com/nananas/Pandocfilters/pandocfilters"
)

var incomment = false

func TestComments(t *testing.T) {
	binary, lookerr := exec.LookPath("pandoc")
	if lookerr != nil {
		t.Skip("No pandoc binary found in PATH")
	}

	pandoc := exec.Command(binary, "-t", "json", "comments_sample.md")

	output, err := pandoc.Output()
	if err != nil {
		t.Skip(err)
	}
	// fmt.Println(string(output))
	// pandocfilters.ToJSONFilter(comment)

	var j interface{}

	if err := json.Unmarshal(output, &j); err != nil {
		t.Error(err)
	}

	format := "html"

	r := reduce(
		func(x interface{}, action pf.Action) interface{} {
			return pf.Walk(x, action, format, "")
		},
		[]pf.Action{comment}, j)

	fmt.Println("...")
	fmt.Println(j)

	fmt.Println("...")
	fmt.Println(r)
}

func comment(k interface{}, v interface{}, format string, meta string) interface{} {

	if k == "RawBlock" {

		fmt.Println("			RAW BLOCK", format)
		s := v.([]interface{})[1].(string)

		if format == "html" {
			if m, err := regexp.MatchString("<!-- BEGIN COMMENT -->", s); err == nil && m {
				incomment = true
				return []string{}
			} else if m, err := regexp.MatchString("<!-- BEGIN COMMENT -->", s); err == nil && m {
				incomment = false
				return []string{}
			}
		}
	}

	if incomment {
		return []string{}
	}

	return nil
}

func reduce(function func(c interface{}, n pf.Action) interface{}, seq []pf.Action, init interface{}) interface{} {

	current := init
	for _, next := range seq {
		current = function(current, next)
	}

	return current

}
