# commentlint

A Go linter that limits comment size. It keeps code self-explanatory and catches the long, redundant comments that AI agents tend to write.

It works as a [golangci-lint](https://golangci-lint.run) module plugin (v2) and as a standalone binary. The module requires Go 1.24.

## Installation

### golangci-lint module plugin

The plugin needs golangci-lint **v2**; v1 is not supported. Module plugins are compiled into a custom golangci-lint binary ([docs](https://golangci-lint.run/plugins/module-plugins/)).

1. Add `.custom-gcl.yml` to your repository:

   ```yaml
   version: v2.14.0 # the golangci-lint version to build
   plugins:
     - module: github.com/prosayfer/comment-lint
       version: latest # or a tagged version
   ```

2. Build the binary. This writes `./custom-gcl`:

   ```sh
   golangci-lint custom
   ```

3. Enable the linter in `.golangci.yml`:

   ```yaml
   version: "2"
   linters:
     enable:
       - commentlint
     settings:
       custom:
         commentlint:
           type: module
           description: Limits comment size relative to the code it describes.
           settings: # optional, see Settings
             funcs:
               complexity-ratio: 0.1
   ```

4. Run `./custom-gcl run ./...`.

### Standalone

```sh
go install github.com/prosayfer/comment-lint/cmd/commentlint@latest
commentlint ./...
```

The standalone binary always uses the default settings.

### Library

`commentlint.NewAnalyzer(commentlint.DefaultSettings())` returns a plain `*analysis.Analyzer`; adjust the `Settings` fields before passing them.

## Settings

Keys go under `linters.settings.custom.commentlint.settings`. Unknown keys, negative numbers and wrong types fail at startup. Setting any limit to `0` disables it.

| Group | Key | Default | Meaning |
|---|---|---|---|
| `funcs` | `max-lines` | 3 | Hard cap on the lines in any one block |
| `funcs` | `ratio` | 0.15 | Allowed comment lines per code line |
| `funcs` | `min-lines` | 1 | Guaranteed allowance (floor) |
| `funcs` | `complexity-ratio` | 0 (off) | Allowed comment lines per point of cognitive complexity; function doc comments only |
| `decls` | `max-lines` | 2 | Hard cap on the lines in any one block |
| `decls` | `ratio` | 1.0 | Allowed comment lines per code line (so a comment can't be longer than its code) |
| `decls` | `min-lines` | 1 | Guaranteed allowance (floor) |

A full configuration with every key at its default:

```yaml
version: "2"
linters:
  enable:
    - commentlint
  settings:
    custom:
      commentlint:
        type: module
        description: Limits comment size relative to the code it describes.
        settings:
          funcs:
            max-lines: 3
            ratio: 0.15
            min-lines: 1
            complexity-ratio: 0 # 0.1 suggested, see below
          decls:
            max-lines: 2
            ratio: 1.0
            min-lines: 1
```

`funcs` covers functions, methods and `init`, and applies to both their doc comments and the comments inside their bodies (including closures, also those assigned to package-level vars). `decls` covers struct types, other named types, type aliases, vars, consts and struct fields.

## Allowance formula

For each comment block:

```
allowed = min(max-lines, max(min-lines, min(ceil(code_lines × ratio), ceil(cognitive_complexity × complexity-ratio))))
```

- Disabled (zero) terms drop out. The complexity term applies only to function doc comments.
- `ceil` absorbs float error, so 20 × 0.15 gives 3, not 4.
- A block is reported when its comment lines exceed `allowed`.
- The diagnostic names the setting that decided the limit. On a tie, the later step in the formula wins: `min-lines` beats `ratio` and `complexity-ratio`, and `max-lines` beats everything.
- Cognitive complexity is computed with [gocognit](https://github.com/uudashr/gocognit) and matches its reports.

A **comment line** has actual text. Directive lines (`//go:…`, `//nolint`, `//export`, `// +build`, `//line`), empty `//` lines and lines that hold only `/*`, `*/` or `*` don't count. A **code line** is a non-blank line that isn't only a comment.

### Suggested complexity ratio: 0.1

With `ceil`, `complexity-ratio: 0.1` gives cognitive complexity ≤10 one line, 11–20 two lines, and 21+ reaches the default 3-line cap. gocognit's common warning threshold of 15 sits in the middle of the 2-line band. `min-lines` still applies, so a trivial function keeps its one-line godoc.

## What a block is measured against

Each comment block is checked on its own, not summed per declaration.

- **Doc comments** are measured against their declaration: a function from its signature to its closing brace; a struct or grouped declaration as a whole, not counting its fields' or specs' comment lines; a field or a spec inside a group on its own.
- **In-body comments** are measured against the node they precede inside the function body: a statement, a case clause, a composite literal element and so on. A dangling comment with nothing after it before the closing brace is measured as one code line.
- **Trailing comments** are measured as one code line, so a single-line trailing comment always fits.
- Commented-out code is an ordinary comment.

## Exemptions and skipped files

Never checked:

- interface types (including constraint interfaces) and the comments on their methods, whether defined or aliased
- named func types, whether defined or aliased
- struct fields whose type is a func type
- package doc comments
- free-floating comments not attached to a declaration or node, such as license headers

Func-typed variables, and fields or variables whose type is an interface, are checked normally.

Generated files (with the standard `Code generated … DO NOT EDIT.` marker) and `_test.go` files are always skipped.

## Diagnostics

One diagnostic per offending block, at the start of the block. It names the block kind (`func doc`, `in-body comment`, `decl doc`, `field comment`, `trailing comment`), the comment lines, the allowance and the binding setting:

```
example.go:3:1: func doc has 3 comment lines, allowed 1 (funcs.min-lines) (commentlint)
```

There are no automatic fixes. To make a justified exception under golangci-lint, put `//nolint:commentlint` as the first line of the block. Directive lines don't count towards the block's size.

## License

[MIT](LICENSE)
