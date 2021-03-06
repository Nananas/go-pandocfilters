
pandocfilters
=============

A Go module for writing [pandoc](http://pandoc.org/) filters

What are pandoc filters?
--------------------------
Pandoc filters
are pipes that read a JSON serialization of the Pandoc AST
from stdin, transform it in some way, and write it to stdout.
They can be used with pandoc (>= 1.12) either using pipes ::

    pandoc -t json -s | ./caps | pandoc -f json

or using the ``--filter`` (or ``-F``) command-line option. ::

    pandoc --filter ./caps -s

For more on pandoc filters, see the pandoc documentation under ``--filter``
and [the tutorial on writing filters](http://johnmacfarlane.net/pandoc/scripting.html).

<!-- Compatibility
----------------
Pandoc 1.16 introduced link and image `attributes` to the existing
`caption` and `target` arguments, requiring a change in pandocfilters
that breaks backwards compatibility. Consequently, you should use:

- pandocfilters version <= 1.2.4 for pandoc versions 1.12--1.15, and
- pandocfilters version >= 1.3.0 for pandoc versions >= 1.16.
 -->

Installing
--------------

Using `git clone`? I still need to figure out `go get`.

Pandocfilters Types
-------------------
In contrast to python, go is a statically typed language. As such, types are needed in every step. Parsing json to go types results in quite ugly `interface{}, []interface{} and map[string]interface{}` types.

The following types are introduced for readability and are a straight conversion from the un-marshalled json types.

The ``pandocfilters`` package uses the following types:

**``type Any interface{}``:**
  For readability and less writing. All objects without clear type are of `Any`type.

**``type List []Any``**:
  A standard array of any types, can be mixed types.

**``type Node map[string]Any``**:
  A map type.

### Type casting
Casting an `Any` type to a known type can be done with the following helper functions:

**``func AsString(a Any) string``**

**``func AsNode(a Any) Node``**

**``func AsList(a Any) List``**

These functions panic if a wrong conversion occurs. See the examples for more use cases.

### Type creation
Creating a new type can be done more easily using the following helper functions:

**``func NewNode() Node``**

**``func NewList(args ...Any) Node``**: This will create a new list using the arguments specified

**``func NewListUnpack(args ...Any) Node``**: Same as the above, except that it will unpack any list in args. This is sometimes useful when wanting to encapsulate an existing list with Any objects.

**``func Empty() List``**: Helper function to create a new List. Is the same as NewList() but more readable, as returning an empty List will delete the object. See `ToJSONFilters()` below.

Available functions
----------------------

The ``pandocfilters`` package exports the following functions:

**``Walk(x Any, action Action, format string, meta Node)``**:
  Walk a tree, applying an action to every object.
  Returns a modified tree.

**``ToJSONFilter(action Action)``**, 
**``ToJSONFilters(actions []Action)``**:
  Converts a list of actions into a filter that reads a JSON-formatted
  pandoc document from stdin, transforms it by walking the tree
  with the actions, and returns a new JSON-formatted pandoc document
  to stdout.
  The argument is a list of functions,
  If the function returns None, the object to which it applies
  will remain unchanged.  If it returns an object, the object will
  be replaced.    If it returns a list, the list will be spliced in to
  the list to which the target object belongs.    (So, returning an
  empty list deletes the object.)

**``Attributes(attrs Node)``**:
  Returns an attribute list, constructed from the
  Node ``attrs``.

How to use
--------------
Most users will only need ``ToJSONFilter``.  Here is a simple example
of its use (from `examples/caps.go`)::

    /*
    Pandoc filter to convert all regular text to uppercase.
    Code, link URLs, etc. are not affected.
    */

    package main
    import (
        "strings"
        . "github.com/nananas/Pandocfilters/pandocfilters"
    )

    func main() {
        ToJSONFilter(caps)
    }

    func caps(k string, v Any, format string, meta Node) Any {
        if k == "Str" {
            return Str(strings.ToUpper(AsString(v)))
        }
        return nil
    }


Examples
--------

The examples subdirectory in the source repository contains the
following filters. These filters should provide a useful starting point
for developing your own pandocfilters.

**``caps.py``**:
    Pandoc filter to convert all regular text to uppercase. Code, link
    URLs, etc. are not affected.

**``comments.py``**:
    Pandoc filter that causes everything between
    ``<!-- BEGIN COMMENT -->`` and ``<!-- END COMMENT -->`` to be ignored.
    The comment lines must appear on lines by themselves, with blank
    lines surrounding

**``deemph.py``**:
    Pandoc filter that causes emphasized text to be displayed in ALL
    CAPS.

**``deflists.py``**:
    Pandoc filter to convert definition lists to bullet lists with the
    defined terms in strong emphasis (for compatibility with standard
    markdown).

**``metavars.py``**:
    Pandoc filter to allow interpolation of metadata fields into a
    document. ``%{fields}`` will be replaced by the field's value, assuming
    it is of the type ``MetaInlines`` or ``MetaString``.

**``myemph.py``**:
    Pandoc filter that causes emphasis to be rendered using the custom
    macro ``\myemph{...}`` rather than ``\emph{...}`` in latex. Other output
    formats are unaffected.

**``theorem.py``**:
    Pandoc filter to convert divs with ``class="theorem"`` to LaTeX theorem
    environments in LaTeX output, and to numbered theorems in HTML
    output.

**TODO** from the python pandocfilters tool:  
``abc.py, gabc.py, graphviz.py, lilypond.py, platuml.py, tikz.py``