# 04: Full allowance formula for func docs

**What to build:** Func doc comments get the full allowance from the spec: allowed = min(max-lines, max(min-lines, ceil(code_lines × ratio))). Code lines are the non-blank, non-comment lines from the func signature to its closing brace. The diagnostic names the setting that decided the limit. When the floor raised the allowance, it names `min-lines`.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Blank lines and comment-only lines inside the func don't count as code lines
- [ ] With the default ratio 0.15, a 13-code-line func allows 2 lines (`ceil(13 × 0.15) = 2`), and a 2-line doc passes while a 3-line doc is flagged naming `funcs.ratio`
- [ ] A tiny func (e.g. 3 code lines) still allows 1 line because of the `min-lines` floor. A 2-line doc there is flagged naming `funcs.min-lines`
- [ ] A long func where the ratio exceeds max-lines is capped by `funcs.max-lines`, and the diagnostic names it
- [ ] `ratio: 0` removes the ratio term, so only max-lines and min-lines apply
- [ ] Boundary cases (exactly at the allowance, and one over) are covered for each limit
