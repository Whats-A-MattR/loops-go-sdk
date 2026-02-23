# Contributing to loops-go-sdk

Thanks for your interest in contributing. This is a community-maintained SDK for the [Loops API](https://loops.so/docs/api).

## How to contribute

- **Bug reports and feature requests:** Open a [GitHub issue](https://github.com/Whats-A-MattR/loops-go-sdk/issues).
- **Code changes:** Open a pull request (PR) against `main`. Follow the guidelines below so your PR can be merged.

## Before you submit a PR

1. **Fork and clone** the repo, then create a branch for your change.

2. **Run the test suite locally:**
   ```bash
   go build ./...
   go test ./... -race -count=1
   go vet ./...
   ```
   All tests must pass. Optional: run fuzzing for a short time:
   ```bash
   go test ./... -fuzz=FuzzClientResponse -fuzztime=20s -count=1
   ```

3. **Keep the SDK aligned with the Loops OpenAPI spec.** Types and endpoints should match [the spec](https://app.loops.so/openapi.json). If you add or change endpoints, update `openapi.json` and ensure `TestOpenAPI_SDKEndpointsExistInSpec` (and any new tests) pass.

4. **Add tests** for new behavior. Prefer:
   - Happy-path tests that assert request/response shape (spec compliance).
   - Negative tests for validation and error handling.
   - Fuzz targets where random input could catch panics or bad error handling.

5. **Code style:** Follow standard Go conventions (`gofmt`, `go vet`). The project uses Go 1.21+.

## PR checks

- CI runs on every PR: build, tests (including race detector), fuzzing, and `go vet`. The PR must be green before merge.
- Maintainers may request changes or ask for clarification.

## Branch protection

The default branch is protected: merges go through PRs, and status checks (CI) must pass. See [.github/BRANCH_PROTECTION.md](.github/BRANCH_PROTECTION.md) for details.

## Questions

Open an issue with the “question” label or reach out to the maintainers.
