# 01: Tracer bullet: plugin with a max-lines check on func docs

**What to build:** A Go module (`github.com/prosayfer/comment-lint`, Go 1.24) that registers `commentlint` as a golangci-lint v2 module plugin. Its constructor takes raw settings, fills in defaults and returns the analyzer. At this point it checks only one thing: a func or method doc comment with more comment lines than `funcs.max-lines` (default 3) gets a diagnostic. This ticket builds the test harness everything later relies on: build through the public constructor, run `analysistest` on `testdata` example files, and assert with `// want`. See `docs/spec.md`.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] The module registers with golangci-lint through `plugin-module-register` under the name `commentlint`, using the load mode that provides type info
- [ ] The constructor accepts a nil or empty settings map and applies every default from the spec (funcs: 3 / 0.15 / 1 / 0; decls: 2 / 1.0 / 1). Strict validation is not needed yet
- [ ] A func or method doc comment with more lines than `funcs.max-lines` gets exactly one diagnostic, placed at the start of the comment
- [ ] A doc comment exactly at the limit gets no diagnostic
- [ ] The diagnostic wording matches the spec, e.g. "func doc has 6 comment lines, allowed 3 (funcs.max-lines)"
- [ ] A settings override of `funcs.max-lines` changes the result, verified through the constructor
- [ ] All tests go through the public constructor and `analysistest`; nothing tests internal helpers directly
