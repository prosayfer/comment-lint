# 07: In-body and trailing comments

**What to build:** Comments inside function bodies are checked on their own, with the `funcs.*` settings. Each in-body block is measured against the next AST node it precedes: a statement (including an entire `if`/`for`/`switch`), a `case` clause, or a composite literal element. A block with nothing after it before the closing `}` is measured as one code line. Comments inside closures are in-body blocks of the enclosing func. Commented-out code is an ordinary comment. Trailing (end-of-line) comments are checked against the line they end, on both funcs and decls.

**Blocked by:** 04, 05

**Status:** ready-for-agent

- [ ] A 2-line comment above a single `i++` is flagged, and a 1-line comment passes
- [ ] A 2-line comment above a 13-code-line `for` loop passes (`ceil(13 × 0.15) = 2`), and a 3-line one is flagged
- [ ] A comment before a `case` clause is measured against that clause
- [ ] A comment inside a composite literal is measured against the element that follows it
- [ ] A multi-line comment right before `}` with nothing after it is flagged under `min-lines`
- [ ] A comment inside a closure is checked with the enclosing func's settings
- [ ] A large block of commented-out code is flagged
- [ ] Single-line trailing comments on statements and fields pass under the defaults
- [ ] Diagnostics use the right block kind ("in-body comment", "trailing comment")
