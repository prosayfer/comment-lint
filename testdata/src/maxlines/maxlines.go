package maxlines

// atLimit has three lines.
// Two.
// Three.
func atLimit() {}

// overLimit has four lines. // want "func doc has 4 comment lines, allowed 3 \\(funcs.max-lines\\)"
// Two.
// Three.
// Four.
func overLimit() {}

type T struct{}

// method has four lines. // want "func doc has 4 comment lines, allowed 3 \\(funcs.max-lines\\)"
// Two.
// Three.
// Four.
func (T) method() {}

// essay has six lines. // want "func doc has 6 comment lines, allowed 3 \\(funcs.max-lines\\)"
// Two.
// Three.
// Four.
// Five.
// Six.
func essay() {}
