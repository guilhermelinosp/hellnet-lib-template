# golang-lib-template

> GitHub template for production-ready Go libraries, pre-configured with the
> canonical Hellnet pattern: CI, linting, pre-commit hooks, dependency
> automation, and an env-first API.

Click **"Use this template"** to scaffold a new Go library in seconds.

## What's included

- **Canonical Hellnet API** — the seeded example demonstrates the exact pattern
  every Hellnet library shares (hellnet-lib-kafka, hellnet-lib-cache,
  hellnet-lib-telemetry, hellnet-lib-database, hellnet-lib-api):
  - configuration via [hellnet-lib-environments]: every option is exposed as a
    `HELLNET_<LIB>_*` environment variable with a shared `HELLNET_*` fallback,
    and `.env` files load automatically (dev only, self-contained);
  - **constructors without `context.Context`** — `New`, `NewFromEnv`, `MustNew`;
    the runtime captures one `context.Background()` internally, so no public
    method takes a `ctx`;
  - options flow `DefaultOptions() → fromEnv(base) → validate()`;
  - versioning stays **v1 forever**: signature changes ship as minor/patch
    releases (never `BREAKING CHANGE`, never a major bump).
- **`.golangci.yml`** — curated linter config (errcheck, staticcheck, gosec, revive, …).
- **Lefthook** pre-commit hooks (`.lefthook.yml`): `go fmt`, `go vet`, `go mod tidy`,
  `golangci-lint`, `gitleaks`.
- **CI** (`.github/workflows`):
  - `pipeline.yml` (main): guard-semver (blocks auto-major) + semantic release +
    lib quality via [github.com/guilhermelinosp/templates].
  - `pr-check.yml` (PR): shellcheck, merge-check, gitleaks, labeler, quality.
  - `codeql.yml`: go + actions matrix.
- **Dependency automation** via Dependabot (`github-actions` + `gomod`).
- **Repo meta**: issue/PR/discussion templates, `CODEOWNERS`, `SECURITY.md`, `CONTRIBUTING.md`, `FUNDING.yml`.

## Quick start

```bash
# 1. create your repo from this template, then:
cd <repo>
go mod edit -module github.com/<you>/<repo>   # replace the module path
go mod edit -module=...                       # and rename the package dir if desired
go mod tidy
```

Rename the seeded API to your library:

1. Rename the package and `envPrefix` in `golanglibtemplate.go`
   (e.g. `HELLNET_KAFKA_`); keep the shared `HELLNET_` fallback.
2. Replace `Options`, `Greet` and the validation rule with your own API.
3. Keep `DefaultOptions()`, `fromEnv`, `New`/`MustNew` and `loadEnvFiles`
   — they are the canonical contract every Hellnet lib exposes.

## Usage (canonical pattern)

```go
package main

import (
	"fmt"

	"github.com/<you>/<repo>"
)

func main() {
	// env-first: HELLNET_<LIB>_* (or shared HELLNET_*), .env loaded for you
	c, err := <repo>.New()
	if err != nil {
		panic(err)
	}

	// or fail fast at startup
	c = <repo>.MustNew()

	msg, err := c.Greet("World")
	if err != nil {
		panic(err)
	}
	fmt.Println(msg) // Hello, World!
}
```

### Environment variables

| Variable | Purpose | Default |
|---|---|---|
| `HELLNET_TEMPLATE_NAME` / `HELLNET_NAME` | greeting target | `World` |
| `HELLNET_TEMPLATE_REPEATS` / `HELLNET_REPEATS` | repeat count | `1` |
| `HELLNET_TEMPLATE_VERBOSE` / `HELLNET_VERBOSE` | debug output | `false` |

`.env` files are loaded automatically by the constructors (conventional
`./.env`, plus parent-directory candidates, **dev environments only**) — no
external loader call needed.

## Develop

```bash
go fmt ./...    # format
go vet ./...    # vet
go test ./...   # tests (add -race for the race detector)
golangci-lint run
```

Install the git hooks once:

```bash
lefthook install
```

## Versioning

Releases and version bumps are derived from [Conventional Commits].
Hellnet libraries stay on **major v1**: a change to a public signature ships as
a minor or patch release — never a `BREAKING CHANGE`, never a major bump.

## License

[Apache 2.0](LICENSE)

[hellnet-lib-environments]: https://github.com/guilhermelinosp/hellnet-lib-environments
[github.com/guilhermelinosp/templates]: https://github.com/guilhermelinosp/templates
[Conventional Commits]: https://www.conventionalcommits.org/