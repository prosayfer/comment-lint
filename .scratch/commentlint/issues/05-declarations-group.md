# 05: Declarations group

**What to build:** Doc comments on struct types, named non-struct types, type aliases, vars, consts and struct fields are checked with the `decls.*` settings (defaults: max-lines 2, ratio 1.0, min-lines 1), using the same formula as func docs.
- Each struct field and each spec in a `var`/`const`/`type` group gets its own budget.
- A struct's or group's doc comment is measured against the whole struct or group, not counting its fields' or specs' comment lines.

Interfaces are not part of this ticket; they are exempted in 06.

**Blocked by:** 04

**Status:** ready-for-agent

- [ ] A 3-line doc on a single-line var or const is flagged, naming `decls.max-lines` or `decls.min-lines`, whichever decided the limit
- [ ] A 2-line doc on a 1-line field is flagged naming `decls.min-lines`, and a 1-line doc passes
- [ ] A struct's doc is measured against the struct's code lines, and the fields' own comment lines don't add to them
- [ ] Each spec in a grouped `const ( ... )` is checked on its own, and the group's doc is checked against the whole group
- [ ] Named non-struct types (`type Status int`) and aliases are checked with `decls.*`
- [ ] Diagnostics use the right block kind ("decl doc", "field comment")
