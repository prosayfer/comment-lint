package formula

// tinyAtFloor has 1 lines.
func tinyAtFloor() int {
	n := 0
	return n
}

// tinyOverFloor has 2 lines. // want "func doc has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
// More.
func tinyOverFloor() int {
	n := 0
	return n
}

// sixAtFloor has 1 lines.
func sixAtFloor() int {
	n := 0
	n++
	n++
	return n
}

// sixOverFloor has 2 lines. // want "func doc has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
// More.
func sixOverFloor() int {
	n := 0
	n++
	n++
	return n
}

// sevenAtRatio has 2 lines.
// More.
func sevenAtRatio() int {
	n := 0
	n++
	n++
	n++
	return n
}

// sevenOverRatio has 3 lines. // want "func doc has 3 comment lines, allowed 2 \\(funcs.ratio\\)"
// More.
// More.
func sevenOverRatio() int {
	n := 0
	n++
	n++
	n++
	return n
}

// thirteenAtRatio has 2 lines.
// More.
func thirteenAtRatio() int {
	n := 0
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	return n
}

// thirteenOverRatio has 3 lines. // want "func doc has 3 comment lines, allowed 2 \\(funcs.ratio\\)"
// More.
// More.
func thirteenOverRatio() int {
	n := 0
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	return n
}

// twentyAtCap has 3 lines.
// More.
// More.
func twentyAtCap() int {
	n := 0
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	return n
}

// twentyOverCap has 4 lines. // want "func doc has 4 comment lines, allowed 3 \\(funcs.max-lines\\)"
// More.
// More.
// More.
func twentyOverCap() int {
	n := 0
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	return n
}

// longOverCap has 4 lines. // want "func doc has 4 comment lines, allowed 3 \\(funcs.max-lines\\)"
// More.
// More.
// More.
func longOverCap() int {
	n := 0
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	n++
	return n
}

// paddedSix has 2 lines. // want "func doc has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
// More.
func paddedSix() int {
	n := 0
	n++
	n++

	// note

	/* another */

	return n
}
