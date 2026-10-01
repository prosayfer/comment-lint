---
title: commentlint — comment size linter
status: ready-for-agent
---

# commentlint — comment size linter

## Problem Statement

AI coding agents write a lot of comments: several paragraphs of godoc on trivial helpers, `// increment counter` above `i++`, a comment before every step of a function body. The author wants self-explanatory code with no obvious or redundant comments. Nothing in the Go linting ecosystem limits how *big* comments are relative to the code they describe. Reviewers catch this by hand, inconsistently. Agents keep doing it because no automated feedback tells them to stop.

## Solution

`commentlint` is a Go analyzer that runs as a golangci-lint v2 module plugin and as a standalone binary. It measures every comment block, compares it with the code it describes, and reports any block that is larger than allowed.

The allowance for a block depends on the size of that code, a hard line cap, a guaranteed minimum, and (optionally) the cognitive complexity of the function. Contract-like declarations (interfaces, func types, callback fields) and package docs are exempt. There are two settings groups, one for functions and one for other declarations, so each can be tuned separately in `.golangci.yml`. The defaults are strict: an average function or declaration may have about one or two lines of comments.

## User Stories

1. As a Go developer, I want a linter that flags oversized comments, so that agent-written code stays free of redundant commentary.
2. As a Go developer, I want to enable the linter through golangci-lint's module plugin system, so that it runs alongside my other linters with one command.
3. As a Go developer, I want to configure the linter in `.golangci.yml` under the custom linter settings, so that all lint configuration lives in one place.
4. As a Go developer, I want a standalone binary, so that I can run the linter without building a custom golangci-lint.
5. As a Go developer, I want sensible strict defaults, so that I get value without writing any configuration.
6. As a Go developer, I want function doc comments limited by a maximum number of lines, so that no function carries an essay.
7. As a Go developer, I want function doc comments limited by a ratio to the function's code lines, so that small functions get proportionally small comments.
8. As a Go developer, I want a guaranteed minimum allowance (one line by default), so that every function and declaration can still have a one-line godoc.
9. As a Go developer, I want comments inside function bodies checked, so that step-by-step narration is caught.
10. As a Go developer, I want an in-body comment measured against the statement it precedes, so that a comment above a large loop may be longer than a comment above a single assignment.
11. As a Go developer, I want in-body comments before `case` clauses, composite literal elements and similar nodes measured against that node, so that every comment has a well-defined subject.
12. As a Go developer, I want a comment with nothing after it before a closing brace measured as if it described one line of code, so that dangling comments stay short.
13. As a Go developer, I want commented-out code treated like any other comment, so that large blocks of dead code get flagged.
14. As a Go developer, I want trailing (end-of-line) comments counted like other comment blocks, so that the rules stay consistent, knowing that a single-line trailing comment always fits the minimum allowance.
15. As a Go developer, I want each comment block checked on its own rather than summed per function, so that each diagnostic points at one specific block.
16. As a Go developer, I want struct types, named non-struct types, type aliases, variables, constants and struct fields in one "declarations" group, so that I tune data-shaped code together.
17. As a Go developer, I want functions, methods and `init` in one "functions" group, so that I tune behavior-shaped code together.
18. As a Go developer, I want a separate max-lines, ratio and minimum for each group, so that declarations and functions can have different budgets.
19. As a Go developer, I want each struct field and each spec in a `var`/`const`/`type` group to have its own budget, so that one heavily commented field doesn't hide behind a large struct.
20. As a Go developer, I want a struct's or group's doc comment measured against the whole struct or group (not counting the field comments), so that a large struct may carry a somewhat larger doc.
21. As a Go developer, I want interface types and their method comments exempt, so that I can document contracts as fully as they need.
22. As a Go developer, I want named func types exempt, so that callback and handler contracts can be documented freely.
23. As a Go developer, I want func-typed struct fields exempt, so that callback hooks can be documented like contracts.
24. As a Go developer, I want func-typed variables and fields whose type is an interface to be linted normally, so that the exemption stays narrow.
25. As a Go developer, I want package doc comments exempt, so that package overviews can be long prose.
26. As a Go developer, I want free-floating comments (license headers, section dividers not attached to anything) ignored, so that the linter only judges comments that describe code.
27. As a Go developer, I want compiler and tool directives (`//go:generate`, `//go:build`, `//nolint`, `//export`, `// +build`) excluded from line counts, so that required directives never count against my budget.
28. As a Go developer, I want empty `//` separator lines and bare `/*` / `*/` marker lines excluded from counts, so that only lines with actual text count.
29. As a Go developer, I want code lines counted as non-blank, non-comment lines, so that blank lines and comments can't inflate the allowance.
30. As a Go developer, I want generated files always skipped, so that I'm never asked to fix code I don't own.
31. As a Go developer, I want `_test.go` files always skipped, so that test narration isn't policed.
32. As a Go developer, I want exported and unexported identifiers treated the same, so that the rules stay simple and predictable.
33. As a Go developer, I want an optional ceiling based on cognitive complexity for function doc comments, so that simple functions get fewer comment lines even when they are long.
34. As a Go developer, I want the complexity ceiling off by default, so that I opt in only when I've decided it suits my codebase.
35. As a Go developer, I want the docs to suggest a complexity ratio at which a function at the common gocognit warning threshold (about 15) gets two lines, so that I have a well-reasoned starting value.
36. As a Go developer, I want the minimum allowance to still apply when the complexity ceiling is on, so that trivial exported functions can keep a one-line godoc.
37. As a Go developer, I want any limit set to zero to be disabled, so that I can turn off individual rules.
38. As a Go developer, I want all enabled limits applied together with the tightest one winning, so that each setting reliably caps the result.
39. As a Go developer, I want one diagnostic per offending block that names the block kind, the actual line count, the allowed count and the setting that decided the limit, so that I know exactly what to change.
40. As a Go developer, I want the linter never to auto-delete comments, so that the fix (usually shortening) stays a human or agent decision.
41. As a Go developer, I want golangci-lint's standard `//nolint:commentlint` suppression to work, so that I can make justified exceptions.
42. As a Go developer, I want invalid configuration (negative values, unknown keys, wrong types) to fail loudly at startup, so that a typo doesn't silently disable the linter.
43. As an AI coding agent, I want precise, actionable diagnostics, so that I can shorten my comments in one pass without guessing.
44. As a maintainer of the linter, I want every rule covered by fixture-based tests at the plugin's public constructor, so that I can refactor internals freely.
45. As a maintainer of the linter, I want CI to build a custom golangci-lint with the plugin, so that I know registration works with real golangci-lint.
46. As a user installing the plugin, I want the module to require only Go 1.24, so that I don't have to adopt a bleeding-edge toolchain.

## Implementation Decisions

**Identity and distribution**
- Linter name: `commentlint`. Module path: `github.com/prosayfer/comment-lint`. Public repository.
- Registers as a golangci-lint v2 module plugin through `github.com/golangci/plugin-module-register`, using the load mode that provides type info (func-type exemptions need it).
- Also ships a standalone command built on the `singlechecker` from `golang.org/x/tools/go/analysis`, with the same defaults.
- `go.mod` declares Go 1.24 and targets the latest golangci-lint v2.x available at implementation time.

**Public surface (the only seam)**
- A plugin constructor takes the raw settings golangci-lint decodes, validates them, applies defaults and returns the analyzer. This constructor is the only public entry point that tests use.
- Settings are decoded strictly: unknown keys, negative numbers and wrong types return an error from the constructor.

**Settings contract** (keys under the custom linter's `settings` in `.golangci.yml`)

| Group | Key | Default | Meaning |
|---|---|---|---|
| `funcs` | `max-lines` | 3 | Hard cap on the lines in any one block |
| `funcs` | `ratio` | 0.15 | Allowed comment lines per code line |
| `funcs` | `min-lines` | 1 | Guaranteed allowance (floor) |
| `funcs` | `complexity-ratio` | 0 (off) | Allowed comment lines per point of cognitive complexity; function doc comments only; suggested value 0.1 |
| `decls` | `max-lines` | 2 | Hard cap on the lines in any one block |
| `decls` | `ratio` | 1.0 | Allowed comment lines per code line (so a comment can't be longer than its code) |
| `decls` | `min-lines` | 1 | Guaranteed allowance (floor) |

- Doc and in-body blocks in the `funcs` group share one set of settings. There are no separate doc/body sub-settings.
- Setting any limit to zero disables it.

**Allowance formula** (per comment block)
- Allowed lines = min(max-lines, max(min-lines, min(ceil(code_lines × ratio), ceil(cognitive_complexity × complexity-ratio)))).
- The complexity term applies only to function doc comments and only when complexity-ratio > 0. Disabled terms drop out of the formula.
- ceil is used for rounding. A block violates the rules when its comment-line count exceeds the allowed lines.
- The limit that ended up binding (max-lines, ratio, complexity-ratio) is recorded so the diagnostic can name it. When the floor decided the allowance, including a tie with the ratio or complexity term, the diagnostic names min-lines.
- Cognitive complexity comes from the `github.com/uudashr/gocognit` package, computed on the function declaration. Its numbers must match gocognit's own reports.

**Block model**
- A *comment block* is one attached comment group: a doc comment, an in-body comment group, or a trailing comment.
- Blocks are evaluated on their own (not summed per declaration).
- *Doc blocks* are measured against their declaration's full extent:
  - a function: from its signature to its closing brace
  - a struct or grouped declaration: the whole construct, not counting its fields' or specs' comment lines
  - a field or a spec inside a group: that field or spec alone
- *In-body blocks* are measured against the next AST node they precede inside a function body: a statement, a case clause, a composite literal element, etc. With no following node before the closing brace, they're measured as one code line.
- *Trailing comments* are measured against the line they end. In practice they always fit within min-lines. They get no special treatment.
- Commented-out code is an ordinary comment.

**Line counting**
- A comment line has actual text. These don't count:
  - directive lines: `//go:` forms, `//nolint`, `//export`, `// +build`, `//line`
  - empty `//` lines
  - lines of a block comment that hold only `/*`, `*/` or `*` decoration (as in `/** … */` blocks)
- A code line is a non-blank line that isn't only a comment.

**Grouping and exemptions**
- `funcs` group: functions, methods, `init`. Closures are not separate declarations; comments inside them are in-body blocks of the enclosing function. A func literal in a package-level var counts as a function body, so its comments are `funcs` in-body blocks. Comments in package-level composite literals outside any func body aren't checked.
- `decls` group: struct types, named non-struct types, type aliases (other than the exempt ones below), vars, consts, struct fields.
- Always exempt:
  - interface types (including generic constraint interfaces) and the doc comments on their methods, whether defined or aliased
  - named func types, whether defined or aliased
  - struct fields whose type is a func type
  - these hold at package level and inside function bodies alike
  - package doc comments
  - free-floating comments not attached to any declaration or node
- Not exempt: func-typed variables, fields or vars whose type is an interface.
- Skipped files: generated files (standard `Code generated ... DO NOT EDIT.` marker) and `_test.go` files, always.
- Exported and unexported identifiers are treated the same.

**Diagnostics**
- One diagnostic per violating block, placed at the start of the block.
- The message states the block kind (func doc, in-body comment, decl doc, field comment, trailing comment), the actual comment-line count, the allowed count and the binding setting. Example shape: "func doc has 6 comment lines, allowed 3 (funcs.max-lines)".
- No suggested fixes.
- Suppression uses golangci-lint's standard `//nolint` handling. The linter itself doesn't implement suppression.

## Testing Decisions

- **A good test** checks only external behavior: given settings and a Go source fixture, which diagnostics appear, where, and with what message. Tests don't reach into internal counting, grouping or formula helpers, so those can be refactored freely.
- **One seam:** every behavioral test goes through the plugin's public constructor with a settings map (or none, for defaults), then runs the returned analyzer with `golang.org/x/tools/go/analysis/analysistest` over `testdata` fixture packages with `// want` expectations.
- **Coverage through that seam:**
  - defaults and each setting, including zero-disables
  - each rule (max-lines, ratio, min-lines floor, complexity ceiling) and which one binds
  - each boundary (exactly at the allowance, one over)
  - ceil rounding
  - each exemption (interfaces and their methods, func types, func-typed fields, package docs, floating comments, directives)
  - skipped files (generated, `_test.go`)
  - in-body attachment (statement, case clause, literal element, dangling before `}`)
  - struct vs field vs grouped-spec budgets
  - trailing comments
  - commented-out code
  - diagnostic wording
- **Invalid config:** constructor-level tests assert errors for unknown keys, negative values and wrong types.
- **Complexity:** fixtures include functions whose gocognit scores are known, covering the suggested ratio's bands (≤10 → 1 line, 11–20 → 2, 21+ → 3 when capped).
- **Integration smoke test (CI only):** a GitHub Actions job builds a custom golangci-lint binary with the plugin via `golangci-lint custom` and lints a small fixture module. It checks registration and settings plumbing, not the rules.
- **Approach:** TDD, red-green-refactor, one rule at a time.
- **Prior art:** none in this repo (it's new). Follow the standard `analysistest` fixture conventions from `golang.org/x/tools` and the module plugin example from golangci-lint's documentation.

## Out of Scope

- Semantic judgments ("this comment is redundant", "this restates the code"). Only size is measured.
- Auto-fixes or suggested fixes.
- Different limits for exported and unexported identifiers.
- Separate limits for doc vs in-body blocks within the functions group.
- Special handling of consecutive trailing comments.
- Linting test files or generated files, or making that configurable.
- Package docs, file-level floating comments and license headers.
- The legacy `.so` Go plugin mode and contributing upstream to golangci-lint.
- Support for golangci-lint v1.

## Further Notes

- Using the plugin needs golangci-lint v2. The author's local install (v1.54.2) must be upgraded.
- Why the suggested complexity ratio is 0.1: with ceil, cognitive complexity ≤10 gets 1 line, 11–20 gets 2 (so the common warning threshold of 15 sits in the middle of the 2-line band), and 21+ reaches the default 3-line cap.
- The declarations ratio of 1.0 mostly enforces "a comment may not be longer than its code". Fields and single-line specs rely on min-lines, and max-lines is the real cap.
