package inbody

func statements(items []int) int {
	n := 0
	// One line above a single statement is fine.
	n++
	// Two lines above a single statement are too many. // want "in-body comment has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
	// Two.
	n++
	// Two lines fit a seven-line loop.
	// Two.
	for _, item := range items {
		if item > 0 {
			n++
		}
		n--
		n++
	}
	// Three lines don't. // want "in-body comment has 3 comment lines, allowed 2 \\(funcs.ratio\\)"
	// Two.
	// Three.
	for _, item := range items {
		if item > 0 {
			n++
		}
		n--
		n++
	}
	return n
}

func cases(x int) int {
	n := 0
	switch x {
	// Two lines fit a seven-line clause.
	// Two.
	case 1:
		n++
		n++
		n++
		n++
		n++
		n++
	// Two lines don't fit a one-line clause. // want "in-body comment has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
	// Two.
	case 2:
	}
	return n
}

func literals() [][]int {
	return [][]int{
		// Two lines fit a seven-line element.
		// Two.
		{
			1,
			2,
			3,
			4,
			5,
		},
		// Two lines don't fit a one-line element. // want "in-body comment has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
		// Two.
		{6},
	}
}

func dangling(ok bool) int {
	n := 0
	if ok {
		n++
		// A dangling comment is measured as one line, // want "in-body comment has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
		// not against the loop after its block.
	}
	for i := 0; i < 3; i++ {
		n++
		n++
		n++
		n++
		n++
	}
	return n
	// Nothing follows this one either. // want "in-body comment has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
	// Two.
}

func closures() func() int {
	return func() int {
		n := 0
		// A closure's comments belong to the enclosing func. // want "in-body comment has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
		// Two.
		n++
		return n
	}
}

func commentedOut() int {
	n := 0
	// for i := 0; i < 10; i++ { // want "in-body comment has 4 comment lines, allowed 1 \\(funcs.min-lines\\)"
	// 	n += i
	// 	n *= 2
	// }
	return n
}

func trailing() int {
	n := 0
	n++ // A single-line trailing comment always fits.
	n++ /* A multi-line trailing comment
	does not. // want "trailing comment has 2 comment lines, allowed 1 \\(funcs.min-lines\\)" */
	return n
}

type fields struct {
	a int // A single-line trailing comment always fits.
	b int /* A multi-line trailing comment
	does not. // want "trailing comment has 2 comment lines, allowed 1 \\(decls.min-lines\\)" */
	c func() /* An exempt field's trailing comment
	may be long. */
}

const (
	x = 1 // A single-line trailing comment always fits.
	y = 2 /* A multi-line trailing comment
	does not. // want "trailing comment has 2 comment lines, allowed 1 \\(decls.min-lines\\)" */
)

func localTypes() {
	type local interface {
		// Methods of local interfaces are contracts too.
		// Two.
		Do()
	}
}

func localContracts() {
	// Handler is a local func type, documented as a contract.
	// Two.
	type Handler func()
	// Alias is a local alias of a func type, also a contract.
	// Two.
	type Alias = func()
	type hooks struct {
		// OnDone is a func-typed field of a local struct.
		// Two.
		OnDone func()
		// Count is not a contract. // want "in-body comment has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
		// Two.
		Count int
	}
}

var packageClosure = func() int {
	n := 0
	// A package-level closure is checked as a func body. // want "in-body comment has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
	// Two.
	n++
	type local struct {
		// Its local struct fields are reported once. // want "in-body comment has 2 comment lines, allowed 1 \\(funcs.min-lines\\)"
		// Two.
		F int
	}
	return n
}
