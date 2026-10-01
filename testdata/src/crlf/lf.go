package crlf

// lf has a 13-line raw string, so 15 code lines allow ceil(15 × 0.15) = 3. // want "func doc has 4 comment lines, allowed 3 \\(funcs.max-lines\\)"
// Two.
// Three.
// Four.
func lf() string {
	return `1
2
3
4
5
6
7
8
9
10
11
12
13`
}
