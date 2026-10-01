package counting

// paragraphs has three text lines.
//
// Two.
//
// Three.
func paragraphs() {}

// overWithDirectives has four text lines. // want "func doc has 4 comment lines, allowed 3 \\(funcs.max-lines\\)"
// Two.
// Three.
// Four.
//
//go:noinline
func overWithDirectives() {}

/*
blockAtLimit has three text lines.
Two.
Three.
*/
func blockAtLimit() {}

/*
blockOver has four text lines.
Two.
Three.
Four. // want "func doc has 4 comment lines, allowed 3 \\(funcs.max-lines\\)"
*/
func blockOver() {}

/**
 * starred has three text lines.
 *
 * Two.
 *
 * Three.
 */
func starred() {}

/**
 * starredOver has four text lines.
 *
 * Two.
 * Three.
 * Four. // want "func doc has 4 comment lines, allowed 3 \\(funcs.max-lines\\)"
 */
func starredOver() {}

/* single is one line. */
// Two.
// Three.
func single() {}

/* inlineOver starts on the marker line.
Two.
Three.
Four. // want "func doc has 4 comment lines, allowed 3 \\(funcs.max-lines\\)" */
func inlineOver() {}

// withDirectives is last because its //line directive remaps positions.
// Two.
// Three.
//
//go:generate echo hi
//nolint:gocritic
//export withDirectives
// +build ignore
//line counting.go:1
func withDirectives() {}
