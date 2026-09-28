package decls

// one is a single-line var with a one-line doc.
var one = 1

// tooLong has three lines on a single-line var. // want "decl doc has 3 comment lines, allowed 1 \\(decls.min-lines\\)"
// Two.
// Three.
var tooLong = 1

// multi has a three-line value, so the ratio allows three and the cap binds. // want "decl doc has 3 comment lines, allowed 2 \\(decls.max-lines\\)"
// Two.
// Three.
var multi = []int{
	1,
}

// Limit is a single const. // want "decl doc has 2 comment lines, allowed 1 \\(decls.min-lines\\)"
// Two.
const Limit = 3

// Status is a named non-struct type. // want "decl doc has 2 comment lines, allowed 1 \\(decls.min-lines\\)"
// Two.
type Status int

// Alias is a type alias. // want "decl doc has 2 comment lines, allowed 1 \\(decls.min-lines\\)"
// Two.
type Alias = Status

// Point is measured against the whole struct.
// Two.
type Point struct {
	// X is fine.
	X int
	// Y is too long. // want "field comment has 2 comment lines, allowed 1 \\(decls.min-lines\\)"
	// Two.
	Y int
}

// Colors is measured against the whole group.
// Two.
const (
	// Red is fine.
	Red = iota
	// Green is too long for one spec. // want "decl doc has 2 comment lines, allowed 1 \\(decls.min-lines\\)"
	// Two.
	Green
)

// Sizes is over even the group's allowance. // want "decl doc has 3 comment lines, allowed 2 \\(decls.max-lines\\)"
// Two.
// Three.
var (
	small = 1
	large = 2
)

type (
	// Grouped is a spec in a type group. // want "decl doc has 2 comment lines, allowed 1 \\(decls.min-lines\\)"
	// Two.
	Grouped int
)

var anonymous struct {
	// field is a field of an anonymous struct. // want "field comment has 2 comment lines, allowed 1 \\(decls.min-lines\\)"
	// Two.
	field int
}
