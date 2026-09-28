package maxlinesoverride

// fiveLines is at the raised limit.
// Two.
// Three.
// Four.
// Five.
func fiveLines() {}

// sixLines is one over. // want "func doc has 6 comment lines, allowed 5 \\(funcs.max-lines\\)"
// Two.
// Three.
// Four.
// Five.
// Six.
func sixLines() {}
