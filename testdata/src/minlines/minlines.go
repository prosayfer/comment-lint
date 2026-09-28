package minlines

// atRaisedFloor has 2 lines.
// More.
func atRaisedFloor() int {
	n := 0
	return n
}

// overRaisedFloor has 3 lines. // want "func doc has 3 comment lines, allowed 2 \\(funcs.min-lines\\)"
// More.
// More.
func overRaisedFloor() int {
	n := 0
	return n
}
