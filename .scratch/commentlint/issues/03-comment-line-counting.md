# 03: Comment line counting rules

**What to build:** Only lines with real text count as comment lines. These are excluded from a block's size:
- directive lines (`//go:` forms, `//nolint`, `//export`, `// +build`, `//line`)
- empty `//` separator lines
- block-comment lines that hold only `/*` or `*/`

Checked through the max-lines results on func docs.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] A func doc of 3 text lines plus a `//go:generate` line and a `//nolint` line is not flagged at `max-lines: 3`
- [ ] Empty `//` lines between paragraphs are not counted
- [ ] In a `/* ... */` doc comment, lines holding only `/*` or `*/` are not counted, and lines with text are
- [ ] A single-line `/* text */` comment counts as one line
- [ ] The diagnostic's reported count reflects the excluded lines
