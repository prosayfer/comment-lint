// Copyright header that is not attached to anything.
// Two.
// Three.
// Four.

// Package exempt has a long package doc.
// Two.
// Three.
// Four.
package exempt

// Reader is an interface, documented as fully as its contract needs.
// Two.
// Three.
// Four.
type Reader interface {
	// Read is a method whose doc is part of the contract.
	// Two.
	// Three.
	// Four.
	Read(p []byte) (int, error)
}

// Number is a generic constraint interface.
// Two.
// Three.
// Four.
type Number interface {
	~int | ~float64
}

// HandlerFunc is a named func type.
// Two.
// Three.
// Four.
type HandlerFunc func(n int) error

// ReaderAlias is an alias of an interface.
// Two.
// Three.
// Four.
type ReaderAlias = Reader

// FuncAlias is an alias of a func type.
// Two.
// Three.
// Four.
type FuncAlias = func()

// Contracts groups only exempt types.
// Two.
// Three.
// Four.
type (
	// Closer is an interface.
	// Two.
	// Three.
	Closer interface{ Close() error }
	// Callback is a func type.
	// Two.
	// Three.
	Callback func()
)

// ---- section divider, separated by a blank line ----
// Two.
// Three.

type Server struct {
	// OnDone is a func-typed callback field.
	// Two.
	// Three.
	OnDone func()
	// Handler has a named func type.
	// Two.
	// Three.
	Handler HandlerFunc
	// Source is an interface-typed field, so it's linted. // want "field comment has 2 comment lines, allowed 1 \\(decls.min-lines\\)"
	// Two.
	Source Reader
	// inline is an interface literal whose methods are exempt.
	inline interface {
		// Do is documented freely.
		// Two.
		// Three.
		Do()
	}
}

// hook is a func-typed var, so it's linted. // want "decl doc has 2 comment lines, allowed 1 \\(decls.min-lines\\)"
// Two.
var hook func()

// Mixed groups an exempt and a linted type. // want "decl doc has 3 comment lines, allowed 2 \\(decls.max-lines\\)"
// Two.
// Three.
type (
	Opener interface{ Open() }
	Size   int
)
