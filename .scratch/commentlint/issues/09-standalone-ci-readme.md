# 09: Standalone binary, CI and README

**What to build:** A standalone `commentlint` command built on the `singlechecker` from `golang.org/x/tools/go/analysis`, with the same defaults, for use without golangci-lint. A GitHub Actions workflow runs the test suite, then builds a custom golangci-lint v2 binary with the plugin via `golangci-lint custom` and lints a small example module to check that registration and settings reach the analyzer. The README documents installation (the golangci-lint v2 requirement, `.custom-gcl.yml`, `.golangci.yml`), the settings table, the allowance formula, the exemptions and the suggested `complexity-ratio: 0.1`.

**Blocked by:** 02

**Status:** ready-for-agent

- [ ] Running the standalone binary on a package prints the same diagnostics as the analyzer tests
- [ ] CI runs `go test ./...` on push and pull request
- [ ] CI builds the custom golangci-lint binary and lints the example module, and the job fails if an expected diagnostic is missing
- [ ] CI passes a non-default setting through `.golangci.yml` and checks that it changes the result
- [ ] The README covers installation, configuration, the formula, the exemptions and the suggested complexity ratio with its reasoning
