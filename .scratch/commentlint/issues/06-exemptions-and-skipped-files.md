# 06: Exemptions and skipped files

**What to build:** These comments are never flagged, however long:
- contract-like declarations: interface types (including generic constraint interfaces) and the doc comments on their methods, named func types, struct fields whose type is a func type
- package doc comments
- free-floating comments not attached to any declaration or node (license headers, section dividers)

Generated files (`Code generated ... DO NOT EDIT.`) and `_test.go` files are skipped entirely.

**Blocked by:** 05

**Status:** ready-for-agent

- [ ] Long docs on an interface type and on its methods produce no diagnostics
- [ ] A long doc on a generic constraint interface produces no diagnostics
- [ ] A long doc on `type HandlerFunc func(...)` produces no diagnostics
- [ ] A long doc on a func-typed struct field produces no diagnostics
- [ ] A long doc on a func-typed var, and on a field whose type is an interface, is still flagged
- [ ] A long package doc and a long license header produce no diagnostics
- [ ] A generated file with oversized comments produces no diagnostics
- [ ] A `_test.go` file with oversized comments produces no diagnostics
