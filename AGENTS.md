## Working agreements

- Read `docs/spec.md` before changing what the analyzer reports or how it is configured or tested.
- Fix the root cause of a problem or design flaw.
- Design interfaces that are hard to misuse.
- Keep the smallest shape that works: inline a helper called from one place or a variable read once, drop a guard whose branch repeats another, and leave every case that needs no special handling to `default`. Extract on the second use, not in anticipation of one.
- Write self-explanatory code. Comment only the non-obvious *why*, in one or two lines.
- No tautological tests.