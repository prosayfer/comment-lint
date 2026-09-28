# 08: Complexity ceiling

**What to build:** When `funcs.complexity-ratio` is above 0, a func doc comment is also capped at `ceil(cognitive_complexity × complexity-ratio)`, so simple funcs get fewer comment lines even if they are long. The formula becomes min(max-lines, max(min-lines, min(ceil(code × ratio), ceil(cognit × complexity-ratio)))). Cognitive complexity comes from the `github.com/uudashr/gocognit` package. The ceiling applies to func docs only and is off by default.

**Blocked by:** 04

**Status:** ready-for-agent

- [ ] With the default (0) the complexity term has no effect
- [ ] With `complexity-ratio: 0.1`, example functions with known gocognit scores behave like this: ≤10 → 1 line, 11–20 → 2, 21+ → 3 (capped by max-lines)
- [ ] A long but trivial func (cognit 0) still allows 1 line because of the `min-lines` floor
- [ ] When the complexity term decides the limit, the diagnostic names `funcs.complexity-ratio`
- [ ] In-body and trailing comments are not affected by the complexity ceiling
- [ ] The complexity values used match gocognit's own reports for the example functions
