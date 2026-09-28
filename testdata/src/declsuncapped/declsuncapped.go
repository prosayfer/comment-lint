package declsuncapped

// AtRatio has four code lines, since the field comments don't count.
// Two.
// Three.
// Four.
type AtRatio struct {
	// A is documented.
	A int
	// B is documented.
	B int
}

// OverRatio has four code lines, since the field comments don't count. // want "decl doc has 5 comment lines, allowed 4 \\(decls.ratio\\)"
// Two.
// Three.
// Four.
// Five.
type OverRatio struct {
	// A is documented.
	A int
	// B is documented.
	B int
}
