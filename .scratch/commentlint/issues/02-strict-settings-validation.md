# 02: Strict settings validation

**What to build:** Invalid configuration fails loudly. The plugin constructor returns an error for unknown keys, negative numbers and wrong value types, in both the `funcs` and `decls` groups. Setting any limit to `0` disables that limit.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] An unknown top-level key or group key (e.g. `funcs.max-line`) makes the constructor return an error that names the key
- [ ] A negative value for any limit returns an error
- [ ] A wrong type (e.g. a string for `ratio`) returns an error
- [ ] A partial config merges with the defaults (e.g. setting only `decls.max-lines` keeps every other default)
- [ ] `max-lines: 0` disables the max-lines cap, verified by a long doc comment producing no max-lines diagnostic
- [ ] Constructor-level tests cover each error case
