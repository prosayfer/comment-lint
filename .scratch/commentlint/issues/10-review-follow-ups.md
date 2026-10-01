# 10: Review follow-ups

**What to build:** Fix what the review of tickets 01–09 (`b335373..a5cc1e5`) found. Items 1–4 are bugs against `docs/spec.md`. Items 5–7 need a decision from the author before any code changes. Items 8–9 are breaches of `AGENTS.md`, and item 10 lists optional clean-ups. Fix each bug test-first through `commentlint.New` + `analysistest` (see `docs/spec.md`, Testing Decisions).

Also: tickets 04 and 07 say `ceil(14 × 0.15) = 2`, but it is 3 (14 × 0.15 = 2.1). The implementation follows the spec formula, and the 2-line band under the defaults is 7–13 code lines. Correct the two tickets.

**Blocked by:** None (can start immediately)

**Status:** done (items 2 and 5–7 decided 2026-10-01)

## Bugs

1. **Contract exemptions inside function bodies.** Spec: named func types and func-typed struct fields are "Always exempt". In `checkInBody`, only comments inside an `InterfaceType` are skipped. Comments on a local `type H func()` or on a func-typed field of a local struct are linted as in-body comments.
   - [x] A long doc on a local `type H func()` produces no diagnostic
   - [x] A long doc on a func-typed field of a local struct produces no diagnostic
2. **Closures in package-level vars are never checked.** Spec: "comments inside them are in-body blocks of the enclosing function", and func-typed vars are linted. `checkInBody` requires the top-level decl to be a `FuncDecl`, so the comments in `var V = func() { … }` are never checked. **Decided:** `funcs`, as in-body blocks, with the literal as the enclosing function. Comments in package-level literals outside any func body stay unchecked. Record this in the spec.
   - [x] An oversized comment inside a package-level func literal is flagged
3. **Config keys match case-insensitively.** Spec: unknown keys return an error. `encoding/json` accepts `{"FUNCS": {"Max-Lines": 5}}` as a valid override.
   - [x] Keys that differ only in case make `New` return an error that names the key
4. **Code lines undercounted in CRLF files.** `codeLines` works out where a raw string ends from `len(lit)`, but the scanner strips `\r` from raw-string literals, so the computed end falls short in CRLF files. A func with a 13-line raw string gets `allowed 2 (funcs.ratio)` with CRLF and 3 with LF.
   - [x] A CRLF fixture gives the same diagnostics as its LF twin

## Decisions needed

5. **Which setting a tie names.** When the ratio term and the floor are equal (e.g. 3 code lines → `ceil(0.45) = 1` = min-lines 1), the diagnostic currently names `min-lines`. The spec says min-lines is named "when the floor raised the allowance", and a tie didn't raise it. Ticket 04's example, the README and the CI grep all rely on the current behaviour.
   - Option A: keep it, and reword the spec as "when the floor decided the allowance, including a tie".
   - Option B: name the ratio on a tie. Update ticket 04's example, the README, `testdata/integration` and CI.
   - **Decided: A.** The README already documents this. Only the spec sentence changes.
6. **Lone `*` lines in `/** … */` blocks.** A line holding only `*` counts as a comment line. Should it be excluded like empty `//` lines?
   - **Decided:** exclude it. Add it to the spec's line-counting list.
   - [x] A `/** … */` doc whose lone `*` lines would push it over the allowance produces no diagnostic
7. **Aliases of interfaces.** `type R = io.Reader` is exempt, because its underlying type is an interface. The spec lists "type aliases" under decls. Should it stay exempt as a contract, or be linted?
   - **Decided:** keep it exempt. The exemption follows the underlying type. Spec: interface and func types are exempt "whether defined or aliased". The "type aliases" in decls are aliases of other types.
   - [x] A long doc on `type H = func()` produces no diagnostic

## Standards (`AGENTS.md`: "inline a helper called from one place or a variable read once")

8. [x] Inline these helpers, each called from one place: `analyzer` into `run` in `commentlint_test.go`, and `isDirective` and `startsBefore` in `commentlint.go`.
9. [x] Inline the local `decls` in `run`, which is read once.

## Optional clean-ups (judgement calls)

10. - **Complexity is easy to forget:** a `rule` can have `complexityRatio > 0` while `complexity` stays unset, which quietly allows 0 lines. Pass complexity as an argument to `allowance` instead.
    - **Duplicated code:** the `ValueSpec` and `TypeSpec` cases in `checkSpec` have the same body, and `check(…Comment, "trailing comment", …, 1)` appears 4 times.
    - **Settings keys spelled twice:** once in the JSON tags and once in the negative-value list.
    - **Names:** `lineSet` actually maps each code line to its first code token, and `funcs` vs `funcDocs` doesn't say what differs.
    - **Done:** `funcDocs` is gone. The checker holds `complexityRatio` and sets it on the func-doc rule together with the complexity, in one statement. Settings keys are now spelled once, because the negative-value check walks the decoded settings. `lineSet` is now `codeLines`.
    - **Skipped:** merging the `ValueSpec`/`TypeSpec` cases. The bodies need each type's own `Doc`/`Comment` fields, so a merged version isn't shorter.
    - **Don't embed `limits` in `funcLimits`:** embedding changes config error paths to `funcs.limits.ratio`, which breaks the invalid-settings test.
